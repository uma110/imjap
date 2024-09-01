package main

import (
	"fmt"

	"github.com/imjap/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	}

	server.Serve()
}
