package middleware

import (
	"errors"
	"net/http"
	"petshop/constants"
	"petshop/delivery/common"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

func GenerateToken(userID int, email string, role string) (string, error) {
	claim := jwt.MapClaims{}
	claim["user_id"] = userID
	claim["email"] = email
	claim["role"] = role
	claim["exp"] = time.Now().Add(time.Hour * 120).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, err := token.SignedString([]byte(constants.SecretKey))

	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}

func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get(echo.HeaderAuthorization)

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)

			if !ok {
				return nil, errors.New("invalid token")
			}

			return []byte(constants.SecretKey), nil
		})

		claim, _ := token.Claims.(jwt.MapClaims)

		role := claim["role"]
		if role != "admin" {
			return c.JSON(http.StatusUnauthorized, common.ErrorResponse(401, "Not admin"))
		}

		return next(c)
	}
}

func ExtractTokenUserID(c echo.Context) int {
	authHeader := c.Request().Header.Get(echo.HeaderAuthorization)

	if !strings.Contains(authHeader, "Bearer") {
		return 0
	}

	tokenString := ""
	arrayToken := strings.Split(authHeader, " ")
	if len(arrayToken) == 2 {
		tokenString = arrayToken[1]
	}

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("invalid token")
		}

		return []byte(constants.SecretKey), nil
	})

	claim, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return 0
	}

	userID := int(claim["user_id"].(float64))

	return userID
}

func ExtractTokenEmail(c echo.Context) string {
	authHeader := c.Request().Header.Get(echo.HeaderAuthorization)

	if !strings.Contains(authHeader, "Bearer") {
		return ""
	}

	tokenString := ""
	arrayToken := strings.Split(authHeader, " ")
	if len(arrayToken) == 2 {
		tokenString = arrayToken[1]
	}

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("invalid token")
		}

		return []byte(constants.SecretKey), nil
	})

	claim, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return ""
	}

	email := claim["email"].(string)

	return email
}
