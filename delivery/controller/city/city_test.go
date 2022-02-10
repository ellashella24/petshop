package city

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

func TestCreateCity(t *testing.T) {
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

	t.Run("1. Error Create City - Error Validate", func(t *testing.T) {
		ev := echo.New()

		requestBody, _ := json.Marshal(CreateCityFormatRequest{
			Name: "city1",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		req.Header.Set("Content-Type", "application/json")
		context := ev.NewContext(req, res)
		context.SetPath("/city")

		cityController := NewCityController(mockCityRepository{})
		err := mw.IsAdmin(cityController.CreateCity())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Can't get the input", response.Message)
	})

	t.Run("2. Error Create City - Can't Create Input", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}

		requestBody, _ := json.Marshal(CreateCityFormatRequest{
			Name: "city1",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/city")

		cityController := NewCityController(mockFalseCityRepository{})

		err := (mw.IsAdmin)(cityController.CreateCity())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.DefaultResponse

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Can't create city", response.Message)
	})

	t.Run("3. Success Create City", func(t *testing.T) {
		e := echo.New()
		e.Validator = &CustomValidator{validator: validator.New()}

		requestBody, _ := json.Marshal(CreateCityFormatRequest{
			Name: "city2",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/city")

		cityController := NewCityController(mockCityRepository{})

		err := (mw.IsAdmin)(cityController.CreateCity())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.DefaultResponse

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})
}

func TestGetAllCity(t *testing.T) {
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

	t.Run("1. Error Get All City Test", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/cities")

		cityController := NewCityController(mockFalseCityRepository{})
		err := (mw.IsAdmin)(cityController.GetAllCity())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "City not found", response.Message)
	})

	t.Run("2. Success Get All City Test", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/cities")

		cityController := NewCityController(mockCityRepository{})
		err := (mw.IsAdmin)(cityController.GetAllCity())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})
}

func TestGetCityProfile(t *testing.T) {
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

	t.Run("1. Error Get City Profile Test", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/city/profile")
		context.SetParamNames("id")
		context.SetParamValues("1")

		cityController := NewCityController(mockFalseCityRepository{})
		err := (mw.IsAdmin)(cityController.GetCityProfile())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "City not found", response.Message)
	})

	t.Run("2. Success Get City Profile Test", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/city")
		context.SetParamNames("id")
		context.SetParamValues("1")

		cityController := NewCityController(mockCityRepository{})
		err := (mw.IsAdmin)(cityController.GetCityProfile())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})
}

func TestUpdateCityProfile(t *testing.T) {
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

	t.Run("1. Error Update City - Error Validate", func(t *testing.T) {
		ev := echo.New()

		requestBody, _ := json.Marshal(UpdateCityFormatRequest{
			Name: "city1 new",
		})

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := ev.NewContext(req, res)
		context.SetPath("/city/profile")
		context.SetParamNames("id")
		context.SetParamValues("1")

		cityController := NewCityController(mockFalseCityRepository{})

		err := (mw.IsAdmin)(cityController.UpdateCityProfile())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.DefaultResponse

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Can't get the input", response.Message)
	})

	t.Run("2. Error Update City - Can't Create Input", func(t *testing.T) {
		requestBody, _ := json.Marshal(UpdateCityFormatRequest{
			Name: "city1 new",
		})

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/city/profile")
		context.SetParamNames("id")
		context.SetParamValues("1")

		cityController := NewCityController(mockFalseCityRepository{})

		err := (mw.IsAdmin)(cityController.UpdateCityProfile())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.DefaultResponse

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Can't update city", response.Message)
	})

	t.Run("3. Success Update City", func(t *testing.T) {
		requestBody, _ := json.Marshal(CreateCityFormatRequest{
			Name: "city1 new",
		})

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/city/profile")
		context.SetParamNames("id")
		context.SetParamValues("1")

		cityController := NewCityController(mockCityRepository{})

		err := (mw.IsAdmin)(cityController.UpdateCityProfile())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.DefaultResponse

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})
}

func TestDeleteCity(t *testing.T) {
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

	t.Run("1. Error Delete City Test", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/city")
		context.SetParamNames("id")
		context.SetParamValues("1")

		cityController := NewCityController(mockFalseCityRepository{})
		err := (mw.IsAdmin)(cityController.DeleteCity())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Can't delete city", response.Message)
	})

	t.Run("2. Success Delete City Test", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/city")
		context.SetParamNames("id")
		context.SetParamValues("1")

		cityController := NewCityController(mockCityRepository{})
		err := (mw.IsAdmin)(cityController.DeleteCity())(context)

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

func (mu mockUserRepository) FindCityByID(cityID int) (entity.City, error) {
	return entity.City{
		ID: 1, Name: "city1"}, nil
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

type mockCityRepository struct{}

func (mc mockCityRepository) GetAllCity() ([]entity.City, error) {
	return []entity.City{
		{ID: 1, Name: "city1"}}, nil
}

func (mc mockCityRepository) GetCityByID(cityID int) (entity.City, error) {
	return entity.City{
		ID: 1, Name: "city1"}, nil
}

func (mu mockCityRepository) CreateCity(newCity entity.City) (entity.City, error) {
	return entity.City{
		ID: 1, Name: "city1"}, nil
}

func (mu mockCityRepository) UpdateCity(cityID int, updatedCity entity.City) (entity.City, error) {
	return entity.City{
		ID: 1, Name: "city1 new"}, nil
}

func (mu mockCityRepository) DeleteCity(cityID int) (entity.City, error) {
	return entity.City{
		ID: 0, Name: ""}, nil
}

type mockFalseCityRepository struct{}

func (mfc mockFalseCityRepository) GetAllCity() ([]entity.City, error) {
	return []entity.City{
		{ID: 0, Name: ""}}, errors.New("can't get cities data")
}

func (mfc mockFalseCityRepository) GetCityByID(cityID int) (entity.City, error) {
	return entity.City{
		ID: 0, Name: ""}, errors.New("can't get city data")
}

func (mfu mockFalseCityRepository) CreateCity(newCity entity.City) (entity.City, error) {
	return entity.City{
		ID: 0, Name: ""}, errors.New("can't create city data")
}

func (mfu mockFalseCityRepository) UpdateCity(cityID int, updatedCity entity.City) (entity.City, error) {
	return entity.City{
		ID: 0, Name: ""}, errors.New("can't update city data")
}

func (mfu mockFalseCityRepository) DeleteCity(cityID int) (entity.City, error) {
	return entity.City{
		ID: 0, Name: ""}, errors.New("can't get city data")
}
