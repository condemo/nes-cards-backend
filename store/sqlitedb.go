package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/condemo/nes-cards-backend/types"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

type SqliteStore struct {
	db *bun.DB
}

func NewSqliteStore() *SqliteStore {
	homeDir, err := os.UserHomeDir()
	logError(err)

	dataDir := path.Join(homeDir, ".local/share/nes-cards")
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		os.Mkdir(dataDir, os.FileMode(0o744))
	}

	dsn := fmt.Sprintf("file:%s/data.db?cache=shared", dataDir)

	sqldb, err := sql.Open(sqliteshim.ShimName, dsn)
	logError(err)

	db := bun.NewDB(sqldb, sqlitedialect.New())

	return &SqliteStore{db: db}
}

func (s *SqliteStore) Init() (*bun.DB, error) {
	// Game Table
	_, err := s.db.NewCreateTable().Model((*types.Game)(nil)).
		IfNotExists().Exec(context.Background())
	if err != nil {
		return nil, err
	}

	// Player Table
	_, err = s.db.NewCreateTable().Model((*types.Player)(nil)).
		IfNotExists().Exec(context.Background())
	if err != nil {
		return nil, err
	}

	// Tower Table
	_, err = s.db.NewCreateTable().Model((*types.Stats)(nil)).
		IfNotExists().Exec(context.Background())
	if err != nil {
		return nil, err
	}

	return s.db, nil
}

func logError(e error) {
	if e != nil {
		log.Fatalln("database error: ", e)
	}
}
