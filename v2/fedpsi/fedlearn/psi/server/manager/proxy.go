package manager

import (
	"fmt"
	"time"

	"fedlearn/psi/model"
	grpcproxy "fedlearn/psi/proxy/grpc"
)

func StartEgressListenerForTask(task *model.Task) error {
	go func() {
		startTCP2GRPCProxy(task)
	}()
	time.Sleep(time.Duration(5) * time.Second)
	return nil
}

func startTCP2GRPCProxy(task *model.Task) error {
	p, err := model.GetPartyByName(task.PartyName)
	if err != nil {
		return err
	}
	targetGrpcAddress := fmt.Sprintf("%s:%d", p.WorkServer, p.WorkPort)
	if _, err := grpcproxy.StartTCP2GrpcServer(task.Uuid, targetGrpcAddress); err != nil {
		return err
	}
	return nil
}
