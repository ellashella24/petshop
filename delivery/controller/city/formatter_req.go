package city

type CreateCityFormatRequest struct {
	Name string `json:"name" form:"name" validate:"required"`
}

type UpdateCityFormatRequest struct {
	Name string `json:"name" form:"name" validate:"required"`
}