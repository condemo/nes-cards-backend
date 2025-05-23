package store

import (
	"context"
	"database/sql"
	"os"

	"github.com/condemo/nes-cards-backend/types"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type PostgresqlStore struct {
	db *bun.DB
}

func NewPostgresqlStore() *PostgresqlStore {
	dsn := os.Getenv("DB_DSN")

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	return &PostgresqlStore{db: db}
}

func (s *PostgresqlStore) Init() (*bun.DB, error) {
	ctx := context.Background()
	// User Table
	_, err := s.db.NewCreateTable().Model((*types.User)(nil)).
		IfNotExists().Exec(ctx)

	// Game Table
	_, err = s.db.NewCreateTable().Model((*types.Game)(nil)).
		IfNotExists().Exec(ctx)
	if err != nil {
		return nil, err
	}

	// Player Table
	_, err = s.db.NewCreateTable().Model((*types.Player)(nil)).
		IfNotExists().Exec(ctx)
	if err != nil {
		return nil, err
	}

	// Tower Table
	_, err = s.db.NewCreateTable().Model((*types.Stats)(nil)).
		IfNotExists().Exec(ctx)
	if err != nil {
		return nil, err
	}
	return s.db, nil
}
