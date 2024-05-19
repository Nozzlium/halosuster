package repository

import (
	"bytes"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/nozzlium/halosuster/internal/constant"
	"github.com/nozzlium/halosuster/internal/model"
	"github.com/nozzlium/halosuster/internal/util"
)

type PatientRepository struct {
	db *pgx.Conn
}

func NewPatientRepository(
	db *pgx.Conn,
) *PatientRepository {
	return &PatientRepository{
		db: db,
	}
}

func (r *PatientRepository) Create(
	ctx context.Context,
	patient model.Patient,
) (model.Patient, error) {
	query := `
    insert into patients 
      (
        identity_number,
        user_id,
        phone_number,
        name,
        birthdate,
        gender,
        identity_card_image_url,
        created_at,
        updated_at
      )
    values 
      (
        $1, $2, $3, $4, $5, $6, $7, $8, $9
      )
  `
	_, err := r.db.Exec(
		ctx,
		query,
		patient.IdentityNumber,
		patient.UserID,
		patient.PhoneNumber,
		patient.Name,
		patient.Birthdate,
		patient.Gender,
		patient.IdentityScanImg,
		patient.CreatedAt,
		patient.UpdatedAt,
	)
	if err != nil {
		return model.Patient{}, err
	}

	return patient, nil
}

func (r *PatientRepository) FindById(
	ctx context.Context,
	id string,
) (model.Patient, error) {
	query := `
    select 
      identity_number,
      phone_number,
      name,
      birthdate,
      gender,
      created_at
    from patients
    where identity_number = $1 and
      deleted_at is null;
  `

	var patient model.Patient
	err := r.db.QueryRow(ctx, query, id).
		Scan(
			&patient.IdentityNumber,
			&patient.PhoneNumber,
			&patient.Name,
			&patient.Birthdate,
			&patient.Gender,
			&patient.CreatedAt,
		)
	if err != nil {
		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return patient, constant.ErrNotFound
		}
		return model.Patient{}, err
	}

	return patient, nil
}

func (r *PatientRepository) FindAll(
	ctx context.Context,
	queries model.PatientQuery,
) ([]model.Patient, error) {
	var query bytes.Buffer
	query.WriteString(`
    select 
      identity_number,
      phone_number,
      name,
      birthdate,
      gender,
      created_at
    from patients
    where 1 = 1
  `)
	queryString, params := util.BuildQueryStringAndParams(
		&query,
		queries.BuildWhereClauses,
		queries.BuildPagination,
		queries.BuildOrderByClause,
		true,
	)
	rows, err := r.db.Query(
		ctx,
		queryString,
		params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	patientData := make(
		[]model.Patient,
		0,
		queries.Limit,
	)
	for rows.Next() {
		var patient model.Patient
		err := rows.
			Scan(
				&patient.IdentityNumber,
				&patient.PhoneNumber,
				&patient.Name,
				&patient.Birthdate,
				&patient.Gender,
				&patient.CreatedAt,
			)
		if err != nil {
			return nil, err
		}

		patientData = append(
			patientData,
			patient,
		)
	}

	return patientData, nil
}
