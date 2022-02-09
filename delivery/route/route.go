package route

import (
"petshop/constant"
	"petshop/delivery/controller/user"
	mw "petshop/delivery/middleware"
  	"petshop/delivery/controller/product"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterPath(e *echo.Echo, userCtrl *user.UserController, pc *product.ProductController) {
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
  
  	e.POST("/product", pc.CreateProduct(), middleware.JWT([]byte("secret123")))
	e.GET("/product", pc.GetAllProduct(), middleware.JWT([]byte("secret123")))
	e.GET("/product/:id", pc.GetProductByID(), middleware.JWT([]byte("secret123")))
	e.GET("/product?store=", pc.GetProductByID(), middleware.JWT([]byte("secret123")))
	e.PUT("/product/:id", pc.UpdateProduct(), middleware.JWT([]byte("secret123")))
	e.DELETE("/product/:id", pc.DeleteProduct(), middleware.JWT([]byte("secret123")))
	e.GET("/stock/product/:id", pc.GetStockHistory(), middleware.JWT([]byte("secret123")))

}
