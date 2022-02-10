package category

type CreateCategoryFormatRequest struct {
	Name string `json:"name" form:"name" validate:"required"`
}

type UpdateCategoryFormatRequest struct {
	Name string `json:"name" form:"name" validate:"required"`
}
