package product

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"petshop/constants"
	"petshop/delivery/common"
	"petshop/delivery/controller/user"

	"petshop/entity"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

var jwtToken = ""

func TestSetup(t *testing.T) {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	requestBody, _ := json.Marshal(
		user.LoginFormatRequest{
			Email:    "admin1@mail.com",
			Password: "admin1",
		},
	)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
	res := httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/json")
	context := e.NewContext(req, res)
	context.SetPath("/users/login")

	userController := user.NewUserController(mockUserRepository{})
	userController.Login()(context)

	var response common.ResponseSuccess

	json.Unmarshal([]byte(res.Body.Bytes()), &response)

	data := (response.Data).(map[string]interface{})

	jwtToken = data["token"].(string)
}

func TestCreateProduct(t *testing.T) {
	t.Run(
		"1. Error create product - can't get image", func(t *testing.T) {
			e := echo.New()
			e.Validator = &CustomValidator{validator: validator.New()}

			requestBody, _ := json.Marshal(
				CreateProductRequestFormat{
					Name: "Whiskas", Price: 100000, Stock: 100, CategoryID: 2,
				},
			)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
			res := httptest.NewRecorder()
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
			req.Header.Set("Content-Type", "application/json")
			context := e.NewContext(req, res)
			context.SetPath("/product")

			productController := NewProductController(mockFalseProductRepository{})
			err := middleware.JWT([]byte(constants.SecretKey))(productController.CreateProduct())(context)

			if err != nil {
				fmt.Println(err)
				return
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Bad Request", response.Message)
		},
	)

	t.Run(
		"2. Error create product - the image isn't in jpg/png format", func(t *testing.T) {
			e := echo.New()
			e.Validator = &CustomValidator{validator: validator.New()}

			path := "example.txt"

			body := new(bytes.Buffer)

			writer := multipart.NewWriter(body)

			writer.WriteField("name", "Whiskas")
			writer.WriteField("price", "100000")
			writer.WriteField("stock", "100")
			writer.WriteField("storeid", "1")
			writer.WriteField("categoryid", "2")

			file, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}

			w, _ := writer.CreateFormFile("file", path)

			io.Copy(w, file)

			writer.Close()

			req := httptest.NewRequest(http.MethodPost, "/", body)
			res := httptest.NewRecorder()
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
			req.Header.Add("Content-Type", writer.FormDataContentType())
			context := e.NewContext(req, res)
			context.SetPath("/product")

			productController := NewProductController(mockFalseProductRepository{})
			err = middleware.JWT([]byte(constants.SecretKey))(productController.CreateProduct())(context)

			if err != nil {
				fmt.Println(err)
				return
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Bad Request", response.Message)
		},
	)

	t.Run(
		"3. Error create product - can't create product", func(t *testing.T) {
			e := echo.New()
			e.Validator = &CustomValidator{validator: validator.New()}

			path := "Cat03.jpg"

			body := new(bytes.Buffer)

			writer := multipart.NewWriter(body)

			writer.WriteField("name", "Whiskas")
			writer.WriteField("price", "100000")
			writer.WriteField("stock", "100")
			writer.WriteField("storeid", "1")
			writer.WriteField("categoryid", "2")

			file, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}

			w, _ := writer.CreateFormFile("file", path)

			io.Copy(w, file)

			writer.Close()

			req := httptest.NewRequest(http.MethodPost, "/", body)
			res := httptest.NewRecorder()
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
			req.Header.Add("Content-Type", writer.FormDataContentType())
			context := e.NewContext(req, res)
			context.SetPath("/product")

			productController := NewProductController(mockFalseProductRepository{})
			err = middleware.JWT([]byte(constants.SecretKey))(productController.CreateProduct())(context)

			if err != nil {
				fmt.Println(err)
				return
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Bad Request", response.Message)
		},
	)

	t.Run(
		"4. Success create product 1 - non grooming product", func(t *testing.T) {
			e := echo.New()
			e.Validator = &CustomValidator{validator: validator.New()}

			path := "Cat03.jpg"

			body := new(bytes.Buffer)

			writer := multipart.NewWriter(body)

			writer.WriteField("name", "Whiskas")
			writer.WriteField("price", "100000")
			writer.WriteField("stock", "100")
			writer.WriteField("storeid", "1")
			writer.WriteField("categoryid", "2")

			file, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}

			w, _ := writer.CreateFormFile("file", path)

			io.Copy(w, file)

			writer.Close()

			req := httptest.NewRequest(http.MethodPost, "/", body)
			res := httptest.NewRecorder()
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
			req.Header.Add("Content-Type", writer.FormDataContentType())
			context := e.NewContext(req, res)
			context.SetPath("/product")

			productController := NewProductController(mockProductRepository{})
			err = middleware.JWT([]byte(constants.SecretKey))(productController.CreateProduct())(context)

			if err != nil {
				fmt.Println(err)
				return
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Successful Operation", response.Message)
		},
	)

	t.Run(
		"5. Success create product 2 - grooming product", func(t *testing.T) {
			e := echo.New()
			e.Validator = &CustomValidator{validator: validator.New()}

			path := "Cat03.jpg"

			body := new(bytes.Buffer)

			writer := multipart.NewWriter(body)

			writer.WriteField("name", "Grooming")
			writer.WriteField("price", "100000")
			writer.WriteField("stock", "100")
			writer.WriteField("storeid", "1")
			writer.WriteField("categoryid", "1")

			file, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}

			w, _ := writer.CreateFormFile("file", path)

			io.Copy(w, file)

			writer.Close()

			req := httptest.NewRequest(http.MethodPost, "/", body)
			res := httptest.NewRecorder()
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
			req.Header.Add("Content-Type", writer.FormDataContentType())
			context := e.NewContext(req, res)
			context.SetPath("/product")

			productController := NewProductController(mockProductRepository{})
			err = middleware.JWT([]byte(constants.SecretKey))(productController.CreateProduct())(context)

			if err != nil {
				fmt.Println(err)
				return
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Successful Operation", response.Message)
		},
	)
}

func TestGetAllProduct(t *testing.T) {
	t.Run(
		"1. Error get all product", func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			res := httptest.NewRecorder()

			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
			context := e.NewContext(req, res)
			context.SetPath("/products")

			productController := NewProductController(mockFalseProductRepository{})
			err := middleware.JWT([]byte(constants.SecretKey))(productController.GetAllProduct())(context)

			if err != nil {
				fmt.Println(err)
				return
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Not Found", response.Message)
		},
	)

	t.Run(
		"2. Success get all product", func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			res := httptest.NewRecorder()

			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
			context := e.NewContext(req, res)
			context.SetPath("/products")

			productController := NewProductController(mockProductRepository{})
			err := middleware.JWT([]byte(constants.SecretKey))(productController.GetAllProduct())(context)

			if err != nil {
				fmt.Println(err)
				return
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Successful Operation", response.Message)
		},
	)
}

func TestGetProductByID(t *testing.T) {
	t.Run(
		"1. Error get product - product id not found", func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			res := httptest.NewRecorder()

			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
			context := e.NewContext(req, res)
			context.SetPath("/product")
			context.SetParamNames("id")
			context.SetParamValues("100000")

			categoryController := NewProductController(mockFalseProductRepository{})
			err := middleware.JWT([]byte(constants.SecretKey))(categoryController.GetProductByID())(context)

			if err != nil {
				fmt.Println(err)
				return
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Not Found", response.Message)
		},
	)

	t.Run(
		"2. Success get product", func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			res := httptest.NewRecorder()

			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
			context := e.NewContext(req, res)
			context.SetPath("/product")
			context.SetParamNames("id")
			context.SetParamValues("2")

			categoryController := NewProductController(mockProductRepository{})
			err := middleware.JWT([]byte(constants.SecretKey))(categoryController.GetProductByID())(context)

			if err != nil {
				fmt.Println(err)
				return
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Successful Operation", response.Message)
		},
	)
}

func TestGetStockHistory(t *testing.T) {
	t.Run(
		"1. Error get product - product id not found", func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			res := httptest.NewRecorder()

			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
			context := e.NewContext(req, res)
			context.SetPath("/stock/product")
			context.SetParamNames("id")
			context.SetParamValues("100000")

			categoryController := NewProductController(mockFalseProductRepository{})
			err := middleware.JWT([]byte(constants.SecretKey))(categoryController.GetStockHistory())(context)

			if err != nil {
				fmt.Println(err)
				return
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Not Found", response.Message)
		},
	)

	t.Run(
		"2. Success get stock history", func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			res := httptest.NewRecorder()

			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
			context := e.NewContext(req, res)
			context.SetPath("/stock/product")
			context.SetParamNames("id")
			context.SetParamValues("2")

			categoryController := NewProductController(mockProductRepository{})
			err := middleware.JWT([]byte(constants.SecretKey))(categoryController.GetStockHistory())(context)

			if err != nil {
				fmt.Println(err)
				return
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Successful Operation", response.Message)
		},
	)
}

func TestGetProductStoreID(t *testing.T) {
	t.Run(
		"1. Error get all product by store id - store id not found", func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			res := httptest.NewRecorder()

			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
			context := e.NewContext(req, res)
			context.SetPath("/products")
			context.SetParamNames("store")
			context.SetParamValues("100000")

			productController := NewProductController(mockFalseProductRepository{})
			err := middleware.JWT([]byte(constants.SecretKey))(productController.GetProductStoreID())(context)

			if err != nil {
				fmt.Println(err)
				return
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Not Found", response.Message)
		},
	)

	t.Run(
		"2. Success get all product by store id", func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			res := httptest.NewRecorder()

			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
			context := e.NewContext(req, res)
			context.SetPath("/product")
			context.SetParamNames("store")
			context.SetParamValues("1")

			productController := NewProductController(mockProductRepository{})
			err := middleware.JWT([]byte(constants.SecretKey))(productController.GetProductStoreID())(context)

			if err != nil {
				fmt.Println(err)
				return
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Successful Operation", response.Message)
		},
	)
}

func TestUpdateProduct(t *testing.T) {
	t.Run(
		"1. Error create product - the image isn't in jpg/png format", func(t *testing.T) {
			e := echo.New()
			e.Validator = &CustomValidator{validator: validator.New()}

			path := "example.txt"

			body := new(bytes.Buffer)

			writer := multipart.NewWriter(body)

			writer.WriteField("name", "Whiskas new")
			writer.WriteField("price", "100000")
			writer.WriteField("stock", "100")
			writer.WriteField("storeid", "1")
			writer.WriteField("categoryid", "2")

			file, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}

			w, _ := writer.CreateFormFile("file", path)

			io.Copy(w, file)

			writer.Close()

			req := httptest.NewRequest(http.MethodPut, "/", body)
			res := httptest.NewRecorder()
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
			req.Header.Add("Content-Type", writer.FormDataContentType())
			context := e.NewContext(req, res)
			context.SetPath("/product")
			context.SetParamNames("id")
			context.SetParamValues("1")

			productController := NewProductController(mockFalseProductRepository{})
			err = middleware.JWT([]byte(constants.SecretKey))(productController.UpdateProduct())(context)

			if err != nil {
				fmt.Println(err)
				return
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Bad Request", response.Message)
		},
	)

	t.Run(
		"2. Error create product - can't create product", func(t *testing.T) {
			e := echo.New()
			e.Validator = &CustomValidator{validator: validator.New()}

			path := "Cat03.jpg"

			body := new(bytes.Buffer)

			writer := multipart.NewWriter(body)

			writer.WriteField("name", "Whiskas new")
			writer.WriteField("price", "100000")
			writer.WriteField("stock", "100")
			writer.WriteField("storeid", "1")
			writer.WriteField("categoryid", "2")

			file, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}

			w, _ := writer.CreateFormFile("file", path)

			io.Copy(w, file)

			writer.Close()

			req := httptest.NewRequest(http.MethodPut, "/", body)
			res := httptest.NewRecorder()
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
			req.Header.Add("Content-Type", writer.FormDataContentType())
			context := e.NewContext(req, res)
			context.SetPath("/product")
			context.SetParamNames("id")
			context.SetParamValues("1")

			productController := NewProductController(mockFalseProductRepository{})
			err = middleware.JWT([]byte(constants.SecretKey))(productController.UpdateProduct())(context)

			if err != nil {
				fmt.Println(err)
				return
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Bad Request", response.Message)
		},
	)

	t.Run(
		"4. Success create product 1 - non grooming product", func(t *testing.T) {
			e := echo.New()
			e.Validator = &CustomValidator{validator: validator.New()}

			path := "Cat03.jpg"

			body := new(bytes.Buffer)

			writer := multipart.NewWriter(body)

			writer.WriteField("name", "Whiskas new")
			writer.WriteField("price", "100000")
			writer.WriteField("stock", "100")

			file, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}

			w, _ := writer.CreateFormFile("file", path)

			io.Copy(w, file)

			writer.Close()

			req := httptest.NewRequest(http.MethodPut, "/", body)
			res := httptest.NewRecorder()
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
			req.Header.Add("Content-Type", writer.FormDataContentType())
			context := e.NewContext(req, res)
			context.SetPath("/product")
			context.SetParamNames("id")
			context.SetParamValues("2")

			productController := NewProductController(mockProductRepository{})
			err = middleware.JWT([]byte(constants.SecretKey))(productController.UpdateProduct())(context)

			if err != nil {
				fmt.Println(err)
				return
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Successful Operation", response.Message)
		},
	)

	t.Run(
		"5. Success create product 2 - grooming product", func(t *testing.T) {
			e := echo.New()
			e.Validator = &CustomValidator{validator: validator.New()}

			path := "Cat03.jpg"

			body := new(bytes.Buffer)

			writer := multipart.NewWriter(body)

			writer.WriteField("name", "Grooming new")
			writer.WriteField("price", "100000")
			writer.WriteField("stock", "1")

			file, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}

			w, _ := writer.CreateFormFile("file", path)

			io.Copy(w, file)

			writer.Close()

			req := httptest.NewRequest(http.MethodPut, "/", body)
			res := httptest.NewRecorder()
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
			req.Header.Add("Content-Type", writer.FormDataContentType())
			context := e.NewContext(req, res)
			context.SetPath("/product")
			context.SetParamNames("id")
			context.SetParamValues("1")

			productController := NewProductController(mockProductRepository{})
			err = middleware.JWT([]byte(constants.SecretKey))(productController.UpdateProduct())(context)

			if err != nil {
				fmt.Println(err)
				return
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Successful Operation", response.Message)
		},
	)
}
func TestDeleteProduct(t *testing.T) {
	t.Run(
		"1. Error get product - product id not found", func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			res := httptest.NewRecorder()

			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
			context := e.NewContext(req, res)
			context.SetPath("/product")
			context.SetParamNames("id")
			context.SetParamValues("100000")

			categoryController := NewProductController(mockFalseProductRepository{})
			err := middleware.JWT([]byte(constants.SecretKey))(categoryController.DeleteProduct())(context)

			if err != nil {
				fmt.Println(err)
				return
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Bad Request", response.Message)
		},
	)

	t.Run(
		"2. Success get product", func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			res := httptest.NewRecorder()

			req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
			context := e.NewContext(req, res)
			context.SetPath("/product")
			context.SetParamNames("id")
			context.SetParamValues("2")

			categoryController := NewProductController(mockProductRepository{})
			err := middleware.JWT([]byte(constants.SecretKey))(categoryController.DeleteProduct())(context)

			if err != nil {
				fmt.Println(err)
				return
			}

			var response common.ResponseSuccess

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, "Successful Operation", response.Message)
		},
	)
}

