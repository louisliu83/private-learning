package scheduler

import (
	"context"
	"fedlearn/psi/common/log"
	"fedlearn/psi/model"
	"fedlearn/psi/server/manager"
)

type TaskStatusScheduler struct {
	tm *manager.TaskMgr
	Ticker
}

func NewTaskStatusScheduler(capacity int) *TaskStatusScheduler {
	return &TaskStatusScheduler{
		tm:     manager.GetTaskMgr(),
		Ticker: *NewTicker(capacity, 5),
	}
}

var _ Scheduler = NewTaskStatusScheduler(DefaultChanCapacity)

func (s *TaskStatusScheduler) Start() {
	s.Ticker.Start()
	go func() {
		for {
			<-s.Ticker.tickChan
			s.schedulerFunc()
		}
	}()
}

func (s *TaskStatusScheduler) schedulerFunc() {
	ctx := contexForScheduler()
	log.Debugln(ctx, "TaskStatusScheduler.schedulerFunc entered")
	acts, err := model.ListActivities()
	if err != nil {
		return
	}
	if len(acts) == 0 {
		return
	}
	for _, act := range acts {
		if act.Status != model.ActivityStatus_failed && act.Status != model.ActivityStatus_Completed {
			log.Debugln(ctx, "check activity status of ", act.Name)
			updateActivityStatus(ctx, act)
		}
	}
}

func updateActivityStatus(ctx context.Context, act *model.Activity) error {
	jobs, err := model.ListJobsOfActivity(act.Uuid)
	if err != nil || len(jobs) == 0 {
		log.Warningf(ctx, "No jobs for activity %s", act.Name)
		return nil
	}

	var completed int = 0
	var running int = 0

	for _, job := range jobs {
		if job.Status == model.JobStatus_Failed {
			act.Status = model.ActivityStatus_failed
			model.UpdateActivity(act)
			return nil
		}
		if job.Status == model.JobStatus_Completed {
			completed += 1
		}
		if job.Status == model.JobStatus_Running {
			running += 1
		}
	}

	if completed == len(jobs) {
		act.Status = model.ActivityStatus_Completed
		model.UpdateActivity(act)
		return nil
	}

	if completed < len(jobs) && running > 0 {
		act.Status = model.ActivityStatus_Running
		model.UpdateActivity(act)
		return nil
	}
	return nil
}
