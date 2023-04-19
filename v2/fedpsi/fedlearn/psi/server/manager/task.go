package manager

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
	"time"

	"fedlearn/psi/api/types"
	"fedlearn/psi/common/config"
	"fedlearn/psi/common/log"
	"fedlearn/psi/common/utils"
	"fedlearn/psi/model"
	"fedlearn/psi/worker"
)

type TaskMgr struct {
}

func (m *TaskMgr) TaskCreate(ctx context.Context, r types.TaskCreateRequestV2) (err error) {
	log.Debugln(ctx, "TaskMgr.TaskCreate is called")
	if r.TaskUID == "" {
		r.TaskUID = utils.UUIDStr()
	}
	// Set the initiator as the Initiator Party
	r.Initiator = config.GetConfig().PartyName
	// Set the initiator as the client mode
	r.Mode = model.TaskRunMode_Client

	if err = m.createTask(ctx, r); err != nil {
		return
	}

	// Creating task in remote party
	remoteClient, err := getRemoteClient(r.PartyDataset.PartyName)
	if err != nil {
		return fmt.Errorf("Remote Error:%v", err)
	}

	r.Mode = GetPeeringMode(r.Mode)
	r.LocalDataset, r.PartyDataset = r.PartyDataset, r.LocalDataset
	ok, err := remoteClient.CreatePartyTask(ctx, r)
	if err != nil {
		return fmt.Errorf("Remote Error:%v", err)
	}
	if !ok {
		return fmt.Errorf("%s", "Remote Unknown Error")
	}

	return nil
}

func (m *TaskMgr) PartyTaskCreate(ctx context.Context, r types.TaskCreateRequestV2) (err error) {
	log.Debugln(ctx, "TaskMgr.PartyTaskCreate is called")
	return m.createTask(ctx, r)
}

