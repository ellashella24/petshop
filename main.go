package main

import (
	"fmt"
	"petshop/config"
	cityCtrl "petshop/delivery/controller/city"
	"petshop/delivery/middleware"
	"petshop/delivery/route"
	cityRepo "petshop/repository/city"
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

	cityRepository := cityRepo.NewCityRepository(db)
	cityController := cityCtrl.NewCityController(cityRepository)

	route.RegisterPath(e, cityController)

	fmt.Println(db)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", config.Port)))
}
