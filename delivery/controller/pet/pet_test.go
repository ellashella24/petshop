package pet

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

func TestCreatePet(t *testing.T) {
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

	t.Run("1. Error create pet - Error validate", func(t *testing.T) {
		ev := echo.New()

		requestBody, _ := json.Marshal(CreatePetFormatRequest{
			Name: "pet1",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		req.Header.Set("Content-Type", "application/json")
		context := ev.NewContext(req, res)
		context.SetPath("/user/pet")

		petController := NewPetController(mockPetRepository{})
		err := mw.IsAdmin(petController.CreatePet())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Can't get the input", response.Message)
	})

	t.Run("2. Error create pet - Can't create pet", func(t *testing.T) {
		requestBody, _ := json.Marshal(CreatePetFormatRequest{
			Name: "pet1",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		context := e.NewContext(req, res)
		context.SetPath("/user/pet")

		petController := NewPetController(mockFalsePetRepository{})
		petController.CreatePet()(context)

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Can't create pet", response.Message)
	})

	t.Run("3. Success create pet", func(t *testing.T) {
		requestBody, _ := json.Marshal(CreatePetFormatRequest{
			Name: "pet1",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		req.Header.Set("Content-Type", "application/json")
		context := e.NewContext(req, res)
		context.SetPath("/user/pet")

		petController := NewPetController(mockPetRepository{})
		err := mw.IsAdmin(petController.CreatePet())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})
}

func TestGetAllPetByUser(t *testing.T) {
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

	t.Run("1. Error Get All Pet By User Test", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/user/pets")

		petController := NewPetController(mockFalsePetRepository{})
		err := (mw.IsAdmin)(petController.GetAllPetByUser())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Pet not found", response.Message)
	})

	t.Run("2. Success Get All Pet By User ID Test", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/cities")

		petController := NewPetController(mockPetRepository{})
		err := (mw.IsAdmin)(petController.GetAllPetByUser())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})
}

func TestGetPetProfile(t *testing.T) {
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

	t.Run("1. Error Get Pet Profile Test", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/user/pet/profile")
		context.SetParamNames("id")
		context.SetParamValues("1")

		petController := NewPetController(mockFalsePetRepository{})
		err := (mw.IsAdmin)(petController.GetPetProfile())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Pet not found", response.Message)
	})

	t.Run("2. Success Get Pet Profile Test", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/user/pet")
		context.SetParamNames("id")
		context.SetParamValues("1")

		petController := NewPetController(mockPetRepository{})
		err := (mw.IsAdmin)(petController.GetPetProfile())(context)

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

		petController := NewPetController(mockFalsePetRepository{})
		petController.GetGroomingStatusByPetID()(context)

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

		petController := NewPetController(mockPetRepository{})
		petController.GetGroomingStatusByPetID()(context)

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
		context.SetPath("/user/grooming_status/pet")
		context.SetParamNames("id")
		context.SetParamValues("1")

		petController := NewPetController(mockFalsePetRepository{})
		petController.UpdateFinalGroomingStatus()(context)

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Can't get grooming status", response.Message)
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
		context.SetPath("/user/grooming_status/pet")
		context.SetParamNames("id")
		context.SetParamValues("1")

		petController := NewPetController(mockPetRepository{})
		petController.UpdateFinalGroomingStatus()(context)

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})
}

func TestUpdatePetProfile(t *testing.T) {
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

	t.Run("1. Error Update Pet - Error Validate", func(t *testing.T) {
		ev := echo.New()

		requestBody, _ := json.Marshal(UpdatePetFormatRequest{
			Name: "pet1 new",
		})

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := ev.NewContext(req, res)
		context.SetPath("/user/pet/profile")
		context.SetParamNames("id")
		context.SetParamValues("1")

		petController := NewPetController(mockFalsePetRepository{})

		err := (mw.IsAdmin)(petController.UpdatePetProfile())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.DefaultResponse

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Can't get the input", response.Message)
	})

	t.Run("2. Error Update City - Can't Create Input", func(t *testing.T) {
		requestBody, _ := json.Marshal(UpdatePetFormatRequest{
			Name: "pet1 new",
		})

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/user/pet/profile")
		context.SetParamNames("id")
		context.SetParamValues("1")

		petController := NewPetController(mockFalsePetRepository{})

		err := (mw.IsAdmin)(petController.UpdatePetProfile())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.DefaultResponse

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Can't update pet", response.Message)
	})

	t.Run("3. Success Update City", func(t *testing.T) {
		requestBody, _ := json.Marshal(CreatePetFormatRequest{
			Name: "pet1 new",
		})

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/city/profile")
		context.SetParamNames("id")
		context.SetParamValues("1")

		cityController := NewPetController(mockPetRepository{})

		err := (mw.IsAdmin)(cityController.UpdatePetProfile())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.DefaultResponse

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})
}

func TestDeletePet(t *testing.T) {
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

	t.Run("1. Error Delete Pet Test", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/user/pet")
		context.SetParamNames("id")
		context.SetParamValues("1")

		cityController := NewPetController(mockFalsePetRepository{})
		err := (mw.IsAdmin)(cityController.DeletePet())(context)

		if err != nil {
			fmt.Println(err)
			return
		}

		var response common.ResponseSuccess

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Can't delete pet", response.Message)
	})

	t.Run("2. Success Delete Pet Test", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		context := e.NewContext(req, res)
		context.SetPath("/user/pet")
		context.SetParamNames("id")
		context.SetParamValues("1")

		petController := NewPetController(mockPetRepository{})
		err := (mw.IsAdmin)(petController.DeletePet())(context)

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

type mockPetRepository struct{}

func (mc mockPetRepository) GetAllPetByUserID(userID int) ([]entity.Pet, error) {
	return []entity.Pet{
		{ID: 1, Name: "pet1", UserID: 1}}, nil
}

func (mc mockPetRepository) GetPetProfileByID(petID, userID int) (entity.Pet, error) {
	return entity.Pet{
		ID: 1, Name: "city1", UserID: 1}, nil
}

func (mc mockPetRepository) GetGroomingStatusByPetID(petID, userID int) (entity.GroomingStatus, error) {
	return entity.GroomingStatus{
		ID: 1, Status: "TELAH DIBAYAR", PetID: 1}, nil
}

func (mc mockPetRepository) UpdateFinalGroomingStatus(petID, userID int) (entity.GroomingStatus, error) {
	return entity.GroomingStatus{
		ID: 1, Status: "SELESAI", PetID: 1}, nil
}

func (mc mockPetRepository) CreatePet(newPet entity.Pet) (entity.Pet, error) {
	return entity.Pet{
		ID: 1, Name: "city1", UserID: 1}, nil
}

func (mc mockPetRepository) UpdatePetProfile(petID, userID int, updatedPet entity.Pet) (entity.Pet, error) {
	return entity.Pet{
		ID: 1, Name: "city1 new", UserID: 1}, nil
}

func (mc mockPetRepository) DeletePet(petID, userID int) (entity.Pet, error) {
	return entity.Pet{
		ID: 0, Name: "", UserID: 0}, nil
}

type mockFalsePetRepository struct{}

func (mfc mockFalsePetRepository) GetAllPetByUserID(userID int) ([]entity.Pet, error) {
	return []entity.Pet{
		{ID: 0, Name: "", UserID: 0}}, errors.New("can't get cities data")
}

func (mfc mockFalsePetRepository) GetPetProfileByID(petID, userID int) (entity.Pet, error) {
	return entity.Pet{
		ID: 0, Name: "", UserID: 0}, errors.New("can't get city data")
}

func (mfc mockFalsePetRepository) GetGroomingStatusByPetID(petID, userID int) (entity.GroomingStatus, error) {
	return entity.GroomingStatus{
		ID: 0, Status: "", PetID: 0}, errors.New("not found grooming status")
}

func (mfc mockFalsePetRepository) UpdateFinalGroomingStatus(petID, userID int) (entity.GroomingStatus, error) {
	return entity.GroomingStatus{
		ID: 0, Status: "", PetID: 0}, errors.New("not found grooming status")
}

func (mfu mockFalsePetRepository) CreatePet(newPet entity.Pet) (entity.Pet, error) {
	return entity.Pet{
		ID: 0, Name: "", UserID: 0}, errors.New("can't create city data")
}

func (mfu mockFalsePetRepository) UpdatePetProfile(petID, userID int, updatedPet entity.Pet) (entity.Pet, error) {
	return entity.Pet{
		ID: 0, Name: "", UserID: 0}, errors.New("can't update city data")
}

func (mfu mockFalsePetRepository) DeletePet(petID, userID int) (entity.Pet, error) {
	return entity.Pet{
		ID: 0, Name: "", UserID: 0}, errors.New("can't get city data")
}
