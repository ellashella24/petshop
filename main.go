package main

import (
	"fmt"
	"petshop/config"
	productController "petshop/delivery/controller/product"
	userCtrl "petshop/delivery/controller/user"
  categoryCtrl "petshop/delivery/controller/category"
	"petshop/delivery/middleware"
	"petshop/delivery/route"
	productRepo "petshop/repository/product"
	userRepo "petshop/repository/user"
  categoryRepo "petshop/repository/category"
	"petshop/util"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	config := config.GetConfig()
	db := util.InitDB(config)
	e := echo.New()

	middleware.LogMiddleware(e)

	e.Validator = &CustomValidator{validator: validator.New()}

	userRepository := userRepo.NewUserRepository(db)
	userController := userCtrl.NewUserController(userRepository)
  
  categoryRepository := categoryRepo.NewCategoryRepository(db)
	categoryController := categoryCtrl.NewCategoryController(categoryRepository)

	productRepo := productRepo.NewProductRepository(db)
	productController := productController.NewProductController(productRepo)

	route.RegisterPath(e, userController, productController, categoryController)

	fmt.Println(db)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", config.Port)))
}
