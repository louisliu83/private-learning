package types

type TProxyStartRequest struct {
	ProxyType string `json:"proxyType"`
	PartyName string `json:"partyName"`
}

type TProxyStopRequest struct {
	ProxyType string `json:"proxyType"`
}

type TProxyStatusResponse struct {
	IngressProxyStatus string `json:"ingress"`
	EgressProxyStatus  string `json:"egress"`
	EgressProxyTarget  string `json:"egressTarget"`
}
