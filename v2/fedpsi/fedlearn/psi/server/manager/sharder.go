package manager

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"fedlearn/psi/common/config"
	"fedlearn/psi/common/log"
	"fedlearn/psi/model"
)

type shardTask struct {
	ctx context.Context
	ds  *model.Dataset
}

type Sharder struct {
	shardTaskPool chan *shardTask
}

var (
	sharder *Sharder = &Sharder{
		shardTaskPool: make(chan *shardTask, 64),
	}
)

func GetSharder() *Sharder {
	return sharder
}

func ShardDataset(ctx context.Context, ds *model.Dataset) error {
	sharder.shardTaskPool <- &shardTask{
		ctx: ctx,
		ds:  ds,
	}
	return nil
}

func (s *Sharder) Run(ctx context.Context) {
	log.Infoln(ctx, "Sharder start to run")
	paralism := config.GetConfig().DataSet.Sharders
	if paralism == 0 {
		paralism = 1
	}

	for i := 1; i <= paralism; i++ {
		log.Infof(ctx, "Sharder %d/%d is running...", i, paralism)
		go func() {
			for {
				t := <-sharder.shardTaskPool
				sharder.shardDataset(t.ctx, t.ds)
			}
		}()
	}
}

// ShardDataset will split the dataset if the dataset's count > MaxLines
func (s *Sharder) shardDataset(ctx context.Context, ds *model.Dataset) error {
	log.Debugln(ctx, "ShardDataset is called")
	if ds.Status != model.DatasetOK {
		return errors.New("Dataset not ready")
	}

	if config.GetConfig().DataSet.MaxLines <= 0 || ds.Count <= config.GetConfig().DataSet.MaxLines {
		ds.Status = model.DatasetAvailable
		ds.Shards = 1
		if err := model.UpdateDataset(ds); err != nil {
			log.Errorf(ctx, "Update dataset %s status to availabe error:%v", ds.Name, err)
			return fmt.Errorf("Update dataset %s status to availabe error:%w", ds.Name, err)
		}
		return nil
	}

	ds.Status = model.DatasetSharding
	if err := model.UpdateDataset(ds); err != nil {
		log.Errorf(ctx, "Update dataset %s status to sharding error:%v", ds.Name, err)
		return fmt.Errorf("Update dataset %s status to sharding error:%w", ds.Name, err)
	}

	filePath := DatasetPath(ds.Md5, ds.Name, 0)
	f, err := os.Open(filePath)
	if err != nil {
		log.Errorf(ctx, "Open file %s error:%v", filePath, err)
		return fmt.Errorf("Open file %s error:%w", filePath, err)
	}
	defer f.Close()
	var lineCount int64 = 0
	var shardIndex int32 = 1
	br := bufio.NewReader(f)
	lines := make([]byte, 0)
	for {
		line, _, err := br.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		lines = append(lines, line...)
		lines = append(lines, '\n')
		lineCount++
		if lineCount == config.GetConfig().DataSet.MaxLines {
			sharder.saveShard(ctx, lines, shardIndex, lineCount, *ds)
			shardIndex = shardIndex + 1
			lineCount = 0
			lines = make([]byte, 0)
		}
	}
	if len(lines) > 0 { // save the last data
		sharder.saveShard(ctx, lines, shardIndex, lineCount, *ds)
	} else {
		// backout this shardIndex if no lines
		shardIndex -= 1
	}

	ds.Shards = shardIndex
	ds.Status = model.DatasetAvailable
	if err := model.UpdateDataset(ds); err != nil {
		log.Errorf(ctx, "Shards dataset %s error:%v", ds.Name, err)
		return fmt.Errorf("Shards dataset %s error:%w", ds.Name, err)
	}
	return nil
}

func (s *Sharder) saveShard(ctx context.Context, lines []byte, shardIndex int32, count int64, ds model.Dataset) error {
	shardFilePath := DatasetPath(ds.Md5, ds.Name, shardIndex)
	ds.Size = int64(len(lines))
	err := ioutil.WriteFile(shardFilePath, lines, 0776)
	if err != nil {
		log.Errorf(ctx, "Save dataset shards %s-%d error:%v", ds.Name, shardIndex, err)
		return err
	}
	ds.Id = 0
	ds.Index = shardIndex
	ds.Count = count
	//ds.Size = int64(cap(lines))
	ds.Status = model.DatasetOK
	// ds.Path = shardFilePath
	if err := model.AddDataset(&ds); err != nil {
		log.Errorf(ctx, "Add dataset shard %s-%d error:%v", ds.Name, shardIndex, err)
		return fmt.Errorf("Add dataset shard %s-%d error:%w", ds.Name, shardIndex, err)
	}
	return nil
}
