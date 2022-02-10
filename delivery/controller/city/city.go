package city

import (
	"net/http"
	"petshop/delivery/common"
	"petshop/entity"
	"petshop/repository/city"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CityController struct {
	cityRepo city.City
}

func NewCityController(cityRepo city.City) *CityController {
	return &CityController{cityRepo}
}

func (cc *CityController) CreateCity() echo.HandlerFunc {
	return func(c echo.Context) error {
		createCityReq := CreateCityFormatRequest{}

		c.Bind(&createCityReq)

		err := c.Validate(&createCityReq)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Can't get the input"))
		}

		newCity := entity.City{}
		newCity.Name = createCityReq.Name

		res, err := cc.cityRepo.CreateCity(newCity)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Can't create city"))
		}

		response := CityFormatResponse{}
		response.ID = res.ID
		response.Name = res.Name

		return c.JSON(http.StatusOK, common.SuccessResponse(response))
	}
}

func (cc *CityController) GetAllCity() echo.HandlerFunc {
	return func(c echo.Context) error {
		res, err := cc.cityRepo.GetAllCity()

		if err != nil || len(res) == 0 {
			return c.JSON(http.StatusNotFound, common.ErrorResponse(404, "City not found"))
		}

		responses := []CityFormatResponse{}
		response := CityFormatResponse{}

		for i := 0; i < len(res); i++ {
			response.ID = res[i].ID
			response.Name = res[i].Name
			responses = append(responses, response)
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(responses))
	}
}

func (cc *CityController) GetCityProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		cityID, _ := strconv.Atoi(c.Param("id"))

		res, err := cc.cityRepo.GetCityByID(cityID)

		if err != nil || res.ID == 0 {
			return c.JSON(http.StatusNotFound, common.ErrorResponse(404, "City not found"))
		}

		response := CityFormatResponse{}

		response.ID = res.ID
		response.Name = res.Name

		return c.JSON(http.StatusOK, common.SuccessResponse(response))
	}
}

func (cc *CityController) UpdateCityProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		cityID, _ := strconv.Atoi(c.Param("id"))

		updateCityReq := UpdateCityFormatRequest{}

		c.Bind(&updateCityReq)

		err := c.Validate(&updateCityReq)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Can't get the input"))
		}

		updatedCity := entity.City{}
		updatedCity.ID = uint(cityID)
		updatedCity.Name = updateCityReq.Name

		res, err := cc.cityRepo.UpdateCity(cityID, updatedCity)

		if err != nil || res.Name == "" {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Can't update city"))
		}

		response := CityFormatResponse{}
		response.ID = uint(cityID)
		response.Name = res.Name

		return c.JSON(http.StatusOK, common.SuccessResponse(response))
	}
}

func (cc *CityController) DeleteCity() echo.HandlerFunc {
	return func(c echo.Context) error {
		cityID, _ := strconv.Atoi(c.Param("id"))

		_, err := cc.cityRepo.DeleteCity(cityID)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Can't delete city"))
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(nil))
	}
}
