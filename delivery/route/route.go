package route

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"petshop/delivery/controller/product"
)

func RegisterPath(e *echo.Echo, pc *product.ProductController) {

	e.POST("/product", pc.CreateProduct(), middleware.JWT([]byte("secret123")))
	e.GET("/product", pc.GetAllProduct(), middleware.JWT([]byte("secret123")))
	e.GET("/product/:id", pc.GetProductByID(), middleware.JWT([]byte("secret123")))
	e.GET("/product?store=", pc.GetProductByID(), middleware.JWT([]byte("secret123")))
	e.PUT("/product/:id", pc.UpdateProduct(), middleware.JWT([]byte("secret123")))
	e.DELETE("/product/:id", pc.DeleteProduct(), middleware.JWT([]byte("secret123")))
	e.GET("/stock/product/:id", pc.GetStockHistory(), middleware.JWT([]byte("secret123")))

}
