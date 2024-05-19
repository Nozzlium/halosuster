package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nozzlium/halosuster/internal/constant"
)

func HandleError(
	ctx *fiber.Ctx,
	err ErrorResponse,
) error {
	switch err.error {
	case constant.ErrNotFound:
		return ctx.Status(fiber.StatusNotFound).
			JSON(fiber.Map{
				"message": err.message,
			})
	case constant.ErrConflict:
		return ctx.Status(fiber.StatusConflict).
			JSON(fiber.Map{
				"message": err.message,
			})
	case constant.ErrUnauthorized:
		return ctx.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{
				"message": err.message,
			})
	case constant.ErrBadInput,
		constant.ErrInvalidBody,
		constant.ErrInsufficientFund,
		constant.ErrInvalidChange,
		constant.ErrInsufficientStock:
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"message": err.message,
			})
	default:
		log.Printf(
			"internal error: %v",
			err.error,
		)
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{
				"message": "internal server error",
			})
	}
}

type ErrorResponse struct {
	error   error
	message string
	detail  string
}

func (e ErrorResponse) Error() string {
	return e.error.Error()
}
