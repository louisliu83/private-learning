package scheduler

import (
	"time"

	"pa.cn/fedlearn/psi/model"

	"pa.cn/fedlearn/psi/log"
	"pa.cn/fedlearn/psi/server/manager"
)

/*
Dataset scheduler check if the dataset exceed its valid time.
If true: expired it and delete all the dataset

If the dataset's expired date is blank, set it is 100 years later

*/

type DatasetScheduler struct {
	dm       *manager.DataSetMgr
	tickChan chan struct{}
}

var _ Scheduler = NewDatasetScheduler()

func NewDatasetScheduler() *DatasetScheduler {
	return &DatasetScheduler{
		dm:       manager.GetDatasetMgr(),
		tickChan: make(chan struct{}, 1024*1024),
	}
}

func (s *DatasetScheduler) Tick() {
	// Server how to notify the scheduler
	s.tickChan <- struct{}{}
}

func (s *DatasetScheduler) Start() {

	go func() {
		for {
			s.tickChan <- struct{}{}
			sleepOneMinute()
		}
	}()

	go func() {
		for {
			<-s.tickChan // check whether has the tick
			s.process()
		}
	}()
}

func sleepOneMinute() {
	time.Sleep(time.Duration(1) * time.Minute)
}

func sleepUntilNext0() {
	now := time.Now()
	next := now.Add(time.Duration(24) * time.Hour)
	next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
	time.Sleep(next.Sub(now))
}

func (s *DatasetScheduler) process() {
	ctx := contexForScheduler()
	log.Debugf(ctx, "DatasetScheduler run: %v\n", time.Now())
	datasetList, err := model.ListDataset()
	if err != nil {
		log.Warningln(ctx, "No dataset found for dataset scheduler to process")
		return
	}

	for _, d := range datasetList {
		if d.ExpiredDate.Before(model.GetAnchorTime()) {
			d.ExpiredDate = time.Now().AddDate(100, 0, 0)
			if err := model.UpdateDataset(d); err != nil {
				log.Errorf(ctx, "Set expired date of dataset %s error %v", d.Name, err)
			} else {
				log.Infof(ctx, "Set expired date of dataset %s as %v", d.Name, d.ExpiredDate)
			}
		} else {
			if !d.IsValid() {
				d.Status = model.DatasetExpired
				if err := model.UpdateDataset(d); err != nil {
					log.Errorf(ctx, "Set expired status of dataset %s error %v", d.Name, err)
				} else {
					log.Infof(ctx, "Set expired status of dataset %s as %v", d.Name, d.Status)
				}
				shards, err := model.GetDatasetShards(d.Name)
				if err != nil {
					log.Warningln(ctx, "No shards found for dataset %s error %v", d.Name, err)
				}

				for _, shard := range shards {
					if shard.Index != 0 {
						shard.Status = model.DatasetExpired
						if err := model.UpdateDataset(shard); err != nil {
							log.Errorf(ctx, "Set expired status of dataset %s shard %d error %v", shard.Name, shard.Index, err)
						} else {
							log.Infof(ctx, "Set expired status of dataset %s  shard %d as %v", shard.Name, shard.Index, shard.Status)
						}
					}
				}
			}
		}
	}
}
