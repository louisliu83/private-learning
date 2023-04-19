package worker

import (
	"pa.cn/fedlearn/psi/config"
	"pa.cn/fedlearn/psi/model"
	"pa.cn/fedlearn/psi/tproxy"
)

func GetAvailableIpAndPort() (string, int32) {
	if config.GetConfig().FeatureGate.TProxy {
		return tproxy.DefaultIngressTargetIP, int32(tproxy.DefaultIngressTargetPort)
	}
	return config.GetConfig().PsiExecutor.PrivateIP, config.GetConfig().PsiExecutor.PrivatePort
}

func GetPublicIPAndPort() (string, int32) {
	if config.GetConfig().FeatureGate.TProxy {
		return tproxy.DefaultEgressListenerIP, int32(tproxy.DefaultEgressListenerPort)
	}
	return config.GetConfig().PsiExecutor.PublicIP, config.GetConfig().PsiExecutor.PublicPort
}

// Get the ip and port to connect when start the psi executor as client mode
// Normally it should be the task.ServerIP and task.ServerPort, but if we use egress proxy, we get the egress ip and port
func GetTargetIpAndPort(task *model.Task) (string, int32) {
	if !config.GetConfig().FeatureGate.TProxy {
		return task.ServerIP, task.ServerPort
	}
	return GetEgressIPAndPortOfParty(task.PartyName)
}

func GetEgressIPAndPortOfParty(partyName string) (string, int32) {
	// TODO: fix this
	ip := config.GetConfig().PsiExecutor.PublicIP
	port := config.GetConfig().PsiExecutor.PublicPort
	return ip, port
}
