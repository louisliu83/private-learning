package worker

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"pa.cn/fedlearn/psi/config"

	"pa.cn/fedlearn/psi/log"
	"pa.cn/fedlearn/psi/model"
	"pa.cn/fedlearn/psi/utils"
)

type TaskExecutor struct {
	BinPath string
}

func NewExecutor(bin string) *TaskExecutor {
	e := &TaskExecutor{
		BinPath: bin,
	}

	if err := e.selfCheck(); err != nil {
		logrus.Errorln("PSI Executor self check failed.", err)
	}

	return e
}

func (e *TaskExecutor) selfCheck() error {
	fileInfo, err := os.Lstat(e.BinPath)
	if err != nil {
		logrus.Errorf("PSI Executor %s not exists.", e.BinPath, err)
		return err
	}
	logrus.Infoln(e.BinPath, "File perm:", fileInfo.Mode().Perm())
	return nil
}

func (e *TaskExecutor) Start(ctx context.Context, task *model.Task) error {

	ds, err := model.GetDatasetByNameAndIndex(task.LocalDSName, task.LocalDSIndex)
	if err != nil {
		log.Errorf(ctx, "DB get dataset %s error %v", task.LocalDSName, err)
		return err
	}

	inputFilePath := utils.TaskDataSetPath(ds.Md5, task.LocalDSName, task.LocalDSIndex)
	resultPath := utils.TaskResultPath(task.Uuid)

	pidChan := make(chan int)

	go func() {
		task.ExecStart = time.Now()
		if err := model.UpdateTask(task); err != nil {
			log.Errorf(ctx, "update task %s exec start time error %v\n", task.Uuid, err)
		}

		if err := UpdateJobExecStart(ctx, task.JobUid, task.ExecStart); err != nil {
			log.Errorf(ctx, "update job %s exec start time error %v\n", task.JobUid, err)
		}

		serverIP, serverPort := task.ServerIP, task.ServerPort
		if task.Mode == model.TaskRunMode_Client {
			serverIP, serverPort = GetTargetIpAndPort(task)
		}

		err = Start(ctx,
			e.BinPath,
			protocolNumber(task.Protocol),
			task.Mode,
			inputFilePath,
			serverIP,
			fmt.Sprintf("%d", serverPort),
			resultPath,
			task.PartyDSCount,
			config.GetConfig().PsiExecutor.ServerTimeout,
			pidChan)

		if err != nil {
			log.Errorf(ctx, "start psi task %s error %v\n", task.Uuid, err)
			task.Status = model.TaskStatus_Failed
			task.ExecEnd = time.Now()
			if err := model.UpdateTask(task); err != nil {
				log.Errorf(ctx, "update task %s status error %v\n", task.Uuid, err)
			}
			if err := UpdateJobStatus(ctx, task.Uuid, task.JobUid); err != nil {
				log.Errorf(ctx, "update job %s status error %v\n", task.JobUid, err)
			}
		} else {
			task.Status = model.TaskStatus_Completed
			task.ExecEnd = time.Now()
			if err := model.UpdateTask(task); err != nil {
				log.Errorf(ctx, "update task %s exec end time error %v\n", task.Uuid, err)
			}
			if err := UpdateJobStatus(ctx, task.Uuid, task.JobUid); err != nil {
				log.Errorf(ctx, "update job %s status error %v\n", task.JobUid, err)
			}
		}
	}()

	pid := <-pidChan
	log.Infof(ctx, "The process id of psi process running in %s mode :%d\n", task.Mode, pid)
	if pid > 0 {
		task.Status = model.TaskStatus_Running
		if err := model.UpdateTask(task); err != nil {
			log.Errorf(ctx, "update task %s status error %v\n", task.Uuid, err)
		}

		pidStr := strconv.Itoa(pid)
		if err := savePIDToFile(pidStr, task.Uuid); err != nil {
			log.Errorf(ctx, "save pid %s to pidfile error %v", pidStr, err)
		}
	} else {
		// error occured
		task.Status = model.TaskStatus_Failed
		task.ExecEnd = time.Now()
		if err := model.UpdateTask(task); err != nil {
			log.Errorf(ctx, "update task %s status error %v\n", task.Uuid, err)
		}
	}
	return nil
}

