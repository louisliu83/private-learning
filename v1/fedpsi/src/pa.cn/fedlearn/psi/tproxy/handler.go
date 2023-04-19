package tproxy

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"pa.cn/fedlearn/psi/prom"

	"pa.cn/fedlearn/psi/log"
)

func (dp *Proxy) MagicStripHandler(src net.Conn, targetAddr string) error {
	defer goCloseConn(src)

	buf := make([]byte, len(MAGIC))
	src.SetDeadline(time.Now().Add(dp.serverWaitDataTimeout()))
	_, err := io.ReadAtLeast(src, buf, len(MAGIC))
	if err != nil {
		log.Errorln(ctx, err)
		return err
	}
	magic := string(buf)

	if magic != MAGIC {
		log.Warningln(ctx, "Unexpected client")
		return errors.New("Unexpected client")
	} else {
		src.SetDeadline(time.Time{})
		log.Infoln(ctx, "Accepted and strip:", magic)

		dst, err := dp.dialTarget()
		if err != nil {
			return err
		}
		defer goCloseConn(dst)
		dp.setKeepAlive(src, dst)
		errc := make(chan error, 1)
		go proxyCopy(errc, src, dst)
		go proxyCopy(errc, dst, src)

		errMsg := <-errc
		log.Errorln(ctx, errMsg)
		return errMsg
	}
}

func (dp *Proxy) MagicPrefixHandler(src net.Conn, targetAddr string) error {
	defer goCloseConn(src)

	dst, err := dp.dialTarget()
	if err != nil {
		return err
	}

	defer goCloseConn(dst)

	dp.setKeepAlive(src, dst)

	src.SetDeadline(time.Time{})
	io.WriteString(dst, MAGIC)

	errc := make(chan error, 1)
	go proxyCopy(errc, src, dst)
	go proxyCopy(errc, dst, src)
	errMsg := <-errc
	log.Errorln(ctx, errMsg)
	return errMsg
}

func (dp *Proxy) dialTarget() (net.Conn, error) {
	ctx := context.Background()
	var cancel context.CancelFunc
	if dp.DialTimeout >= 0 {
		ctx, cancel = context.WithTimeout(ctx, dp.dialTimeout())
	}

	var defaultDialer = new(net.Dialer)
	dst, err := defaultDialer.DialContext(ctx, "tcp", dp.TargetAddr)
	if cancel != nil {
		cancel()
	}
	if err != nil {
		log.Errorln(ctx, err)
		return nil, err
	}
	return dst, nil
}

func (dp *Proxy) setKeepAlive(src, dst net.Conn) error {
	if c, ok := dst.(*net.TCPConn); ok {
		c.SetKeepAlive(true)
		c.SetKeepAlivePeriod(dp.keepAlivePeriod())
	}
	if c, ok := src.(*net.TCPConn); ok {
		c.SetKeepAlive(true)
		c.SetKeepAlivePeriod(dp.keepAlivePeriod())
	}
	return nil
}

// proxyCopy is the function that copies bytes around.
// It's a named function instead of a func literal so users get
// named goroutines in debug goroutine stack dumps.
func proxyCopy(errc chan<- error, dst, src net.Conn) {
	n, err := io.Copy(dst, src)
	connectionTag := fmt.Sprintf("%s->%s", src.RemoteAddr().String(), dst.RemoteAddr().String())
	prom.AddProxyDataBytes(connectionTag, n)
	errc <- err
}

func goCloseConn(c net.Conn) {
	go c.Close()
}
