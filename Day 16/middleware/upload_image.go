package middleware

import (
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func UploadFile(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		file, err := c.FormFile("uploadImage")

		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		src, err := file.Open()

		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		defer src.Close()

		tempFile, err := os.CreateTemp("upload", "image-*.png") // upload/image-3e10e160.png
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		defer tempFile.Close()

		if _, err = io.Copy(tempFile, src); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		data := tempFile.Name() // upload/image-3e10e160.png
		fileName := data[7:]    // image-3e10e160.png

		c.Set("dataFile", fileName)

		return next(c)
	}
}