package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"fedlearn/psi/api"
	"fedlearn/psi/common/config"
	"fedlearn/psi/common/utils"

	"github.com/sirupsen/logrus"
)

type Ticker struct {
	tickChan          chan struct{}
	intervalInSeconds int
}

func NewTicker(capacity int, intervalInSeconds int) *Ticker {
	return &Ticker{
		tickChan:          make(chan struct{}, capacity),
		intervalInSeconds: intervalInSeconds,
	}
}

func (t *Ticker) Tick() {
	t.tickChan <- struct{}{}
}

func (t *Ticker) Start() {
	go func() { //generate ticks every 5 seconds
		for {
			t.Tick()
			time.Sleep(time.Duration(t.intervalInSeconds) * time.Second)
		}
	}()
}

/**
Schedulers implemented
*/

type Scheduler interface {
	Start()
	Tick()
}

var (
	activityScheduler   Scheduler
	taskStatusScheduler Scheduler
	jobScheduler        Scheduler
	taskScheduler       Scheduler
	taskRerunScheduler  Scheduler
	datasetScheduler    Scheduler
	sharderScheduler    Scheduler
	downloadSheduler    Scheduler
)

var (
	DefaultChanCapacity int = 1024 * 64
)

func GetTaskStatusScheduler() Scheduler {
	var once sync.Once
	once.Do(func() {
		logrus.Infoln("New Task Status Scheduler")
		taskStatusScheduler = NewTaskStatusScheduler(DefaultChanCapacity)
	})
	return taskStatusScheduler
}

func GetActivityScheduler() Scheduler {
	var once sync.Once
	once.Do(func() {
		logrus.Infoln("New Activity Scheduler")
		activityScheduler = NewActivityScheduler()
	})
	return activityScheduler
}

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
		logrus.Infoln("New Dataset Validator Scheduler")
		datasetScheduler = NewDatasetScheduler()
	})
	return datasetScheduler
}

// GetSharderScheduler return sigleton Sharder Scheduler
func GetSharderScheduler() Scheduler {
	var once sync.Once
	once.Do(func() {
		logrus.Infoln("New DatasetSharder")
		sharderScheduler = NewSharderScheduler()
	})
	return sharderScheduler
}

// GetDownloadScheduler return sigleton http downloader Scheduler
func GetDownloadScheduler() Scheduler {
	var once sync.Once
	once.Do(func() {
		logrus.Infoln("New DatasetHTTPDownloader")
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
