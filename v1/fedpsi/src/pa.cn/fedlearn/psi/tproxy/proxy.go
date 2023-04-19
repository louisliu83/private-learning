package tproxy

import (
	"errors"
	"net"
	"time"

	"pa.cn/fedlearn/psi/log"
)

const (
	MAGIC = "!!PAB##PSI!!"
)

type ConnHandler func(src net.Conn, targetAddr string) error

type Proxy struct {
	running               bool
	acceptingChan         chan struct{}
	ListenAddr            string
	TargetAddr            string
	ConnHandler           ConnHandler
	DialTimeout           time.Duration
	KeepAlivePeriod       time.Duration
	ServerWaitDataTimeout time.Duration
	_                     bool
}

func New(la, ta string) *Proxy {
	return &Proxy{
		running:    false,
		ListenAddr: la,
		TargetAddr: ta,
	}
}

func NewIngressProxy(la, ta string) *Proxy {
	p := New(la, ta)
	p.SetConnHandler(p.MagicStripHandler)
	setConfig(p)
	return p
}

func NewEgressProxy(la, ta string) *Proxy {
	p := New(la, ta)
	p.SetConnHandler(p.MagicPrefixHandler)
	setConfig(p)
	return p
}

func (p *Proxy) SetDialTimeout(d time.Duration) {
	p.DialTimeout = d
}

func (p *Proxy) SetServerWaitDataTimeout(d time.Duration) {
	p.ServerWaitDataTimeout = d
}

func (p *Proxy) SetKeepAlivePeriod(d time.Duration) {
	p.KeepAlivePeriod = d
}

func (p *Proxy) SetConnHandler(h ConnHandler) {
	p.ConnHandler = h
}

func (p *Proxy) Start() error {
	if p.running {
		return errors.New("already running")
	}

	p.running = true
	p.acceptingChan = make(chan struct{}, 1)

	tcpAddr, err := net.ResolveTCPAddr("tcp", p.ListenAddr)
	if err != nil {
		return err
	}

	l, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}
	defer l.Close()

	type accepted struct {
		conn *net.TCPConn
		err  error
	}

	for p.running {
		c := make(chan accepted, 1)
		go func() {
			conn, err := l.AcceptTCP()
			c <- accepted{conn, err}
		}()
		select {
		case a := <-c:
			if a.err != nil {
				log.Errorln(ctx, err)
			} else {
				if err := a.conn.SetNoDelay(true); err != nil {
					log.Errorln(ctx, err)
				}
				if p.ConnHandler != nil {
					go p.ConnHandler(a.conn, p.TargetAddr)
				}
			}

		case <-p.acceptingChan:
			return nil
		}
	}
	return nil
}

func (p *Proxy) Stop() error {
	if !p.running {
		return errors.New("already stopped")
	}
	p.running = false
	close(p.acceptingChan)
	return nil
}

func (p *Proxy) dialTimeout() time.Duration {
	if p.DialTimeout > 0 {
		return p.DialTimeout
	}
	return 15 * time.Second
}

func (p *Proxy) keepAlivePeriod() time.Duration {
	if p.KeepAlivePeriod > 0 {
		return p.KeepAlivePeriod
	}
	return 120 * time.Second
}

func (p *Proxy) serverWaitDataTimeout() time.Duration {
	if p.ServerWaitDataTimeout > 0 {
		return p.ServerWaitDataTimeout
	}
	return 15 * time.Second
}
