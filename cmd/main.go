package main

import (
	"log"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/bootstrap"
)

func main() {
	app := bootstrap.InitApp()
	log.Println("Server running on port 8008")
	app.Listen(":8008")
}