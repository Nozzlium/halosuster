package model

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/nozzlium/halosuster/internal/constant"
	"github.com/nozzlium/halosuster/internal/util"
)

type OrderBy string

const (
	Asc  OrderBy = "asc"
	Desc OrderBy = "desc"
)

func (o OrderBy) IsValid() bool {
	switch o {
	case Asc, Desc:
		return true
	default:
		return false
	}
}

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

func (u *User) ToUserDataResponseBody() (UserDataResponseBody, error) {
	employeeIDInt, err := strconv.ParseUint(
		u.EmployeeID,
		10,
		64,
	)
	if err != nil {
		return UserDataResponseBody{}, err
	}
	return UserDataResponseBody{
		UserID: u.ID.String(),
		NIP:    employeeIDInt,
		Name:   u.Name,
		CreatedAt: util.ToISO8601(
			u.CreatedAt,
		),
	}, nil
}

func (u *User) ToUserRegisterResponseBody() (UserRegisterResponseBody, error) {
	employeeIDInt, err := strconv.ParseUint(
		u.EmployeeID,
		10,
		64,
	)
	if err != nil {
		return UserRegisterResponseBody{}, err
	}

	return UserRegisterResponseBody{
		UserID: u.ID.String(),
		NIP:    employeeIDInt,
		Name:   u.Name,
	}, nil
}

func (u *User) ToNurseResponseBody() (NurseRegisterResponseBody, error) {
	employeeIDInt, err := strconv.ParseUint(
		u.EmployeeID,
		10,
		64,
	)
	if err != nil {
		return NurseRegisterResponseBody{}, err
	}

	return NurseRegisterResponseBody{
		UserID: u.ID.String(),
		NIP:    employeeIDInt,
		Name:   u.Name,
	}, nil
}

