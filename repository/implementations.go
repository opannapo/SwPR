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

func (r *Repository) UserGet(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error) {
	//TODO implement me
	panic("implement me")
}
