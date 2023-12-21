// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error)
	UserCreate(ctx context.Context, input UserCreate) (idResult int64, err error)
	UserGetByPhone(ctx context.Context, phone string) (output *UserGet, err error)
	LoginAttemptCreate(ctx context.Context, input LoginAttemptCreate) (idResult int64, err error)
}
