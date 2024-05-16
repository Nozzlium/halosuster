package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nozzlium/halosuster/internal/constant"
	"github.com/nozzlium/halosuster/internal/model"
)

type AuthRepository struct {
	db *pgx.Conn
}

func NewAuthRepository(
	db *pgx.Conn,
) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) Save(
	ctx context.Context,
	user model.User,
) (model.User, error) {
	query := `
    insert into users
    (
      id,
      employee_id,
      name,
      password,
      identity_card_image_url,
      created_at,
      updated_at
    ) values (
      $1,
      $2,
      $3,
      $4,
      $5,
      $6,
      $7
    )
  `
	_, err := r.db.Exec(
		ctx,
		query,
		user.ID,
		user.EmployeeID,
		user.Name,
		user.Password,
		user.IdentityCardImageURL,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *AuthRepository) FindById(
	ctx context.Context,
	id uuid.UUID,
) (model.User, error) {
	query := `
    select
      id,
      name,
      employee_id,
      password
    from users
    where
      id = $1;
  `

	var user model.User
	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(&user.ID, &user.Name, &user.EmployeeID, &user.Password)
	if err != nil {
		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return model.User{}, constant.ErrNotFound
		}
		return model.User{}, err
	}

	return user, nil
}

func (r *AuthRepository) FindByEmployeeId(
	ctx context.Context,
	employeeId uint64,
) (model.User, error) {
	query := `
    select
      id,
      name,
      employee_id,
      password
    from users
    where
      employee_id = $1;
  `

	var user model.User
	err := r.db.QueryRow(
		ctx,
		query,
		employeeId,
	).Scan(&user.ID, &user.Name, &user.EmployeeID, &user.Password)
	if err != nil {
		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return model.User{}, constant.ErrNotFound
		}
		return model.User{}, err
	}

	return user, nil
}

func (r *AuthRepository) EditPassword(
	ctx context.Context,
	user model.User,
) (model.User, error) {
	query := `
    update users
    set password = $1
    where id = $2
  `
	_, err := r.db.Exec(
		ctx,
		query,
		user.Password,
		user.ID,
	)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
