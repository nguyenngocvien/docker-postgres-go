package gapi

import (
	"fmt"

	db "github.com/viennn/docker-postgres-go/db/sqlc"
	"github.com/viennn/docker-postgres-go/token"
	"github.com/viennn/docker-postgres-go/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTmaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}