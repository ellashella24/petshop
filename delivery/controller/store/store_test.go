package store

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"petshop/delivery/common"
	"petshop/delivery/controller/user"
	mw "petshop/delivery/middleware"
	"petshop/entity"
	"testing"
	"time"

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

func TestCreateStore(t *testing.T) {
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

	t.Run("1. Error create store - Error validate", func(t *testing.T) {
		ev := echo.New()

		requestBody, _ := json.Marshal(CreateStoreFormatRequest{
			Name:   "store1",
			CityID: 1,
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		req.Header.Set("Content-Type", "application/json")
		context := ev.NewContext(req, res)
		context.SetPath("/user/store")

		storeController := NewStoreController(mockStoreRepository{})
		err := mw.IsAdmin(storeController.CreateStore())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Can't get the input", response.Message)
	})

	t.Run("2. Error create store - Can't create store", func(t *testing.T) {
		requestBody, _ := json.Marshal(CreateStoreFormatRequest{
			Name:   "store1",
			CityID: 1,
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		req.Header.Set("Content-Type", "application/json")
		context := e.NewContext(req, res)
		context.SetPath("/user/store")

		storeController := NewStoreController(mockFalseStoreRepository{})
		err := mw.IsAdmin(storeController.CreateStore())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Can't create store", response.Message)
	})

	t.Run("3. Success create store", func(t *testing.T) {
		requestBody, _ := json.Marshal(CreateStoreFormatRequest{
			Name:   "store2",
			CityID: 1,
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		req.Header.Set("Content-Type", "application/json")
		context := e.NewContext(req, res)
		context.SetPath("/user/store")

		storeController := NewStoreController(mockStoreRepository{})
		err := mw.IsAdmin(storeController.CreateStore())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})
}

func TestGetAllStoreByUser(t *testing.T) {
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

	t.Run("1. Error Get All Store By User Test", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/user/stores")

		storeController := NewStoreController(mockFalseStoreRepository{})
		err := (mw.IsAdmin)(storeController.GetAllStoreByUser())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Store not found", response.Message)
	})

	t.Run("2. Success Get All Store By User ID Test", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/cities")

		storeController := NewStoreController(mockStoreRepository{})
		err := (mw.IsAdmin)(storeController.GetAllStoreByUser())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})
}

func TestGetStoreProfile(t *testing.T) {
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

	t.Run("1. Error Get Store Profile Test", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/user/store/profile")
		context.SetParamNames("id")
		context.SetParamValues("1")

		storeController := NewStoreController(mockFalseStoreRepository{})
		err := (mw.IsAdmin)(storeController.GetStoreProfile())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Store not found", response.Message)
	})

	t.Run("2. Success Get Store Profile Test", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/user/store")
		context.SetParamNames("id")
		context.SetParamValues("1")

		storeController := NewStoreController(mockStoreRepository{})
		err := (mw.IsAdmin)(storeController.GetStoreProfile())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})
}

func TestGetGroomingStatusByPetID(t *testing.T) {
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

	t.Run("1. Error get grooming status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		context := e.NewContext(req, res)
		context.SetPath("/user/grooming_status/pet")
		context.SetParamNames("id")
		context.SetParamValues("1")

		storeController := NewStoreController(mockFalseStoreRepository{})
		storeController.GetGroomingStatusByPetID()(context)

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Can't get grooming status", response.Message)
	})

	t.Run("2. Success get grooming status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		context := e.NewContext(req, res)
		context.SetPath("/user/grooming_status/pet")
		context.SetParamNames("id")
		context.SetParamValues("1")

		storeController := NewStoreController(mockStoreRepository{})
		storeController.GetGroomingStatusByPetID()(context)

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})
}

