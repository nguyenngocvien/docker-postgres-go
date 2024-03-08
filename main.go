package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/viennn/docker-postgres-go/app/api"
	db "github.com/viennn/docker-postgres-go/app/db/sqlc"
	"github.com/viennn/docker-postgres-go/app/mail"
	"github.com/viennn/docker-postgres-go/app/util"
	"github.com/viennn/docker-postgres-go/app/worker"
	"golang.org/x/sync/errgroup"
)

var interruptSignal = []os.Signal{
	os.Interrupt,
    syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignal...)
	defer stop()

	connPool, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(connPool)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	// taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	waitGroup, ctx := errgroup.WithContext(ctx)

	runTaskProcessor(ctx, waitGroup, *config, redisOpt, store)
	runGinServer(*config, store)

	err = waitGroup.Wait()
	if err!= nil {
        log.Fatal().Err(err).Msg("error from wait group")
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

func runTaskProcessor(ctx context.Context, waitGroup *errgroup.Group,config util.Config, redisOpt asynq.RedisClientOpt, db db.Store) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, db, mailer)
	
	log.Info().Msg("start task processor")

	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start task processor")
	}

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown task processor")

		taskProcessor.Shutdown()
		log.Info().Msg("task processor is stopped")

		return nil
	})
}

func runGinServer(config util.Config, store db.Store){
	server, err := api.NewServer(config, store)
    if err!= nil {
        log.Fatal().Err(err).Msg("cannot create server")
    }

    err = server.Start(config.HTTPServerAddress)
    if err!= nil {
        log.Fatal().Err(err).Msg("cannot start server")
    }
}