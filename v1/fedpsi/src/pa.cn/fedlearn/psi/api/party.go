package api

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
