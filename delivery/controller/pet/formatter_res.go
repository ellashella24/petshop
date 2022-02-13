package pet

type PetFormatResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	UserID uint   `json:"user_id"`
}

type GroomingStatusResponse struct {
	ID     uint   `json:"id"`
	PetID  uint   `json:"pet_id"`
	Status string `json:"status"`
}