func TestUpdateGroomingStatus(t *testing.T) {
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

	t.Run("1. Error get grooming status", func(t *testing.T) {
		body := new(bytes.Buffer)

		writer := multipart.NewWriter(body)

		writer.WriteField("pet_id", "1000")
		writer.WriteField("store_id", "1000")

		writer.Close()

		req := httptest.NewRequest(http.MethodPut, "/", body)
		res := httptest.NewRecorder()
		req.Header.Add("Content-Type", writer.FormDataContentType())

		context := e.NewContext(req, res)
		context.SetPath("/store/grooming_status/pet")
		context.SetParamNames("id")
		context.SetParamValues("1")

		storeController := NewStoreController(mockFalseStoreRepository{})
		storeController.UpdateGroomingStatus()(context)

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Can't update grooming status", response.Message)
	})

	t.Run("2. Success get grooming status", func(t *testing.T) {
		body := new(bytes.Buffer)

		writer := multipart.NewWriter(body)

		writer.WriteField("pet_id", "1")
		writer.WriteField("store_id", "1")

		writer.Close()

		req := httptest.NewRequest(http.MethodPut, "/", body)
		res := httptest.NewRecorder()
		req.Header.Add("Content-Type", writer.FormDataContentType())

		context := e.NewContext(req, res)
		context.SetPath("/store/grooming_status/pet")
		context.SetParamNames("id")
		context.SetParamValues("1")

		storeController := NewStoreController(mockStoreRepository{})
		storeController.UpdateGroomingStatus()(context)

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})
}

func TestUpdateStoreProfile(t *testing.T) {
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

	t.Run("1. Error Update Store - Error Validate", func(t *testing.T) {
		ev := echo.New()

		requestBody, _ := json.Marshal(UpdateStoreFormatRequest{
			Name: "store1 new",
		})

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := ev.NewContext(req, res)
		context.SetPath("/user/store/profile")
		context.SetParamNames("id")
		context.SetParamValues("1")

		storeController := NewStoreController(mockFalseStoreRepository{})

		err := (mw.IsAdmin)(storeController.UpdateStoreProfile())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.DefaultResponse

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Can't get the input", response.Message)
	})

	t.Run("2. Error Update Store - Can't Update Store", func(t *testing.T) {
		requestBody, _ := json.Marshal(UpdateStoreFormatRequest{
			Name:   "store1 new",
			CityID: 1,
		})

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/user/store/profile")
		context.SetParamNames("id")
		context.SetParamValues("1")

		storeController := NewStoreController(mockFalseStoreRepository{})

		err := (mw.IsAdmin)(storeController.UpdateStoreProfile())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.DefaultResponse

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Can't update store", response.Message)
	})

	t.Run("3. Success Update Store", func(t *testing.T) {
		requestBody, _ := json.Marshal(CreateStoreFormatRequest{
			Name:   "store1 new",
			CityID: 1,
		})

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/city/profile")
		context.SetParamNames("id")
		context.SetParamValues("1")

		cityController := NewStoreController(mockStoreRepository{})

		err := (mw.IsAdmin)(cityController.UpdateStoreProfile())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.DefaultResponse

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})
}

