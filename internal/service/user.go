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
	userRepository *repository.UserRepository
	salt           int
	secret         string
}

func NewUserService(
	userRepository *repository.UserRepository,
	salt int,
	secret string,
) *UserService {
	return &UserService{
		userRepository: userRepository,
		salt:           salt,
		secret:         secret,
	}
}

func (s *UserService) Register(
	ctx context.Context,
	user model.User,
) (model.UserRegisterResponseBody, error) {
	savedUser, err := s.userRepository.FindByEmployeeId(
		ctx,
		user.EmployeeID,
	)
	if err != nil {
		if !errors.Is(
			err,
			constant.ErrNotFound,
		) {
			return model.UserRegisterResponseBody{}, err
		}
	}

	if savedUser.EmployeeID == user.EmployeeID {
		return model.UserRegisterResponseBody{}, constant.ErrConflict
	}

	id, err := uuid.NewV7()
	if err != nil {
		return model.UserRegisterResponseBody{}, err
	}
	currentTime := util.Now()
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		s.salt,
	)
	if err != nil {
		return model.UserRegisterResponseBody{}, err
	}

	user.ID = id
	user.Password = string(
		hashedPassword,
	)
	user.CreatedAt = currentTime
	user.UpdatedAt = currentTime

	result, err := s.userRepository.Save(
		ctx,
		user,
	)
	if err != nil {
		return model.UserRegisterResponseBody{}, err
	}

	accessToken, err := generateJwtToken(
		s.secret,
		result,
	)
	if err != nil {
		return model.UserRegisterResponseBody{}, err
	}

	userResponseBody := result.ToUserRegisterResponseBody()
	userResponseBody.AccessToken = accessToken

	return userResponseBody, nil
}

func (s *UserService) Login(
	ctx context.Context,
	user model.User,
) (model.UserRegisterResponseBody, error) {
	savedUser, err := s.userRepository.FindByEmployeeId(
		ctx,
		user.EmployeeID,
	)
	if err != nil {
		return model.UserRegisterResponseBody{}, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(savedUser.Password),
		[]byte(user.Password),
	)
	if err != nil {
		return model.UserRegisterResponseBody{}, constant.ErrBadInput
	}

	accessToken, err := generateJwtToken(
		s.secret,
		user,
	)
	if err != nil {
		return model.UserRegisterResponseBody{}, err
	}

	userResponseBody := savedUser.ToUserRegisterResponseBody()
	userResponseBody.AccessToken = accessToken

	return userResponseBody, nil
}

func (s *UserService) FindAll(
	ctx context.Context,
	queries model.SearchUserQuery,
) ([]model.UserDataResponseBody, error) {
	users, err := s.userRepository.FindAll(
		ctx,
		queries,
	)
	if err != nil {
		return nil, err
	}

	usersData := make(
		[]model.UserDataResponseBody,
		0,
		len(users),
	)
	for _, user := range users {
		usersData = append(
			usersData,
			user.ToUserDataResponseBody(),
		)
	}

	return usersData, nil
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
	err := util.ValidateUserEmployeeID(
		employeeId,
	)
	if err != nil {
		return model.NurseRegisterResponseBody{}, constant.ErrUnauthorized
	}

	savedNurse, err := s.userRepository.FindByEmployeeId(
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

	result, err := s.userRepository.Save(
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
) (model.UserRegisterResponseBody, error) {
	savedUser, err := s.userRepository.FindByEmployeeId(
		ctx,
		user.EmployeeID,
	)
	if err != nil {
		return model.UserRegisterResponseBody{}, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(savedUser.Password),
		[]byte(user.Password),
	)
	if err != nil {
		return model.UserRegisterResponseBody{}, constant.ErrBadInput
	}

	accessToken, err := generateJwtToken(
		s.secret,
		user,
	)
	if err != nil {
		return model.UserRegisterResponseBody{}, err
	}

	userResponseBody := savedUser.ToUserRegisterResponseBody()
	userResponseBody.AccessToken = accessToken

	return userResponseBody, nil
}

func (s *UserService) GrantNurseAccess(
	ctx context.Context,
	user model.User,
) error {
	employeeId := ctx.Value("employeeId").(uint64)
	err := util.ValidateUserEmployeeID(
		employeeId,
	)
	if err != nil {
		return constant.ErrUnauthorized
	}

	savedNurse, err := s.userRepository.FindById(
		ctx,
		user.ID,
	)
	if err != nil {
		return err
	}

	err = util.ValidateNurseEmployeeID(
		savedNurse.EmployeeID,
	)
	if err != nil {
		return constant.ErrBadInput
	}

	hashedPassBytes, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		s.salt,
	)
	if err != nil {
		return err
	}

	_, err = s.userRepository.EditPassword(
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

func (s *UserService) Update(
	ctx context.Context,
	user model.User,
) (model.User, error) {
	saved, err := s.userRepository.FindByEmployeeId(
		ctx,
		user.EmployeeID,
	)
	if err != nil {
		if !errors.Is(
			err,
			constant.ErrNotFound,
		) {
			return user, err
		}
	}

	if saved.EmployeeID == user.EmployeeID {
		return user, constant.ErrConflict
	}

	edited, err := s.userRepository.Edit(
		ctx,
		user,
	)
	if err != nil {
		return model.User{}, err
	}

	return edited, nil
}
