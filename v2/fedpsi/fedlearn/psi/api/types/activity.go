package types

type Activity struct {
	Uuid            string   `json:"uuid"`
	Name            string   `json:"name"`
	SendID          string   `json:"send_id"`
	Title           string   `json:"title"`
	InitParty       string   `json:"initiator"`        // initiator party who create the activity
	FollowerParty   string   `json:"follower"`         // the party collaberte with initiator party
	Dataset         []string `json:"dataset"`          // all dataset's names in psi system
	FollowerDataset []string `json:"follower_dataset"` // all dataset's names in psi system
	Status          string   `json:"status"`
}

type ActivityCreateRequest struct {
	Activity
}

type ActivityAttachDataRequest struct {
	Uuid    string   `json:"uuid"`
	Dataset []string `json:"dataset"`
}

type ActivityDeleteRequest struct {
	Uuid string `json:"uuid"`
}

type ActivityStartRequest struct {
	Activity
}

type ActivityConfirmRequest struct {
	Uuid            string   `json:"uuid"`
	FollowerDataset []string `json:"follower_dataset"` // all dataset in psi system
}
