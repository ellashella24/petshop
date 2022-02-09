package main

import (
	"fmt"
	"petshop/config"
	petCtrl "petshop/delivery/controller/pet"
	"petshop/delivery/middleware"
	"petshop/delivery/route"
	petRepo "petshop/repository/pet"
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

	petRepository := petRepo.NewPetRepository(db)
	petController := petCtrl.NewPetController(petRepository)

	route.RegisterPath(e, petController)

	fmt.Println(db)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", config.Port)))
}
