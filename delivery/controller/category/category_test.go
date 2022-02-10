package category

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"petshop/delivery/common"
	"petshop/delivery/controller/user"
	mw "petshop/delivery/middleware"
	"petshop/entity"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
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

func TestCreateCategory(t *testing.T) {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	requestBody, _ := json.Marshal(user.LoginFormatRequest{
		Email:    "admin1@mail.com",
		Password: "admin1",
	})

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

	t.Run("1. Error Create Category - Error Validate", func(t *testing.T) {
		ev := echo.New()

		requestBody, _ := json.Marshal(CreateCategoryFormatRequest{
			Name: "category1",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		req.Header.Set("Content-Type", "application/json")
		context := ev.NewContext(req, res)
		context.SetPath("/category")

		categoryController := NewCategoryController(mockCategoryRepository{})
		err := mw.IsAdmin(categoryController.CreateCategory())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Can't get the input", response.Message)
	})

	t.Run("2. Error Create Category - Can't Create Input", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}

		requestBody, _ := json.Marshal(CreateCategoryFormatRequest{
			Name: "category1",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/category")

		categoryController := NewCategoryController(mockFalseCategoryRepository{})

		err := (mw.IsAdmin)(categoryController.CreateCategory())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.DefaultResponse

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Can't create category", response.Message)
	})

	t.Run("3. Success Create Category", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}

		requestBody, _ := json.Marshal(CreateCategoryFormatRequest{
			Name: "category2",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/category")

		categoryController := NewCategoryController(mockCategoryRepository{})

		err := (mw.IsAdmin)(categoryController.CreateCategory())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.DefaultResponse

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})
}

func TestGetAllCategory(t *testing.T) {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	requestBody, _ := json.Marshal(user.LoginFormatRequest{
		Email:    "admin1@mail.com",
		Password: "admin1",
	})

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

	t.Run("1. Error Get All Category Test", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/cities")

		categoryController := NewCategoryController(mockFalseCategoryRepository{})
		err := (mw.IsAdmin)(categoryController.GetAllCategory())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Category not found", response.Message)
	})

	t.Run("2. Success Get All Category Test", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/cities")

		categoryController := NewCategoryController(mockCategoryRepository{})
		err := (mw.IsAdmin)(categoryController.GetAllCategory())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})
}

func TestGetCategoryProfile(t *testing.T) {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	requestBody, _ := json.Marshal(user.LoginFormatRequest{
		Email:    "admin1@mail.com",
		Password: "admin1",
	})

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

	t.Run("1. Error Get Category Profile Test", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/category/profile")
		context.SetParamNames("id")
		context.SetParamValues("1")

		categoryController := NewCategoryController(mockFalseCategoryRepository{})
		err := (mw.IsAdmin)(categoryController.GetCategoryProfile())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Category not found", response.Message)
	})

	t.Run("2. Success Get Category Profile Test", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/category")
		context.SetParamNames("id")
		context.SetParamValues("1")

		categoryController := NewCategoryController(mockCategoryRepository{})
		err := (mw.IsAdmin)(categoryController.GetCategoryProfile())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})
}

func TestUpdateCategoryProfile(t *testing.T) {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	requestBody, _ := json.Marshal(user.LoginFormatRequest{
		Email:    "admin1@mail.com",
		Password: "admin1",
	})

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

	t.Run("1. Error Update Category - Error Validate", func(t *testing.T) {
		ev := echo.New()

		requestBody, _ := json.Marshal(UpdateCategoryFormatRequest{
			Name: "category1 new",
		})

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := ev.NewContext(req, res)
		context.SetPath("/category/profile")
		context.SetParamNames("id")
		context.SetParamValues("1")

		categoryController := NewCategoryController(mockFalseCategoryRepository{})

		err := (mw.IsAdmin)(categoryController.UpdateCategoryProfile())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.DefaultResponse

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Can't get the input", response.Message)
	})

	t.Run("2. Error Update Category - Can't Create Input", func(t *testing.T) {
		requestBody, _ := json.Marshal(UpdateCategoryFormatRequest{
			Name: "category1 new",
		})

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/category/profile")
		context.SetParamNames("id")
		context.SetParamValues("1")

		categoryController := NewCategoryController(mockFalseCategoryRepository{})

		err := (mw.IsAdmin)(categoryController.UpdateCategoryProfile())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.DefaultResponse

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Can't update category", response.Message)
	})

	t.Run("3. Success Update Category", func(t *testing.T) {
		requestBody, _ := json.Marshal(CreateCategoryFormatRequest{
			Name: "category1 new",
		})

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/category/profile")
		context.SetParamNames("id")
		context.SetParamValues("1")

		categoryController := NewCategoryController(mockCategoryRepository{})

		err := (mw.IsAdmin)(categoryController.UpdateCategoryProfile())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.DefaultResponse

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})
}

