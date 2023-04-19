package manager

import (
	"context"
	"errors"
	"fedlearn/psi/api/types"
	"fedlearn/psi/common/config"
	"fedlearn/psi/common/log"
	"fedlearn/psi/common/utils"
	"fedlearn/psi/model"
	"fedlearn/psi/service"
	"fmt"
	"strings"
)

/***
1. Initiator party create an activity, the activity is in "created" state .
2. Initiator party attach all dataset needed to the activity .
3. Initiator start the activity, the activity is in "waiting for confirm", and the activity will be created in follower party
4. Follower party can attach all dataset needed to the activity .
5. Follower party confirm the activity
*/

func (m *TaskMgr) ActivityCreate(ctx context.Context, r types.ActivityCreateRequest) (err error) {
	log.Debugln(ctx, "TaskMgr.ActivityCreate is called")
	if r.Uuid == "" {
		r.Uuid = utils.UUIDStr()
	}

	if r.InitParty == "" {
		r.InitParty = config.GetConfig().PartyName
	}

	if r.FollowerParty == "" {
		log.Errorf(ctx, "empty follower party")
		return fmt.Errorf("empty follower party")
	}
	r.Status = model.ActivityStatus_Created
	err = m.createActivity(ctx, r)
	if err != nil {
		return err
	}

	return nil
}

func (m *TaskMgr) ActivityStart(ctx context.Context, activityUID string) (err error) {
	log.Debugln(ctx, "TaskMgr.ActivityStart is called")
	o, err := model.GetActivityByUuid(activityUID)
	if err != nil {
		return err
	}
	o.Status = model.ActivityStatus_WaitingPartyConfirm
	if err = model.UpdateActivity(o); err != nil {
		return err
	}

	// call party activity create
	a := ToAPIActivity(o)
	remoteClient, err := service.GetRemoteClient(o.FollowerParty)
	if err != nil {
		return
	}

	acr := types.ActivityCreateRequest{
		Activity: *a,
	}

	if ok, err := remoteClient.CreatePartyActivity(ctx, acr); err != nil {
		return err
	} else if !ok {
		return errors.New("create party activity error")
	}
	return nil
}

func (m *TaskMgr) PartyActivityCreate(ctx context.Context, r types.ActivityCreateRequest) (err error) {
	log.Debugln(ctx, "TaskMgr.PartyActivityCreate is called")
	r.Status = model.ActivityStatus_WaitingPartyConfirm
	return m.createActivity(ctx, r)
}

func (m *TaskMgr) createActivity(ctx context.Context, r types.ActivityCreateRequest) (err error) {
	a := &model.Activity{
		Uuid:          r.Uuid,
		Name:          r.Name,
		Title:         r.Title,
		SendID:        r.SendID,
		Desc:          "",
		InitParty:     r.InitParty,
		FollowerParty: r.FollowerParty,
		Status:        r.Status,
	}

	initiatorData := make([]string, 0)
	if len(r.Dataset) != 0 {
		initiatorData = append(initiatorData, r.Dataset...)
	}
	a.InitiatorData = strings.Join(initiatorData, ",")

	err = model.AddActivity(a)
	if err != nil {
		log.Errorf(ctx, "add activity %s failed:%v", r.Name, err)
		return err
	}
	return nil
}

func (m *TaskMgr) ActivityDelete(ctx context.Context, r types.ActivityDeleteRequest) error {
	log.Debugln(ctx, "TaskMgr.ActivityDelete is called")
	return model.DeleteActivityByUUID(r.Uuid)
}

func (m *TaskMgr) ActivityList(ctx context.Context) (as []*types.Activity, err error) {
	log.Debugln(ctx, "TaskMgr.ActivityList is called")

	activities, err := model.ListActivities()
	if err != nil {
		log.Errorf(ctx, "List activities failed:%v", err)
		return nil, err
	}

	for _, act := range activities {
		a := ToAPIActivity(act)
		as = append(as, a)
	}
	return as, err
}

func (m *TaskMgr) ActivityAttachData(ctx context.Context, r types.ActivityAttachDataRequest) (err error) {
	log.Debugln(ctx, "TaskMgr.ActivityAddData is called")
	a, err := model.GetActivityByUuid(r.Uuid)
	if err != nil {
		log.Errorf(ctx, "No activity %s", r.Uuid)
		return fmt.Errorf("No activity %s", r.Uuid)
	}
	if config.GetConfig().PartyName == a.InitParty {
		a.InitiatorData = strings.Join(r.Dataset, ",")
	} else {
		a.FollowerData = strings.Join(r.Dataset, ",")
	}
	err = model.UpdateActivity(a)
	return
}

func (m *TaskMgr) ActivityConfirm(ctx context.Context, r types.ActivityConfirmRequest) (err error) {
	log.Debugln(ctx, "TaskMgr.PartyActivityConfirm is called")

	a, err := model.GetActivityByUuid(r.Uuid)
	if err != nil {
		log.Errorf(ctx, "Get activity %s failed:%v", r.Uuid, err)
		return err
	}

	a.Status = model.ActivityStatus_Confirmed
	err = model.UpdateActivity(a)
	if err != nil {
		log.Errorf(ctx, "Add activity %s failed:%v", a.Name, err)
		return err
	}

	remoteClient, err := service.GetRemoteClient(a.InitParty)
	if err != nil {
		log.Errorln(ctx, err)
		return
	}
	acr := types.ActivityConfirmRequest{
		Uuid:            a.Uuid,
		FollowerDataset: strings.Split(a.FollowerData, ","),
	}
	if ok, err := remoteClient.ConfirmPartyActivity(ctx, acr); !ok || err != nil {
		return err
	}
	a.Status = model.ActivityStatus_Running
	if err := model.UpdateActivity(a); err != nil {
		return err
	}
	return nil
}

func (m *TaskMgr) PartyActivityConfirm(ctx context.Context, r types.ActivityConfirmRequest) (err error) {
	log.Debugln(ctx, "TaskMgr.PartyActivityConfirm is called")
	a, err := model.GetActivityByUuid(r.Uuid)
	if err != nil {
		log.Errorf(ctx, "Get activity %s failed:%v", r.Uuid, err)
		return err
	}
	a.FollowerData = strings.Join(r.FollowerDataset, ",")
	a.Status = model.ActivityStatus_Confirmed
	err = model.UpdateActivity(a)
	if err != nil {
		log.Errorf(ctx, "Add activity %s failed:%v", a.Name, err)
		return err
	}
	return nil
}
