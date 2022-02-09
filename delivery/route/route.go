package route

import (
	"petshop/constant"
	"petshop/delivery/controller/pet"
	mw "petshop/delivery/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterPath(e *echo.Echo, petCtrl *pet.PetController) {
	eAuth := e.Group("")
	eAuth.Use(middleware.JWT([]byte(constant.SecretKey)))
	eAuthAdmin := eAuth.Group("")
	eAuthAdmin.Use(mw.IsAdmin)

	// Pet Route Path
	eAuth.POST("/user/pet", petCtrl.CreatePet())
	eAuth.GET("/user/pets", petCtrl.GetAllPetByUser())
	eAuth.GET("/user/pet/profile/:id", petCtrl.GetPetProfile())
	eAuth.PUT("/user/pet/profile/:id", petCtrl.UpdatePetProfile())
	eAuth.DELETE("/user/pet/:id", petCtrl.DeletePet())
}
