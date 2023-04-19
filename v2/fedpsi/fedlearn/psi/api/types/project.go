package types

import "time"

type Project struct {
	Id            uint64    `json:"id"`
	Uuid          string    `json:"uuid"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	Desc          string    `json:"desc"`
	InitParty     string    `json:"initiator"` // initiator party who create the project
	FollowerParty string    `json:"follower"`  // the party collaberte with initiator party
	Status        string    `json:"status"`
	Creator       string    `json:"creator"`
	UpdateUser    string    `json:"update_user"`
	Created       time.Time `json:"created"`
	Updated       time.Time `json:"Updated"`
}

type ProjectCreateRequest struct {
	Project
}

type ProjectUpdateRequest struct {
	Project
}

type ProjectDeleteRequest struct {
	Id string `json:"id"`
}

type ProjectGetRequest struct {
	Id string `json:"id"`
}
