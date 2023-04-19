package scheduler

import (
	"context"
	"time"

	"fedlearn/psi/api/types"
	"fedlearn/psi/common/log"
	"fedlearn/psi/common/portmgr"
	"fedlearn/psi/model"
	"fedlearn/psi/server/manager"
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
	log.Infof(ctx, "Schedule task %s", t.Uuid)
	port := portmgr.PSIServerPortManager.AcquireAvailablePort(t.Uuid)
	if port <= 0 {
		log.Warningf(ctx, "no available server port for task %s", t.Uuid)
	} else {
		log.Infof(ctx, "get available port for task %s,%d\n", t.Uuid, port)

		taskStartRequest := types.TaskStartRequest{
			TaskUID: t.Uuid,
		}
		// Start the task
		if err := s.tm.TaskStart(ctx, taskStartRequest); err != nil {
			log.Errorf(ctx, "%v", NewTaskScheduleError(t.Uuid, err))
		}

	}
}
