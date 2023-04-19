package proxy

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"fedlearn/psi/common/portmgr"
	"fedlearn/psi/proxy/grpc/proto"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	TASK_UUID_KEY = "TaskUUID"
)

var (
	tcp2grpcProxyMap  = map[string]*TCP2GrpcServer{}
	tcp2grpcProxyLock sync.Mutex
)

// TCP2GrpcServer to proxy TCP traffic to gRPC
type TCP2GrpcServer struct {
	tcpServerPort     int
	tcpServerAddress  string
	targetGrpcAddress string
	taskUID           string
	listener          net.Listener
}

// NewTCP2GrpcServer constructs a TCP2GrpcServer
func NewTCP2GrpcServer(tcpServerPort int, targetGrpcAddress, taskUID string) *TCP2GrpcServer {
	tcpServerAddress := fmt.Sprintf("%s:%d", "127.0.0.1", tcpServerPort)
	return &TCP2GrpcServer{
		tcpServerPort:     tcpServerPort,
		tcpServerAddress:  tcpServerAddress,
		targetGrpcAddress: targetGrpcAddress,
		taskUID:           taskUID,
	}
}

func (s *TCP2GrpcServer) handleTCPConn(tcpConn net.Conn) error {
	logrus.Infoln("handle tcp connection, target to:", s.targetGrpcAddress)

	grpcConn, err := grpc.Dial(s.targetGrpcAddress, grpc.WithInsecure())
	if err != nil {
		logrus.Errorf("connect to grpc %s failed %v", s.targetGrpcAddress, err)
		return fmt.Errorf("connect to grpc %s failed %w", s.targetGrpcAddress, err)
	}
	defer grpcConn.Close()
	md := metadata.Pairs(TASK_UUID_KEY, s.taskUID)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	grpcClient := proto.NewTunnelServiceClient(grpcConn)
	stream, err := grpcClient.Tunnel(ctx)

	if err != nil {
		logrus.Errorln("gRPC client tunnel remote service failed ", err)
		return fmt.Errorf("gRPC client tunnel remote service failed %w", err)
	}

	// Gets data from remote gRPC server and proxy to TCP client
	go func() {
		for {
			chunk, err := stream.Recv()
			if err != nil {
				logrus.Errorf("Recv from grpc target %s terminated: %v", s.targetGrpcAddress, err)
				return
			}
			logrus.Infof("Sending %d bytes to TCP client", len(chunk.Data))
			tcpConn.Write(chunk.Data)
		}
	}()

	// Gets data from TCP client and proxy to remote gRPC server
	func() {
		for {
			tcpData := make([]byte, 64*1024)
			bytesRead, err := tcpConn.Read(tcpData)

			if err == io.EOF {
				logrus.Infoln("Connection finished")
				return
			}
			if err != nil {
				logrus.Errorln("Read from tcp error: ", err)
			}
			logrus.Infof("Sending %d bytes to gRPC server", bytesRead)
			if err := stream.Send(&proto.Chunk{Data: tcpData[0:bytesRead]}); err != nil {
				logrus.Errorln("Failed to send gRPC data: ", err)
			}
		}
	}()

	return nil
}

// Start Starts the server
func (s *TCP2GrpcServer) Start() error {
	ln, err := net.Listen("tcp", s.tcpServerAddress)
	if err != nil {
		logrus.Errorf("listen on tcp %s error: %v", s.tcpServerAddress, err)
		return err
	}

	s.listener = ln
	defer ln.Close()
	logrus.Infof("run TCPServer on %s target to %s for task %s", s.tcpServerAddress, s.targetGrpcAddress, s.taskUID)
	//
	tcp2grpcProxyLock.Lock()
	tcp2grpcProxyMap[s.taskUID] = s
	tcp2grpcProxyLock.Unlock()
	//
	for {
		conn, err := ln.Accept()
		if err != nil {
			logrus.Errorln("TCP listener error:", err)
			return err
		}
		go s.handleTCPConn(conn)
	}
}

func (s *TCP2GrpcServer) Stop() error {
	defer func() {
		portmgr.PSIClientPortManager.ReleasePort(s.tcpServerPort)
	}()

	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func StartTCP2GrpcServer(taskUID, targetGrpcAddress string) (int, error) {
	//Stop the server if exists
	StopTCP2GrpcServer(taskUID)
	port := portmgr.PSIClientPortManager.AcquireAvailablePortWithRetry(taskUID, 6)
	if port < 0 {
		return port, errors.New("no available port for psi client tcp proxy")
	}
	tcp2grpcServer := NewTCP2GrpcServer(port, targetGrpcAddress, taskUID)
	tcp2grpcServer.Start()
	logrus.Infof("tcp2grpc server for task %s started on: %d, target to:%s", taskUID, port, targetGrpcAddress)
	return port, nil
}

func StopTCP2GrpcServer(taskUID string) error {
	if s, ok := tcp2grpcProxyMap[taskUID]; ok && s != nil {
		tcp2grpcProxyLock.Lock()
		err := s.Stop()
		tcp2grpcProxyLock.Unlock()
		return err
	}
	return nil
}
