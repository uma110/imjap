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

	e.POST("/upload", uploadFile)

	e.Logger.Fatal(e.Start(":8011"))
}

// e.POST("/upload", uploadFile)
func uploadFile(c echo.Context) error {
	return file.UploadHandler(c)
}
