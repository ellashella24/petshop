package route

import (
	"petshop/constant"
	"petshop/delivery/controller/store"
	mw "petshop/delivery/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterPath(e *echo.Echo, storeCtrl *store.StoreController) {
	eAuth := e.Group("")
	eAuth.Use(middleware.JWT([]byte(constant.SecretKey)))
	eAuthAdmin := eAuth.Group("")
	eAuthAdmin.Use(mw.IsAdmin)

	// Store Route Path
	eAuth.POST("/user/store", storeCtrl.CreateStore())
	eAuth.GET("/user/stores", storeCtrl.GetAllStoreByUser())
	eAuth.GET("/user/store/profile/:id", storeCtrl.GetStoreProfile())
	eAuth.PUT("/user/store/profile/:id", storeCtrl.UpdateStoreProfile())
	eAuth.DELETE("/user/store/:id", storeCtrl.DeleteStore())
}
