// This file contains types that are used in the repository layer.
package repository

import "database/sql"

type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	Name string
}

type UserCreate struct {
	FullName string
	Password string
	Phone    string
}

type UserGet struct {
	Id        int64
	FullName  string
	Password  string
	Phone     string
	CreatedAt sql.NullString
	UpdatedAt sql.NullString
}

type LoginAttemptCreate struct {
	UserID int64
	Status bool
}
