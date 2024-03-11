package standarderrors

import "github.com/gofiber/fiber/v2"

var NotFound = fiber.ErrNotFound
var AlreadyExist = fiber.NewError(fiber.StatusConflict, "already exist")
var InsufficientBalance = fiber.NewError(fiber.StatusBadRequest, "insufficient balance")
var InvalidSignature = fiber.NewError(fiber.StatusBadRequest, "invalid signature")
