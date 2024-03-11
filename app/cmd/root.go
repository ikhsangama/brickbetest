package cmd

import (
	"brickbetest/config"
	"brickbetest/internal/postgres"
	"brickbetest/internal/sqs"
	"brickbetest/pkg/client/bank"
	"brickbetest/pkg/repository"
	"brickbetest/pkg/service"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"net/http"
)

type App struct {
	accountService  *service.AccountService
	transferService *service.TransferService
	bankClient      *bank.Client
}

func createServices(db *gorm.DB, publisher *sqs.Publisher) (
	*service.AccountService,
	*service.TransferService,
) {
	transferRepo := repository.NewTransferRepository(db)
	merchantRepo := repository.NewMerchantRepository(db)
	ledgerRepo := repository.NewLedgerRepository(db)
	bankClient := bank.NewClient(config.BankBaseUrl(), &http.Client{})
	accountService := service.NewAccountService(bankClient)
	transferService := service.NewTransferService(
		accountService,
		bankClient,
		transferRepo,
		merchantRepo,
		ledgerRepo,
		publisher,
	)
	return accountService, transferService
}

func initApp() *App {
	db, err := postgres.CreateConnection()
	if err != nil {
		log.Fatalf("db create connection failed: %s", err)
	}

	publisher := sqs.NewSqsPublisher(config.GetSqsUrl())
	accountService, transferService := createServices(db, publisher)

	return &App{
		accountService:  accountService,
		transferService: transferService,
	}
}
