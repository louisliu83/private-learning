package scheduler

import (
	"context"
	"fedlearn/psi/api/types"
	"fedlearn/psi/common/config"
	"fedlearn/psi/common/log"
	"fedlearn/psi/model"
	"fedlearn/psi/server/manager"
	"fmt"
	"strings"
)

type ActivityScheduler struct {
	tm *manager.TaskMgr
	Ticker
}

func NewActivityScheduler() *ActivityScheduler {
	return &ActivityScheduler{
		tm:     manager.GetTaskMgr(),
		Ticker: *NewTicker(1024*1024, 120),
	}
}

var _ Scheduler = NewActivityScheduler()

func (s *ActivityScheduler) Start() {
	s.Ticker.Start()

	go func() {
		for {
			<-s.Ticker.tickChan
			s.scheduleFunc()
		}
	}()
}

func (s *ActivityScheduler) scheduleFunc() error {
	t, err := model.GetOldestConfirmedActivity()
	if err == nil && t != nil {
		ctx := contexForScheduler()
		if t.InitParty == config.GetConfig().PartyName {
			s.scheduleClientConfirmed(ctx, t)
		}
	}
	return nil
}

// scheduleClientConfirmed: activity initiator party schedule an confirmed activity, create job from the activity
func (s *ActivityScheduler) scheduleClientConfirmed(ctx context.Context, t *model.Activity) {
	log.Infof(ctx, "initiator party schedule confirmed activity %s", t.Name)
	initDataset := strings.Split(t.InitiatorData, ",")
	if len(initDataset) == 0 {
		return
	}
	followDataset := strings.Split(t.FollowerData, ",")
	if len(followDataset) == 0 {
		return
	}
	for i, id := range initDataset {
		for j, fd := range followDataset {
			localDataset := types.TaskDataset{
				PartyName: t.InitParty,
				DSName:    id,
			}
			partyDataset := types.TaskDataset{
				PartyName: t.FollowerParty,
				DSName:    fd,
			}
			jcr := types.JobSubmitRequest{
				ActUID:       t.Uuid,
				Initiator:    t.InitParty,
				Protocol:     "OT",
				Name:         fmt.Sprintf("Act-%s-%d-%d", t.Name, i, j),
				Desc:         fmt.Sprintf("Job auto generated from activity:%s for dataset(%s,%s)", t.Name, localDataset.DSName, partyDataset.DSName),
				LocalDataset: localDataset,
				PartyDataset: partyDataset,
			}
			s.tm.JobSubmit(ctx, jcr)
		}
	}

	t.Status = model.ActivityStatus_Running
	if err := model.UpdateActivity(t); err != nil {
		log.Errorf(ctx, "Update activity %s error:%v", t.Name, err)
	}
	return
}
