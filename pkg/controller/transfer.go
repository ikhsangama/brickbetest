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
	group.Post("/", ctrl.TransferRequest)
	group.Get("/:id", ctrl.GetTransfer)
}

func (ctrl *transferController) TransferRequest(ctx *fiber.Ctx) (err error) {
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

func (ctrl *transferController) GetTransfer(ctx *fiber.Ctx) error {
	transferId := ctx.Params("id")
	resBody, err := ctrl.transferService.GetTransfer(ctx.Context(), transferId)
	if err != nil {
		return err
	}
	return ctx.JSON(resBody)
}
