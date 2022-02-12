package store

type StoreFormatResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	CityID uint   `json:"city_id"`
	UserID uint   `json:"user_id"`
}
