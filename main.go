package main

import (
	"database/sql"

	"github.com/golang-migrate/migrate"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"github.com/viennn/docker-postgres-go/api"
	db "github.com/viennn/docker-postgres-go/db/sqlc"
	"github.com/viennn/docker-postgres-go/mail"
	"github.com/viennn/docker-postgres-go/util"
	"github.com/viennn/docker-postgres-go/worker"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:123456@localhost:5432/simplebank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load config")
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	go runTaskProcessor(redisOpt, store, config)

	server, err := api.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}
}

func runDBMigration(migrateURL string, dbSource string) {
	migration, err := migrate.New(migrateURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migrate instance")
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up")
	}

	log.Info().Msg("db migrated successfully")
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, db db.Store, config util.Config) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, db, mailer)
	log.Info().Msg("Running task processor")

	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start task processor")
	}
}
