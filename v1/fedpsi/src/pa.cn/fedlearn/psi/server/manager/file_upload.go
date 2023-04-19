package manager

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"pa.cn/fedlearn/psi/config"
	"pa.cn/fedlearn/psi/log"

	"pa.cn/fedlearn/psi/api/types"
	"pa.cn/fedlearn/psi/model"
	"pa.cn/fedlearn/psi/utils"
)

type UploaderManager struct {
}

func (mgr *UploaderManager) CheckChunk(ctx context.Context, r types.ChunkCheckRequest) (bool, error) {
	log.Debugln(ctx, "UploaderManager.CheckChunk is called")
	chunkPath := ChunkPath(r.MD5, r.Chunk)
	f, err := os.Lstat(chunkPath)
	if err != nil {
		log.Warningf(ctx, "Chunk %d of %s does not exist.\n", r.Chunk, r.MD5)
		return false, err
	}
	if f.Size() != r.ChunkSize {
		log.Warningf(ctx, "Chunk %d of %s file size does not match.\n", r.Chunk, r.MD5)
		return false, err
	}
	return true, nil
}

func (mgr *UploaderManager) UploadChunk(ctx context.Context, r types.ChunkUploadRequest) error {
	log.Debugln(ctx, "UploaderManager.UploadChunk is called")
	chunkPath := ChunkPath(r.MD5, r.Chunk)
	baseDir := path.Dir(chunkPath)
	if err := os.MkdirAll(baseDir, 0776); err != nil {
		log.Errorf(ctx, "Failed to mkdir %s:%v", baseDir, err)
		return err
	}
	if len(r.FileData) == 0 {
		log.Errorf(ctx, "Upload empty chunk %s:%d\n", r.MD5, r.Chunk)
		return fmt.Errorf("Empty data of chunk %s:%d", r.MD5, r.Chunk)
	}
	if err := ioutil.WriteFile(chunkPath, r.FileData, 0776); err != nil {
		log.Errorf(ctx, "Error to write chunk %s:%v", chunkPath, err)
		return err
	}
	return nil
}

func (mgr *UploaderManager) MergeChunk(ctx context.Context, r types.ChunkMergeRequest) error {
	log.Debugln(ctx, "UploaderManager.MergeChunk is called")
	validDays := r.ValidDays
	if validDays == 0 {
		validDays = config.GetDatasetValidDays()
	}
	expiredDate := time.Now().AddDate(0, 0, int(validDays))
	dataset := &model.Dataset{
		Name:        r.Name,
		Desc:        r.Description,
		Type:        r.Type,
		Index:       int32(0), // index 0 means all data here, not shards
		BizContext:  r.BizContext,
		ExpiredDate: expiredDate,
	}

	allData := make([]byte, 0)
	if r.ChunkCount == 0 {
		r.ChunkCount = 1
	}
	for i := 0; i < int(r.ChunkCount); i++ {
		chunkData, err := ioutil.ReadFile(ChunkPath(r.MD5, int32(i)))
		if err != nil {
			log.Errorf(ctx, "Read chunk %s:%d failed:%v", r.MD5, i, err)
			return err
		}
		allData = append(allData, chunkData...)
	}

	md5Str := utils.MD5(allData)
	if md5Str != r.MD5 {
		log.Errorf(ctx, "Md5 does not match, expect:%s, but %s\n", r.MD5, md5Str)
		return fmt.Errorf("Md5 does not match: %s", r.MD5)
	}
	targetPath := DatasetPath(r.MD5, r.Name, dataset.Index)
	if err := saveToFile(ctx, targetPath, allData); err != nil {
		log.Errorf(ctx, "Write file %s error %v", targetPath, err)
		return err
	}

	removeChunkDirPath := ChunkDirPath(r.MD5)
	if err := os.Rename(removeChunkDirPath, fmt.Sprintf("delete_%s", removeChunkDirPath)); err != nil {
		log.Warningf(ctx, "Rename %s failed %v \n", removeChunkDirPath, err)
	}

	dataset.Md5 = md5Str
	lineCount, _, _, err := utils.FileMetaInfo(targetPath)
	if err != nil {
		log.Errorf(ctx, "Error to get %s file lines:%v", targetPath, err)
		return err
	}
	dataset.Count = lineCount
	dataset.Size = int64(len(allData))
	dataset.Status = model.DatasetOK
	if err := model.AddDataset(dataset); err != nil {
		return err
	}
	go func() {
		if err := ShardDataset(ctx, dataset); err != nil {
			log.Errorf(ctx, "Sharding dataset %s error:%v", dataset.Name, err)
		}
	}()
	return nil
}

func saveToFile(ctx context.Context, targetPath string, data []byte) error {
	baseDir := path.Dir(targetPath)
	if err := os.MkdirAll(baseDir, 0776); err != nil {
		log.Errorf(ctx, "Failed to mkdir %s:%v", baseDir, err)
		return err
	}
	outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY, 0776)
	if err != nil {
		log.Errorf(ctx, "Open file %s failed:%v", targetPath, err)
		return err
	}
	defer outFile.Close()
	if _, err = outFile.Write(data); err != nil {
		log.Errorf(ctx, "Write file %s failed:%v", targetPath, err)
		return err
	}
	return nil
}
