package middleware

import (
	"brickbetest/internal/standarderrors"
	"bytes"
	"errors"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"io"
	"net/http"
	"strings"
)

func OpenApiValidator(router routers.Router) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req, err := http.NewRequest(c.Method(), c.OriginalURL(), bytes.NewReader(c.Body()))
		if err != nil {
			log.Fatalf("Error while creating request: %s", err)
		}
		req.Header.Set("Content-Type", c.Get("Content-Type"))

		route, pathParams, err := router.FindRoute(req)

		if err != nil {
			if errors.Is(err, standarderrors.NotFound) {
				return c.Status(http.StatusNotFound).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParams,
			Route:      route,
			Options:    &openapi3filter.Options{MultiError: true},
		}

		err = openapi3filter.ValidateRequest(c.Context(), requestValidationInput)
		if err != nil {
			var validationError openapi3.MultiError
			if ok := errors.As(err, &validationError); ok {
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{
					"error": strings.Split(validationError.Error(), " | "),
				})
			}
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Call the next Fiber handler
		err = c.Next()
		if err != nil {
			return err
		}

		// Response body
		res := c.Response()
		bodyReader := bytes.NewReader(res.Body())

		headers := http.Header{}
		res.Header.VisitAll(func(key, value []byte) {
			headers.Add(string(key), string(value))
		})

		// Create http.Response
		httpRes := &http.Response{
			StatusCode:    res.StatusCode(),
			Proto:         "HTTP/1.1",
			ProtoMajor:    1,
			ProtoMinor:    1,
			Body:          io.NopCloser(bodyReader),
			ContentLength: int64(bodyReader.Len()),
			Header:        headers,
		}

		responseValidationInput := &openapi3filter.ResponseValidationInput{
			RequestValidationInput: requestValidationInput,
			Status:                 httpRes.StatusCode,
			Header:                 httpRes.Header,
			Body:                   httpRes.Body,
		}

		if err := openapi3filter.ValidateResponse(c.Context(), responseValidationInput); err != nil {
			var validationError openapi3.MultiError
			ok := errors.As(err, &validationError)
			if ok {
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{
					"error": strings.Split(validationError.Error(), " | "),
				})
			}
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return nil
	}
}
