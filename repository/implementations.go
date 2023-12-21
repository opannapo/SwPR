package repository

import (
	"context"
	"log"
)

func (r *Repository) GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT name FROM test WHERE id = $1", input.Id).Scan(&output.Name)
	if err != nil {
		return
	}
	return
}

func (r *Repository) UserCreate(ctx context.Context, input UserCreate) (idResult int64, err error) {
	sql := `
		INSERT INTO users (full_name, password, phone)
		VALUES ($1, $2, $3)
		RETURNING id`

	err = r.Db.QueryRowContext(ctx, sql, input.FullName, input.Password, input.Phone).Scan(&idResult)
	if err != nil {
		log.Println("error ", err)
	}

	return
}

func (r *Repository) UserGetByPhone(ctx context.Context, phone string) (result *UserGet, err error) {
	result = &UserGet{}
	err = r.Db.QueryRowContext(ctx, "SELECT * FROM users WHERE phone = $1", phone).
		Scan(
			&result.Id,
			&result.FullName,
			&result.Password,
			&result.Phone,
			&result.CreatedAt,
			&result.UpdatedAt,
		)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (r *Repository) LoginAttemptCreate(ctx context.Context, input LoginAttemptCreate) (idResult int64, err error) {
	sql := `
		INSERT INTO login_attempt (user_id, status)
		VALUES ($1, $2)
		RETURNING id`

	err = r.Db.QueryRowContext(ctx, sql, input.UserID, input.Status).Scan(&idResult)
	if err != nil {
		log.Println("error ", err)
	}

	return
}
