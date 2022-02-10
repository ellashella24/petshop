package service

import (
	"github.com/jlaffaye/ftp"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"mime/multipart"
	"net/http"
	"petshop/constants"
	"petshop/delivery/common"
	"petshop/entity"
	"strings"
	"time"
)

func Upload(c echo.Context, userID string, file *multipart.FileHeader) (entity.Product, error) {
	var products entity.Product

	src, err := file.Open()
	if err != nil {
		return products, err
	}
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		return products, err
	}

	src.Seek(0, 0)

	contentType := http.DetectContentType(buffer)
	if contentType != "image/png" && contentType != "image/jpeg" {

		return products, c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	conn, err := ftp.Dial(constants.FTP_ADDR)
	err = conn.Login(constants.FTP_USERNAME, constants.FTP_PASSWORD)
	if err != nil {
		log.Fatal(err.Error())
	}
	time := time.Now().String()

	time = strings.ReplaceAll(time, " ", "")

	destinationFile := userID + "." + time + ".png"

	err = conn.Stor("./images/"+destinationFile, src)
	if err != nil {
		log.Fatal(err.Error())
	}

	products.ImageURL = "http://naufalhibatullah.com/images/" + destinationFile

	return products, nil
}
func Delete(file string) error {

	conn, err := ftp.Dial(constants.FTP_ADDR)
	err = conn.Login(constants.FTP_USERNAME, constants.FTP_PASSWORD)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = conn.Delete(file)

	if err != nil {
		return err
	}

	return nil
}

func DirReset() {
	conn, err := ftp.Dial(constants.FTP_ADDR)
	err = conn.Login(constants.FTP_USERNAME, constants.FTP_PASSWORD)
	if err != nil {
		log.Fatal(err.Error())
	}
	conn.RemoveDirRecur("images")
	conn.MakeDir("images")
}
