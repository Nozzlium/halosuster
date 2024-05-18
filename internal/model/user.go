package model

import (
	"fmt"
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
	EmployeeID           uint64
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

func (u *User) ToUserDataResponseBody() UserDataResponseBody {
	return UserDataResponseBody{
		UserID: u.ID.String(),
		NIP:    u.EmployeeID,
		Name:   u.Name,
		CreatedAt: util.ToISO8601(
			u.CreatedAt,
		),
	}
}

func (u *User) ToUserRegisterResponseBody() UserRegisterResponseBody {
	return UserRegisterResponseBody{
		UserID: u.ID.String(),
		NIP:    u.EmployeeID,
		Name:   u.Name,
	}
}

func (u *User) ToNurseResponseBody() NurseRegisterResponseBody {
	return NurseRegisterResponseBody{
		UserID: u.ID.String(),
		NIP:    u.EmployeeID,
		Name:   u.Name,
	}
}

type UserRegisterRequestBody struct {
	NIP      uint64 `json:"nip"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (body *UserRegisterRequestBody) IsValid() error {
	err := util.ValidateUserEmployeeID(
		body.NIP,
	)
	if err != nil {
		return err
	}

	if nameLen := len(body.Name); nameLen < 5 ||
		nameLen > 50 {
		return constant.ErrBadInput
	}

	if passwordLen := len(body.Password); passwordLen < 5 ||
		passwordLen > 33 {
		return constant.ErrBadInput
	}

	return nil
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

func (body *UserLoginBody) IsValid() error {
	err := util.ValidateUserEmployeeID(
		body.NIP,
	)
	if err != nil {
		return err
	}

	if passLen := len(body.Password); passLen < 5 ||
		passLen > 33 {
		return constant.ErrBadInput
	}

	return nil
}

type NurseLoginBody struct {
	NIP      uint64 `json:"nip"`
	Password string `json:"password"`
}

func (body *NurseLoginBody) IsValid() error {
	err := util.ValidateNurseEmployeeID(
		body.NIP,
	)
	if err != nil {
		return err
	}

	if passLen := len(body.Password); passLen < 5 ||
		passLen > 33 {
		return constant.ErrBadInput
	}

	return nil
}

type NurseRegisterRequestBody struct {
	NIP                 uint64 `json:"nip"`
	Name                string `json:"name"`
	IdentityCardScanImg string `json:"identityCardScanImg"`
}

func (body *NurseRegisterRequestBody) IsValid() error {
	err := util.ValidateNurseEmployeeID(
		body.NIP,
	)
	if err != nil {
		return err
	}

	if nameLen := len(body.Name); nameLen < 5 ||
		nameLen > 50 {
		return constant.ErrBadInput
	}

	if !util.ValidateURL(
		body.IdentityCardScanImg,
	) {
		return constant.ErrBadInput
	}

	return nil
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
			"name ilike %%$%d%%",
		)
		params = append(params, q.Name)
	}

	if q.NIP > 0 {
		clauses = append(
			clauses,
			"employee_id >= $%d",
			"employee_id <= $%d",
		)
		lower, upper := formEmployeeIdWildcard(
			q.NIP,
		)
		params = append(
			params,
			lower,
			upper,
		)
	}

	switch q.Role {
	case "it":
		clauses = append(
			clauses,
			"employee_id >= $%d",
			"employee_id <= $%d",
		)
		lower, upper := formEmployeeIdWildcard(
			615,
		)
		params = append(
			params,
			lower,
			upper,
		)

	case "nurse":
		clauses = append(
			clauses,
			"employee_id >= $%d",
			"employee_id <= $%d",
		)
		lower, upper := formEmployeeIdWildcard(
			303,
		)
		params = append(
			params,
			lower,
			upper,
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

func formEmployeeIdWildcard(
	employeeId uint64,
) (uint64, uint64) {
	var minIT uint64 = 6150000000000
	var minNurse uint64 = 3030000000000

	upper := 0
	for employeeId < minIT || employeeId < minNurse {
		employeeId *= 10
		upper = (upper * 10) + 9
		if employeeId >= minNurse {
			return employeeId, (employeeId + uint64(upper))
		}
	}

	return employeeId, (employeeId + uint64(upper))
}

type UserDataResponseBody struct {
	UserID    string `json:"userId"`
	NIP       uint64 `json:"nip"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}
