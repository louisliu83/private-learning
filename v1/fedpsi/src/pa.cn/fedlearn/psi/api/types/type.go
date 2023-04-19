package types

type ServerInfo struct {
	Version   string   `json:"version"`
	PartyName string   `json:"partyName"`
	Protocols []string `json:"protocols"`
	Status    string   `json:"status"`
}

type ConfigInfo struct {
	Listen                string `json:"tproxy.listen"`
	Target                string `json:"tproxy.target"`
	DialTimeout           int64  `json:"tproxy.dialTimeout"`
	KeepAlivePeriod       int64  `json:"tproxy.keepAlivePeriod"`
	ServerWaitDataTimeout int64  `json:"tproxy.serverWaitDataTimeout"`
	PublicIP              string `json:"psi.publicIP"`
	PrivateIP             string `json:"psi.privateIP"`
	PublicPort            int32  `json:"psi.publicPort"`
	PrivatePort           int32  `json:"psi.privatePort"`
	BinPath               string `json:"psi.binpath"`
}
