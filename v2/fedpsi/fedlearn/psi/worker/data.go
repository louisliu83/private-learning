package worker

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"fedlearn/psi/common/log"
	"fedlearn/psi/common/utils"
	"fedlearn/psi/model"
)

// MergeJobIntersect merge all tasks' intersect together
// At last, if the update the job's intersect count in db in client side
// TODO, how to sync the intersect file and job intersect count to the server side
func MergeJobIntersect(ctx context.Context, tasks []*model.Task, job *model.Job) error {
	// only merge job intersect in client side
	if job.Mode == model.TaskRunMode_Server {
		return nil
	}
	jobIntersectPath := utils.JobIntersectPath(job.Uuid)
	jobPath := filepath.Dir(jobIntersectPath)
	if err := os.MkdirAll(jobPath, os.ModePerm); err != nil {
		log.Errorf(ctx, "Create job intersect dir error:%v", err)
		return err
	}

	jobFile, err := os.Create(jobIntersectPath)
	if err != nil {
		log.Errorf(ctx, "Create job intersect file error:%v", err)
		return err
	}
	defer jobFile.Close()

	for _, task := range tasks {
		taskResultInfoPath := utils.TaskResultPath(task.Uuid)
		data, err := ioutil.ReadFile(taskResultInfoPath)
		if err != nil {
			log.Errorf(ctx, "Read result file of %s error:%v", task.Uuid, err)
			return err
		}

		log.Debugf(ctx, "Write from %s to %s (%d) bytes", task.Uuid, job.Uuid, len(data))

		if _, err := jobFile.Write(data); err != nil {
			log.Errorf(ctx, "Write file %s to error:%v", jobIntersectPath, err)
		}
	}

	UpdateJobIntersectCount(ctx, job)
	return nil
}

func UpdateJobIntersectCount(ctx context.Context, job *model.Job) {
	if job.Mode == model.TaskRunMode_Server {
		return
	}

	jobIntersectPath := utils.JobIntersectPath(job.Uuid)

	lineCount, _, _, err := utils.FileMetaInfo(jobIntersectPath)
	if err != nil {
		log.Errorf(ctx, "Get intersect line count of %s error:%v", jobIntersectPath, err)
	}

	job.IntersectCount = lineCount
	if err := model.UpdateJob(job); err != nil {
		log.Errorf(ctx, "Update job intersect file %s line count failed:%v", jobIntersectPath, err)
	}
}
