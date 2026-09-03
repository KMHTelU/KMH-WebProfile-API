package main

import (
	"database/sql"
	"errors"
	"os"

	"github.com/KMHTelU/KMH-WebProfile-API/configs"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/generated"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/handlers"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/mailer"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/repositories"
	"github.com/KMHTelU/KMH-WebProfile-API/internal/services"
	"github.com/KMHTelU/KMH-WebProfile-API/routes"
	"github.com/KMHTelU/KMH-WebProfile-API/utils"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/yokeTH/gofiber-scalar/scalar/v3"
)

var (
	config  *configs.Config
	db      *sql.DB
	queries *generated.Queries
	repo    *repositories.Repository
	service *services.Service
	cleaner *utils.TokenCleaner
	handler *handlers.Handler
	route   *routes.Routes
)

func init() {
	conf, err := configs.LoadConfig()
	if err != nil {
		log.Infof("Failed to load configuration: %v", err)
	}
	config = conf
	db = configs.ConnectDatabase(config.DSN())
	queries = generated.New(db)

	cleaner = utils.InitializeTokenCleaner(config.JWTSecret, config.JWTRefreshSecret)
	cld, err := cloudinary.NewFromParams(config.CloudinaryCloudName, config.CloudinaryAPIKey, config.CloudinaryAPISecret)
	if err != nil {
		log.Infof("Failed to initialize Cloudinary: %v", err)
	}

	mailClient := mailer.New(mailer.Config{
		Host:       config.SMTPHost,
		Port:       config.SMTPPort,
		Username:   config.SMTPUsername,
		Password:   config.SMTPPassword,
		FromEmail:  config.SMTPFromEmail,
		FromName:   config.SMTPFromName,
		Encryption: config.SMTPEncryption,
		AppName:    config.AppName,
	})

	repo = repositories.InitializeRepository(db, queries)
	service = services.InitializeService(repo, cleaner, cld, mailClient, services.PasswordResetConfig{
		FrontendURL: config.FrontendURL,
		TokenTTL:    config.PasswordResetTokenTTL,
		MaxPerHour:  config.PasswordResetMaxPerHour,
	})
	handler = handlers.InitializeHandler(service)
	route = routes.InitializeRoutes(handler)
}

func main() {
	defer db.Close()
	app := fiber.New(fiber.Config{
		CaseSensitive:      true,
		StrictRouting:      true,
		EnableIPValidation: true,
		StructValidator: &utils.Validator{
			Validator: validator.New(),
		},
		ServerHeader: "KMH Tel-U",
		AppName:      "🔥 KMH Tel-U Profile Web API v" + config.Version,
		// Semua error yang lolos dari handler (termasuk *fiber.Error yang
		// dikembalikan mentah) diformat ke amplop JSON standar, sehingga
		// frontend selalu bisa menampilkan pesan yang bermakna — bukan
		// sekadar "request failed with status code 500".
		ErrorHandler: func(c fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := "Terjadi kesalahan pada server"
			var fiberErr *fiber.Error
			if errors.As(err, &fiberErr) {
				code = fiberErr.Code
				if fiberErr.Message != "" {
					message = fiberErr.Message
				}
			}
			return utils.RespondWithError(c, code, message)
		},
	})

	// Recover: cegah satu panic di handler menjatuhkan seluruh server.
	app.Use(recover.New())

	// CORS: dibutuhkan agar frontend (dev/prod) yang beda origin bisa memanggil API.
	// Auth memakai Bearer token (bukan cookie), jadi origin "*" aman.
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	route.SetupRoutes(app)

	swaggerBytes, err := os.ReadFile("./docs/swagger.json")
	if err != nil {
		log.Infof("Failed to read Swagger file: %v", err)
	}

	fileContentString := string(swaggerBytes)

	app.Get("/docs/*", scalar.New(scalar.Config{
		BasePath:          "/",
		FileContentString: fileContentString,
		Path:              "/docs",
		Title:             "KMH Tel-U Profile Web API Docs v" + config.Version,
		Theme:             scalar.ThemeKepler,
	}))

	log.Fatal(app.Listen(":" + config.ServerPort))

}
