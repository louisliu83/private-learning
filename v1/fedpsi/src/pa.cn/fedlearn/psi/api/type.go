package api

type ServerInfo struct {
	PartyName string   `json:"partyName"`
	Protocols []string `json:"protocols"`
	Status    string   `json:"status"`
}
