package manager

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"pa.cn/fedlearn/psi/api/types"
	"pa.cn/fedlearn/psi/config"
	"pa.cn/fedlearn/psi/log"
	"pa.cn/fedlearn/psi/model"
	"pa.cn/fedlearn/psi/utils"
	"pa.cn/fedlearn/psi/worker"
)

// JobSubmit submit a job to be scheduled.
// As the dataset of the job cannot exceeds max lines, JobScheduler will split the job to tasks accord the dataset count
func (m *TaskMgr) JobSubmit(ctx context.Context, r types.JobSubmitRequest) (err error) {
	log.Debugln(ctx, "TaskMgr.JobSubmit is called")
	// Create Task in local db
	if r.JobUID == "" {
		r.JobUID = utils.UUIDStr()
	}

	// Set the initiator as the Initiator Party
	r.Initiator = config.GetConfig().PartyName
	// Set the initiator as the client mode
	r.Mode = model.TaskRunMode_Client

	r.Protocol = strings.ToUpper(r.Protocol)
	if r.Protocol != worker.DiffieHellman && r.Protocol != worker.OT {
		return fmt.Errorf("task protocol must be [OT|DH], not %s", r.Protocol)
	}

	if err := m.submitJob(ctx, r); err != nil {
		log.Errorf(ctx, "submitJobFailed:%v", err)
		return err
	}

	remoteParty := r.PartyDataset.PartyName
	remoteClient, err := getRemoteClient(remoteParty)
	if err != nil {
		log.Errorf(ctx, "Get party client for %s error:%v", remoteParty, err)
		return fmt.Errorf("Get party client for %s error:%w", remoteParty, err)
	}
	r.Mode = GetPeeringMode(r.Mode)
	r.LocalDataset, r.PartyDataset = r.PartyDataset, r.LocalDataset
	ok, err := remoteClient.PartyJobSubmit(ctx, r)
	if err != nil {
		log.Errorf(ctx, "remote party submit job error:%v", err)
		return fmt.Errorf("remote party submit job error:%w", err)
	}
	if !ok {
		return fmt.Errorf("%s", "Remote Unknown Error")
	}
	return nil
}

func (m *TaskMgr) PartyJobSubmit(ctx context.Context, r types.JobSubmitRequest) (err error) {
	log.Debugln(ctx, "TaskMgr.PartyJobSubmit is called")
	return m.submitJob(ctx, r)
}

func (m *TaskMgr) submitJob(ctx context.Context, r types.JobSubmitRequest) (err error) {
	localDSCount, err := getDatasetCount(ctx, r.LocalDataset.PartyName, r.LocalDataset.DSName, int32(0))
	if err != nil {
		log.Warningf(ctx, "Get dataset count failed:%v", err)
	}
	partyDSCount, err := getDatasetCount(ctx, r.PartyDataset.PartyName, r.PartyDataset.DSName, int32(0))
	if err != nil {
		log.Warningf(ctx, "Get dataset count failed:%v", err)
	}
	job := &model.Job{
		Initiator:    r.Initiator,
		Mode:         r.Mode,
		Status:       model.JobStatus_WaitingPartyConfirm,
		Uuid:         r.JobUID,
		Protocol:     r.Protocol,
		Name:         r.Name,
		Desc:         r.Desc,
		LocalName:    r.LocalDataset.PartyName,
		LocalDSName:  r.LocalDataset.DSName,
		LocalDSCount: localDSCount,
		PartyName:    r.PartyDataset.PartyName,
		PartyDSName:  r.PartyDataset.DSName,
		PartyDSCount: partyDSCount,
	}
	if err := model.AddJob(job); err != nil {
		log.Errorf(ctx, "Add job %s error %v", job.Uuid, err)
		return fmt.Errorf("Add job %s error %w", job.Uuid, err)
	}

	jobDir := utils.JobPath(job.Uuid)
	if err = os.MkdirAll(jobDir, os.ModePerm); err != nil {
		log.Errorf(ctx, "Make dir  %s error %v", jobDir, err)
		return fmt.Errorf("Make dir  %s error %w", jobDir, err)
	}
	return nil
}

func (m *TaskMgr) JobList(ctx context.Context) (jobList []*types.Job, err error) {
	log.Debugln(ctx, "TaskMgr.JobList is called")
	jobList = make([]*types.Job, 0)
	jobs, err := model.ListJobs()
	if err != nil {
		return jobList, err
	}
	for _, j := range jobs {
		job := toAPIJob(j)
		jobList = append(jobList, job)
	}
	return
}

func (m *TaskMgr) JobConfirm(ctx context.Context, r types.JobConfirmRequest) (err error) {
	log.Debugln(ctx, "TaskMgr.JobConfirm is called")
	job, err := model.GetJobByUuid(r.JobUID)
	if err != nil || job == nil {
		log.Errorf(ctx, "Load job %s failed %v \n", r.JobUID, err)
		return fmt.Errorf("No this job %s", r.JobUID)
	}
	remoteClient, err := getRemoteClient(job.PartyName)
	if err != nil {
		log.Errorf(ctx, "Get party client for %s error:%v", job.PartyName, err)
		return fmt.Errorf("Get party client for %s error:%w", job.PartyName, err)
	}
	ok, err := remoteClient.ConfirmPartyJob(ctx, r)
	if err == nil && ok {
		job.Status = model.JobStatus_Confirmed
		err = model.UpdateJob(job)
		if err != nil {
			log.Errorf(ctx, "Set job %s to ready failed:%v\n", r.JobUID, err)
		}
	} else {
		log.Errorf(ctx, "%s, OK=%v, Error:%v", r.JobUID, ok, err)
	}
	return nil
}

