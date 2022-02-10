package pet

type CreatePetFormatRequest struct {
	Name string `json:"name" form:"name" validate:"required"`
}

type UpdatePetFormatRequest struct {
	Name string `json:"name" form:"name" validate:"required"`
}
