package route

import (
	"petshop/constant"
	"petshop/delivery/controller/user"
	mw "petshop/delivery/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterPath(e *echo.Echo, userCtrl *user.UserController) {
	eAuth := e.Group("")
	eAuth.Use(middleware.JWT([]byte(constant.SecretKey)))
	eAuthAdmin := eAuth.Group("")
	eAuthAdmin.Use(mw.IsAdmin)

	// User Route Path
	e.POST("/user/register", userCtrl.Register())
	e.POST("/user/login", userCtrl.Login())

	eAuth.GET("/user/profile", userCtrl.GetUserProfile())
	eAuth.PUT("/user/profile", userCtrl.UpdateProfile())
	eAuth.DELETE("/user", userCtrl.DeleteAccount())
}
