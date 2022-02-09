package route

import (
	"petshop/constant"
	"petshop/delivery/controller/city"
	mw "petshop/delivery/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterPath(e *echo.Echo, cityCtrl *city.CityController) {
	eAuth := e.Group("")
	eAuth.Use(middleware.JWT([]byte(constant.SecretKey)))
	eAuthAdmin := eAuth.Group("")
	eAuthAdmin.Use(mw.IsAdmin)

	// City Route Path
	eAuthAdmin.POST("/city", cityCtrl.CreateCity())
	eAuthAdmin.GET("/cities", cityCtrl.GetAllCity())
	eAuthAdmin.GET("/city/profile/:id", cityCtrl.GetCityProfile())
	eAuthAdmin.PUT("/city/profile/:id", cityCtrl.UpdateCityProfile())
	eAuthAdmin.DELETE("/city/:id", cityCtrl.DeleteCity())
}
