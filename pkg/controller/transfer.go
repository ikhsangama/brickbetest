package controller

import (
	"brickbetest/model"
	"brickbetest/pkg/service"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
)

type transferController struct {
	transferService service.TransferService
}

func NewTransferController(
	group fiber.Router,
	transferService *service.TransferService,
) {
	ctrl := &transferController{
		transferService: *transferService,
	}

	ctrl.registerEndpoints(group)
}

func (ctrl *transferController) registerEndpoints(group fiber.Router) {
	group.Post("/", ctrl.Transfer)
}

func (ctrl *transferController) Transfer(ctx *fiber.Ctx) (err error) {
	var transferReqBody model.TransferReqBody
	bodyBytes := ctx.Body()
	err = json.Unmarshal(bodyBytes, &transferReqBody)
	if err != nil {
		return fiber.ErrBadRequest
	}
	resBody, err := ctrl.transferService.Transfer(ctx.Context(), transferReqBody)
	if err != nil {
		return err
	}
	return ctx.JSON(resBody)
}
