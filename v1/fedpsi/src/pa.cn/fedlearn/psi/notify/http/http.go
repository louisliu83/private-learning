package http

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"

	"pa.cn/fedlearn/psi/client/httpc"
	"pa.cn/fedlearn/psi/log"
	"pa.cn/fedlearn/psi/notify"
)

type HttpNotifier struct {
	TargetPostURL string
	Token         string
	User          string
	Password      string
}

var _ notify.Notifier = &HttpNotifier{}

func (n *HttpNotifier) Notify(ctx context.Context, msg string) error {
	headers := map[string]string{}
	if n.Token != "" {
		headers["Authorization"] = fmt.Sprintf("%s %s", "Bearer", n.Token)
	} else {
		userPass := bytes.NewBufferString(fmt.Sprintf("%s:%s", n.User, n.Password)).Bytes()
		headers["Authorization"] = fmt.Sprintf("%s %s", "Basic", base64.StdEncoding.EncodeToString(userPass))
	}

	_, err := httpc.DoPostWithJson(n.TargetPostURL, headers, bytes.NewBufferString(msg).Bytes())
	if err != nil {
		log.Errorln(ctx, err)
	}
	return err
}
