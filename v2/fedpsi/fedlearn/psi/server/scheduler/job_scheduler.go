package scheduler

import (
	"context"
	"time"

	"fedlearn/psi/api/types"
	"fedlearn/psi/common/log"
	"fedlearn/psi/model"
	"fedlearn/psi/server/manager"
	"fedlearn/psi/service"
)

type JobScheduler struct {
	tm       *manager.TaskMgr
	dm       *manager.DataSetMgr
	tickChan chan struct{}
}

var _ Scheduler = NewJobScheduler()

func NewJobScheduler() *JobScheduler {
	return &JobScheduler{
		tm:       manager.GetTaskMgr(),
		dm:       manager.GetDatasetMgr(),
		tickChan: make(chan struct{}, 1024*1024),
	}
}

func (s *JobScheduler) Tick() {
	// Server how to notify the scheduler
	s.tickChan <- struct{}{}
}

func (s *JobScheduler) Start() {

	go func() { //generate ticks every 60 seconds

		for {
			time.Sleep(time.Duration(5) * time.Second)
			s.tickChan <- struct{}{}
		}
	}()

	go func() {
		for {
			<-s.tickChan // check whether has tick
			j, err := model.GetOldestConfirmedJob()
			if err == nil && j != nil {
				ctx := contexForScheduler()
				if j.Mode == model.TaskRunMode_Client {
					s.scheduleClientModeJob(ctx, j)
				}
				if j.Mode == model.TaskRunMode_Server {
					s.scheduleServerModeJob(ctx, j)
				}
			}
		}
	}()
}

func (s *JobScheduler) scheduleServerModeJob(ctx context.Context, j *model.Job) {
	log.Debugf(ctx, "Schedule Server job: %s\n", j.Uuid)

	lDatasetService := service.GetDatasetService(j.LocalName)
	localDS, err1 := lDatasetService.GetDataset(ctx, j.LocalName, j.LocalDSName)
	if err1 != nil {
		log.Errorf(ctx, "Get local dataset shardsNum error:%v", NewJobScheduleError(j.Uuid, err1))
		j.Status = model.JobStatus_Failed
		if err := model.UpdateJob(j); err != nil {
			log.Errorf(ctx, "Update job error:%v", NewJobScheduleError(j.Uuid, err))
		}
		return
	}
	localShards := int64(localDS.ShardsNum)
	if localShards == 0 {
		localShards = 1
	}
	log.Debugf(ctx, "Dataset %s has %d shards\n", j.LocalDSName, localShards)

	pDatasetService := service.GetDatasetService(j.PartyName)
	partyDS, err2 := pDatasetService.GetDataset(ctx, j.PartyName, j.PartyDSName)
	if err2 != nil {
		log.Errorf(ctx, "Get remote dataset shardsNum error:%v", NewJobScheduleError(j.Uuid, err2))
		j.Status = model.JobStatus_Failed
		if err := model.UpdateJob(j); err != nil {
			log.Errorf(ctx, "Update job error:%v", NewJobScheduleError(j.Uuid, err))
		}
		return
	}
	partyShards := int64(partyDS.ShardsNum)
	if partyShards == 0 {
		partyShards = 1
	}
	log.Debugf(ctx, "Dataset %s has %d shards\n", j.PartyDSName, partyShards)

	// localShards := (j.LocalDSCount-1)/config.GetConfig().DataSet.MaxLines + 1
	// partyShards := (j.PartyDSCount-1)/config.GetConfig().DataSet.MaxLines + 1

	var k int
	for k = 0; k < 12; k++ {
		tasks, err := model.ListTasksOfJob(j.Uuid)

		if err != nil {
			log.Errorf(ctx, "%v", NewJobScheduleError(j.Uuid, err))
			time.Sleep(time.Duration(60) * time.Second)
			continue
		}

		if int64(len(tasks)) == localShards*partyShards {
			log.Infof(ctx, "All tasks of job %s generated, total %d\n", j.Uuid, len(tasks))
			j.Status = model.JobStatus_Running
			if err := model.UpdateJob(j); err != nil {
				log.Errorf(ctx, "Update job error:%v", NewJobScheduleError(j.Uuid, err))
			}
			for _, t := range tasks {
				r := types.TaskConfirmRequest{
					TaskUID: t.Uuid,
				}
				if err := s.tm.TaskConfirm(ctx, r); err != nil {
					log.Infof(ctx, "Confirm task %s failed %v\n", t.Uuid, NewJobScheduleError(j.Uuid, err))
				}
			}
			break
		} else {
			log.Warningf(ctx, "Tasks of job %s generated %d tasks, expects %d time %d\n", j.Uuid, len(tasks), localShards*partyShards, k)
			time.Sleep(time.Duration(60) * time.Second)
		}

	}

	if k == 12 {
		j.Status = model.JobStatus_Failed
		if err := model.UpdateJob(j); err != nil {
			log.Errorf(ctx, "Update job error:%v", NewJobScheduleError(j.Uuid, err))
		}
	}

}

func (s *JobScheduler) scheduleClientModeJob(ctx context.Context, j *model.Job) {
	log.Debugf(ctx, "Schedule Client job: %s\n", j.Uuid)
	localDatashards, err := s.dm.DataSetShardsList(ctx, j.LocalDSName)
	if err != nil {
		log.Errorf(ctx, "Load dataset shards error %v", NewJobScheduleError(j.Uuid, err))
		return
	}
	client, err := service.GetRemoteClient(j.PartyName)
	if err != nil {
		log.Errorf(ctx, "%v", NewJobScheduleError(j.Uuid, err))
		return
	}
	partyDatashards, err := client.PartyDataShards(ctx, j.PartyDSName)
	if err != nil {
		log.Errorf(ctx, "Get dataset shards error %v", NewJobScheduleError(j.Uuid, err))
		return
	}
	for _, ld := range localDatashards {
		for _, rd := range partyDatashards {
			reqv2 := genTaskCreateRequestV2(j, ld, rd)
			err := s.tm.TaskCreate(ctx, reqv2)
			if err != nil {
				log.Errorf(ctx, "%v", NewJobScheduleError(j.Uuid, err))
				return
			}
		}
	}
	j.Status = model.JobStatus_Running
	if err := model.UpdateJob(j); err != nil {
		log.Errorf(ctx, "Update job error:%v", NewJobScheduleError(j.Uuid, err))
	}
}

func genTaskCreateRequestV2(j *model.Job, ld, rd types.Dataset) types.TaskCreateRequestV2 {
	reqv2 := types.TaskCreateRequestV2{
		Initiator: j.Initiator,
		Mode:      j.Mode,
		Protocol:  j.Protocol,
		Name:      j.Name,
		Desc:      j.Desc,
		JobUID:    j.Uuid,
	}
	reqv2.LocalDataset = types.TaskDataset{
		PartyName: j.LocalName,
		DSName:    ld.Name,
		DSIndex:   ld.Index,
	}
	reqv2.PartyDataset = types.TaskDataset{
		PartyName: j.PartyName,
		DSName:    rd.Name,
		DSIndex:   rd.Index,
	}
	return reqv2
}