type mockUserRepository struct{}

func (mu mockUserRepository) FindCityByID(productID int) (entity.City, error) {
	return entity.City{
		ID: 1, Name: "city",
	}, nil
}

func (mu mockUserRepository) GetAllUser() ([]entity.User, error) {
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte("admin1"), bcrypt.MinCost)

	return []entity.User{
		{ID: 1, Name: "admin1", Email: "admin1@mail.com", Password: string(hashPassword), Role: "admin"},
	}, nil
}

func (mu mockUserRepository) GetUserByID(userID int) (entity.User, error) {
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte("admin1"), bcrypt.MinCost)

	return entity.User{
		ID: 1, Name: "admin1", Email: "admin1@mail.com", Password: string(hashPassword), Role: "admin",
	}, nil
}

func (mu mockUserRepository) GetUserByEmail(email string) (entity.User, error) {
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte("admin1"), bcrypt.MinCost)

	return entity.User{
		ID: 1, Name: "admin1", Email: "admin1@mail.com", Password: string(hashPassword), Role: "admin",
	}, nil
}

func (mu mockUserRepository) CreateUser(newUser entity.User) (entity.User, error) {
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte("admin1"), bcrypt.MinCost)

	return entity.User{
		ID: 1, Name: "admin1", Email: "admin1@mail.com", Password: string(hashPassword), Role: "admin",
	}, nil
}

