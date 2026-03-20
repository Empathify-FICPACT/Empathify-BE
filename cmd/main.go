// @title           Empathify API
// @version         1.0
// @description     API untuk platform terapi dan pelatihan keterampilan sosial
// @termsOfService  http://swagger.io/terms/

// @contact.name   Empathify Team
// @contact.email  empathify@gmail.com

// @host      empathify-be-staging.fly.dev
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan token dengan format: Bearer {token}

package main

import (
	"log"

	_ "github.com/Empathify-FICPACT/Empathify-BE/docs"

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