package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/nozzlium/halosuster/internal/constant"
	"github.com/nozzlium/halosuster/internal/model"
	"github.com/nozzlium/halosuster/internal/service"
)

type PatientHandler struct {
	patientService *service.PatientService
}

func NewPatientHandler(
	patientService *service.PatientService,
) *PatientHandler {
	return &PatientHandler{
		patientService: patientService,
	}
}

func (h *PatientHandler) Create(
	ctx *fiber.Ctx,
) error {
	var body model.PatientRegisterBody
	err := ctx.BodyParser(&body)
	if err != nil {
		err = constant.ErrBadInput
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: "invalid body",
				detail: fmt.Sprintf(
					"patient register; failed to parse request body %v",
					err,
				),
			},
		)
	}

	patientModel, err := body.IsValid()
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: "invalid body",
				detail: fmt.Sprintf(
					"patient register; failed to parse request body %v",
					err,
				),
			},
		)
	}

	patientData, err := h.patientService.Create(
		ctx.Context(),
		patientModel,
	)
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: "invalid body",
				detail: fmt.Sprintf(
					"nurse register; failed to register %v",
					err,
				),
			},
		)
	}

	return ctx.Status(fiber.StatusCreated).
		JSON(
			fiber.Map{
				"message": "success",
				"data":    patientData,
			},
		)
}

func (h *PatientHandler) FindAll(
	ctx *fiber.Ctx,
) error {
	var queries model.PatientQuery
	ctx.QueryParser(&queries)
	queries.Offset = ctx.QueryInt(
		"offset",
		0,
	)
	queries.Limit = ctx.QueryInt(
		"limit",
		5,
	)

	data, err := h.patientService.FindAll(
		ctx.Context(),
		queries,
	)
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: err.Error(),
				detail: fmt.Sprintf(
					"find users; error finding patients: %v",
					err,
				),
			},
		)
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
		"data":    data,
	})
}
