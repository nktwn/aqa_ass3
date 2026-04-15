package user

import (
	"context"
	"database/sql"
	"diploma/modules/auth/model"
	"diploma/pkg/client/db"
	"errors"
	"time"

	converter "diploma/modules/auth/repository/user/converter"
	modelRepo "diploma/modules/auth/repository/user/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

const (
	tableName = "users"

	idColumn             = "id"
	nameColumn           = "name"
	phoneNumberColumn    = "phone_number"
	hashedPasswordColumn = "hashed_password"
	roleColumn           = "role"
	createdAtColumn      = "created_at"
	updatedAtColumn      = "updated_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) *repo {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, user *model.AuthUser) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, phoneNumberColumn, hashedPasswordColumn, roleColumn).
		Values(user.Info.Name, user.Info.PhoneNumber, user.HashedPassword, user.Info.Role).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) GetById(ctx context.Context, id int64) (*model.User, error) {
	builder := sq.Select(idColumn, nameColumn, phoneNumberColumn, roleColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{idColumn: id})

	query, args, err := builder.ToSql()

	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	var user modelRepo.User
	var userInfo modelRepo.UserInfo
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&user.ID, &userInfo.Name, &userInfo.PhoneNumber, &userInfo.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNoRows
		}
		return nil, err
	}
	user.Info = &userInfo

	return converter.ToUserFromRepo(&user), nil
}

func (r *repo) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*model.AuthUser, error) {
	builder := sq.Select(idColumn, nameColumn, phoneNumberColumn, roleColumn, hashedPasswordColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{phoneNumberColumn: phoneNumber})

	query, args, err := builder.ToSql()

	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	var user modelRepo.AuthUser
	var userInfo modelRepo.UserInfo
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&user.ID, &userInfo.Name, &userInfo.PhoneNumber, &userInfo.Role, &user.HashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNoRows
		}
		return nil, err
	}
	user.Info = &userInfo

	return converter.ToAuthUserFromRepo(&user), nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	builder := sq.Delete(tableName).
		Where(sq.Eq{idColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository.Delete",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) Update(ctx context.Context, id int64, info *model.UserInfo) error {
	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id}).
		Set(updatedAtColumn, time.Now())

	if info.Name != "" {
		builder = builder.Set(nameColumn, info.Name)
	}
	if info.PhoneNumber != "" {
		builder = builder.Set(phoneNumberColumn, info.PhoneNumber)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository.Update",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	return nil
}
