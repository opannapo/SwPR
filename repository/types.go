// This file contains types that are used in the repository layer.
package repository

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
