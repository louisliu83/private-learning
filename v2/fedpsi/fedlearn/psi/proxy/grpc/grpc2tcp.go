package proxy

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"fedlearn/psi/common/portmgr"

	"fedlearn/psi/proxy/grpc/proto"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
)

func newGrpcServer(withCustomParam bool) *grpc.Server {
	if !withCustomParam {
		return grpc.NewServer()
	}

	kaep := keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
		PermitWithoutStream: true,            // Allow pings even when there are no active streams
	}
	kasp := keepalive.ServerParameters{
		// 当连接处于idle的时长超过 MaxConnectionIdle时，服务端就发送GOAWAY，关闭连接。默认值为无限大
		//连接的最大空闲时长。当超过这个时间时，服务端会向客户端发送GOAWAY帧，关闭空闲的连接，节省连接数。
		MaxConnectionIdle: 3600 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
		// 一个连接只能使用 MaxConnectionAge 这么长的时间，服务端就会关闭这个连接。
		//一个连接可以使用的时间。当一个连接已经使用了超过这个值的时间时，服务端就要强制关闭连接了。如果客户端仍然要连接服务端，可以重新发起连接。这时连接将进入半关闭状态，不再接收新的流。
		MaxConnectionAge: 60 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
		// 服务端优雅关闭连接时长
		//当服务端决定关闭一个连接时，如果有RPC在进行，会等待MaxConnectionAgeGrace时间，让已经存在的流可以正常处理完毕。
		MaxConnectionAgeGrace: 10 * time.Second, // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
		// 这个时间是服务端用来ping 客户端的。默认值为2小时
		Time: 10 * time.Second, // Ping the client if it is idle for 5 seconds to ensure the connection is still active
		// 默认值为20秒
		Timeout: 30 * time.Second, // Wait 1 second for the ping ack before assuming the connection is dead
	}
	return grpc.NewServer(grpc.KeepaliveEnforcementPolicy(kaep), grpc.KeepaliveParams(kasp))
}

// Grpc2TCPServer A server to proxy grpc traffic to TCP
type Grpc2TCPServer struct {
	proto.UnimplementedTunnelServiceServer
	address string
	server  *grpc.Server
}

// NewGrpc2TCPServer constructs a Grpc2TCP server
func NewGrpc2TCPServer(address string) *Grpc2TCPServer {
	return &Grpc2TCPServer{
		address: address,
	}
}

// Start starts the Grpc2TCP server
func (s *Grpc2TCPServer) Start() error {
	ln, err := net.Listen("tcp", s.address)
	if err != nil {
		logrus.Errorln("failed to listen: ", err)
		return err
	}

	// Starts a gRPC server and register services
	grpcServer := newGrpcServer(false)
	s.server = grpcServer
	proto.RegisterTunnelServiceServer(grpcServer, s)

	logrus.Infof("Starting gRPC server on: %s", s.address)
	if err := grpcServer.Serve(ln); err != nil {
		logrus.Errorln("Unable to start gRPC server:", err)
		return err
	}
	return nil
}

// Stop stop the Grpc2TCP server
func (s *Grpc2TCPServer) Stop() {
	if s.server != nil {
		s.server.Stop()
	}
}

// Tunnel the implementation of gRPC Tunnel service
func (s *Grpc2TCPServer) Tunnel(stream proto.TunnelService_TunnelServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context()) // get context from stream
	if !ok {
		return errors.New("no task bound to this stream")
	}

	taskUIDS := md.Get(TASK_UUID_KEY)
	if len(taskUIDS) == 0 {
		return errors.New("no task bound to this stream")
	}
	taskUID := taskUIDS[0]
	logrus.Infof("accept stream bound to task %s", taskUID)

	targetTCPPort := portmgr.GetPortOfTask(taskUID)
	if targetTCPPort <= 0 {
		return errors.New("no available port")
	}
	targetTCPAddress := fmt.Sprintf("127.0.0.1:%d", targetTCPPort)

	tcpConnection, err := net.Dial("tcp", targetTCPAddress)
	if err != nil {
		logrus.Errorf("dail to tcp target %s error: %v", targetTCPAddress, err)
		return fmt.Errorf("dail to tcp target %s error: %w", targetTCPAddress, err)
	}

	logrus.Infof("task %s connected to %s", taskUID, targetTCPAddress)
	// Makes sure the connection gets closed
	defer func() {
		logrus.Infoln("connection closed to ", targetTCPAddress)
		tcpConnection.Close()
	}()

	// Gets data from gRPC client and proxy TCP server
	errReadChan := make(chan error)
	go func() {
		for {
			chunk, err := stream.Recv()
			if chunk != nil {
				data := chunk.Data
				if len(data) > 0 {
					bytesWrote, err := tcpConnection.Write(data)
					if err != nil {
						errReadChan <- fmt.Errorf("error while sending TCP data: %v", err)
						return
					} else {
						logrus.Infof("sending %d bytes to tcp server", bytesWrote)
					}
				}
			} else {
				if err == io.EOF {
					errReadChan <- err
					return
				}
				if err != nil {
					errReadChan <- fmt.Errorf("error while receiving gRPC data: %v", err)
					return
				}
			}
		}
	}()

	// Gets data from remote TCP server and proxy to gRPC client
	errWriteChan := make(chan error)
	go func() {
		buff := make([]byte, 64*1024)
		for {
			bytesRead, err := tcpConnection.Read(buff)
			if bytesRead > 0 {
				if err = stream.Send(&proto.Chunk{Data: buff[0:bytesRead]}); err != nil {
					errWriteChan <- fmt.Errorf("error while sending gRPC data: %v", err)
					return
				} else {
					logrus.Infof("sending %d bytes to gRPC client", bytesRead)
				}
			}

			if err != nil {
				if err == io.EOF {
					errWriteChan <- io.EOF
					logrus.Infoln("remote TCP connection closed")
					return
				}
				errWriteChan <- fmt.Errorf("error while receiving TCP data: %v", err)
				return
			}
		}
	}()

	// Blocking copy around
	returnedError := <-errWriteChan
	if returnedError != io.EOF {
		logrus.Errorln("copy data from TCP server to gRPC client error:", returnedError)
	} else {
		logrus.Infoln("copy data from TCP server to gRPC client:", returnedError)
	}
	returnedError = <-errReadChan
	if returnedError != io.EOF {
		logrus.Errorln("copy data from gRPC client to TCP server error:", returnedError)
	} else {
		logrus.Infoln("copy data from gRPC client to TCP server:", returnedError)
	}

	return returnedError
}
