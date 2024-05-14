package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/nozzlium/halosuster/internal/model"
	"github.com/nozzlium/halosuster/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(
	userService *service.UserService,
) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Register(
	ctx *fiber.Ctx,
) error {
	var body model.UserRequestBody
	err := ctx.BodyParser(&body)
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: "invalid body",
				detail: fmt.Sprintf(
					"user register; failed to parse request body %v",
					err,
				),
			},
		)
	}

	if !body.IsValid() {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: "invalid body",
				detail: fmt.Sprintf(
					"user register; invalid body: %v",
					err,
				),
			},
		)
	}

	data, err := h.userService.Register(
		ctx.UserContext(),
		model.User{
			EmployeeID: body.NIP,
			Name:       body.Name,
			Password:   body.Password,
		},
	)
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: "invalid body",
				detail: fmt.Sprintf(
					"user register; failed to parse request body %v",
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

func (h *UserHandler) Login(
	ctx *fiber.Ctx,
) error {
	var body model.UserLoginBody
	err := ctx.BodyParser(&body)
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: "invalid body",
				detail: fmt.Sprintf(
					"user login; failed to parse request body %v",
					err,
				),
			},
		)
	}

	data, err := h.userService.Login(
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
					"user login; failed to login %v",
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
