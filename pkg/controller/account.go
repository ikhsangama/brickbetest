package controller

import (
	"brickbetest/pkg/service"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type accountController struct {
	accountService service.AccountService
}

func NewAccountController(
	group fiber.Router,
	accountService *service.AccountService,
) {
	ctrl := &accountController{
		accountService: *accountService,
	}

	ctrl.registerEndpoints(group)
}

func (ctrl *accountController) registerEndpoints(group fiber.Router) {
	group.Get("/validate", ctrl.AccountValidation)
}

func (ctrl *accountController) AccountValidation(ctx *fiber.Ctx) (err error) {
	bankCode := ctx.Query("bank_code")
	accountNumber := ctx.Query("account_number")
	resBody, err := ctrl.accountService.Validate(ctx.Context(), bankCode, accountNumber)
	if err != nil {
		return err
	}
	return ctx.Status(http.StatusOK).JSON(resBody)
}
