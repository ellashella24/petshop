package route

import (
	"petshop/constant"
	"petshop/delivery/controller/category"
	mw "petshop/delivery/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterPath(e *echo.Echo, categoryCtrl *category.CategoryController) {
	eAuth := e.Group("")
	eAuth.Use(middleware.JWT([]byte(constant.SecretKey)))
	eAuthAdmin := eAuth.Group("")
	eAuthAdmin.Use(mw.IsAdmin)

	// Category Route Path
	eAuthAdmin.POST("/category", categoryCtrl.CreateCategory())
	eAuthAdmin.GET("/categories", categoryCtrl.GetAllCategory())
	eAuthAdmin.GET("/category/profile/:id", categoryCtrl.GetCategoryProfile())
	eAuthAdmin.PUT("/category/profile/:id", categoryCtrl.UpdateCategoryProfile())
	eAuthAdmin.DELETE("/category/:id", categoryCtrl.DeleteCategory())
}
