package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/viennn/docker-postgres-go/app/api/account"
	"github.com/viennn/docker-postgres-go/app/api/user"
	db "github.com/viennn/docker-postgres-go/app/db/sqlc"
	"github.com/viennn/docker-postgres-go/app/token"
	"github.com/viennn/docker-postgres-go/app/util"
)

const ApiPrefix = "/api/v1"

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
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

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", util.ValidCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	routerDefault := gin.Default()
	routerGroup := routerDefault.Group("/")
	user.Routes(routerGroup)
	account.Routes(routerGroup)

	server.router = routerDefault
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