func (mu mockUserRepository) UpdateUser(userID int, updatedUser entity.User) (entity.User, error) {
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte("admin1"), bcrypt.MinCost)

	return entity.User{
		ID: 1, Name: "admin1 new", Email: "admin1@mail.com", Password: string(hashPassword), Role: "admin",
	}, nil
}

func (mu mockUserRepository) DeleteUser(userID int) (entity.User, error) {
	// hashPassword, _ := bcrypt.GenerateFromPassword([]byte("admin1"), bcrypt.MinCost)

	return entity.User{
		ID: 0, Name: "", Email: "", Password: "", Role: "",
	}, nil
}

type mockProductRepository struct{}

func (mc mockProductRepository) GetAllProduct() ([]entity.Product, error) {
	return []entity.Product{
		{ID: 2, Name: "Whiskas", Price: 100000, Stock: 100, StoreID: 1, CategoryID: 2},
	}, nil
}

func (mc mockProductRepository) GetProductByID(productID int) (entity.Product, error) {
	return entity.Product{
		ID: 2, Name: "Whiskas", Price: 100000, Stock: 100, StoreID: 1, CategoryID: 2,
	}, nil
}

func (mc mockProductRepository) CreateProduct(userID int, newProduct entity.Product) (entity.Product, error) {
	if newProduct.CategoryID == 1 {
		return entity.Product{
			ID: 1, Name: "Grooming", Price: 100000, StoreID: 1, CategoryID: 1,
		}, nil
	}
	return entity.Product{
		ID: 2, Name: "Whiskas", Price: 100000, Stock: 100, StoreID: 1, CategoryID: 2,
	}, nil
}

