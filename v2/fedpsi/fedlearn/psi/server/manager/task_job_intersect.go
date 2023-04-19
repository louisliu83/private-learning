package manager

import (
	"context"
	"fmt"

	"fedlearn/psi/api/types"
	"fedlearn/psi/common/log"
	"fedlearn/psi/model"
)

func (m *TaskMgr) JobListByDataset(ctx context.Context, dsName string) (jobList []*types.Job, err error) {
	log.Debugln(ctx, "TaskMgr.JobListByDataset is called")
	jobList = make([]*types.Job, 0)
	jobs, err := model.ListJobsByLocalDataset(dsName)
	if err != nil {
		return jobList, err
	}
	for _, j := range jobs {
		job := toAPIJob(j)
		jobList = append(jobList, job)
	}
	return
}

func (m *TaskMgr) JobListByDatasetBizContext(ctx context.Context, bizContext string) (jobList []*types.Job, err error) {
	log.Debugln(ctx, "TaskMgr.JobListByDatasetBizContext is called")

	jobList = make([]*types.Job, 0)
	ds, err := model.GetDatasetByBizContext(bizContext)
	if err != nil {
		return jobList, err
	}

	jobs, err := model.ListJobsByLocalDataset(ds.Name)
	if err != nil {
		return jobList, err
	}

	for _, j := range jobs {
		job := toAPIJob(j)
		jobList = append(jobList, job)
	}

	return
}

func (m *TaskMgr) IntersectsOfDatasetBizContext(ctx context.Context, bizContext string) (intersectListResponse *types.IntersectListResponse, err error) {
	log.Debugln(ctx, "TaskMgr.IntersectsOfDatasetBizContext is called")
	intersectListResponse = &types.IntersectListResponse{}
	ds, err := model.GetDatasetByBizContext(bizContext)
	if err != nil {
		return intersectListResponse, err
	}
	intersectListResponse.DatasetName = ds.Name
	intersectListResponse.DatasetDesc = ds.Desc
	intersectListResponse.DatasetBizContext = ds.BizContext
	intersectListResponse.Intersects = make([]string, 0)

	jobList, err := m.JobListByDatasetBizContext(ctx, bizContext)
	if err != nil {
		return intersectListResponse, err
	}

	for _, job := range jobList {
		if job.Status == model.TaskStatus_Completed {
			intersectURL := fmt.Sprintf("http://[API-Server-Address]/apis/v2/job/intersect/download?job_uuid=%s&jwt=[YOUR-TOKEN-ISSUED]", job.Uuid)
			intersectListResponse.Intersects = append(intersectListResponse.Intersects, intersectURL)
		} else {
			return intersectListResponse, fmt.Errorf("jobs are running, not completed.")
		}
	}
	return intersectListResponse, nil
}
