package main

import (
	"fmt"
	"petshop/config"
	userCtrl "petshop/delivery/controller/user"
	"petshop/delivery/middleware"
	"petshop/delivery/route"
	userRepo "petshop/repository/user"
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

	route.RegisterPath(e, userController)

	fmt.Println(db)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", config.Port)))
}
