package controller

import (
	"fmt"
	"net/http"

	"github.com/imjap/internal/model"
	"github.com/imjap/internal/service"
	"github.com/labstack/echo/v4"
)

type ImageController struct {
	ImageService service.ImageService
}

func (ic *ImageController) UploadFile(c echo.Context) error {
	res, err := ic.ImageService.UploadFile(c)
	if !res && err != nil {
		errResp := model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, errResp)
	}

	return c.String(http.StatusOK, "uploaded success")
}

func (ic *ImageController) GetFiles(c echo.Context) error {
	files, err := ic.ImageService.GetFiles()
	if err != nil {
		errResp := model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, errResp)
	}
	return c.JSON(http.StatusOK, files)
}

func (ic *ImageController) GetFile(c echo.Context) error {
	filename := c.Param("name")
	data, mimetype, err := ic.ImageService.GetFile(filename)
	if err != nil {
		errResp := model.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, errResp)
	}

	c.Response().Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Response().Writer.Header().Set("Content-Type", mimetype)

	return c.Blob(http.StatusOK, mimetype, data)
}
