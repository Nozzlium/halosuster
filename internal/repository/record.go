package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nozzlium/halosuster/internal/constant"
	"github.com/nozzlium/halosuster/internal/model"
)

type RecordRepository struct {
	db *pgx.Conn
}

func NewRecordRepository(
	db *pgx.Conn,
) *RecordRepository {
	return &RecordRepository{
		db: db,
	}
}

func (r *RecordRepository) Create(
	ctx context.Context,
	record model.Record,
) (model.Record, error) {
	query := `
    insert into records
    (
      id,
      identity_number,
      user_id,
      symptomps,
      medications,
      created_at,
      updated_at
    ) values (
      $1, $2, $3, $4, $5, $6, $7
    )
  `
	_, err := r.db.Exec(ctx, query,
		record.ID,
		record.IdentityNumber,
		record.UserID,
		record.Symptomps,
		record.Medications,
		record.CreatedAt,
		record.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return model.Record{}, constant.ErrNotFound
			}
		}
		return model.Record{}, err
	}

	return record, nil
}
