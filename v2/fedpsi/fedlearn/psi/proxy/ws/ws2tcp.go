package wsproxy

import (
	"context"
	"net"
	"net/http"
	"time"

	"fedlearn/psi/common/log"

	"fedlearn/psi/proxy"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024 * 1024 * 4,
		WriteBufferSize: 1024 * 1024 * 4,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugln(ctx, "Accept websocket message")
	if wsConn, err := upgrader.Upgrade(w, r, nil); err != nil {
		log.Errorf(context.Background(), "websocket handler error:%v", err)
		return
	} else {
		handleWSConn(wsConn, proxy.DefaultIngressTarget)
	}
}

func dialTcpTarget(targetTCP string) (net.Conn, error) {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Duration(5)*time.Second)

	var defaultDialer = new(net.Dialer)

	dst, err := defaultDialer.DialContext(ctx, "tcp", targetTCP)
	if cancel != nil {
		cancel()
	}

	if err != nil {
		log.Errorf(ctx, "dail tcp address %s error: %v", targetTCP, err)
		return nil, err
	}
	return dst, nil
}

func handleWSConn(wsConn *websocket.Conn, targetTCP string) error {
	tcpConn, err := dialTcpTarget(targetTCP)
	if err != nil {
		log.Errorf(ctx, "dail tcp address %s error: %v", targetTCP, err)
		return err
	}
	copyAround(wsConn, tcpConn)
	return nil
}
