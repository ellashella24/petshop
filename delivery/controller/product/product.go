package product

import (
	"net/http"
	"petshop/delivery/common"
	"petshop/delivery/middleware"
	"petshop/entity"
	"petshop/repository/product"
	"petshop/service"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type ProductController struct {
	productRepo product.Product
}

func NewProductController(productRepo product.Product) *ProductController {
	return &ProductController{productRepo}
}

func (uc *ProductController) GetAllProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		res, err := uc.productRepo.GetAllProduct()

		if err != nil || len(res) == 0 {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		}

		formatProduct := ProductFormatResponse{}
		formatProducts := []ProductFormatResponse{}

		for i := 0; i < len(res); i++ {
			formatProduct.ID = res[i].ID
			formatProduct.Name = res[i].Name
			formatProduct.Price = res[i].Price
			formatProduct.ImageUrl = res[i].ImageURL
			formatProduct.Stock = res[i].Stock

			formatProducts = append(formatProducts, formatProduct)
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(formatProducts))
	}
}

func (uc *ProductController) GetProductByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		productID, _ := strconv.Atoi(c.Param("id"))

		res, err := uc.productRepo.GetProductByID(productID)

		if err != nil || res.ID == 0 {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		}

		formatProduct := ProductFormatResponse{
			ID:       res.ID,
			Name:     res.Name,
			Price:    res.Price,
			Stock:    res.Stock,
			ImageUrl: res.ImageURL,
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(formatProduct))
	}
}

func (uc *ProductController) GetStockHistory() echo.HandlerFunc {
	return func(c echo.Context) error {
		productID, _ := strconv.Atoi(c.Param("id"))

		res, err := uc.productRepo.GetStockHistory(productID)

		if err != nil || len(res) == 0 {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		}

		stockResponse := StockFormatResponse{}
		stockResponses := []StockFormatResponse{}

		for i := 0; i < len(res); i++ {
			stockResponse.ID = res[i].ID
			stockResponse.ProductID = res[i].ProductID
			stockResponse.Stock = res[i].Stock

			stockResponses = append(stockResponses, stockResponse)
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(stockResponses))
	}
}

func (uc *ProductController) GetProductStoreID() echo.HandlerFunc {
	return func(c echo.Context) error {
		storeID, _ := strconv.Atoi(c.QueryParam("store"))

		res, err := uc.productRepo.GetProductByStoreID(storeID)

		if err != nil || len(res) == 0 {
			return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
		}

		formatProduct := ProductFormatResponse{}
		formatProducts := []ProductFormatResponse{}

		for i := 0; i < len(res); i++ {
			formatProduct.ID = res[i].ID
			formatProduct.Name = res[i].Name
			formatProduct.Price = res[i].Price
			formatProduct.Stock = res[i].Stock
			formatProducts = append(formatProducts, formatProduct)
		}

		// return c.JSON(http.StatusOK, common.SuccessResponse(formatProduct))
		return c.JSON(http.StatusOK, common.SuccessResponse(formatProducts))
	}
}

func (uc *ProductController) CreateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := middleware.ExtractTokenUserID(c)
		CreateProductReq := CreateProductRequestFormat{}
		c.Bind(&CreateProductReq)

		file, err := c.FormFile("file")

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		imageURL, err := service.Upload(c, strconv.Itoa(userId), file)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		newProduct := entity.Product{}
		newProduct.Name = CreateProductReq.Name
		newProduct.Price = CreateProductReq.Price
		newProduct.StoreID = CreateProductReq.StoreID
		newProduct.ImageURL = imageURL.ImageURL
		newProduct.CategoryID = CreateProductReq.CategoryID
		if CreateProductReq.CategoryID != 1 {
			newProduct.Stock = CreateProductReq.Stock
		}

		res, err := uc.productRepo.CreateProduct(userId, newProduct)

		if err != nil {
			destination := strings.ReplaceAll(res.ImageURL, "http://naufalhibatullah.com/images/", "")
			service.Delete(destination)
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}
		formatProduct := ProductFormatResponse{}
		if res.CategoryID != 1 {
			formatProduct = ProductFormatResponse{
				ID:       res.ID,
				Name:     res.Name,
				Price:    res.Price,
				Stock:    res.Stock,
				ImageUrl: res.ImageURL,
			}
		} else {
			formatProduct = ProductFormatResponse{
				ID:       res.ID,
				Name:     res.Name,
				Price:    res.Price,
				ImageUrl: res.ImageURL,
			}
		}
		return c.JSON(http.StatusOK, common.SuccessResponse(formatProduct))
	}
}

func (uc *ProductController) UpdateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		productID, _ := strconv.Atoi(c.Param("id"))
		userId := middleware.ExtractTokenUserID(c)

		productUpdateReq := UpdateProductRequestFormat{}

		c.Bind(&productUpdateReq)

		file, _ := c.FormFile("file")
		updateFile := ""
		if file != nil {
			file, err := service.Upload(c, strconv.Itoa(userId), file)
			if err != nil {
				return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
			}
			updateFile = file.ImageURL
		}

		updatedProduct := entity.Product{}
		updatedProduct.Name = productUpdateReq.Name
		updatedProduct.Price = productUpdateReq.Price
		updatedProduct.Stock = productUpdateReq.Stock
		updatedProduct.ImageURL = updateFile

		res, err := uc.productRepo.UpdateProduct(productID, updatedProduct)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		formatProduct := ProductFormatResponse{}

		if res.CategoryID != 1 {
			formatProduct = ProductFormatResponse{
				ID:       uint(productID),
				Name:     res.Name,
				Price:    res.Price,
				Stock:    res.Stock,
				ImageUrl: res.ImageURL,
			}
		} else {
			formatProduct = ProductFormatResponse{
				ID:       uint(productID),
				Name:     res.Name,
				Price:    res.Price,
				ImageUrl: res.ImageURL,
			}
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(formatProduct))
	}
}

func (uc *ProductController) DeleteProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		productID, _ := strconv.Atoi(c.Param("id"))

		_, err := uc.productRepo.DeleteProduct(productID)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
		}

		return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
	}
}
