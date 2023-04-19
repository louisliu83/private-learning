package scheduler

import (
	"pa.cn/fedlearn/psi/server/manager"
)

type DatasetDownloader struct {
}

var _ Scheduler = NewDatasetDownloader()

func NewDatasetDownloader() *DatasetDownloader {
	return &DatasetDownloader{}
}

func (s *DatasetDownloader) Tick() {
	// No need impl
}

func (s *DatasetDownloader) Start() {
	ctx := contexForScheduler()
	manager.GetHttpDownloader().Run(ctx)
}