func (mc mockProductRepository) UpdateProduct(productID int, updatedProduct entity.Product) (entity.Product, error) {
	if updatedProduct.CategoryID == 1 {
		return entity.Product{
			ID: 1, Name: "Grooming new", Price: 100000, StoreID: 1, CategoryID: 1,
		}, nil
	}
	return entity.Product{
		ID: 2, Name: "Whiskas new", Price: 100000, Stock: 100, StoreID: 1, CategoryID: 2,
	}, nil
}

func (mu mockProductRepository) DeleteProduct(productID int) (entity.Product, error) {
	return entity.Product{
		ID: 0, Name: "",
	}, nil
}

func (mc mockProductRepository) GetProductByStoreID(storeID int) ([]entity.Product, error) {
	return []entity.Product{
		{ID: 2, Name: "Whiskas", Price: 100000, Stock: 100, StoreID: 1, CategoryID: 2},
	}, nil
}

func (mc mockProductRepository) GetStockHistory(productID int) ([]entity.StockHistory, error) {
	return []entity.StockHistory{
		{ID: 1, ProductID: 2, Stock: 100},
	}, nil
}

type mockFalseProductRepository struct{}

func (mfc mockFalseProductRepository) GetAllProduct() ([]entity.Product, error) {
	return []entity.Product{
		{ID: 0, Name: "", Price: 0, Stock: 0, StoreID: 0, CategoryID: 0},
	}, errors.New("can't get product data")
}

