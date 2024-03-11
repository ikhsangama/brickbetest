package cmd

import (
	"brickbetest/middleware"
	"brickbetest/pkg/controller"
	"brickbetest/pkg/service"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/legacy"
	"github.com/gofiber/fiber/v2"
	"log"
)

func StartRestServer() {
	f := fiber.New(fiber.Config{ErrorHandler: middleware.ErrorHandler})

	spec := LoadOpenApiSpec("docs/openapi.yaml")
	router := CreateRouter(spec)
	f.Use(middleware.OpenApiValidator(router))

	app := initApp()

	SetupRoutes(f, app.transferService, app.accountService)

	StartServer(f)
}

func LoadOpenApiSpec(filepath string) *openapi3.T {
	spec, err := openapi3.NewLoader().LoadFromFile(filepath)
	if err != nil {
		panic(fmt.Sprintf("failed to load openapi %v", err))
	}
	return spec
}

func CreateRouter(spec *openapi3.T) routers.Router {
	router, err := legacy.NewRouter(spec)
	if err != nil {
		log.Fatalf("Error while creating router: %s", err)
	}
	return router
}

func SetupRoutes(app *fiber.App, transferService *service.TransferService, accountService *service.AccountService) {
	account := app.Group("/v1/account")
	controller.NewAccountController(account, accountService)
	transfer := app.Group("/v1/transfer")
	controller.NewTransferController(transfer, transferService)
	callback := app.Group("/v1/callback")
	controller.NewCallbackController(callback, *transferService)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
}

func StartServer(app *fiber.App) {
	err := app.Listen(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
