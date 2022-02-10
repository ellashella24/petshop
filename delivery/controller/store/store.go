package store

import (
	"net/http"
	"petshop/delivery/common"
	"petshop/delivery/middleware"
	"petshop/entity"
	"petshop/repository/store"
	"strconv"

	"github.com/labstack/echo/v4"
)

type StoreController struct {
	storeRepo store.Store
}

func NewStoreController(storeRepo store.Store) *StoreController {
	return &StoreController{storeRepo}
}

func (sc *StoreController) CreateStore() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := middleware.ExtractTokenUserID(c)

		createStoreReq := CreateStoreFormatRequest{}

		c.Bind(&createStoreReq)

		err := c.Validate(&createStoreReq)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Can't get the input"))
		}

		newStore := entity.Store{}
		newStore.Name = createStoreReq.Name
		newStore.CityID = createStoreReq.CityID
		newStore.UserID = uint(userID)

		res, err := sc.storeRepo.CreateStore(newStore)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Can't create store"))
		}

		response := StoreFormatResponse{}
		response.ID = res.ID
		response.Name = res.Name
		response.CityID = res.CityID
		response.UserID = res.UserID

		return c.JSON(http.StatusOK, common.SuccessResponse(response))
	}
}

func (sc *StoreController) GetAllStoreByUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := middleware.ExtractTokenUserID(c)

		res, err := sc.storeRepo.GetAllStoreByUserID(userID)

		if err != nil || len(res) == 0 {
			return c.JSON(http.StatusNotFound, common.ErrorResponse(404, "Store not found"))
		}

		responses := []StoreFormatResponse{}
		response := StoreFormatResponse{}

		for i := 0; i < len(res); i++ {
			response.ID = res[i].ID
			response.Name = res[i].Name
			response.CityID = res[i].CityID
			response.UserID = res[i].UserID

			responses = append(responses, response)
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(responses))
	}
}

func (sc *StoreController) GetStoreProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		storeID, _ := strconv.Atoi(c.Param("id"))

		userID := middleware.ExtractTokenUserID(c)

		res, err := sc.storeRepo.GetStoreProfile(storeID, userID)

		if err != nil || res.ID == 0 {
			return c.JSON(http.StatusNotFound, common.ErrorResponse(404, "Store not found"))
		}

		response := StoreFormatResponse{}

		response.ID = res.ID
		response.Name = res.Name
		response.CityID = res.CityID
		response.UserID = res.UserID

		return c.JSON(http.StatusOK, common.SuccessResponse(response))
	}
}

func (sc *StoreController) UpdateStoreProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		storeID, _ := strconv.Atoi(c.Param("id"))

		userID := middleware.ExtractTokenUserID(c)

		updateStoreReq := UpdateStoreFormatRequest{}

		c.Bind(&updateStoreReq)

		err := c.Validate(&updateStoreReq)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Can't get the input"))
		}

		updatedStore := entity.Store{}
		updatedStore.Name = updateStoreReq.Name
		updatedStore.CityID = updateStoreReq.CityID
		updatedStore.UserID = uint(userID)

		res, err := sc.storeRepo.UpdateStoreProfile(storeID, userID, updatedStore)

		if err != nil || res.Name == "" {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Can't update store"))
		}

		response := StoreFormatResponse{}
		response.ID = uint(storeID)
		response.Name = res.Name
		response.CityID = res.CityID
		response.UserID = uint(userID)

		return c.JSON(http.StatusOK, common.SuccessResponse(response))
	}
}

func (sc *StoreController) DeleteStore() echo.HandlerFunc {
	return func(c echo.Context) error {
		storeID, _ := strconv.Atoi(c.Param("id"))

		userID := middleware.ExtractTokenUserID(c)

		_, err := sc.storeRepo.DeleteStore(storeID, userID)

		if err != nil {
			return c.JSON(http.StatusNotFound, common.ErrorResponse(400, "Can't delete store"))
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(nil))
	}
}
