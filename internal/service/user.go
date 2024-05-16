package service

import (
	"context"
	"encoding/base64"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nozzlium/halosuster/internal/constant"
	"github.com/nozzlium/halosuster/internal/model"
	"github.com/nozzlium/halosuster/internal/repository"
	"github.com/nozzlium/halosuster/internal/util"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	authRepository *repository.AuthRepository
	salt           int
	secret         string
}

func NewUserService(
	authRepository *repository.AuthRepository,
	salt int,
	secret string,
) *UserService {
	return &UserService{
		authRepository: authRepository,
		salt:           salt,
		secret:         secret,
	}
}

func (s *UserService) Register(
	ctx context.Context,
	user model.User,
) (model.UserResponseBody, error) {
	savedUser, err := s.authRepository.FindByEmployeeId(
		ctx,
		user.EmployeeID,
	)
	if err != nil {
		if !errors.Is(
			err,
			constant.ErrNotFound,
		) {
			return model.UserResponseBody{}, err
		}
	}

	if savedUser.EmployeeID == user.EmployeeID {
		return model.UserResponseBody{}, constant.ErrConflict
	}

	id, err := uuid.NewV7()
	if err != nil {
		return model.UserResponseBody{}, err
	}
	currentTime := util.Now()
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		s.salt,
	)
	if err != nil {
		return model.UserResponseBody{}, err
	}

	user.ID = id
	user.Password = string(
		hashedPassword,
	)
	user.CreatedAt = currentTime
	user.UpdatedAt = currentTime

	result, err := s.authRepository.Save(
		ctx,
		user,
	)
	if err != nil {
		return model.UserResponseBody{}, err
	}

	accessToken, err := generateJwtToken(
		s.secret,
		result,
	)
	if err != nil {
		return model.UserResponseBody{}, err
	}

	userResponseBody := result.ToUserResponseBody()
	userResponseBody.AccessToken = accessToken

	return userResponseBody, nil
}

func (s *UserService) Login(
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
		return model.UserResponseBody{}, constant.ErrBadInput
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

func generateJwtToken(
	secret string,
	user model.User,
) (string, error) {
	token := jwt.New(
		jwt.SigningMethodHS256,
	)

	claims := token.Claims.(jwt.MapClaims)
	userID := base64.RawStdEncoding.EncodeToString(
		[]byte(user.ID.String()),
	)
	employeeId := strconv.FormatUint(
		user.EmployeeID,
		16,
	)
	log.Println(employeeId)
	claims["si"] = userID
	claims["ut"] = employeeId
	claims["exp"] = time.Now().
		Add(time.Hour * 72).
		Unix()

	t, err := token.SignedString(
		[]byte(secret),
	)
	if err != nil {
		return "", err
	}

	return t, nil
}

func (s *UserService) RegisterNurse(
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

func (s *UserService) LoginNurse(
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
		return model.UserResponseBody{}, constant.ErrBadInput
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

func (s *UserService) GrantNurseAccess(
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
