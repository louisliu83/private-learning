package api

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
