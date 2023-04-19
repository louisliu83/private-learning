package wsproxy

import (
	"context"

	"fedlearn/psi/proxy"
)

var (
	ctx context.Context = proxy.ContexForProxy("websocket")
)