func TestDeleteCategory(t *testing.T) {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	requestBody, _ := json.Marshal(user.LoginFormatRequest{
		Email:    "admin1@mail.com",
		Password: "admin1",
	})

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

	t.Run("1. Error Delete Category Test", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/category")
		context.SetParamNames("id")
		context.SetParamValues("1")

		categoryController := NewCategoryController(mockFalseCategoryRepository{})
		err := (mw.IsAdmin)(categoryController.DeleteCategory())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Can't delete category", response.Message)
	})

	t.Run("2. Success Delete Category Test", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/category")
		context.SetParamNames("id")
		context.SetParamValues("1")

		categoryController := NewCategoryController(mockCategoryRepository{})
		err := (mw.IsAdmin)(categoryController.DeleteCategory())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})
}

type mockUserRepository struct{}

func (mu mockUserRepository) FindCityByID(categoryID int) (entity.City, error) {
	return entity.City{
		ID: 1, Name: "city"}, nil
}

func (mu mockUserRepository) GetAllUser() ([]entity.User, error) {
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte("admin1"), bcrypt.MinCost)

	return []entity.User{
		{ID: 1, Name: "admin1", Email: "admin1@mail.com", Password: string(hashPassword), Role: "admin"}}, nil
}

func (mu mockUserRepository) GetUserByID(userID int) (entity.User, error) {
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte("admin1"), bcrypt.MinCost)

	return entity.User{
		ID: 1, Name: "admin1", Email: "admin1@mail.com", Password: string(hashPassword), Role: "admin"}, nil
}

func (mu mockUserRepository) GetUserByEmail(email string) (entity.User, error) {
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte("admin1"), bcrypt.MinCost)

	return entity.User{
		ID: 1, Name: "admin1", Email: "admin1@mail.com", Password: string(hashPassword), Role: "admin"}, nil
}

func (mu mockUserRepository) CreateUser(newUser entity.User) (entity.User, error) {
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte("admin1"), bcrypt.MinCost)

	return entity.User{
		ID: 1, Name: "admin1", Email: "admin1@mail.com", Password: string(hashPassword), Role: "admin"}, nil
}

func (mu mockUserRepository) UpdateUser(userID int, updatedUser entity.User) (entity.User, error) {
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte("admin1"), bcrypt.MinCost)

	return entity.User{
		ID: 1, Name: "admin1 new", Email: "admin1@mail.com", Password: string(hashPassword), Role: "admin"}, nil
}

func (mu mockUserRepository) DeleteUser(userID int) (entity.User, error) {
	// hashPassword, _ := bcrypt.GenerateFromPassword([]byte("admin1"), bcrypt.MinCost)

	return entity.User{
		ID: 0, Name: "", Email: "", Password: "", Role: ""}, nil
}

type mockCategoryRepository struct{}

func (mc mockCategoryRepository) GetAllCategory() ([]entity.Category, error) {
	return []entity.Category{
		{ID: 1, Name: "category1"}}, nil
}

func (mc mockCategoryRepository) GetCategoryByID(categoryID int) (entity.Category, error) {
	return entity.Category{
		ID: 1, Name: "category1"}, nil
}

func (mu mockCategoryRepository) CreateCategory(newCategory entity.Category) (entity.Category, error) {
	return entity.Category{
		ID: 1, Name: "category1"}, nil
}

func (mu mockCategoryRepository) UpdateCategory(categoryID int, updatedCategory entity.Category) (entity.Category, error) {
	return entity.Category{
		ID: 1, Name: "category1 new"}, nil
}

func (mu mockCategoryRepository) DeleteCategory(categoryID int) (entity.Category, error) {
	return entity.Category{
		ID: 0, Name: ""}, nil
}

type mockFalseCategoryRepository struct{}

func (mfc mockFalseCategoryRepository) GetAllCategory() ([]entity.Category, error) {
	return []entity.Category{
		{ID: 0, Name: ""}}, errors.New("can't get cities data")
}

func (mfc mockFalseCategoryRepository) GetCategoryByID(categoryID int) (entity.Category, error) {
	return entity.Category{
		ID: 0, Name: ""}, errors.New("can't get category data")
}

func (mfu mockFalseCategoryRepository) CreateCategory(newCategory entity.Category) (entity.Category, error) {
	return entity.Category{
		ID: 0, Name: ""}, errors.New("can't create category data")
}

func (mfu mockFalseCategoryRepository) UpdateCategory(categoryID int, updatedCategory entity.Category) (entity.Category, error) {
	return entity.Category{
		ID: 0, Name: ""}, errors.New("can't update category data")
}

func (mfu mockFalseCategoryRepository) DeleteCategory(categoryID int) (entity.Category, error) {
	return entity.Category{
		ID: 0, Name: ""}, errors.New("can't get category data")
}
