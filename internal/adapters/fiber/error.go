package fiber_adapter

import (
	"github.com/gofiber/fiber/v2"
	"github.com/x0k/skillrock-tasks-service/internal/shared"
)

func BadRequest(err error) error {
	return &fiber.Error{
		Code:    fiber.ErrBadRequest.Code,
		Message: err.Error(),
	}
}

func ServiceError(err *shared.ServiceError) error {
	if err.Expected {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: err.Msg,
		}
	} else {
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Msg,
		}
	}
}

func SpecificServiceError(err *shared.ServiceError, code int) error {
	return &fiber.Error{
		Code:    code,
		Message: err.Msg,
	}
}
