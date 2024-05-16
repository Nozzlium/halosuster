package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nozzlium/halosuster/internal/constant"
	"github.com/nozzlium/halosuster/internal/model"
	"github.com/nozzlium/halosuster/internal/repository"
	"github.com/nozzlium/halosuster/internal/util"
	"golang.org/x/crypto/bcrypt"
)

type NurseService struct {
	authRepository *repository.AuthRepository
	salt           int
	secret         string
}

func NewNurseService(
	authRepository *repository.AuthRepository,
) *NurseService {
	return &NurseService{
		authRepository: authRepository,
	}
}

func (s *NurseService) Register(
	ctx context.Context,
	user model.User,
) (model.NurseRegisterResponseBody, error) {
	employeeId := ctx.Value("employeeId").(uint64)
	if !util.ValidateUserEmployeeID(
		employeeId,
	) {
		return model.NurseRegisterResponseBody{}, constant.ErrUnauthorized
	}

	savedNurse, err := s.authRepository.FindByEmployeeId(
		ctx,
		user.EmployeeID,
	)
	if err != nil {
		if !errors.Is(
			err,
			constant.ErrNotFound,
		) {
			return model.NurseRegisterResponseBody{}, err
		}
	}

	if savedNurse.EmployeeID == user.EmployeeID {
		return model.NurseRegisterResponseBody{}, constant.ErrConflict
	}

	id, err := uuid.NewV7()
	if err != nil {
		return model.NurseRegisterResponseBody{}, err
	}
	currentTime := util.Now()

	user.ID = id
	user.CreatedAt = currentTime
	user.UpdatedAt = currentTime

	result, err := s.authRepository.Save(
		ctx,
		user,
	)
	if err != nil {
		return model.NurseRegisterResponseBody{}, err
	}

	return result.ToNurseResponseBody(), nil
}

func (s *NurseService) Login(
	ctx context.Context,
	user model.User,
) (model.UserResponseBody, error) {
	savedUser, err := s.authRepository.FindByEmployeeId(
		ctx,
		user.EmployeeID,
	)
	if err != nil {
		return model.UserResponseBody{}, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(savedUser.Password),
		[]byte(user.Password),
	)
	if err != nil {
		return model.UserResponseBody{}, err
	}

	accessToken, err := generateJwtToken(
		s.secret,
		user,
	)
	if err != nil {
		return model.UserResponseBody{}, err
	}

	userResponseBody := savedUser.ToUserResponseBody()
	userResponseBody.AccessToken = accessToken

	return userResponseBody, nil
}

func (s *NurseService) GiveAccess(
	ctx context.Context,
	user model.User,
) error {
	employeeId := ctx.Value("employeeId").(uint64)
	if !util.ValidateUserEmployeeID(
		employeeId,
	) {
		return constant.ErrUnauthorized
	}

	savedNurse, err := s.authRepository.FindById(
		ctx,
		user.ID,
	)
	if err != nil {
		return err
	}

	if !util.ValidateNurseEmployeeID(
		savedNurse.EmployeeID,
	) {
		return constant.ErrBadInput
	}

	hashedPassBytes, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		s.salt,
	)
	if err != nil {
		return err
	}

	_, err = s.authRepository.EditPassword(
		ctx,
		model.User{
			ID: user.ID,
			Password: string(
				hashedPassBytes,
			),
		},
	)
	if err != nil {
		return err
	}

	return nil
}
