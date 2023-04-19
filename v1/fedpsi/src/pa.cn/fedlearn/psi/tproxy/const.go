package tproxy

import (
	"context"
	"fmt"
)

var (
	ctx context.Context = context.Background()
)

var (
	//Default Ingress
	DefaultIngressTargetIP   = "127.0.0.1"
	DefaultIngressTargetPort = 17766
	DefaultIngressTarget     = fmt.Sprintf("%s:%d", DefaultIngressTargetIP, DefaultIngressTargetPort)

	//Default Egress
	DefaultEgressListenerIP   = "127.0.0.1"
	DefaultEgressListenerPort = 17767
	DefaultEgressListener     = fmt.Sprintf("%s:%d", DefaultEgressListenerIP, DefaultEgressListenerPort)
)
