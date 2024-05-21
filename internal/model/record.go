package model

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/nozzlium/halosuster/internal/constant"
	"github.com/nozzlium/halosuster/internal/util"
)

type Record struct {
	ID             uuid.UUID
	IdentityNumber string
	UserID         uuid.UUID
	Symptomps      string
	Medications    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
}

type RecordRegisterBody struct {
	IdentityNumber uint64 `json:"identityNumber"`
	Symptomps      string `json:"symptoms"`
	Medications    string `json:"medications"`
}

func (body *RecordRegisterBody) IsValid() (Record, error) {
	var record Record
	identityNumberString := strconv.FormatUint(
		body.IdentityNumber,
		10,
	)
	err := util.ValidateIdentityNumber(
		identityNumberString,
	)
	if err != nil {
		return record, err
	}
	record.IdentityNumber = identityNumberString

	if sympLen := len(body.Symptomps); sympLen < 1 ||
		sympLen > 2000 {
		return record, constant.ErrBadInput
	}
	record.Symptomps = body.Symptomps

	if medLen := len(body.Medications); medLen < 1 ||
		medLen > 2000 {
		return record, constant.ErrBadInput
	}
	record.Medications = body.Medications

	return record, nil
}

type RecordPatientBody struct {
	IdentityNumber      uint64 `json:"identityNumber"`
	PhoneNumber         string `json:"phoneNumber"`
	Name                string `json:"name"`
	Birthdate           string `json:"birthdate"`
	Gender              string `json:"gender"`
	IdentityCardScanImg string `json:"identityCardScanImg"`
}

type RecordUserBody struct {
	NIP    uint64 `json:"nip"`
	Name   string `json:"name"`
	UserID string `json:"userId"`
}

type RecordResponseBody struct {
	Symptomps      string            `json:"symptomps"`
	Medications    string            `json:"medications"`
	CreatedAt      string            `json:"createdAt"`
	IdentityDetail RecordPatientBody `json:"identityDetail"`
	CreatedBy      RecordUserBody    `json:"createdBy"`
}
