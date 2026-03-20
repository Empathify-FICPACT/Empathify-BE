package main

import (
	"log"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/bootstrap"
	"github.com/Empathify-FICPACT/Empathify-BE/config"
)

func main() {
	config.Load()

	db := bootstrap.NewDB()
	defer db.Close()

	app := bootstrap.NewApp(db)

	log.Printf("Server running on port %s", config.App.AppPort)
	if err := app.Listen(":" + config.App.AppPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}