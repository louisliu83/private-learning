package scheduler

import (
	"context"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"pa.cn/fedlearn/psi/api"
	"pa.cn/fedlearn/psi/config"
	"pa.cn/fedlearn/psi/utils"
)

type Scheduler interface {
	Start()
	Tick()
}

var (
	jobScheduler       Scheduler
	taskScheduler      Scheduler
	taskRerunScheduler Scheduler
	datasetScheduler   Scheduler
	sharderScheduler   Scheduler
	downloadSheduler   Scheduler
)

// GetJobScheduler return sigleton Job Scheduler
func GetJobScheduler() Scheduler {
	var once sync.Once
	once.Do(func() {
		logrus.Infoln("New Job Scheduler")
		jobScheduler = NewJobScheduler()
	})
	return jobScheduler
}

// GetTaskScheduler return sigleton Task Scheduler
func GetTaskScheduler() Scheduler {
	var once sync.Once
	once.Do(func() {
		logrus.Infoln("New Task Scheduler")
		taskScheduler = NewTaskScheduler()
	})
	return taskScheduler
}

// GetTaskRerunScheduler return sigleton Task Scheduler
func GetTaskRerunScheduler() Scheduler {
	var once sync.Once
	once.Do(func() {
		logrus.Infoln("New Task Rerun Scheduler")
		taskRerunScheduler = NewTaskRerunScheduler()
	})
	return taskRerunScheduler
}

// GetDatasetScheduler return sigleton Dataset Scheduler
func GetDatasetScheduler() Scheduler {
	var once sync.Once
	once.Do(func() {
		logrus.Infoln("New Dataset Scheduler")
		datasetScheduler = NewDatasetScheduler()
	})
	return datasetScheduler
}

// GetSharderScheduler return sigleton Sharder Scheduler
func GetSharderScheduler() Scheduler {
	var once sync.Once
	once.Do(func() {
		logrus.Infoln("New Sharder Scheduler")
		sharderScheduler = NewSharderScheduler()
	})
	return sharderScheduler
}

// GetDownloadScheduler return sigleton http downloader Scheduler
func GetDownloadScheduler() Scheduler {
	var once sync.Once
	once.Do(func() {
		logrus.Infoln("New HTTPDownloader Scheduler")
		downloadSheduler = NewDatasetDownloader()
	})
	return downloadSheduler
}

func contexForScheduler() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, api.Trace_ID, fmt.Sprintf("Trace_Sched_%s", utils.UUIDStr()))
	ctx = context.WithValue(ctx, api.ReqHeader_PSIUserID, "SysScheduler")
	ctx = context.WithValue(ctx, api.ReqHeader_PSIUserParty, config.GetConfig().PartyName)
	return ctx
}
