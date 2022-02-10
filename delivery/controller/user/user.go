package user

import (
	"net/http"
	"petshop/delivery/common"
	"petshop/delivery/middleware"
	"petshop/entity"
	"petshop/repository/user"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	userRepo user.User
}

func NewUserController(userRepo user.User) *UserController {
	return &UserController{userRepo}
}

func (uc *UserController) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		userRegisterReq := RegisterFormatRequest{}

		c.Bind(&userRegisterReq)

		err := c.Validate(&userRegisterReq)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Can't get the input"))
		}

		//_, err = uc.userRepo.FindCityByID(int(userRegisterReq.CityID))
		//
		//if err != nil {
		//	return c.JSON(http.StatusBadRequest, common.ErrorResponse(404, "City not found"))
		//}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(userRegisterReq.Password), bcrypt.MinCost)

		newUser := entity.User{}
		newUser.Name = userRegisterReq.Name
		newUser.Email = userRegisterReq.Email
		newUser.Password = string(hashedPassword)
		newUser.CityID = userRegisterReq.CityID

		res, err := uc.userRepo.CreateUser(newUser)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Can't create user"))
		}

		response := UserFormatResponse{}
		response.ID = res.ID
		response.Name = res.Name
		response.Email = res.Email
		response.CityID = res.CityID

		return c.JSON(http.StatusOK, common.SuccessResponse(response))
	}
}

func (uc *UserController) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		userLoginReq := LoginFormatRequest{}

		c.Bind(&userLoginReq)

		err := c.Validate(&userLoginReq)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Can't get the input"))
		}

		res, err := uc.userRepo.GetUserByEmail(userLoginReq.Email)

		if err != nil || res.ID == 0 {
			return c.JSON(http.StatusNotFound, common.ErrorResponse(404, "User not found"))
		}

		err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(userLoginReq.Password))

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Password doesn't match"))
		}

		token, _ := middleware.GenerateToken(int(res.ID), res.Email, res.Role)

		response := LoginFormatResponse{}
		response.Name = res.Name
		response.Email = res.Email
		response.Token = token

		return c.JSON(http.StatusOK, common.SuccessResponse(response))
	}
}

func (uc *UserController) GetUserProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := middleware.ExtractTokenUserID(c)

		res, err := uc.userRepo.GetUserByID(userID)

		if err != nil || res.ID == 0 {
			return c.JSON(http.StatusNotFound, common.ErrorResponse(404, "User not found"))
		}

		response := UserFormatResponse{}
		response.ID = res.ID
		response.Name = res.Name
		response.Email = res.Email
		response.CityID = res.CityID

		return c.JSON(http.StatusOK, common.SuccessResponse(response))
	}
}

func (uc *UserController) UpdateProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := middleware.ExtractTokenUserID(c)

		userUpdateReq := UpdateFormatRequest{}

		c.Bind(&userUpdateReq)

		err := c.Validate(&userUpdateReq)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Can't get the input"))
		}

		//city, err := uc.userRepo.FindCityByID(int(userUpdateReq.CityID))
		//
		//if err != nil || city.ID == 0 {
		//	return c.JSON(http.StatusBadRequest, common.ErrorResponse(404, "City not found"))
		//}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(userUpdateReq.Password), bcrypt.MinCost)

		updatedUser := entity.User{}
		updatedUser.ID = uint(userID)
		updatedUser.Name = userUpdateReq.Name
		updatedUser.Email = userUpdateReq.Email
		updatedUser.Password = string(hashedPassword)
		updatedUser.CityID = userUpdateReq.CityID

		res, err := uc.userRepo.UpdateUser(userID, updatedUser)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Can't update user"))
		}

		response := UserFormatResponse{}
		response.ID = res.ID
		response.Name = res.Name
		response.Email = res.Email
		response.CityID = res.CityID

		return c.JSON(http.StatusOK, common.SuccessResponse(response))
	}
}

func (uc *UserController) DeleteAccount() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := middleware.ExtractTokenUserID(c)

		_, err := uc.userRepo.DeleteUser(userID)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Can't delete account"))
		}

		return c.JSON(http.StatusBadRequest, common.SuccessResponse(nil))
	}
}
