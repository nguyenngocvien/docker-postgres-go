package main

import (
	"database/sql"
	"log"

	"github.com/viennn/docker-postgres-go/api"
	db "github.com/viennn/docker-postgres-go/db/sqlc"
	"github.com/viennn/docker-postgres-go/util"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:123456@localhost:5432/simplebank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	config, err := util.LoadConfig(".")
	if err!= nil {
        log.Fatal("Cannot load config: ", err)
    }
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err!= nil {
        log.Fatal("cannot create server: ", err)
    }

	err = server.Start(config.ServerAddress)
	if err!= nil {
        log.Fatal("cannot start server: ", err)
    }
}
