package manager

import (
	"fmt"
	"path/filepath"

	"pa.cn/fedlearn/psi/config"
	"pa.cn/fedlearn/psi/model"
)

func GetPeeringMode(mode string) string {
	if mode == model.TaskRunMode_Server {
		return model.TaskRunMode_Client
	}
	return model.TaskRunMode_Server
}

func ChunkPath(md5 string, chunk int32) string {
	return filepath.Join(config.GetConfig().DataSet.Dir, TMPDIR, md5, fmt.Sprintf("chunk_%d", chunk))
}

func ChunkDirPath(md5 string) string {
	return filepath.Join(config.GetConfig().DataSet.Dir, TMPDIR, md5)
}

func DatasetPathWithoutMD5(name string, index int32) string {
	return filepath.Join(config.GetConfig().DataSet.Dir, TMPDIR, fmt.Sprintf("%s_%d_tmp", name, index))
}

func DatasetPath(md5 string, name string, index int32) string {
	return filepath.Join(config.GetConfig().DataSet.Dir, md5, fmt.Sprintf("%s_%d", name, index))
}
