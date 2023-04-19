package scheduler

import (
	"context"
	"time"

	"pa.cn/fedlearn/psi/utils"

	"pa.cn/fedlearn/psi/api/types"
	"pa.cn/fedlearn/psi/config"
	"pa.cn/fedlearn/psi/log"
	"pa.cn/fedlearn/psi/model"
	"pa.cn/fedlearn/psi/server/manager"
)

type TaskScheduler struct {
	tm       *manager.TaskMgr
	tickChan chan struct{}
}

func NewTaskScheduler() *TaskScheduler {
	return &TaskScheduler{
		tm:       manager.GetTaskMgr(),
		tickChan: make(chan struct{}, 1024*1024),
	}
}

var _ Scheduler = NewTaskScheduler()

func (s *TaskScheduler) Tick() {
	// Server how to notify the scheduler
	s.tickChan <- struct{}{}
}

func (s *TaskScheduler) Start() {

	go func() { //generate ticks every 5 seconds
		for {
			s.tickChan <- struct{}{}
			time.Sleep(time.Duration(5) * time.Second)
		}
	}()

	go func() {
		for {
			<-s.tickChan // check whether has tick
			t, err := model.GetOldestTask()
			if err == nil && t != nil {
				ctx := contexForScheduler()
				s.schedule(ctx, t)
			}
		}
	}()
}

func (s *TaskScheduler) schedule(ctx context.Context, t *model.Task) {
	log.Infof(ctx, "Schedule task %s\n", t.Uuid)
	if utils.CheckPortInUse(t.ServerPort) {
		log.Warningf(ctx, "Task %s port %d in use\n", t.Uuid, t.ServerPort)
	} else {
		log.Debugf(ctx, "Task %s port %d is OK\n", t.Uuid, t.ServerPort)
		if config.GetConfig().StartOnceConfirm {
			taskStartRequest := types.TaskStartRequest{
				TaskUID: t.Uuid,
			}
			if err := s.tm.TaskStart(ctx, taskStartRequest); err != nil {
				log.Errorf(ctx, "%v", NewTaskScheduleError(t.Uuid, err))
			}
		}
	}
}
