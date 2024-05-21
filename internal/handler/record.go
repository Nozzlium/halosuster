package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/nozzlium/halosuster/internal/constant"
	"github.com/nozzlium/halosuster/internal/model"
	"github.com/nozzlium/halosuster/internal/service"
)

type RecordHandler struct {
	recordService *service.RecordService
}

func NewRecordHandler(
	recordService *service.RecordService,
) *RecordHandler {
	return &RecordHandler{
		recordService: recordService,
	}
}

func (h *RecordHandler) Create(
	ctx *fiber.Ctx,
) error {
	var body model.RecordRegisterBody
	err := ctx.BodyParser(&body)
	if err != nil {
		err = constant.ErrBadInput
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: "invalid body",
				detail: fmt.Sprintf(
					"medical record register; failed to parse request body %v",
					err,
				),
			},
		)
	}

	recordModel, err := body.IsValid()
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: "invalid body",
				detail: fmt.Sprintf(
					"record register; failed to parse request body %v",
					err,
				),
			},
		)
	}

	_, err = h.recordService.Create(
		ctx.Context(),
		recordModel,
	)
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: "invalid body",
				detail: fmt.Sprintf(
					"record register; failed to parse request body %v",
					err,
				),
			},
		)
	}

	return ctx.Status(fiber.StatusCreated).
		JSON(fiber.Map{
			"message": "success",
		})
}
