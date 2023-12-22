package handler

import (
	"swpr/repository"
	"swpr/util"
)

type Server struct {
	Repository   repository.RepositoryInterface
	PasswordUtil util.PasswordInterface
}

type NewServerOptions struct {
	Repository   repository.RepositoryInterface
	PasswordUtil util.PasswordInterface
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		Repository:   opts.Repository,
		PasswordUtil: opts.PasswordUtil,
	}
}
