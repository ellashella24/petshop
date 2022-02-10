package store

type CreateStoreFormatRequest struct {
	Name   string `json:"name" form:"name" validate:"required"`
	CityID uint   `json:"city_id" form:"city_id" validate:"required"`
}

type UpdateStoreFormatRequest struct {
	Name   string `json:"name" form:"name" validate:"required"`
	CityID uint   `json:"city_id" form:"city_id" validate:"required"`
}