func TestDeleteStore(t *testing.T) {
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

	t.Run("1. Error Delete Store Test", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/user/store")
		context.SetParamNames("id")
		context.SetParamValues("1")

		cityController := NewStoreController(mockFalseStoreRepository{})
		err := (mw.IsAdmin)(cityController.DeleteStore())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Can't delete store", response.Message)
	})

	t.Run("2. Success Delete Store Test", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/user/store")
		context.SetParamNames("id")
		context.SetParamValues("1")

		storeController := NewStoreController(mockStoreRepository{})
		err := (mw.IsAdmin)(storeController.DeleteStore())(context)

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

type mockStoreRepository struct{}

func (mu mockStoreRepository) FindCityByID(cityID int) (entity.City, error) {
	return entity.City{
		ID: 1, Name: "city1"}, nil
}

func (mc mockStoreRepository) GetAllStoreByUserID(userID int) ([]entity.Store, error) {
	return []entity.Store{
		{ID: 1, Name: "store1", UserID: 1}}, nil
}

func (mc mockStoreRepository) GetStoreProfile(storeID, userID int) (entity.Store, error) {
	return entity.Store{
		ID: 1, Name: "city1", UserID: 1}, nil
}

func (mc mockStoreRepository) GetListTransactionByStoreID(storeID int) ([]entity.Transaction, error) {
	return []entity.Transaction{
		{ID: 1, UserID: 1, InvoiceID: "Invoice-1", PaymentMethod: "Bank Transfer", PaidAt: time.Now(), TotalPrice: 100000, PaymentStatus: "Paid"}}, nil
}

func (mc mockStoreRepository) GetGroomingStatusByPetID(storeID, petID int) (entity.GroomingStatus, error) {
	return entity.GroomingStatus{
		ID: 1, Status: "TELAH DIBAYAR", PetID: 1}, nil
}

func (mc mockStoreRepository) UpdateGroomingStatus(storeID, petID int) (entity.GroomingStatus, error) {
	return entity.GroomingStatus{
		ID: 1, Status: "PROSES GROOMING", PetID: 1}, nil
}

func (mu mockStoreRepository) CreateStore(newStore entity.Store) (entity.Store, error) {
	return entity.Store{
		ID: 1, Name: "city1", UserID: 1}, nil
}

func (mu mockStoreRepository) UpdateStoreProfile(storeID, userID int, updatedStore entity.Store) (entity.Store, error) {
	return entity.Store{
		ID: 1, Name: "city1 new", UserID: 1}, nil
}

func (mu mockStoreRepository) DeleteStore(storeID, userID int) (entity.Store, error) {
	return entity.Store{
		ID: 0, Name: "", UserID: 0}, nil
}

type mockFalseStoreRepository struct{}

func (mfc mockFalseStoreRepository) FindCityByID(cityID int) (entity.City, error) {
	return entity.City{
		ID: 0, Name: ""}, nil
}

func (mfc mockFalseStoreRepository) GetAllStoreByUserID(userID int) ([]entity.Store, error) {
	return []entity.Store{
		{ID: 0, Name: "", UserID: 0}}, errors.New("can't get cities data")
}

func (mfc mockFalseStoreRepository) GetStoreProfile(storeID, userID int) (entity.Store, error) {
	return entity.Store{
		ID: 0, Name: "", UserID: 0}, errors.New("can't get city data")
}

func (mc mockFalseStoreRepository) GetGroomingStatusByPetID(storeID, petID int) (entity.GroomingStatus, error) {
	return entity.GroomingStatus{
		ID: 0, Status: "", PetID: 0}, errors.New("not found grooming status")
}

func (mc mockFalseStoreRepository) GetListTransactionByStoreID(storeID int) ([]entity.Transaction, error) {
	return []entity.Transaction{
		{ID: 0, UserID: 0, InvoiceID: "", PaymentMethod: "", TotalPrice: 0, PaymentStatus: ""}}, errors.New("not found transactions")
}

func (mc mockFalseStoreRepository) UpdateGroomingStatus(storeID, petID int) (entity.GroomingStatus, error) {
	return entity.GroomingStatus{
		ID: 0, Status: "", PetID: 0}, errors.New("not found grooming status")
}

func (mfu mockFalseStoreRepository) CreateStore(newStore entity.Store) (entity.Store, error) {
	return entity.Store{
		ID: 0, Name: "", UserID: 0}, errors.New("can't create city data")
}

func (mfu mockFalseStoreRepository) UpdateStoreProfile(storeID, userID int, updatedStore entity.Store) (entity.Store, error) {
	return entity.Store{
		ID: 0, Name: "", UserID: 0}, errors.New("can't update city data")
}

func (mfu mockFalseStoreRepository) DeleteStore(storeID, userID int) (entity.Store, error) {
	return entity.Store{
		ID: 0, Name: "", UserID: 0}, errors.New("can't get city data")
}
