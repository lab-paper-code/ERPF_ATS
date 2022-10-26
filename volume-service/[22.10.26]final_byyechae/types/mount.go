package types


type Mount struct {
	MountPath   string `json:"path"`
	IP          string `json:"ip"`
	ID          string `json:"id"`
	//Description string `json:"description,omitempty"`
}
