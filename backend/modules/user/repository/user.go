package repository

import (
	"context"
	"time"

	model "diploma/modules/user/model"
	"diploma/pkg/client/db"

	sq "github.com/Masterminds/squirrel"
)

const (
	uTableName         = "users"
	uIdColumn          = "id"
	uNameColumn        = "name"
	uPhoneNumberColumn = "phone_number"
)

func (r *repo) UserByID(ctx context.Context, id int64) (model.User, error) {
	builder := sq.Select(uIdColumn, uNameColumn, uPhoneNumberColumn, "role").
		PlaceholderFormat(sq.Dollar).
		From(uTableName).
		Where(sq.Eq{uIdColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return model.User{}, err
	}

	q := db.Query{
		Name:     "user_repository.GetUserByID",
		QueryRaw: query,
	}

	var user model.User
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Role)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *repo) UpdateUser(ctx context.Context, user model.User) error {
	builder := sq.Update(uTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{uIdColumn: user.ID}).
		Set("updated_at", time.Now())

	if user.Name != "" {
		builder = builder.Set(uNameColumn, user.Name)
	}
	if user.PhoneNumber != "" {
		builder = builder.Set(uPhoneNumberColumn, user.PhoneNumber)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository.UpdateUser",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return err
}

// func (r *repo) DeleteUser(ctx context.Context, id int64) error {
// 	builder := sq.Delete(uTableName).
// 		PlaceholderFormat(sq.Dollar).
// 		Where(sq.Eq{uIdColumn: id})

// 	query, args, err := builder.ToSql()
// 	if err != nil {
// 		return err
// 	}

// 	q := db.Query{
// 		Name:     "user_repository.DeleteUser",
// 		QueryRaw: query,
// 	}

// 	_, err = r.db.DB().ExecContext(ctx, q, args...)
// 	return err
// }
