package service

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
	employeeId := base64.RawStdEncoding.EncodeToString(
		[]byte(user.EmployeeID),
	)
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
