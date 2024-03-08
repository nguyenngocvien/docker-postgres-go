package db

import "github.com/jackc/pgx/v5/pgxpool"

type Store interface {
	Querier
}

type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

var StoreDB Store

func NewStore(connPool *pgxpool.Pool) Store {
	StoreDB = &SQLStore{
		connPool: connPool,
        Queries: New(connPool),
	}
	return StoreDB
}
