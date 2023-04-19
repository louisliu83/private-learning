package scheduler

import (
	"context"
	"time"

	"fedlearn/psi/api/types"
	"fedlearn/psi/common/log"
	"fedlearn/psi/model"
	"fedlearn/psi/server/manager"
)

type TaskRerunScheduler struct {
	tm             *manager.TaskMgr
	tickChan       chan struct{}
	taskRerunCache map[string]int
}

func NewTaskRerunScheduler() *TaskRerunScheduler {
	return &TaskRerunScheduler{
		tm:             manager.GetTaskMgr(),
		tickChan:       make(chan struct{}, 1024*1024),
		taskRerunCache: make(map[string]int, 0),
	}
}

var _ Scheduler = NewTaskRerunScheduler()

func (s *TaskRerunScheduler) Tick() {
	// Server how to notify the scheduler
	s.tickChan <- struct{}{}
}

func (s *TaskRerunScheduler) Start() {

	go func() { //generate ticks every 5 seconds
		for {
			s.tickChan <- struct{}{}
			time.Sleep(time.Duration(120) * time.Second)
		}
	}()

	go func() {
		for {
			<-s.tickChan // check whether has tick
			ts, err := model.GetFailedTasks()
			if err == nil && ts != nil {
				for i := 0; i < len(ts); i++ { //loop all failed tasks
					if v, ok := s.taskRerunCache[ts[i].Uuid]; ok {
						if v >= 5 {
							continue // has rerun 5 times in this life, skip
						} else {
							// increase rerun times
							s.taskRerunCache[ts[i].Uuid] = v + 1
						}
					} else {
						// set rerun times
						s.taskRerunCache[ts[i].Uuid] = 1
					}

					ctx := contexForScheduler()
					s.schedule(ctx, ts[i])

					time.Sleep(time.Second * time.Duration(30))
				}
			}
		}
	}()
}

func (s *TaskRerunScheduler) schedule(ctx context.Context, t *model.Task) {
	log.Infof(ctx, "Schedule failed task %s\n", t.Uuid)

	taskRerunRequest := types.TaskRerunRequest{
		TaskUID: t.Uuid,
	}
	if err := s.tm.TaskRerun(ctx, taskRerunRequest); err != nil {
		log.Errorf(ctx, "%v", NewTaskScheduleError(t.Uuid, err))
	}

}