func (e *TaskExecutor) Stop(ctx context.Context, task *model.Task) error {
	pidStr, err := getPIDOfTask(task.Uuid)
	if err != nil {
		log.Errorf(ctx, "get pid of task %s error %v", task.Uuid, err)
		return err
	}

	log.Infof(ctx, "Stop the task %s with PID %s", task.Uuid, pidStr)
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return err
	}

	Stop(ctx, pid)
	return nil
}

func ResultInfo(ctx context.Context, taskUUID string) ([]byte, error) {
	resultInfoPath := utils.TaskResultPath(taskUUID)
	data, err := ioutil.ReadFile(resultInfoPath)
	if err != nil {
		log.Errorf(ctx, "Failed to read the task result (%s)\n", taskUUID)
		return nil, errors.New("Read result file err")
	}
	return data, nil
}

func UpdateJobExecStart(ctx context.Context, jobUUID string, startTime time.Time) error {
	zeroTime := time.Time{}
	if jobUUID == "" {
		log.Errorf(ctx, "there is no job %s", jobUUID)
		return nil
	}

	job, err := model.GetJobByUuid(jobUUID)
	if err != nil {
		log.Errorf(ctx, "there is no job %s", jobUUID)
		return err
	}

	if job.ExecStart == zeroTime {
		job.ExecStart = startTime
		if err := model.UpdateJob(job); err != nil {
			log.Errorf(ctx, "update job exec start time failed:%v", jobUUID, err)
		}
	}

	return nil
}

// UpdateJobStatus set the job status if all subtasks of job all completed.
func UpdateJobStatus(ctx context.Context, taskUUID, jobUUID string) error {
	if jobUUID == "" {
		log.Infof(ctx, "task %s has no job %s", taskUUID, jobUUID)
		return nil
	}
	job, err := model.GetJobByUuid(jobUUID)
	if err != nil {
		log.Errorf(ctx, "Get job %s from db failed %v", jobUUID, err)
		return err
	}
	tasks, err := model.ListTasksOfJob(jobUUID)
	if err != nil {
		log.Errorf(ctx, "List tasks of job %s from db failed %v", jobUUID, err)
		return err
	}
	tasksCount := len(tasks)
	for _, task := range tasks {
		if task.Status == model.TaskStatus_Failed {
			job.Status = model.JobStatus_Failed
			if err := model.UpdateJob(job); err != nil {
				log.Errorf(ctx, "Update job %s status failed %v", jobUUID, err)
			}
			break
		}
		if task.Status == model.TaskStatus_Completed {
			tasksCount--
		}
	}
	if tasksCount == 0 { //All tasks of job completed.
		// Merge job first
		if err := MergeJobIntersect(ctx, tasks, job); err != nil {
			log.Errorf(ctx, "Merge job %s intersect error:%v", job.Uuid, err)
		}

		job.Status = model.JobStatus_Completed
		job.ExecEnd = time.Now()
		if err := model.UpdateJob(job); err != nil {
			log.Errorf(ctx, "Update job %s error:%v", job.Uuid, err)
		}

	}
	return nil
}

func savePIDToFile(pid, taskUID string) error {
	pidFilePath := utils.TaskPIDPath(taskUID)
	f, err := os.OpenFile(pidFilePath, os.O_WRONLY|os.O_CREATE, 0776)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(pid); err != nil {
		return err
	}
	return nil
}

func getPIDOfTask(taskUID string) (string, error) {
	pidPath := utils.TaskPIDPath(taskUID)
	data, err := ioutil.ReadFile(pidPath)
	if err != nil {
		return "", err
	}
	pidStr := string(data)
	return pidStr, nil
}
