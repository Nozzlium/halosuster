package repository

import (
	"bytes"
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nozzlium/halosuster/internal/constant"
	"github.com/nozzlium/halosuster/internal/model"
	"github.com/nozzlium/halosuster/internal/util"
)

type UserRepository struct {
	db *pgx.Conn
}

func NewUserRepository(
	db *pgx.Conn,
) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Save(
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

func (r *UserRepository) FindById(
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
      id = $1 and
      deleted_at is null;
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

func (r *UserRepository) FindAll(
	ctx context.Context,
	searchQuery model.SearchUserQuery,
) ([]model.User, error) {
	var query bytes.Buffer
	query.WriteString(`
    select 
      id,
      employee_id,
      name,
      created_at
      from users
    where 1 = 1
  `)
	queryString, params := util.BuildQueryStringAndParams(
		&query,
		searchQuery.BuildWhereClauses,
		searchQuery.BuildPagination,
		searchQuery.BuildOrderByClause,
		true,
	)
	log.Println(
		"query hasil",
		queryString,
		params,
	)

	rows, err := r.db.Query(
		ctx,
		queryString,
		params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make(
		[]model.User,
		0,
		searchQuery.Limit,
	)
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID,
			&user.EmployeeID,
			&user.Name,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) FindByEmployeeId(
	ctx context.Context,
	employeeId string,
) (model.User, error) {
	query := `
    select
      id,
      name,
      employee_id,
      password
    from users
    where
      employee_id = $1 and
      deleted_at is null;
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

func (r *UserRepository) EditPassword(
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

func (r *UserRepository) Edit(
	ctx context.Context,
	user model.User,
) (model.User, error) {
	query := `
    update users
    set employee_id = $1, name = $2
    where id = $3
  `
	_, err := r.db.Exec(
		ctx,
		query,
		user.EmployeeID,
		user.Name,
		user.ID,
	)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *UserRepository) SetDeletedAt(
	ctx context.Context,
	user model.User,
) (model.User, error) {
	query := `
    update users
    set deleted_at = $1
    where id = $2
  `
	_, err := r.db.Exec(
		ctx,
		query,
		user.DeletedAt,
		user.ID,
	)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