func (m *TaskMgr) PartyJobConfirm(ctx context.Context, r types.JobConfirmRequest) (err error) {
	log.Debugln(ctx, "TaskMgr.PartyJobConfirm is called")
	return m.confirmJob(ctx, r.JobUID)
}

func (m *TaskMgr) confirmJob(ctx context.Context, jobUID string) (err error) {
	job, err := model.GetJobByUuid(jobUID)
	if err != nil || job == nil {
		log.Errorf(ctx, "Load job %s failed %v \n", jobUID, err)
		return fmt.Errorf("No this job %s", jobUID)
	}
	job.Status = model.JobStatus_Confirmed
	if err := model.UpdateJob(job); err != nil {
		return err
	}
	return nil
}

func (m *TaskMgr) JobIntersectRead(ctx context.Context, jobUID string) (data []byte, err error) {
	log.Debugln(ctx, "TaskMgr.JobIntersectRead is called")
	job, err := model.GetJobByUuid(jobUID)
	if err != nil {
		log.Errorf(ctx, "Get job %s error %v", jobUID, err)
		return nil, fmt.Errorf("Get job %s error %w", jobUID, err)
	}

	if job.Status != model.TaskStatus_Completed {
		return nil, fmt.Errorf("Job %s is not completed", jobUID)
	}

	//If intersect file exists ...
	if utils.IsJobIntersectExists(jobUID) {
		data, err = ioutil.ReadFile(utils.JobIntersectPath(jobUID))
		if err == nil {
			return data, nil
		}
		log.Errorf(ctx, "Read job %s intersect file error:%v", jobUID, err)
		return nil, fmt.Errorf("Read job %s intersect file error:%w", jobUID, err)
	}

	if job.Mode == model.TaskRunMode_Server {
		// get result info from client side
		log.Infoln(ctx, "Read result from client side")
		remoteClient, err := getRemoteClient(job.PartyName)
		if err != nil {
			return nil, fmt.Errorf("Remote Error:%v", err)
		}

		var intersectData []byte
		if intersect, err := remoteClient.JobIntersectPartyResult(ctx, jobUID); err != nil {
			return nil, err
		} else {
			intersectData = bytes.NewBufferString(intersect).Bytes()
		}

		if err := ioutil.WriteFile(utils.JobIntersectPath(jobUID), intersectData, 0776); err != nil {
			log.Warningln(ctx, "Write intersect file failed.")
		}
		return intersectData, nil
	}

	return make([]byte, 0), errors.New("Unknown Error")
}

func (m *TaskMgr) JobStop(ctx context.Context, r types.JobStopRequest) (err error) {
	log.Debugln(ctx, "TaskMgr.JobStop is called")
	jobUID := r.JobUID
	job, err := model.GetJobByUuid(jobUID)
	if err != nil {
		log.Errorf(ctx, "Get job %s error %v", jobUID, err)
		return fmt.Errorf("Get job %s error %w", jobUID, err)
	}

	// if the job can be stopped
	if !model.JobCanStop(job.Status) {
		log.Warningf(ctx, "Cannot stop job in %s", job.Status)
		return fmt.Errorf("Cannot stop job in %s", job.Status)
	}

	tasks, err := model.ListTasksOfJob(jobUID)
	if err != nil {
		log.Errorf(ctx, "List tasks of job %s failed:%v", jobUID, err)
		return fmt.Errorf("List tasks of job %s failed:%v", jobUID, err)
	}

	for _, task := range tasks {
		go func(taskUID string) {
			taskStopRequest := types.TaskStopRequest{
				TaskUID: taskUID,
			}
			if err := m.TaskStop(ctx, taskStopRequest); err != nil {
				log.Errorf(ctx, "Stop task %s job %s failed:%v", taskUID, jobUID, err)
			}
		}(task.Uuid)
	}

	job.Status = model.JobStatus_Cancel
	if err := model.UpdateJob(job); err != nil {
		log.Errorf(ctx, "Update job %s failed:%v", jobUID, err)
	}
	return nil
}

func (m *TaskMgr) JobDel(ctx context.Context, r types.JobDelRequest) (err error) {
	log.Debugln(ctx, "TaskMgr.JobDel is called")
	jobUID := r.JobUID
	job, err := model.GetJobByUuid(jobUID)

	if err != nil {
		log.Errorf(ctx, "Get job %s error %v", jobUID, err)
		return fmt.Errorf("Get job %s error %w", jobUID, err)
	}

	if job.Status == model.TaskStatus_Running {
		log.Errorf(ctx, "Cannot delele running job %s", jobUID)
		return fmt.Errorf("Cannot delele running job %s", jobUID)
	}

	if err = model.DeleteJob(job); err != nil {
		log.Errorf(ctx, "Delete job %s error %v", jobUID, err)
		return fmt.Errorf("Delete job %s error %w", jobUID, err)
	}

	if err = model.DeleteTasksOfJob(jobUID); err != nil {
		log.Errorf(ctx, "Delete tasks of job %s error %v", jobUID, err)
		return fmt.Errorf("Delete task of job %s error %w", jobUID, err)
	}
	// TODO: remove tasks/jobs in file system
	return nil
}
