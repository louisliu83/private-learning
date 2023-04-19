package manager

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"fedlearn/psi/api/types"
	"fedlearn/psi/common/config"
	"fedlearn/psi/common/log"
	"fedlearn/psi/common/utils"
	"fedlearn/psi/model"
)

func (mgr *UploaderManager) FileLocalCopy(ctx context.Context, r types.FileLocalCopyRequest) error {
	log.Debugln(ctx, "UploaderManager.FileLocalCopy is called")
	if r.FilePath == "" {
		return errors.New("Empty file path")
	}
	// check wether file path exists

	validDays := r.ValidDays
	if validDays == 0 {
		validDays = config.GetDatasetValidDays()
	}
	expiredDate := time.Now().AddDate(0, 0, int(validDays))

	dataset := &model.Dataset{
		Name:        r.Name,
		Desc:        r.Description,
		Type:        r.Type,
		Index:       int32(0),
		BizContext:  r.BizContext,
		ExpiredDate: expiredDate,
	}

	if err := model.AddDataset(dataset); err != nil {
		log.Errorf(ctx, "add dataset %s into db failed:%v\n", r.Name, err)
		return fmt.Errorf("add dataset %s into db failed:%w\n", r.Name, err)
	}

	go func() error {
		md5Str := utils.MD5OfFile(r.FilePath)
		dataset.Md5 = md5Str

		lineCount, fileSize, _, err := utils.FileMetaInfo(r.FilePath)
		if err != nil {
			log.Errorf(ctx, "get metadata of file %s failed:%v\n", r.FilePath, err)
			dataset.Status = model.DatasetCopyFailed
			model.UpdateDataset(dataset)
			return fmt.Errorf("get metadata of file %s failed:%w\n", r.FilePath, err)
		}

		dataset.Count = lineCount
		dataset.Size = fileSize
		dataset.Status = model.DatasetOK

		fileReader, err := os.Open(r.FilePath)
		if err != nil {
			log.Errorf(ctx, "Open %s failed:%v\n", r.FilePath, err)
			dataset.Status = model.DatasetCopyFailed
			model.UpdateDataset(dataset)
			return err
		}
		targetPath := DatasetPath(md5Str, r.Name, dataset.Index)
		baseDir := path.Dir(targetPath)
		if err = os.MkdirAll(baseDir, os.ModePerm); err != nil {
			log.Errorf(ctx, "Make dir  %s error %v", baseDir, err)
			return err
		}
		fileWriter, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR, 0776)
		if err != nil {
			log.Errorf(ctx, "Open %s failed:%v\n", targetPath, err)
			dataset.Status = model.DatasetCopyFailed
			model.UpdateDataset(dataset)
			return err
		}
		defer fileReader.Close()
		defer fileWriter.Close()

		fileBuf := make([]byte, 256*1024*1024)
		if _, err := io.CopyBuffer(fileWriter, fileReader, fileBuf); err != nil {
			log.Errorf(ctx, "Copyfile from %s to %s failed:%v\n", r.FilePath, targetPath, err)
			if err != io.EOF {
				dataset.Status = model.DatasetCopyFailed
				model.UpdateDataset(dataset)
				return fmt.Errorf("Copyfile from %s to %s failed:%v\n", r.FilePath, targetPath, err)
			}
		}

		if err := model.UpdateDataset(dataset); err != nil {
			log.Errorf(ctx, "update dataset %s into db failed:%v\n", r.Name, err)
			return fmt.Errorf("update dataset %s into db failed:%w\n", r.Name, err)
		}

		dataset, err = model.GetDatasetByNameAndIndex(r.Name, int32(0))
		if err != nil {
			log.Errorf(ctx, "load dataset %s from db failed:%v\n", r.Name, err)
			return fmt.Errorf("load dataset %s from db failed:%w\n", r.Name, err)
		}

		if err := ShardDataset(ctx, dataset); err != nil {
			log.Errorf(ctx, "Sharding dataset %s error:%v", dataset.Name, err)
			return err
		}
		return nil
	}()

	return nil
}
