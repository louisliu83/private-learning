package proxy

import (
	"context"
	"fmt"

	"fedlearn/psi/api"
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

type TraceKey string

func ContexForProxy(name string) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, TraceKey(api.Trace_ID), fmt.Sprintf("Proxy_%s", name))
	return ctx
}
