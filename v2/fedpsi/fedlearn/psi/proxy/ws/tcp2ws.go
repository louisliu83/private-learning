package wsproxy

import (
	"crypto/tls"
	"net"

	"fedlearn/psi/common/log"
	"fedlearn/psi/proxy"

	"github.com/gorilla/websocket"
)

type WrappedTCPConn struct {
	conn *net.TCPConn
	err  error
}

type TCP2WSProxy struct {
	target      string
	running     bool
	runningChan chan struct{}
}

func NewTCP2WSProxy(target string) *TCP2WSProxy {
	return &TCP2WSProxy{
		target:      target,
		running:     false,
		runningChan: nil,
	}
}

func (p *TCP2WSProxy) IsRunning() bool {
	return p.running
}

func (p *TCP2WSProxy) GetTarget() string {
	return p.target
}

func (p *TCP2WSProxy) Start() error {
	p.runningChan = make(chan struct{}, 1)
	p.running = true

	tcpAddr, err := net.ResolveTCPAddr("tcp", proxy.DefaultEgressListener)
	if err != nil {
		log.Errorf(ctx, "resolve tcp address %s error:%v", proxy.DefaultEgressListener, err)
		return err
	}

	l, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Errorf(ctx, "listen at tcp address %s error:%v", tcpAddr.String(), err)
		return err
	}

	for p.running {
		wcc := make(chan WrappedTCPConn, 1)

		go func() {
			conn, err := l.AcceptTCP()
			wcc <- WrappedTCPConn{conn, err}
		}()

		select {
		case wc := <-wcc:
			if wc.err != nil {
				log.Errorln(ctx, "accept tcp connectioin error:", err)
			} else {
				if err := wc.conn.SetNoDelay(true); err != nil {
					log.Errorln(ctx, "set tcp connection NoDelay error:", err)
				}
				handleTCPConn(wc.conn, p.target)
			}

		case <-p.runningChan:
			return nil
		}
	}
	return nil
}

func (p *TCP2WSProxy) Stop() error {
	if !p.running {
		return nil
	}
	p.running = false
	if p.runningChan != nil {
		close(p.runningChan)
	}
	return nil
}

func dialWS(wsAddr string) (*websocket.Conn, error) {
	d := websocket.Dialer{
		TLSClientConfig: &tls.Config{
			RootCAs:            nil,
			InsecureSkipVerify: true,
		},
	}
	wsConn, _, err := d.DialContext(ctx, wsAddr, nil)
	if err != nil {
		log.Errorf(ctx, "dial web socket address %s error:%v", wsAddr, err)
		return nil, err
	}
	return wsConn, nil
}

func handleTCPConn(c net.Conn, wsAddr string) error {
	ws, err := dialWS(wsAddr)
	if err != nil {
		log.Errorf(ctx, "dial web socket address %s error:%v", wsAddr, err)
		return err
	}
	copyAround(ws, c)
	return nil
}
