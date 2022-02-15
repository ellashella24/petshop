package store

type CreateStoreFormatRequest struct {
	Name   string `json:"name" form:"name" validate:"required"`
	CityID uint   `json:"city_id" form:"city_id" validate:"required"`
}

type UpdateStoreFormatRequest struct {
	Name   string `json:"name" form:"name" `
	CityID uint   `json:"city_id" form:"city_id"`
}

type GetGroomingStatusFormatRequest struct {
	StoreID uint `json:"store_id" form:"store_id"`
	PetID   uint `json:"pet_id" form:"pet_id"`
}

type UpdateGroomingStatusFormatRequest struct {
	StoreID uint `json:"store_id" form:"store_id"`
	PetID   uint `json:"pet_id" form:"pet_id"`
}