func (mfc mockFalseProductRepository) GetProductByID(productID int) (entity.Product, error) {
	return entity.Product{
		ID: 0, Name: "", Price: 0, Stock: 0, StoreID: 0, CategoryID: 0,
	}, errors.New("can't get product data")
}

func (mfc mockFalseProductRepository) CreateProduct(userID int, newProduct entity.Product) (entity.Product, error) {
	return entity.Product{
		ID: 0, Name: "", Price: 0, Stock: 0, StoreID: 0, CategoryID: 0,
	}, errors.New("can't create product data")
}

func (mfc mockFalseProductRepository) UpdateProduct(productID int, updatedProduct entity.Product) (
	entity.Product, error,
) {
	return entity.Product{
		ID: 0, Name: "", Price: 0, Stock: 0, StoreID: 0, CategoryID: 0,
	}, errors.New("can't update product data")
}

func (mfc mockFalseProductRepository) DeleteProduct(productID int) (entity.Product, error) {
	return entity.Product{
		ID: 0, Name: "", Price: 0, Stock: 0, StoreID: 0, CategoryID: 0,
	}, errors.New("can't get product data")
}

func (mfc mockFalseProductRepository) GetProductByStoreID(storeID int) ([]entity.Product, error) {
	return []entity.Product{
		{ID: 0, Name: "", Price: 0, Stock: 0, StoreID: 0, CategoryID: 0},
	}, errors.New("can't get product data")
}

func (mfc mockFalseProductRepository) GetStockHistory(productID int) ([]entity.StockHistory, error) {
	return []entity.StockHistory{
		{ID: 0, ProductID: 0, Stock: 0},
	}, errors.New("can't get stock product history data")
}
