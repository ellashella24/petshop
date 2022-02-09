package user

type RegisterFormatRequest struct {
	Name     string `json:"name" form:"name" validate:"required"`
	Email    string `json:"email" form:"email" validate:"email,required"`
	Password string `json:"password" form:"password" validate:"required"`
	CityID   uint   `json:"city_id" form:"city_id" validate:"required"`
}

type LoginFormatRequest struct {
	Email    string `json:"email" form:"email" validate:"email,required"`
	Password string `json:"password" form:"password" validate:"required"`
}

type UpdateFormatRequest struct {
	Name     string `json:"name" form:"name" validate:"required"`
	Email    string `json:"email" form:"email" validate:"email,required"`
	Password string `json:"password" form:"password" validate:"required"`
	CityID   uint   `json:"city_id" form:"city_id" validate:"required"`
}
