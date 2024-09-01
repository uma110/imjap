package server

import (
	"github.com/imjap/internal/controller"
	"github.com/imjap/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Serve() {
	imageService := service.ImageService{}
	imageController := controller.ImageController{ImageService: imageService}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", imageController.GetFiles)
	e.GET("/:name", imageController.GetFile)
	e.POST("/upload", imageController.UploadFile)

	e.Logger.Fatal(e.Start(":8011"))
}