func (m *TaskMgr) createTask(ctx context.Context, r types.TaskCreateRequestV2) (err error) {
	r.Mode = strings.ToLower(r.Mode)
	if r.Mode != model.TaskRunMode_Server && r.Mode != model.TaskRunMode_Client {
		return fmt.Errorf("task run mode must be [server|client], not %s", r.Mode)
	}

	r.Protocol = strings.ToUpper(r.Protocol)
	if r.Protocol != worker.DiffieHellman && r.Protocol != worker.OT {
		return fmt.Errorf("task protocol must be [OT|DH], not %s", r.Protocol)
	}

	t := &model.Task{
		Initiator: r.Initiator,
		Mode:      r.Mode,
		Status:    model.TaskStatus_WaitingPartyConfirm,
		Uuid:      r.TaskUID,
		JobUid:    r.JobUID,
		Protocol:  r.Protocol,
		Name:      r.Name,
		Desc:      r.Desc,
	}

	t.LocalName = r.LocalDataset.PartyName
	t.LocalDSName = r.LocalDataset.DSName
	t.LocalDSIndex = r.LocalDataset.DSIndex
	if count1, err := getDatasetCount(ctx, t.LocalName, t.LocalDSName, t.LocalDSIndex); err != nil {
		log.Warningf(ctx, "Failed to get dataset(%s:%s):%v\n", t.LocalDSName, t.LocalDSIndex, err)
	} else {
		t.LocalDSCount = count1
	}

	t.PartyName = r.PartyDataset.PartyName
	t.PartyDSName = r.PartyDataset.DSName
	t.PartyDSIndex = r.PartyDataset.DSIndex
	if count2, err := getDatasetCount(ctx, t.PartyName, t.PartyDSName, t.PartyDSIndex); err != nil {
		log.Warningf(ctx, "Failed to get dataset(%s:%s):%v\n", t.PartyDSName, t.PartyDSIndex, err)
	} else {
		t.PartyDSCount = count2
	}

	if err = model.AddTask(t); err != nil {
		return err
	}

	taskDir := utils.TaskPath(t.Uuid)
	if err = os.MkdirAll(taskDir, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (m *TaskMgr) TaskConfirm(ctx context.Context, r types.TaskConfirmRequest) (err error) {
	log.Debugln(ctx, "TaskMgr.TaskConfirm is called")
	task, err := model.GetTaskByUuid(r.TaskUID)
	if err != nil {
		return err
	}
	if config.GetConfig().PartyName == task.Initiator {
		return errors.New("Task should be confirmed by the party, not the initiator")
	}
	task.Status = model.TaskStatus_Created
	err = model.UpdateTask(task)
	remoteClient, err := getRemoteClient(task.PartyName)
	if err != nil {
		return fmt.Errorf("Remote Error:%v", err)
	}
	r.PartyDSCount = task.LocalDSCount
	ok, err := remoteClient.ConfirmPartyTask(ctx, r)
	if err == nil && ok {
		task.Status = model.TaskStatus_Ready
		err = model.UpdateTask(task)
		if err != nil {
			log.Errorf(ctx, "Set task %s to ready failed:%v\n", r.TaskUID, err)
		}
	} else {
		log.Errorf(ctx, "%s, OK=%v, Error:%v", r.TaskUID, ok, err)
	}
	return
}

func (m *TaskMgr) PartyTaskConfirm(ctx context.Context, r types.TaskConfirmRequest) (err error) {
	log.Debugln(ctx, "TaskMgr.PartyTaskConfirm is called")
	task, err := model.GetTaskByUuid(r.TaskUID)
	if err != nil {
		return err
	}
	remoteClient, err := getRemoteClient(task.PartyName)
	if err != nil {
		return fmt.Errorf("Remote Error:%v", err)
	}
	remoteTask, err := remoteClient.GetPartyTask(ctx, task.Uuid)
	if err != nil {
		return err
	}
	if remoteTask.Status != model.TaskStatus_Created {
		return errors.New("party has not confirm the task")
	}
	task.Status = model.TaskStatus_Created
	task.PartyDSCount = r.PartyDSCount
	task.ServerIP = r.PSIServerIP
	task.ServerPort = r.PSIServerPort
	err = model.UpdateTask(task)
	return
}

func (m *TaskMgr) TaskGet(ctx context.Context, r types.TaskGetRequest) (task *types.Task, err error) {
	log.Debugln(ctx, "TaskMgr.TaskGet is called")
	t, err := model.GetTaskByUuid(r.TaskUID)
	if err != nil {
		return nil, err
	}
	task = &types.Task{
		Uuid:         t.Uuid,
		JobUID:       t.JobUid,
		Initiator:    t.Initiator,
		Mode:         t.Mode,
		Status:       t.Status,
		Name:         t.Name,
		Desc:         t.Desc,
		LocalName:    t.LocalName,
		LocalDSName:  t.LocalDSName,
		LocalDSIndex: t.LocalDSIndex,
		LocalDSCount: t.LocalDSCount,
		PartyName:    t.PartyName,
		PartyDSName:  t.PartyDSName,
		PartyDSIndex: t.PartyDSIndex,
		PartyDSCount: t.PartyDSCount,
	}
	return task, nil
}

func (m *TaskMgr) TaskListByPage(ctx context.Context, r types.TaskListRequest) (taskList []*types.Task, count int64, pageCount int64, err error) {
	log.Debugln(ctx, "TaskMgr.TaskList is called")
	taskList = make([]*types.Task, 0)

	listFunc := func() ([]*model.Task, int64, error) {
		return model.ListTasksByPage(r.PageNum, r.PageSize)
	}

	if r.JobUID != "" {
		listFunc = func() ([]*model.Task, int64, error) {
			return model.ListTasksOfJobByPage(r.JobUID, r.PageNum, r.PageSize)
		}
	}

	tasks := make([]*model.Task, 0)
	tasks, count, err = listFunc()
	if err != nil {
		return taskList, 0, 0, err
	}

	for _, t := range tasks {
		apiTask := &types.Task{
			Uuid:         t.Uuid,
			JobUID:       t.JobUid,
			Initiator:    t.Initiator,
			Mode:         t.Mode,
			Name:         t.Name,
			Desc:         t.Desc,
			Status:       t.Status,
			LocalName:    t.LocalName,
			LocalDSName:  t.LocalDSName,
			LocalDSIndex: t.LocalDSIndex,
			LocalDSCount: t.LocalDSCount,
			PartyName:    t.PartyName,
			PartyDSName:  t.PartyDSName,
			PartyDSIndex: t.PartyDSIndex,
			PartyDSCount: t.PartyDSCount,
		}
		taskList = append(taskList, apiTask)
	}
	pageCount = int64(math.Ceil((float64(count) / float64(r.PageSize))))
	return taskList, count, pageCount, nil
}

func (m *TaskMgr) taskClientStartLocally(ctx context.Context, task *model.Task) {
	task.Status = model.TaskStatus_ClientWaiting
	model.UpdateTask(task)

	psiConfig := config.GetConfig().PsiExecutor
	te := worker.NewExecutor(psiConfig.BinPath)
	remoteClient, err := getRemoteClient(task.PartyName)
	if err != nil {
		log.Errorf(ctx, "Remote Error %s:%v\n", task.Uuid, err)
		return
	}

	times := 2 * 60
	var i int
	for i := 0; i < times; i++ {
		time.Sleep(time.Duration(500) * time.Millisecond)

		remoteTaskInfo, err := remoteClient.GetPartyTask(ctx, task.Uuid)
		if err == nil && remoteTaskInfo.Status == model.TaskStatus_Running {
			break
		}
		if remoteTaskInfo.Status == model.TaskStatus_Failed || remoteTaskInfo.Status == model.TaskStatus_Cancel {
			log.Errorf(ctx, "Server task %s failed, set client task failed. \n", task.Uuid)
			task.Status = model.TaskStatus_Failed
			model.UpdateTask(task)
		}
	}

	if i >= times {
		log.Errorf(ctx, "Timeout:server Party %s is not running, cannot start client \n", task.PartyName)
	}

	log.Infof(ctx, "Will start task %s as the client mode", task.Uuid)

	if err := StartEgressListenerForTask(task); err != nil {
		log.Errorf(ctx, "Start egress listener failed.")
	}

	te.Start(ctx, task)
}

func (m *TaskMgr) TaskStart(ctx context.Context, r types.TaskStartRequest) (err error) {
	log.Debugln(ctx, "TaskMgr.TaskStart is called")
	task, err := model.GetTaskByUuid(r.TaskUID)
	if err != nil {
		log.Errorf(ctx, "no task %s, %v", r.TaskUID, err)
		return fmt.Errorf("no task %s, %v", r.TaskUID, err)
	}

	if task.Status == model.TaskStatus_Running || task.Status == model.TaskStatus_ClientWaiting {
		log.Warnf(ctx, "Task %s is already running", task.Uuid)
		return fmt.Errorf("%s", "Task is already running.")
	}

	if task.Status == model.TaskStatus_WaitingPartyConfirm {
		log.Warnf(ctx, "Task %s needs to be confirmed before running", task.Uuid)
		return fmt.Errorf("%s", "Task has not been confirmed.")
	}

	if task.Status != model.TaskStatus_Created {
		m.cleanTask(ctx, task)
	}

	task.Status = model.TaskStatus_Init
	if err := model.UpdateTask(task); err != nil {
		log.Errorf(ctx, "update task %s error: %v", task.Uuid, err)
		return fmt.Errorf("update task %s error: %v", task.Uuid, err)
	}

	psiConfig := config.GetConfig().PsiExecutor
	te := worker.NewExecutor(psiConfig.BinPath)
	// if it is initiator
	if task.Mode == model.TaskRunMode_Server {
		// Start Server task locally, and after the task is running, start the client task
		log.Infof(ctx, "Will start the task %s as the server mode", task.Uuid)
		go func() {
			if err := te.Start(ctx, task); err != nil {
				log.Errorf(ctx, "Start the task %s as the server mode failed:%v", task.Uuid, err)
			}
		}()

		// Start the task in party
		remoteClient, err := getRemoteClient(task.PartyName)
		if err != nil {
			log.Errorf(ctx, "get remote client %s error:%v", task.PartyName, err)
			return fmt.Errorf("get remote client %s error:%v", task.PartyName, err)
		}

		go func() {
			totalTry := 3
			round := 1
			for round <= totalTry {
				ok, err := remoteClient.StartPartyTask(ctx, r)
				if err == nil && ok {
					log.Infof(ctx, "(%d/%d)Start the task %s as client mode in remote party ok", round, totalTry, task.Uuid)
					break
				}
				if err != nil {
					log.Warningf(ctx, "(%d/%d)Start the task %s as client mode in remote party failed:%v", round, totalTry, task.Uuid, err)
				}
				if !ok {
					log.Warningf(ctx, "(%d/%d)Start the task %s as client mode in remote party failed:%v", round, totalTry, task.Uuid, ok)
				}
				round++
				time.Sleep(time.Duration(round) * time.Second)
			}

			if round > totalTry {
				log.Errorf(ctx, "Start the task %s as client mode in remote party failed at last", task.Uuid)
				if err := te.Stop(ctx, task); err != nil {
					log.Errorf(ctx, "Stop the task %s error:%v", task.Uuid, err)
				}
				task.Status = model.TaskStatus_Failed
				model.UpdateTask(task)
			}
		}()

	} else {
		remoteClient, err := getRemoteClient(task.PartyName)
		if err != nil {
			return fmt.Errorf("Remote Error:%v", err)
		}
		// Start Server task remotely

		go func() {
			ok, err := remoteClient.StartPartyTask(ctx, r)
			if err == nil && ok {
				log.Infof(ctx, "Start the task %s as server mode in remote party ok", task.Uuid)
			}
			if err != nil {
				log.Warningf(ctx, "Start the task %s as server mode in remote party failed:%v", task.Uuid, err)
			}
			if !ok {
				log.Warningf(ctx, "Start the task %s as server mode in remote party failed:%v", task.Uuid, ok)
			}
		}()

		// Query the pair task status, if task status is Running
		go m.taskClientStartLocally(ctx, task)
	}
	return nil
}

func (m *TaskMgr) PartyTaskStart(ctx context.Context, r types.TaskStartRequest) (err error) {
	log.Debugln(ctx, "TaskMgr.PartyTaskStart is called")
	task, err := model.GetTaskByUuid(r.TaskUID)
	if err != nil {
		return err
	}

	if task.Status == model.TaskStatus_Running || task.Status == model.TaskStatus_ClientWaiting {
		log.Warnf(ctx, "Task %s is already running.\n", task.Uuid)
		return fmt.Errorf("%s", "Task is already running.")
	}

	if task.Status == model.TaskStatus_WaitingPartyConfirm {
		log.Warnf(ctx, "Task %s needs to be confirmed before running.\n", task.Uuid)
		return fmt.Errorf("%s", "Task has not been confirmed.")
	}

	if task.Status != model.TaskStatus_Created {
		m.cleanTask(ctx, task)
	}

	task.Status = model.TaskStatus_Init
	model.UpdateTask(task)

	psiConfig := config.GetConfig().PsiExecutor
	te := worker.NewExecutor(psiConfig.BinPath)
	// if it is initiator
	if task.Mode == model.TaskRunMode_Server {
		// Start Server task locally, and after the task is running, start the client task
		log.Infof(ctx, "Will start the task %s as the server mode", task.Uuid)
		go te.Start(ctx, task)

	} else {
		// Query the pair task status, if task status is Running
		go m.taskClientStartLocally(ctx, task)
	}
	return nil
}

func (m *TaskMgr) cleanTask(ctx context.Context, task *model.Task) {
	if err := os.Remove(utils.TaskPIDPath(task.Uuid)); err != nil {
		log.Warningf(ctx, "Remove %s failed: %v", utils.TaskPIDPath(task.Uuid), err)
	}
	if err := os.Remove(utils.TaskResultPath(task.Uuid)); err != nil {
		log.Warningf(ctx, "Remove %s failed: %v", utils.TaskResultPath(task.Uuid), err)
	}
	if task.Mode == model.TaskRunMode_Server {
		if err := os.Remove(utils.TaskIntersectPath(task.Uuid)); err != nil {
			log.Warningf(ctx, "Remove %s failed: %v", utils.TaskIntersectPath(task.Uuid), err)
		}
	}
}

func (m *TaskMgr) TaskStop(ctx context.Context, r types.TaskStopRequest) (err error) {
	log.Debugln(ctx, "TaskMgr.TaskStop is called")
	task, err := model.GetTaskByUuid(r.TaskUID)
	if err != nil {
		return err
	}

	err = m.taskStop(ctx, r)
	if err != nil {
		return err
	}

	remoteClient, err := getRemoteClient(task.PartyName)
	if err != nil {
		return fmt.Errorf("Remote Error:%v", err)
	}
	ok, err := remoteClient.StopPartyTask(ctx, r)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("Stop party task %s error", r.TaskUID)
	}
	return nil
}

func (m *TaskMgr) PartyTaskStop(ctx context.Context, r types.TaskStopRequest) (err error) {
	log.Debugln(ctx, "TaskMgr.TaskStop is called")
	return m.taskStop(ctx, r)
}

func (m *TaskMgr) taskStop(ctx context.Context, r types.TaskStopRequest) (err error) {
	task, err := model.GetTaskByUuid(r.TaskUID)
	if err != nil {
		return err
	}
	psiConfig := config.GetConfig().PsiExecutor
	te := worker.NewExecutor(psiConfig.BinPath)
	log.Infof(ctx, "Stop task %s \n", task.Uuid)
	if model.TaskCanStop(task.Status) {
		// Stop task internal
		te.Stop(ctx, task)
		// Update the db state
		task.Status = model.TaskStatus_Cancel
		err = model.UpdateTask(task)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *TaskMgr) TaskIntersectRead(ctx context.Context, r types.TaskIntersectionDownloadRequest) (data []byte, err error) {
	log.Debugln(ctx, "TaskMgr.TaskIntersectRead is called")
	task, err := model.GetTaskByUuid(r.TaskUID)
	if err != nil {
		log.Errorf(ctx, "Load task %s error %v", r.TaskUID, err)
		return nil, fmt.Errorf("Load task %s error %w", r.TaskUID, err)
	}
	//If intersect file exists ...
	if utils.IsTaskIntersectExists(task.Uuid) {
		data, err = ioutil.ReadFile(utils.TaskIntersectPath(task.Uuid))
		if err == nil {
			return data, nil
		}
		log.Errorf(ctx, "Read task %s intersect file error:%v", task.Uuid, err)
		return nil, fmt.Errorf("Read task %s intersect file error:%w", task.Uuid, err)
	}

	var intersectData []byte
	if task.Mode == model.TaskRunMode_Server {
		// get result info from client side
		log.Infoln(ctx, "Read result from client side")
		remoteClient, err := getRemoteClient(task.PartyName)
		if err != nil {
			return nil, fmt.Errorf("Remote Error:%v", err)
		}

		if intersect, err := remoteClient.TaskIntersectPartyResult(ctx, task.Uuid); err != nil {
			return nil, err
		} else {
			intersectData = bytes.NewBufferString(intersect).Bytes()
		}

	}
	if task.Mode == model.TaskRunMode_Client {
		if intersect, err := worker.ResultInfo(ctx, r.TaskUID); err != nil {
			return nil, err
		} else {
			intersectData = intersect
		}
	}

	if err := ioutil.WriteFile(utils.TaskIntersectPath(task.Uuid), intersectData, 0776); err != nil {
		log.Warningln(ctx, "Write intersect file failed.")
	}

	return intersectData, nil
}

func (m *TaskMgr) TaskRerun(ctx context.Context, r types.TaskRerunRequest) (err error) {
	log.Debugln(ctx, "TaskMgr.TaskRerun is called")
	task, err := model.GetTaskByUuid(r.TaskUID)
	if err != nil {
		log.Errorf(ctx, "Load task %s error %v", r.TaskUID, err)
		return fmt.Errorf("Load task %s error %w", r.TaskUID, err)
	}

	remoteClient, err := getRemoteClient(task.PartyName)
	if err != nil {
		return fmt.Errorf("Remote Error:%v", err)
	}

	if task.Mode == model.TaskRunMode_Server {
		success, err := remoteClient.RerunPartyTask(ctx, r)
		if err != nil || !success {
			return fmt.Errorf("Rerun task %s failed", r.TaskUID)
		}
		err = m.taskRerun(ctx, r)
	} else {
		err = m.taskRerun(ctx, r)
		if err != nil {
			return
		}
		success, err := remoteClient.RerunPartyTask(ctx, r)
		if err != nil || !success {
			return fmt.Errorf("Rerun task %s failed", r.TaskUID)
		}
	}

	return
}

func (m *TaskMgr) PartyTaskRerun(ctx context.Context, r types.TaskRerunRequest) (err error) {
	log.Debugln(ctx, "TaskMgr.PartyTaskRerun is called")
	err = m.taskRerun(ctx, r)
	return err
}

func (m *TaskMgr) taskRerun(ctx context.Context, r types.TaskRerunRequest) (err error) {
	log.Debugln(ctx, "TaskMgr.taskRerun is called")

	task, err := model.GetTaskByUuid(r.TaskUID)
	if err != nil {
		log.Errorf(ctx, "Load task %s error %v", r.TaskUID, err)
		return fmt.Errorf("Load task %s error %w", r.TaskUID, err)
	}

	if err := m.cleanJobStateWhenRerunTask(ctx, task); err != nil {
		log.Errorf(ctx, "Clean state of job %s when rerun task %s error %v", task.JobUid, task.Uuid, err)
	}

	if task.Mode == model.TaskRunMode_Server {
		task.Status = model.TaskStatus_Ready
		err = model.UpdateTask(task)
	} else {
		task.Status = model.TaskStatus_Created
		err = model.UpdateTask(task)
	}

	return err
}

func (m *TaskMgr) cleanJobStateWhenRerunTask(ctx context.Context, task *model.Task) error {
	if task == nil {
		return nil
	}
	job, err := model.GetJobByUuid(task.JobUid)
	if err != nil {
		log.Errorf(ctx, "Load job of task %s error %v", task.Uuid, err)
		return fmt.Errorf("Load job of task %s error %w", task.Uuid, err)
	}

	job.Status = model.JobStatus_Running
	if err := model.UpdateJob(job); err != nil {
		return err
	}

	if task.Mode == model.TaskRunMode_Server && utils.IsJobIntersectExists(task.JobUid) {
		// if this is server side and has a intersect file ,remove it
		jobIntersectPath := utils.JobIntersectPath(task.JobUid)
		if err := os.Remove(jobIntersectPath); err != nil {
			log.Errorf(ctx, "Delete job %s intersect failed when rerun task %s :%v", task.JobUid, task.Uuid, err)
			return err
		}
	}

	if job.ActUid == "" {
		return nil
	}

	act, err := model.GetActivityByUuid(job.ActUid)
	if err != nil {
		log.Errorf(ctx, "No activity for job %s", job.Uuid)
		return err
	}

	act.Status = model.ActivityStatus_Running
	model.UpdateActivity(act)

	return nil
}
