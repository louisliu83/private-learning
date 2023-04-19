package tproxy

import (
	"sync"
	"time"

	"pa.cn/fedlearn/psi/config"
	"pa.cn/fedlearn/psi/log"
)

var (
	psiProxy     *PSITProxy
	psiProxyLock = sync.Mutex{}
)

func GetTProxy() *PSITProxy {
	psiProxyLock.Lock()
	defer psiProxyLock.Unlock()
	if psiProxy == nil {
		psiProxy = newPSIProxy()
	}
	return psiProxy
}

// PSITProxy definition
type PSITProxy struct {
	ingressProxy *Proxy
	egressProxy  *Proxy
}

func newPSIProxy() *PSITProxy {
	tp := &PSITProxy{}

	tp.ingressProxy = NewIngressProxy(config.GetConfig().TProxy.Listen, DefaultIngressTarget)
	tp.egressProxy = NewEgressProxy(DefaultEgressListener, config.GetConfig().TProxy.Target)

	return tp
}

func (p *PSITProxy) Start() {
	if !config.GetConfig().TProxy.DisableServer {
		p.StartServerProxy()
	}
	if !config.GetConfig().TProxy.DisableClient {
		p.StartClientProxy()
	}
}

func (p *PSITProxy) Stop() {
	p.StopServerProxy()
	p.StopClientProxy()
}

func (p *PSITProxy) GetServerProxy() *Proxy {
	return p.ingressProxy
}

func (p *PSITProxy) SetServerProxy(sp *Proxy) {
	p.ingressProxy = sp
}

func (p *PSITProxy) IsServerProxyRunning() bool {
	return p.ingressProxy != nil && p.ingressProxy.running
}

func (p *PSITProxy) GetClientProxy() *Proxy {
	return p.egressProxy
}

func (p *PSITProxy) SetClientProxy(cp *Proxy) {
	p.egressProxy = cp
}

func (p *PSITProxy) IsClientProxyRunning() bool {
	return p.egressProxy != nil && p.egressProxy.running
}

func (p *PSITProxy) StartServerProxy() {
	if p.ingressProxy != nil {
		log.Infof(ctx, "Start ingress proxy, serving at %s, will target to %s \n", p.ingressProxy.ListenAddr, p.ingressProxy.TargetAddr)
		go func() {
			log.Errorln(ctx, p.ingressProxy.Start())
		}()
	}
}

func (p *PSITProxy) StopServerProxy() error {
	if p.ingressProxy != nil {
		return p.ingressProxy.Stop()
	}
	return nil

}

func (p *PSITProxy) StartClientProxy() {
	if p.egressProxy != nil {
		log.Infof(ctx, "Start outgress proxy, serving at %s, will target to %s \n", p.egressProxy.ListenAddr, p.egressProxy.TargetAddr)
		go func() {
			log.Errorln(ctx, p.egressProxy.Start())
		}()
	}
}

func (p *PSITProxy) StopClientProxy() error {
	if p.egressProxy != nil {
		return p.egressProxy.Stop()
	}
	return nil
}

func setConfig(pxy *Proxy) {
	if pxy == nil {
		return
	}
	pxy.SetDialTimeout(
		time.Duration(config.GetConfig().TProxy.DialTimeout) * time.Second)
	pxy.SetKeepAlivePeriod(
		time.Duration(config.GetConfig().TProxy.KeepAlivePeriod) * time.Second)
	pxy.SetServerWaitDataTimeout(
		time.Duration(config.GetConfig().TProxy.ServerWaitDataTimeout) * time.Second)
}
