package tproxy

import (
	"sync"
	"time"
)

var (
	restartClientLock = sync.Mutex{}
)

func RestartClientProxy(la, ta string) {
	restartClientLock.Lock()
	defer restartClientLock.Unlock()
	GetTProxy().StopClientProxy()
	time.Sleep(time.Duration(10) * time.Second)
	ncp := NewEgressProxy(la, ta)
	GetTProxy().SetClientProxy(ncp)
	GetTProxy().StartClientProxy()
}
