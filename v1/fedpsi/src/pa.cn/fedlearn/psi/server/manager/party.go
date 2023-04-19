package manager

import (
	"context"
	"fmt"

	"pa.cn/fedlearn/psi/api/types"
	"pa.cn/fedlearn/psi/log"
	"pa.cn/fedlearn/psi/model"
)

type PartyMgr struct {
}

func (m *PartyMgr) PartyRegister(ctx context.Context, r *types.PartyRegisterRequest) error {
	log.Debugln(ctx, "PartyMgr.PartyRegister is called")
	p := model.Party{
		Name:             r.Name,
		Scheme:           r.Scheme,
		ControllerServer: r.ControllerServer,
		ControllerPort:   r.ControllerPort,
		WorkServer:       r.WorkServer,
		WorkPort:         r.WorkPort,
		Token:            r.Token,
	}
	if p.WorkServer == "" {
		p.WorkServer = r.ControllerServer
	}
	return model.AddParty(&p)
}

func (m *PartyMgr) PartyUpdate(ctx context.Context, r *types.PartyUpdateRequest) error { // Only update token
	log.Debugln(ctx, "PartyMgr.PartyUpdate is called")
	p, err := model.GetPartyByName(r.Name)
	if err != nil {
		log.Errorf(ctx, "No party named %s: %v", r.Name, err)
		return fmt.Errorf("No party named %s: %w", r.Name, err)
	}
	if p == nil {
		log.Errorf(ctx, "No party named %s", r.Name)
		return fmt.Errorf("No party named %s", r.Name)
	}
	p.Token = r.Token
	err = model.UpdateParty(p)
	if err != nil {
		log.Errorf(ctx, "Update party %s token error:%v", r.Name, err)
		return fmt.Errorf("Update party %s token error:%w", r.Name, err)
	}
	return nil
}

func (m *PartyMgr) PartyList(ctx context.Context) []*types.PartyInfo {
	log.Debugln(ctx, "PartyMgr.PartyList is called")
	apiParties := make([]*types.PartyInfo, 0)
	ps, err := model.ListParties()
	if err != nil {
		log.Errorf(ctx, "List paryies error %v", err)
		return apiParties
	}
	for _, r := range ps {
		apiP := types.PartyInfo{
			Name:             r.Name,
			Scheme:           r.Scheme,
			ControllerServer: r.ControllerServer,
			ControllerPort:   r.ControllerPort,
			WorkServer:       r.WorkServer,
			WorkPort:         r.WorkPort,
		}
		apiParties = append(apiParties, &apiP)
	}
	return apiParties
}

func (m *PartyMgr) PartyDel(ctx context.Context, name string) error {
	log.Debugln(ctx, "PartyMgr.PartyDel is called")
	p, err := model.GetPartyByName(name)
	if err != nil {
		log.Errorf(ctx, "No party %s %v\n", name, err)
		return err
	}
	if err = model.DeleteParty(p); err != nil {
		log.Errorf(ctx, "Delete party %s error:%v\n", name, err)
	}
	return err
}

func (m *PartyMgr) PartyReady(ctx context.Context, name string) error {
	log.Debugln(ctx, "PartyMgr.PartyReady is called")
	remoteClient, err := getRemoteClient(name)
	if err != nil {
		log.Errorf(ctx, "Get party client for %s error:%v", name, err)
		return fmt.Errorf("Get party client for %s error:%w", name, err)
	}
	si, err := remoteClient.PartyStatus(ctx)
	if err != nil {
		log.Errorf(ctx, "Get party status for %s error:%v", name, err)
		return err
	}
	log.Infof(ctx, "Party %s is OK", si.PartyName)
	return err
}
