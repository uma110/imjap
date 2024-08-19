package main

import (
	"fmt"

	"github.com/imjap/pkg/file"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Response struct {
	Status  int
	Message string
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", getFiles)
	e.GET("/:name", getFile)
	e.POST("/upload", uploadFile)

	e.Logger.Fatal(e.Start(":8011"))
}

// e.GET("/", getFiles)
func getFiles(c echo.Context) error {
	return file.GetFilesHandler(c)
}

// e.GET("/file", getFile)
func getFile(c echo.Context) error {
	return file.GetFileHandler(c)
}

// e.POST("/upload", uploadFile)
func uploadFile(c echo.Context) error {
	return file.UploadHandler(c)
}
