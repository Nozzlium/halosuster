package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nozzlium/halosuster/internal/constant"
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
	var body model.UserRegisterRequestBody
	err := ctx.BodyParser(&body)
	if err != nil {
		err = constant.ErrBadInput
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

	userModel, err := body.IsValid()
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: err.Error(),
				detail: fmt.Sprintf(
					"user register; invalid body: %v",
					err,
				),
			},
		)
	}

	data, err := h.userService.Register(
		ctx.UserContext(),
		userModel,
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
		err = constant.ErrBadInput
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

	userModel, err := body.IsValid()
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: err.Error(),
				detail: fmt.Sprintf(
					"user login; invalid body: %v",
					err,
				),
			},
		)
	}

	data, err := h.userService.Login(
		ctx.UserContext(),
		userModel,
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

func (h *UserHandler) RegisterNurse(
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

	userModel, err := body.IsValid()
	if err != nil {
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

	data, err := h.userService.RegisterNurse(
		ctx.Context(),
		userModel,
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

func (h *UserHandler) LoginNurse(
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

	userModel, err := body.IsValid()
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: err.Error(),
				detail: fmt.Sprintf(
					"nurse login; invalid body: %v",
					err,
				),
			},
		)
	}

	data, err := h.userService.LoginNurse(
		ctx.UserContext(),
		userModel,
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

func (h *UserHandler) GrantNurseAccess(
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

	err = h.userService.GrantNurseAccess(
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

	return ctx.JSON(fiber.Map{
		"message": "success",
	})
}

func (h *UserHandler) FindAll(
	ctx *fiber.Ctx,
) error {
	var queries model.SearchUserQuery
	ctx.QueryParser(&queries)
	queries.Offset = ctx.QueryInt(
		"offset",
		0,
	)
	queries.Limit = ctx.QueryInt(
		"limit",
		5,
	)

	data, err := h.userService.FindAll(
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
					"find users; error finding users: %v",
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

func (h *UserHandler) Update(
	ctx *fiber.Ctx,
) error {
	userIdString := ctx.Params("userId")
	var body model.NurseEditRequestBody
	err := ctx.BodyParser(&body)
	if err != nil {
		err = constant.ErrBadInput
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: "invalid body",
				detail: fmt.Sprintf(
					"nurse edit; failed to parse request body %v",
					err,
				),
			},
		)
	}

	user, err := body.IsValid()
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: "invalid body",
				detail: fmt.Sprintf(
					"nurse edit; invalid body: %v",
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
	user.ID = userId

	_, err = h.userService.UpdateNurse(
		ctx.Context(),
		user,
	)
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: "failed to edit",
				detail: fmt.Sprintf(
					"nurse edit; failed to edit %v",
					err,
				),
			},
		)
	}

	return ctx.JSON(fiber.Map{
		"message": "success",
	})
}

func (h *UserHandler) Delete(
	ctx *fiber.Ctx,
) error {
	userIdString := ctx.Params("userId")

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
					"nurse delete; failed to parse userID  %v",
					err,
				),
			},
		)
	}

	_, err = h.userService.DeleteNurse(
		ctx.Context(),
		userId,
	)
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: "failed to edit",
				detail: fmt.Sprintf(
					"nurse delete; failed to edit %v",
					err,
				),
			},
		)
	}

	return ctx.JSON(
		fiber.Map{"message": "success"},
	)
}
