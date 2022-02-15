package pet

import (
	"net/http"
	"petshop/delivery/common"
	"petshop/delivery/middleware"
	"petshop/entity"
	"petshop/repository/pet"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PetController struct {
	petRepo pet.Pet
}

func NewPetController(petRepo pet.Pet) *PetController {
	return &PetController{petRepo}
}

func (pc *PetController) CreatePet() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := middleware.ExtractTokenUserID(c)

		createNewPetReq := CreatePetFormatRequest{}

		c.Bind(&createNewPetReq)

		err := c.Validate(&createNewPetReq)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Can't get the input"))
		}

		newPet := entity.Pet{}
		newPet.Name = createNewPetReq.Name
		newPet.UserID = uint(userID)

		res, err := pc.petRepo.CreatePet(newPet)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Can't create pet"))
		}

		response := PetFormatResponse{}
		response.ID = res.ID
		response.Name = res.Name
		response.UserID = res.UserID

		return c.JSON(http.StatusOK, common.SuccessResponse(response))
	}
}

func (pc *PetController) GetAllPetByUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := middleware.ExtractTokenUserID(c)

		res, err := pc.petRepo.GetAllPetByUserID(userID)

		if err != nil || len(res) == 0 {
			return c.JSON(http.StatusNotFound, common.ErrorResponse(404, "Pet not found"))
		}

		responses := []PetFormatResponse{}
		response := PetFormatResponse{}

		for i := 0; i < len(res); i++ {
			response.ID = res[i].ID
			response.Name = res[i].Name
			response.UserID = res[i].UserID

			responses = append(responses, response)
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(responses))
	}
}

func (pc *PetController) GetPetProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		petID, _ := strconv.Atoi(c.Param("id"))

		userID := middleware.ExtractTokenUserID(c)

		res, err := pc.petRepo.GetPetProfileByID(petID, userID)

		if err != nil || res.ID == 0 {
			return c.JSON(http.StatusNotFound, common.ErrorResponse(404, "Pet not found"))
		}

		response := PetFormatResponse{}
		response.ID = res.ID
		response.Name = res.Name
		response.UserID = res.UserID

		return c.JSON(http.StatusOK, common.SuccessResponse(response))
	}
}

func (pc *PetController) UpdatePetProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		petID, _ := strconv.Atoi(c.Param("id"))

		userID := middleware.ExtractTokenUserID(c)

		updatePetReq := UpdatePetFormatRequest{}

		c.Bind(&updatePetReq)

		err := c.Validate(&updatePetReq)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Can't get the input"))
		}

		updatedPet := entity.Pet{}
		updatedPet.Name = updatePetReq.Name

		res, err := pc.petRepo.UpdatePetProfile(petID, userID, updatedPet)

		if err != nil || res.Name == "" {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Can't update pet"))
		}

		response := PetFormatResponse{}
		response.ID = uint(petID)
		response.Name = res.Name
		response.UserID = uint(userID)

		return c.JSON(http.StatusOK, common.SuccessResponse(response))
	}
}

func (pc *PetController) DeletePet() echo.HandlerFunc {
	return func(c echo.Context) error {
		petID, _ := strconv.Atoi(c.Param("id"))

		userID := middleware.ExtractTokenUserID(c)

		_, err := pc.petRepo.DeletePet(petID, userID)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Can't delete pet"))
		}

		return c.JSON(http.StatusOK, common.SuccessResponse(nil))

	}
}

func (pc *PetController) GetGroomingStatusByPetID() echo.HandlerFunc {
	return func(c echo.Context) error {
		getGroomingReq := GetGroomingStatusFormatRequest{}
		userID := middleware.ExtractTokenUserID(c)

		c.Bind(&getGroomingReq)

		c.Validate(&getGroomingReq)

		res, err := pc.petRepo.GetGroomingStatusByPetID(int(getGroomingReq.PetID), userID)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Can't get grooming status"))
		}

		response := GroomingStatusResponse{}
		response.ID = res.ID
		response.PetID = res.PetID
		response.Status = res.Status

		return c.JSON(http.StatusOK, common.SuccessResponse(response))
	}
}

func (pc *PetController) UpdateFinalGroomingStatus() echo.HandlerFunc {
	return func(c echo.Context) error {
		getGroomingReq := GetGroomingStatusFormatRequest{}
		userID := middleware.ExtractTokenUserID(c)

		c.Bind(&getGroomingReq)

		c.Validate(&getGroomingReq)

		res, err := pc.petRepo.UpdateFinalGroomingStatus(int(getGroomingReq.PetID), userID)

		if err != nil {
			return c.JSON(http.StatusBadRequest, common.ErrorResponse(400, "Can't get grooming status"))
		}

		response := GroomingStatusResponse{}
		response.ID = res.ID
		response.PetID = res.PetID
		response.Status = res.Status

		return c.JSON(http.StatusOK, common.SuccessResponse(response))
	}
}
