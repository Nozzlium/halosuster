package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nozzlium/halosuster/internal/constant"
	"github.com/nozzlium/halosuster/internal/model"
	"github.com/nozzlium/halosuster/internal/service"
)

type NurseHandler struct {
	nurseService *service.NurseService
}

func NewNurseHandler(
	nurseService *service.NurseService,
) *NurseHandler {
	return &NurseHandler{
		nurseService: nurseService,
	}
}

func (h *NurseHandler) Register(
	ctx *fiber.Ctx,
) error {
	var body model.NurseRegisterRequestBody
	err := ctx.BodyParser(&body)
	if err != nil {
		err = constant.ErrBadInput
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: "invalid body",
				detail: fmt.Sprintf(
					"nurse register; failed to parse request body %v",
					err,
				),
			},
		)
	}

	if !body.IsValid() {
		err = constant.ErrBadInput
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: "invalid body",
				detail: fmt.Sprintf(
					"nurse register; invalid body: %v",
					err,
				),
			},
		)
	}

	fmt.Println(
		"ini lhooo masbroooo",
		ctx.Locals("employeeId"),
	)
	data, err := h.nurseService.Register(
		ctx.Context(),
		model.User{
			EmployeeID:           body.NIP,
			Name:                 body.Name,
			IdentityCardImageURL: body.IdentityCardScanImg,
		},
	)
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: "invalid body",
				detail: fmt.Sprintf(
					"nurse register; failed to parse request body %v",
					err,
				),
			},
		)
	}

	return ctx.Status(fiber.StatusCreated).
		JSON(fiber.Map{
			"message": "success",
			"data":    data,
		})
}

func (h *NurseHandler) Login(
	ctx *fiber.Ctx,
) error {
	var body model.NurseLoginBody
	err := ctx.BodyParser(&body)
	if err != nil {
		err = constant.ErrBadInput
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: "invalid body",
				detail: fmt.Sprintf(
					"nurse login; failed to parse request body %v",
					err,
				),
			},
		)
	}

	if !body.IsValid() {
		err := constant.ErrBadInput
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: "invalid body",
				detail: fmt.Sprintf(
					"nurse login; invalid body: %v",
					err,
				),
			},
		)
	}

	data, err := h.nurseService.Login(
		ctx.UserContext(),
		model.User{
			EmployeeID: body.NIP,
			Password:   body.Password,
		},
	)
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: "failed to login",
				detail: fmt.Sprintf(
					"nurse login; failed to login %v",
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

func (h *NurseHandler) GiveAccess(
	ctx *fiber.Ctx,
) error {
	userIdString := ctx.Params("userId")
	var body model.NurseGiveAccessRequestBody
	err := ctx.BodyParser(&body)
	if err != nil ||
		userIdString == "" {
		err = constant.ErrBadInput
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: "invalid body",
				detail: fmt.Sprintf(
					"nurse register; failed to parse request body %v",
					err,
				),
			},
		)
	}

	userId, err := uuid.ParseBytes(
		[]byte(userIdString),
	)
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: "error parsing user ID",
				detail: fmt.Sprintf(
					"nurse access; failed to parse userID  %v",
					err,
				),
			},
		)
	}

	if !body.IsValid() {
		err = constant.ErrBadInput
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: "invalid body",
				detail: fmt.Sprintf(
					"nurse register; invalid body: %v",
					err,
				),
			},
		)
	}

	fmt.Println(
		"ini lhooo masbroooo",
		ctx.Locals("employeeId"),
	)
	err = h.nurseService.GiveAccess(
		ctx.Context(),
		model.User{
			ID:       userId,
			Password: body.Password,
		},
	)
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: "invalid body",
				detail: fmt.Sprintf(
					"nurse register; failed to parse request body %v",
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
