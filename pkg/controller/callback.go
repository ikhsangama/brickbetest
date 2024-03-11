package controller

import (
	"brickbetest/internal/signature"
	"brickbetest/internal/standarderrors"
	"brickbetest/model"
	"brickbetest/pkg/service"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type callbackController struct {
	transferService service.TransferService
}

func NewCallbackController(
	group fiber.Router,
	transferService service.TransferService,
) {
	ctrl := &callbackController{
		transferService: transferService,
	}

	ctrl.registerEndpoints(group)
}

func (ctrl *callbackController) registerEndpoints(group fiber.Router) {
	group.Post("/transfer", ctrl.HandlePaymentCallback)
}

func (ctrl *callbackController) HandlePaymentCallback(ctx *fiber.Ctx) (err error) {
	sign := ctx.GetReqHeaders()["Signature"]
	valid := signature.Verify(sign)
	if !valid {
		return standarderrors.InvalidSignature
	}
	var callbackReqBody model.CallbackReqBody
	bodyBytes := ctx.Body()
	err = json.Unmarshal(bodyBytes, &callbackReqBody)
	if err != nil {
		return fiber.ErrBadRequest
	}
	resBody, err := ctrl.transferService.HandleTransferCallback(ctx.Context(), callbackReqBody)
	if err != nil {
		return err
	}
	return ctx.Status(http.StatusOK).JSON(resBody)
}
