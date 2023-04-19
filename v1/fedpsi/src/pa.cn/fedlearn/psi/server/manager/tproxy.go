package manager

import (
	"context"
	"fmt"

	"pa.cn/fedlearn/psi/model"

	"pa.cn/fedlearn/psi/api/types"
	"pa.cn/fedlearn/psi/log"
	"pa.cn/fedlearn/psi/tproxy"
)

type TProxyManager struct {
}

func (p *TProxyManager) Start(ctx context.Context, r types.TProxyStartRequest) error {
	log.Debugln(ctx, "TProxyManager.Start is called")
	if r.ProxyType == "server" {
		tproxy.GetTProxy().StartServerProxy()
	} else if r.ProxyType == "client" {
		party, err := model.GetPartyByName(r.PartyName)
		if err != nil {
			return err
		}
		go func() {
			log.Infof(ctx, "Restart client proxy %s target to %s", tproxy.DefaultEgressListener, fmt.Sprintf("%s:%d", party.WorkServer, party.WorkPort))
			tproxy.RestartClientProxy(
				tproxy.DefaultEgressListener,
				fmt.Sprintf("%s:%d", party.WorkServer, party.WorkPort))
		}()
	}
	return nil
}

func (p *TProxyManager) Stop(ctx context.Context, r types.TProxyStopRequest) error {
	log.Debugln(ctx, "TProxyManager.Stop is called")
	if r.ProxyType == "server" {
		tproxy.GetTProxy().StopServerProxy()
	} else if r.ProxyType == "client" {
		tproxy.GetTProxy().StopClientProxy()
	}
	return nil
}

const (
	ProxyON  = "ON"
	ProxyOFF = "OFF"
)

func (p *TProxyManager) Status(ctx context.Context) *types.TProxyStatusResponse {
	t := &types.TProxyStatusResponse{}
	if tproxy.GetTProxy().IsServerProxyRunning() {
		t.IngressProxyStatus = ProxyON
	} else {
		t.IngressProxyStatus = ProxyOFF
	}

	if tproxy.GetTProxy().IsClientProxyRunning() {
		t.EgressProxyStatus = ProxyON
		t.EgressProxyTarget = getPartyNameByTargetAddress(tproxy.GetTProxy().GetClientProxy().TargetAddr)
	} else {
		t.EgressProxyStatus = ProxyOFF
	}
	return t
}

func getPartyNameByTargetAddress(targetAdd string) string {
	parties, err := model.ListParties()
	if err != nil {
		return "UnknownParty"
	}
	for _, party := range parties {
		partyAdd := fmt.Sprintf("%s:%d", party.WorkServer, party.WorkPort)
		if partyAdd == targetAdd {
			return party.Name
		}
	}
	return "UnknownParty"
}
