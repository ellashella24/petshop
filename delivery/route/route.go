package route

import (
	"petshop/constant"
	"petshop/delivery/controller/category"
	"petshop/delivery/controller/city"
	"petshop/delivery/controller/pet"
	"petshop/delivery/controller/product"
	"petshop/delivery/controller/transaction"
	"petshop/delivery/controller/user"
	mw "petshop/delivery/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterPath(e *echo.Echo, userCtrl *user.UserController, pc *product.ProductController, categoryCtrl *category.CategoryController, tc *transaction.TransactionController, cityCtrl *city.CityController, petCtrl *pet.PetController) {
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

	// City Route Path
	eAuthAdmin.POST("/city", cityCtrl.CreateCity())
	eAuthAdmin.GET("/cities", cityCtrl.GetAllCity())
	eAuthAdmin.GET("/city/profile/:id", cityCtrl.GetCityProfile())
	eAuthAdmin.PUT("/city/profile/:id", cityCtrl.UpdateCityProfile())
	eAuthAdmin.DELETE("/city/:id", cityCtrl.DeleteCity())
	// User Route Path
	e.POST("/user/register", userCtrl.Register())
	e.POST("/user/login", userCtrl.Login())

	eAuth.GET("/user/profile", userCtrl.GetUserProfile())
	eAuth.PUT("/user/profile", userCtrl.UpdateProfile())
	eAuth.DELETE("/user", userCtrl.DeleteAccount())

	// Category Route Path
	eAuthAdmin.POST("/category", categoryCtrl.CreateCategory())
	eAuthAdmin.GET("/categories", categoryCtrl.GetAllCategory())
	eAuthAdmin.GET("/category/profile/:id", categoryCtrl.GetCategoryProfile())
	eAuthAdmin.PUT("/category/profile/:id", categoryCtrl.UpdateCategoryProfile())
	eAuthAdmin.DELETE("/category/:id", categoryCtrl.DeleteCategory())

	// Product Route Path
	e.POST("/product", pc.CreateProduct(), middleware.JWT([]byte("secret123")))
	e.GET("/product", pc.GetAllProduct(), middleware.JWT([]byte("secret123")))
	e.GET("/product/:id", pc.GetProductByID(), middleware.JWT([]byte("secret123")))
	e.GET("/product?store=", pc.GetProductByID(), middleware.JWT([]byte("secret123")))
	e.PUT("/product/:id", pc.UpdateProduct(), middleware.JWT([]byte("secret123")))
	e.DELETE("/product/:id", pc.DeleteProduct(), middleware.JWT([]byte("secret123")))
	e.GET("/stock/product/:id", pc.GetStockHistory(), middleware.JWT([]byte("secret123")))

	e.POST("/transaction", tc.Create(), middleware.JWT([]byte("secret123")))
	e.POST("/callback", tc.Callback())

}
