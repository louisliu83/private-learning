package manager

import (
	"context"
	"errors"
	"fedlearn/psi/api"
	"fedlearn/psi/api/types"
	"fedlearn/psi/common/config"
	"fedlearn/psi/common/log"
	"fedlearn/psi/common/utils"
	"fedlearn/psi/model"
	"fedlearn/psi/service"
	"fmt"
	"math"
)

type ProjectMgr struct {
}

func (m *ProjectMgr) ProjectCreate(ctx context.Context, r *types.ProjectCreateRequest) error {
	log.Debugln(ctx, "PrjectMgr.ProjectAdd is called")
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
	user := fmt.Sprintf("%s", ctx.Value(api.ReqHeader_PSIUserID))
	r.Creator = user
	r.UpdateUser = user
	p := model.Project{
		Uuid:          r.Uuid,
		Name:          r.Name,
		Type:          r.Type,
		Desc:          r.Desc,
		InitParty:     r.InitParty,
		FollowerParty: r.FollowerParty,
		Status:        r.Status,
		Creator:       r.Creator,
		UpdateUser:    r.UpdateUser,
		Created:       r.Created,
		Updated:       r.Updated,
	}
	err := model.AddProject(&p)
	if err != nil {
		return err
	}

	// call party project create
	remoteClient, err := service.GetRemoteClient(p.FollowerParty)
	if err != nil {
		return err
	}

	if ok, err := remoteClient.CreatePartyProject(ctx, *r); err != nil {
		return err
	} else if !ok {
		return errors.New("create party project error")
	}
	return nil
}

func (m *ProjectMgr) PartyProjectCreate(ctx context.Context, r *types.ProjectCreateRequest) error {
	log.Debugln(ctx, "PrjectMgr.PartyProjectAdd is called")
	p := model.Project{
		Uuid:          r.Uuid,
		Name:          r.Name,
		Type:          r.Type,
		Desc:          r.Desc,
		InitParty:     r.FollowerParty,
		FollowerParty: r.InitParty,
		Status:        r.Status,
		Creator:       r.Creator,
		UpdateUser:    r.UpdateUser,
		Created:       r.Created,
		Updated:       r.Updated,
	}
	err := model.AddProject(&p)
	if err != nil {
		return err
	}
	return nil
}

func (m *ProjectMgr) ProjectUpdate(ctx context.Context, r *types.ProjectUpdateRequest) error {
	p, err := model.GetProjectById(r.Id)
	if err != nil {
		log.Errorf(ctx, "No Project Id %d: %v", r.Id, err)
		return fmt.Errorf("No Project Id %d: %w", r.Id, err)
	}
	if p == nil {
		log.Errorf(ctx, "No Project Id %d", r.Id)
		return fmt.Errorf("No Project Id %d", r.Id)
	}
	p.Name = r.Name
	p.Type = r.Type
	p.Desc = r.Desc
	p.InitParty = r.InitParty
	p.FollowerParty = r.FollowerParty
	p.UpdateUser = fmt.Sprintf("%s", ctx.Value(api.ReqHeader_PSIUserID))
	err = model.UpdateProject(p)
	if err != nil {
		log.Errorf(ctx, "Update project %d error:%v", r.Id, err)
		return fmt.Errorf("Update project %d error:%w", r.Id, err)
	}
	return nil
}

func (m *ProjectMgr) ProjectList(ctx context.Context, name string) []*types.Project {
	log.Debugln(ctx, "ProjectMgr.ProjectList is called")
	projects := make([]*types.Project, 0)
	ps, err := model.ListProjects(name)
	if err != nil {
		log.Errorf(ctx, "List projects error %v", err)
		return projects
	}
	for _, r := range ps {
		p := ToAPIProject(r)
		projects = append(projects, p)
	}
	return projects
}

func (m *ProjectMgr) ProjectDel(ctx context.Context, id uint64) error {
	log.Debugln(ctx, "ProjectMgr.ProjectDel is called")
	p, err := model.GetProjectById(id)
	if err != nil {
		log.Errorf(ctx, "No project %d %v\n", id, err)
		return err
	}
	jobs, err := model.ListJobsOfActivity(p.Uuid)
	if err != nil {
		log.Errorf(ctx, "List jobs of project %s %v\n", p.Uuid, err)
		return err
	}
	if len(jobs) > 0 {
		return fmt.Errorf("该项目下有匹配任务，无法删除。")
	}

	if err = model.DeleteProject(p); err != nil {
		log.Errorf(ctx, "Delete project %s error:%v\n", p.Uuid, err)
	}
	return err
}

func (m *ProjectMgr) ProjectGet(ctx context.Context, id uint64) (a *types.Project, err error) {
	log.Debugln(ctx, "ProjectMgr.ProjectGet is called")
	p, err := model.GetProjectById(id)
	if err != nil {
		log.Errorf(ctx, "No project %d %v\n", id, err)
		return nil, err
	}
	pro := ToAPIProject(p)
	return pro, nil
}

func (m *ProjectMgr) ProjectJobsGet(ctx context.Context, puid string, pageNum int, pageSize int) ([]*types.Job, int64, int64) {
	log.Debugln(ctx, "ProjectMgr.ProjectJobsGet is called")
	jobs := make([]*types.Job, 0)
	js, count, err := model.ListJobsOfActivityByPage(puid, pageNum, pageSize)
	if err != nil {
		log.Errorf(ctx, "List jobs by project uuid %d error %v", puid, err)
		return jobs, 0, 0
	}
	for _, r := range js {
		j := toAPIJob(r)
		jobs = append(jobs, j)
	}
	pageCount := int64(math.Ceil((float64(count) / float64(pageSize))))
	return jobs, count, pageCount
}
