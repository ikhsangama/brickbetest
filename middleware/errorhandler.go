package middleware

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"log"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
			_ = c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "An internal server error occurred",
			})
		}
	}()

	var e *fiber.Error
	if errors.As(err, &e) {
		log.Printf("A Fiber error occurred: %s", e.Error())
		return c.Status(e.Code).JSON(fiber.Map{
			"error": e.Message,
		})
	}

	log.Printf("An error occurred: %s", err.Error())
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "An internal server error occurred",
	})
}
