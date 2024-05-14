package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/nozzlium/halosuster/internal/util"
)

type User struct {
	ID                   uuid.UUID
	EmployeeID           string
	Name                 string
	Password             string
	IdentityCardImageURL string
	CreatedBy            uuid.UUID
	UpdatedBy            uuid.UUID
	DeletedBy            uuid.UUID
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            time.Time
}

func (u *User) ToUserResponseBody() UserResponseBody {
	return UserResponseBody{
		UserID: u.ID.String(),
		NIP:    u.EmployeeID,
		Name:   u.Name,
	}
}

type UserRequestBody struct {
	NIP      string `json:"nip"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (body *UserRequestBody) IsValid() bool {
	if !util.ValidateUserEmployeeID(
		body.NIP,
	) {
		return false
	}

	if nameLen := len(body.Name); nameLen < 5 ||
		nameLen > 50 {
		return false
	}

	if passwordLen := len(body.Password); passwordLen < 5 ||
		passwordLen > 33 {
		return false
	}

	return true
}

type UserResponseBody struct {
	UserID      string `json:"userId"`
	NIP         string `json:"nip"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

type UserLoginBody struct {
	NIP      string `json:"nip"`
	Password string `json:"password"`
}
