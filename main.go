package main

import (
	"database/sql"

	"github.com/KMHTelU/KMH-WebProfile-API/configs"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/handlers"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/repositories"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/services"
	"github.com/KMHTelU/KMH-WebProfile-API/routes"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

var (
	config  *configs.Config
	db      *sql.DB
	queries *generated.Queries
	repo    *repositories.Repository
	service *services.Service
	cleaner *utils.TokenCleaner
)

func init() {
	conf, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	config = conf
	db = configs.ConnectDatabase(config.DatabaseURL)
	queries = generated.New(db)

	cleaner = utils.InitializeTokenCleaner(config.JWTSecret)

	repo = repositories.InitializeRepository(db, queries)
	service = services.InitializeService(repo, cleaner)
	handlers.InitializeHandler(service)
}

func main() {
	defer db.Close()
	app := fiber.New(fiber.Config{
		CaseSensitive:      true,
		StrictRouting:      true,
		EnableIPValidation: true,
		ServerHeader:       "KMH Tel-U",
		AppName:            "KMH Tel-U Profile Web API v" + config.Version,
	})

	routes.SetupRoutes(app)

	app.Hooks().OnPostShutdown(func(err error) error {
		if err != nil {
			log.Infof("Shutdown error: %v", err)
		} else {
			log.Info("Server shutdown completed successfully")
		}
		return nil
	})
	go app.Listen(":" + config.ServerPort)

	app.Shutdown()
}
