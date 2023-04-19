package manager

import (
	"context"
	"errors"
	"os"
	"path"
	"strings"
	"time"

	"pa.cn/fedlearn/psi/api/types"
	"pa.cn/fedlearn/psi/client/httpc"
	"pa.cn/fedlearn/psi/config"
	"pa.cn/fedlearn/psi/log"
	"pa.cn/fedlearn/psi/model"
	"pa.cn/fedlearn/psi/utils"
)

func (mgr *UploaderManager) FilePull(ctx context.Context, r types.FilePullRequest) error {
	log.Debugln(ctx, "UploaderManager.FilePull is called")
	if r.URL == "" {
		return errors.New("Empty file url")
	}
	if !strings.HasPrefix(r.URL, "http") && !strings.HasPrefix(r.URL, "https") {
		return errors.New("Only support http/https when download dataset file")
	}

	validDays := r.ValidDays
	if validDays == 0 {
		validDays = config.GetDatasetValidDays()
	}
	expiredDate := time.Now().AddDate(0, 0, int(validDays))

	dataset := &model.Dataset{
		Name:        r.Name,
		Desc:        r.Description,
		Type:        r.Type,
		URL:         r.URL,
		Status:      model.DatasetDownloading,
		Index:       int32(0),
		BizContext:  r.BizContext,
		ExpiredDate: expiredDate,
	}

	if err := model.AddDataset(dataset); err != nil {
		return err
	}

	dataset, err := model.GetDatasetByNameAndIndex(r.Name, int32(0))
	if err != nil {
		return err
	}

	// go download(ctx, dataset)
	GetHttpDownloader().DownloadDataset(ctx, dataset)

	return nil
}

func download(ctx context.Context, dataset *model.Dataset) error {
	var md5Str string
	var targetPath string
	var err error
	if config.GetConfig().DataSet.DownMethod > 0 {
		md5Str, targetPath, err = downloadBigMemory(ctx, dataset)
	} else {
		md5Str, targetPath, err = downloadCopy(ctx, dataset)
	}

	if err != nil {
		return err
	}

	dataset.Md5 = md5Str

	lineCount, fileSize, _, err := utils.FileMetaInfo(targetPath)
	if err != nil {
		log.Errorf(ctx, "Download %s filemetainfo failed:%v\n", dataset.URL, err)
		dataset.Status = model.DatasetDownloadFailed
		model.UpdateDataset(dataset)
		return err
	}
	dataset.Count = lineCount
	dataset.Size = fileSize
	dataset.Status = model.DatasetOK

	if err := model.UpdateDataset(dataset); err != nil {
		log.Errorf(ctx, "Download %s update db failed:%v\n", dataset.URL, err)
		return err
	}

	if err := ShardDataset(ctx, dataset); err != nil {
		log.Errorf(ctx, "Sharding dataset %s error:%v", dataset.Name, err)
		return err
	}
	return nil
}

func downloadBigMemory(ctx context.Context, dataset *model.Dataset) (md5 string, filepath string, err error) {
	headers := map[string]string{
		"User-Agent": "psi client",
	}

	data, err := httpc.DoGetBig(dataset.URL, headers)
	if err != nil {
		log.Errorf(ctx, "Download %s failed:%v\n", dataset.URL, err)
		dataset.Status = model.DatasetDownloadFailed
		model.UpdateDataset(dataset)
		return "", "", err
	}

	md5Str := utils.MD5(data)

	targetPath := DatasetPath(md5Str, dataset.Name, dataset.Index)
	if err = saveToFile(ctx, targetPath, data); err != nil {
		log.Errorf(ctx, "Download %s save failed:%v\n", dataset.URL, err)
		dataset.Status = model.DatasetDownloadFailed
		model.UpdateDataset(dataset)
		return "", "", err
	}

	return md5Str, targetPath, nil
}

func downloadCopy(ctx context.Context, dataset *model.Dataset) (md5 string, filepath string, err error) {
	headers := map[string]string{
		"User-Agent": "psi client",
	}

	tmpTargetFilePath := DatasetPathWithoutMD5(dataset.Name, dataset.Index)
	if err := httpc.DownloadAndSaveFile(dataset.URL, headers, tmpTargetFilePath); err != nil {
		log.Errorf(ctx, "Download %s failed:%v\n", dataset.URL, err)
		dataset.Status = model.DatasetDownloadFailed
		model.UpdateDataset(dataset)
		return "", "", err
	}

	md5Str := utils.MD5OfFile(tmpTargetFilePath)

	targetPath := DatasetPath(md5Str, dataset.Name, dataset.Index)

	baseDir := path.Dir(targetPath)
	if err := os.MkdirAll(baseDir, 0776); err != nil {
		log.Errorf(ctx, "Download %s mkdir %s failed:%v", dataset.URL, baseDir, err)
		dataset.Status = model.DatasetDownloadFailed
		model.UpdateDataset(dataset)
		return "", "", err
	}

	if err := os.Rename(tmpTargetFilePath, targetPath); err != nil {
		log.Errorf(ctx, "Download %s rename failed:%v\n", dataset.URL, err)
		dataset.Status = model.DatasetDownloadFailed
		model.UpdateDataset(dataset)
		return "", "", err
	}
	return md5Str, targetPath, nil
}