type UserRegisterRequestBody struct {
	NIP      uint64 `json:"nip"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (body *UserRegisterRequestBody) IsValid() (User, error) {
	var user User
	employeeIdString := strconv.FormatUint(
		body.NIP,
		10,
	)
	err := util.ValidateUserEmployeeID(
		employeeIdString,
	)
	if err != nil {
		return user, err
	}
	user.EmployeeID = employeeIdString

	if nameLen := len(body.Name); nameLen < 5 ||
		nameLen > 50 {
		return user, constant.ErrBadInput
	}
	user.Name = body.Name

	if passwordLen := len(body.Password); passwordLen < 5 ||
		passwordLen > 33 {
		return user, constant.ErrBadInput
	}
	user.Password = body.Password

	return user, nil
}

type UserRegisterResponseBody struct {
	UserID      string `json:"userId"`
	NIP         uint64 `json:"nip"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

type UserLoginBody struct {
	NIP      uint64 `json:"nip"`
	Password string `json:"password"`
}

func (body *UserLoginBody) IsValid() (User, error) {
	var user User
	employeeIdString := strconv.FormatUint(
		body.NIP,
		10,
	)
	err := util.ValidateUserEmployeeID(
		employeeIdString,
	)
	if err != nil {
		return user, err
	}
	user.EmployeeID = employeeIdString

	if passLen := len(body.Password); passLen < 5 ||
		passLen > 33 {
		return user, constant.ErrBadInput
	}
	user.Password = body.Password

	return user, nil
}

type NurseLoginBody struct {
	NIP      uint64 `json:"nip"`
	Password string `json:"password"`
}

func (body *NurseLoginBody) IsValid() (User, error) {
	var user User
	employeeIdString := strconv.FormatUint(
		body.NIP,
		10,
	)
	err := util.ValidateGeneralEmployeeID(
		employeeIdString,
	)
	if err != nil {
		return user, err
	}
	user.EmployeeID = employeeIdString

	if passLen := len(body.Password); passLen < 5 ||
		passLen > 33 {
		return user, constant.ErrBadInput
	}
	user.Password = body.Password

	return user, nil
}

type NurseRegisterRequestBody struct {
	NIP                 uint64 `json:"nip"`
	Name                string `json:"name"`
	IdentityCardScanImg string `json:"identityCardScanImg"`
}

func (body *NurseRegisterRequestBody) IsValid() (User, error) {
	var user User
	employeeIdString := strconv.FormatUint(
		body.NIP,
		10,
	)
	err := util.ValidateGeneralEmployeeID(
		employeeIdString,
	)
	if err != nil {
		return user, err
	}
	user.EmployeeID = employeeIdString

	if nameLen := len(body.Name); nameLen < 5 ||
		nameLen > 50 {
		return user, constant.ErrBadInput
	}
	user.Name = body.Name

	if !util.ValidateURL(
		body.IdentityCardScanImg,
	) {
		return user, constant.ErrBadInput
	}
	user.IdentityCardImageURL = body.IdentityCardScanImg

	return user, nil
}

type NurseGiveAccessRequestBody struct {
	Password string `json:"password"`
}

func (body *NurseGiveAccessRequestBody) IsValid() bool {
	if passLen := len(body.Password); passLen < 5 ||
		passLen > 33 {
		return false
	}

	return true
}

type NurseEditRequestBody struct {
	NIP  uint64 `json:"nip"`
	Name string `json:"name"`
}

func (body *NurseEditRequestBody) IsValid() (User, error) {
	var user User
	employeeIdString := strconv.FormatUint(
		body.NIP,
		10,
	)
	err := util.ValidateGeneralEmployeeID(
		employeeIdString,
	)
	if err != nil {
		return user, err
	}
	user.EmployeeID = employeeIdString

	if nameLen := len(body.Name); nameLen < 5 ||
		nameLen > 50 {
		return user, constant.ErrBadInput
	}
	user.Name = body.Name

	return user, nil
}

type NurseRegisterResponseBody struct {
	UserID string `json:"userId"`
	NIP    uint64 `json:"nip"`
	Name   string `json:"name"`
}

type SearchUserQuery struct {
	UserID    string  `query:"userId"`
	Name      string  `query:"name"`
	NIP       uint64  `query:"nip"`
	Role      string  `query:"role"`
	CreatedAt OrderBy `query:"createdAt"`
	Offset    int
	Limit     int
}

func (q *SearchUserQuery) BuildWhereClauses() ([]string, []interface{}) {
	clauses := make([]string, 0, 4)
	params := make([]interface{}, 0, 4)

	if q.UserID != "" {
		clauses = append(
			clauses,
			"id = $%d",
		)
		params = append(
			params,
			q.UserID,
		)
	}

	if q.Name != "" {
		clauses = append(
			clauses,
			`name ilike '%%' || $%d || '%%'`,
		)
		params = append(params, q.Name)
	}

	employeeIdString := strconv.FormatUint(
		q.NIP,
		10,
	)
	if employeeIdString != "0" &&
		employeeIdString != "" {
		clauses = append(
			clauses,
			`employee_id like $%d || '%%'`,
		)
		params = append(
			params,
			employeeIdString,
		)
	}

	switch q.Role {
	case "it":
		clauses = append(
			clauses,
			`employee_id like $%d || '%%'`,
		)

		params = append(
			params,
			"615",
		)

	case "nurse":
		clauses = append(
			clauses,
			"employee_id like $%d || '%%'",
		)
		params = append(
			params,
			"303",
		)
	}

	return clauses, params
}

func (q SearchUserQuery) BuildPagination() (string, []interface{}) {
	var params []interface{}

	limit := 5
	offset := 0
	if q.Limit > 0 {
		limit = q.Limit
	}
	if q.Offset > 0 {
		offset = q.Offset
	}
	params = append(
		params,
		limit,
		offset,
	)

	return "limit $%d offset $%d", params
}

func (q SearchUserQuery) BuildOrderByClause() []string {
	var sqlClause []string

	if q.CreatedAt != "" ||
		OrderBy(
			q.CreatedAt,
		).IsValid() {
		sqlClause = append(
			sqlClause,
			fmt.Sprintf(
				"created_at %s",
				q.CreatedAt,
			),
		)
	} else {
		sqlClause = append(
			sqlClause,
			"created_at desc",
		)
	}

	return sqlClause
}

type UserDataResponseBody struct {
	UserID    string `json:"userId"`
	NIP       uint64 `json:"nip"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}
