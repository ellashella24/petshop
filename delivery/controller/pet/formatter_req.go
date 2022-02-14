package pet

type CreatePetFormatRequest struct {
	Name string `json:"name" form:"name" validate:"required"`
}

type UpdatePetFormatRequest struct {
	Name string `json:"name" form:"name" validate:"required"`
}

type GetGroomingStatusFormatRequest struct {
	PetID  uint `json:"pet_id" form:"pet_id"`
	UserID uint `json:"user_id" form:"user_id"`
}

type UpdateGroomingStatusFormatRequest struct {
	PetID  uint `json:"pet_id" form:"pet_id"`
	UserID uint `json:"user_id" form:"user_id"`
}
