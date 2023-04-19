package manager

import (
	"fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"pa.cn/fedlearn/psi/model"
)

func saveShard(lines []byte, shardIndex int32, count int64, ds model.Dataset) error {
	shardFilePath := DatasetPath(ds.Md5, ds.Name, shardIndex)
	ds.Size = int64(len(lines))
	err := ioutil.WriteFile(shardFilePath, lines, 0776)
	if err != nil {
		logrus.Errorf("Save dataset shards %s-%d error:%v", ds.Name, shardIndex, err)
		return err
	}
	ds.Id = 0
	ds.Index = shardIndex
	ds.Count = count
	//ds.Size = int64(cap(lines))
	ds.Status = model.DatasetOK
	// ds.Path = shardFilePath
	if err := model.AddDataset(&ds); err != nil {
		logrus.Errorf("Add dataset shard %s-%d error:%v", ds.Name, shardIndex, err)
		return fmt.Errorf("Add dataset shard %s-%d error:%w", ds.Name, shardIndex, err)
	}
	return nil
}
