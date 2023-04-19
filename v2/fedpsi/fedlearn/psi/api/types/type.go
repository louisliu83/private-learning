package types

type ServerInfo struct {
	Version   string   `json:"version"`
	PartyName string   `json:"partyName"`
	Protocols []string `json:"protocols"`
	Status    string   `json:"status"`
}

type ConfigInfo struct {
}

type User struct {
	UserName    string `json:"username"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName"`
	Party       string `json:"party"`
}

type UserRegisterRequest struct {
	User
}

type TokenGenResponse struct {
	Token string `json:"token"`
}

type PartyInfo struct {
	Name             string `json:"name"`
	Scheme           string `json:"scheme"`
	ControllerServer string `json:"controllerServer"`
	ControllerPort   int32  `json:"controllerPort"`
	WorkServer       string `json:"workServer"`
	WorkPort         int32  `json:"workPort"`
	Token            string `json:"token"`
}

type PartyRegisterRequest struct {
	PartyInfo
}

type PartyUpdateRequest struct {
	PartyInfo
}

type Channel struct {
	ChannelCode    string `json:"channel_code"`
	ChannelName    string `json:"channel"`
	SubchannelCode string `json:"subchannel_code"`
	SubchannelName string `json:"subchannel"`
	Party          string `json:"party"`
	Desc           string `json:"desc"`
}

type ChannelCreateRequest struct {
	Channel
}

type PageInfo struct {
	Items     interface{} `json:"items"`
	Count     int64       `json:"count"`
	PageCount int64       `json:"pageCount"`
	PageNum   int         `json:"pageNum"`
	PageSize  int         `json:"pageSize"`
}
