package portmgr

import (
	"sync"
	"time"

	"fedlearn/psi/common/utils"

	"github.com/sirupsen/logrus"
)

const (
	EmptyTask = "EmptyTask"
)

var (
	PSIServerPortManager *PortManager
	PSIClientPortManager *PortManager
)

var (
	taskPortMap     = map[string]int{}
	taskPortMapLock sync.Mutex
)

func addTaskPort(taskUID string, port int) {
	taskPortMapLock.Lock()
	defer taskPortMapLock.Unlock()
	taskPortMap[taskUID] = port
}

func rmTaskPort(taskUID string) {
	taskPortMapLock.Lock()
	defer taskPortMapLock.Unlock()
	delete(taskPortMap, taskUID)
}

func GetPortOfTask(taskUID string) int {
	p := taskPortMap[taskUID]
	return p
}

type PortManager struct {
	portMapLock  sync.Mutex
	portUsingMap map[int]string
}

func NewPortManager(portStart, portNum int) *PortManager {
	portUsingMap := map[int]string{}

	for i := 0; i < portNum; i++ {
		portUsingMap[portStart+i] = EmptyTask
	}

	portManager := &PortManager{
		portUsingMap: portUsingMap,
	}
	return portManager
}

func (pm *PortManager) ReleasePort(p int) {
	releasePort := func() bool {
		pm.portMapLock.Lock()
		defer pm.portMapLock.Unlock()
		taskUID, ok := pm.portUsingMap[p]
		if !ok { // no task for this port, return
			rmTaskPort(taskUID)
			return true
		}

		if !utils.CheckPortInUse(int32(p)) {
			pm.portUsingMap[p] = EmptyTask
			logrus.Infoln("relase port:", p)
			rmTaskPort(taskUID)
			return true
		}
		return false
	}

	for i := 1; i <= 10; i++ {
		if releasePort() {
			return
		}
		time.Sleep(time.Duration(2) * time.Second)
	}
}

func (pm *PortManager) AcquireAvailablePort(taskUID string) int {
	pm.portMapLock.Lock()
	defer pm.portMapLock.Unlock()
	for p, t := range pm.portUsingMap {
		if t == EmptyTask {
			if !utils.CheckPortInUse(int32(p)) {
				logrus.Infof("get available port %d for task %s", p, taskUID)
				pm.portUsingMap[p] = taskUID
				addTaskPort(taskUID, p)
				return p
			} else {
				logrus.Infoln("port is in use:", p)
			}
		} else {
			logrus.Infoln("port is in use:", p)
		}
	}
	logrus.Warningf("no available port for task %s", taskUID)
	return -1
}

func (pm *PortManager) AcquireAvailablePortWithRetry(taskUID string, retry int) int {
	for i := 1; i <= retry; i++ {
		if p := pm.AcquireAvailablePort(taskUID); p > 0 {
			return p
		}
		time.Sleep(time.Duration(10) * time.Second)
	}
	logrus.Warningf("no available port for task %s after retry %d times", taskUID, retry)
	return -1
}
