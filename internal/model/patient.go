package model

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/nozzlium/halosuster/internal/constant"
	"github.com/nozzlium/halosuster/internal/util"
)

type Patient struct {
	IdentityNumber  string
	UserID          uuid.UUID
	PhoneNumber     string
	Name            string
	Birthdate       time.Time
	Gender          string
	IdentityScanImg string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       time.Time
}

type PatientRegisterBody struct {
	IdentityNumber      uint64 `json:"identityNumber"`
	PhoneNumber         string `json:"phoneNumber"`
	Name                string `json:"name"`
	Birthdate           string `json:"birthdate"`
	Gender              string `json:"gender"`
	IdentityCardScanImg string `json:"identityCardScanImg"`
}

func (body *PatientRegisterBody) IsValid() (Patient, error) {
	var patient Patient
	identityNumberString := strconv.FormatUint(
		body.IdentityNumber,
		10,
	)
	err := util.ValidateIdentityNumber(
		identityNumberString,
	)
	if err != nil {
		return patient, err
	}
	patient.IdentityNumber = identityNumberString

	err = util.ValidatePhoneNumber(
		body.PhoneNumber,
	)
	if err != nil {
		return patient, err
	}
	patient.PhoneNumber = body.PhoneNumber

	if nameLen := len(body.Name); nameLen < 3 ||
		nameLen > 30 {
		return patient, constant.ErrBadInput
	}
	patient.Name = body.Name

	if body.Birthdate == "" {
		return patient, constant.ErrBadInput
	}
	birthdate, err := time.Parse(
		"2006-01-02T15:04:05.999Z",
		body.Birthdate,
	)
	if err != nil {
		return patient, constant.ErrBadInput
	}
	patient.Birthdate = birthdate

	if body.Gender != "male" &&
		body.Gender != "female" {
		return patient, constant.ErrBadInput
	}
	patient.Gender = body.Gender

	if !util.ValidateURL(
		body.IdentityCardScanImg,
	) {
		return patient, constant.ErrBadInput
	}
	patient.IdentityScanImg = body.IdentityCardScanImg

	return patient, nil
}

type PatientResponseBody struct {
	IdentityNumber uint64 `json:"identityNumber"`
	PhoneNumber    string `json:"phoneNumber"`
	Name           string `json:"name"`
	Birthdate      string `json:"birthDate"`
	Gender         string `json:"gender"`
	CreatedAt      string `json:"createdAt"`
}

func (patient *Patient) ToResponseBody() (PatientResponseBody, error) {
	var patientBody PatientResponseBody
	userIdUint, err := strconv.ParseUint(
		patient.IdentityNumber,
		10,
		64,
	)
	if err != nil {
		return patientBody, err
	}

	return PatientResponseBody{
		IdentityNumber: userIdUint,
		PhoneNumber:    patient.PhoneNumber,
		Name:           patient.Name,
		Birthdate: util.ToISO8601(
			patient.Birthdate,
		),
		Gender: patient.Gender,
		CreatedAt: util.ToISO8601(
			patient.CreatedAt,
		),
	}, nil
}

type PatientQuery struct {
	IdentityNumber string `query:"identityNumber"`
	Name           string `query:"name"`
	PhoneNumber    string `query:"phoneNumber"`
	CreatedAt      string `query:"createdAt"`
	Offset         int
	Limit          int
}

func (q *PatientQuery) BuildWhereClauses() ([]string, []interface{}) {
	clauses := make([]string, 0, 4)
	params := make([]interface{}, 0, 4)

	if q.IdentityNumber != "" {
		clauses = append(
			clauses,
			"identity_number = $%d",
		)
		params = append(
			params,
			q.IdentityNumber,
		)
	}

	if q.Name != "" {
		clauses = append(
			clauses,
			`name ilike '%%' || $%d || '%%'`,
		)
		params = append(params, q.Name)
	}

	if q.PhoneNumber != "" {
		phoneNumberString := "+" + q.PhoneNumber
		clauses = append(
			clauses,
			`phone_number like $%d || '%%'`,
		)
		params = append(
			params,
			phoneNumberString,
		)
	}

	return clauses, params
}

func (q *PatientQuery) BuildPagination() (string, []interface{}) {
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

func (q *PatientQuery) BuildOrderByClause() []string {
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
