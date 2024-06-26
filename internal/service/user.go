package service

import (
	"context"
	"encoding/base64"
	"errors"
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

	userResponseBody, err := result.ToUserRegisterResponseBody()
	if err != nil {
		return model.UserRegisterResponseBody{}, err
	}
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
		savedUser,
	)
	if err != nil {
		return model.UserRegisterResponseBody{}, err
	}

	userResponseBody, err := savedUser.ToUserRegisterResponseBody()
	if err != nil {
		return model.UserRegisterResponseBody{}, err
	}
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

	usersDataCol := make(
		[]model.UserDataResponseBody,
		0,
		len(users),
	)
	for _, user := range users {
		userData, err := user.ToUserDataResponseBody()
		if err != nil {
			return nil, err
		}

		usersDataCol = append(
			usersDataCol,
			userData,
		)
	}

	return usersDataCol, nil
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
	employeeId := user.EmployeeID
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
	employeeId := ctx.Value("employeeId").(string)
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

	return result.ToNurseResponseBody()
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
		savedUser,
	)
	if err != nil {
		return model.UserRegisterResponseBody{}, err
	}

	userResponseBody, err := savedUser.ToUserRegisterResponseBody()
	if err != nil {
		return model.UserRegisterResponseBody{}, err
	}
	userResponseBody.AccessToken = accessToken

	return userResponseBody, nil
}

func (s *UserService) GrantNurseAccess(
	ctx context.Context,
	user model.User,
) error {
	employeeId := ctx.Value("employeeId").(string)
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

	err = util.ValidateGeneralEmployeeID(
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

func (s *UserService) UpdateNurse(
	ctx context.Context,
	user model.User,
) (model.User, error) {
	employeeId := ctx.Value("employeeId").(string)
	err := util.ValidateUserEmployeeID(
		employeeId,
	)
	if err != nil {
		return model.User{}, err
	}

	existingUser, err := s.userRepository.FindById(
		ctx,
		user.ID,
	)
	if err != nil {
		return model.User{}, err
	}

	err = util.ValidateIsANurse(
		existingUser.EmployeeID,
	)
	if err != nil {
		return model.User{}, err
	}

	_, err = s.userRepository.FindByEmployeeId(
		ctx,
		user.EmployeeID,
	)
	if err != nil {
		if !errors.Is(
			err,
			constant.ErrNotFound,
		) {
			return model.User{}, err
		}
	} else {
		return model.User{}, constant.ErrConflict
	}

	existingUser.EmployeeID = user.EmployeeID
	existingUser.Name = user.Name

	saved, err := s.userRepository.Edit(
		ctx,
		existingUser,
	)
	if err != nil {
		return model.User{}, err
	}

	return saved, nil
}

func (s *UserService) DeleteNurse(
	ctx context.Context,
	id uuid.UUID,
) (model.User, error) {
	employeeId := ctx.Value("employeeId").(string)
	err := util.ValidateUserEmployeeID(
		employeeId,
	)
	if err != nil {
		return model.User{}, err
	}

	existingUser, err := s.userRepository.FindById(
		ctx,
		id,
	)
	if err != nil {
		return model.User{}, err
	}

	err = util.ValidateIsANurse(
		existingUser.EmployeeID,
	)
	if err != nil {
		return model.User{}, err
	}

	_, err = s.userRepository.SetDeletedAt(
		ctx,
		model.User{
			ID:        id,
			DeletedAt: time.Now(),
		},
	)
	if err != nil {
		return model.User{}, err
	}

	return model.User{}, nil
}
