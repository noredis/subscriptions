package httpext

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type FiberError struct {
	Error  string       `json:"error"`
	Fields []FieldError `json:"fields,omitempty"`
}

type FieldError struct {
	Field       string `json:"field"`
	Description string `json:"description"`
}

func Error(c *fiber.Ctx, code int, err string) error {
	return c.Status(code).JSON(FiberError{Error: err})
}

func ValidationError(c *fiber.Ctx, vErrs validator.ValidationErrors) error {
	fieldErrors := make([]FieldError, 0)

	for _, fErr := range vErrs {
		fieldErrors = append(fieldErrors, FieldError{
			Field:       fErr.Field(),
			Description: mapTagToMessage(fErr),
		})
	}

	return c.Status(http.StatusUnprocessableEntity).JSON(FiberError{
		Error:  "validation error",
		Fields: fieldErrors,
	})
}

func mapTagToMessage(fErr validator.FieldError) string {
	switch fErr.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fErr.Field())
	case "gte":
		return fmt.Sprintf("%s should be more than %s", fErr.Field(), fErr.Param())
	case "uuid":
		return fmt.Sprintf("%s should be uuid", fErr.Field())
	case "date_format":
		return fmt.Sprintf("%s must match 'MM-YYYY'", fErr.Field())
	default:
		return fErr.Tag()
	}
}
