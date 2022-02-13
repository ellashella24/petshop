package category

import (
	"net/http"
	"petshop/delivery/common"
	"petshop/entity"
	"petshop/repository/category"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CategoryController struct {
	categoryRepo category.Category
}

func NewCategoryController(categoryRepo category.Category) *CategoryController {
	return &CategoryController{categoryRepo}
}

func (cc *CategoryController) CreateCategory() echo.HandlerFunc {
	return func(c echo.Context) error {
		createCategoryReq := CreateCategoryFormatRequest{}

		c.Bind(&createCategoryReq)

		err := c.Validate(&createCategoryReq)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Can't get the input"))
		}

		newCategory := entity.Category{}
		newCategory.Name = createCategoryReq.Name

		res, err := cc.categoryRepo.CreateCategory(newCategory)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Can't create category"))
		}

		response := CategoryFormatResponse{}
		response.ID = res.ID
		response.Name = res.Name

		return c.JSON(http.StatusOK, common.SuccessResponse(response))
	}
}

func (cc *CategoryController) GetAllCategory() echo.HandlerFunc {
	return func(c echo.Context) error {
		res, err := cc.categoryRepo.GetAllCategory()

		if err != nil || len(res) == 0 {
			return c.JSON(http.StatusNotFound, common.ErrorResponse(404, "Category not found"))
		}

		responses := []CategoryFormatResponse{}
		response := CategoryFormatResponse{}

		for i := 0; i < len(res); i++ {
			response.ID = res[i].ID
			response.Name = res[i].Name
			responses = append(responses, response)
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(responses))
	}
}

func (cc *CategoryController) GetCategoryProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		categoryID, _ := strconv.Atoi(c.Param("id"))

		res, err := cc.categoryRepo.GetCategoryByID(categoryID)

		if err != nil || res.ID == 0 {
			return c.JSON(http.StatusNotFound, common.ErrorResponse(404, "Category not found"))
		}

		response := CategoryFormatResponse{}

		response.ID = res.ID
		response.Name = res.Name

		return c.JSON(http.StatusOK, common.SuccessResponse(response))
	}
}

func (cc *CategoryController) UpdateCategoryProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		categoryID, _ := strconv.Atoi(c.Param("id"))

		updateCategoryReq := UpdateCategoryFormatRequest{}

		c.Bind(&updateCategoryReq)

		err := c.Validate(&updateCategoryReq)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Can't get the input"))
		}

		updatedCategory := entity.Category{}
		updatedCategory.ID = uint(categoryID)
		updatedCategory.Name = updateCategoryReq.Name

		res, err := cc.categoryRepo.UpdateCategory(categoryID, updatedCategory)

		if err != nil || res.Name == "" {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Can't update category"))
		}

		response := CategoryFormatResponse{}
		response.ID = uint(categoryID)
		response.Name = res.Name

		return c.JSON(http.StatusOK, common.SuccessResponse(response))
	}
}

func (cc *CategoryController) DeleteCategory() echo.HandlerFunc {
	return func(c echo.Context) error {
		categoryID, _ := strconv.Atoi(c.Param("id"))

		_, err := cc.categoryRepo.DeleteCategory(categoryID)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Can't delete category"))
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(nil))
	}
}
