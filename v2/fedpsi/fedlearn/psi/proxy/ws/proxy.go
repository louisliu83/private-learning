package wsproxy

import (
	"net"

	"fedlearn/psi/common/log"

	"github.com/gorilla/websocket"
)

// copyAround copy data stream between websocket connection and tcp connection
func copyAround(wsConn *websocket.Conn, tcpConn net.Conn) {
	go func() { // copy data from websocket connection to tcp connection
		for {
			if t, data, err := wsConn.ReadMessage(); err != nil {
				log.Errorf(ctx, "ws-read from websocket connection error:%v", err)
				break
			} else {
				if t == websocket.BinaryMessage {
					if _, err := tcpConn.Write(data); err != nil {
						log.Errorf(ctx, "ws-write to tcp connection error:%v", err)
						break
					}
				}
			}
		}
	}()

	go func() { // copy data from tcp connection to websocket connection
		for {
			buf := make([]byte, 4*1024*1024)
			var n int
			var err error
			if n, err = tcpConn.Read(buf); err != nil {
				log.Errorf(ctx, "ws-read from tcp connection error:%v", err)
				break
			}
			if err = wsConn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
				log.Errorf(ctx, "ws-write to websocket connection error:%v", err)
				break
			}
		}
	}()
}
