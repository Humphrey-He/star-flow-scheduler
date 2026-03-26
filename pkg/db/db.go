package db

import (
	"context"
	"database/sql"
	"fmt"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/Humphrey-He/star-flow-scheduler/pkg/ent"
	"github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	Client *ent.Client
	SQL    *sql.DB
}

func Open(dsn string) (*DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("postgres dsn is empty")
	}

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	drv := entsql.OpenDB(dialect.Postgres, sqlDB)
	client := ent.NewClient(ent.Driver(drv))

	return &DB{Client: client, SQL: sqlDB}, nil
}

func Close(db *DB) error {
	if db == nil {
		return nil
	}
	if err := db.Client.Close(); err != nil {
		return err
	}
	return db.SQL.Close()
}

func Transact(ctx context.Context, client *ent.Client, fn func(ctx context.Context, tx *ent.Tx) error) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	if err := fn(ctx, tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func init() {
	stdlib.RegisterDefaultPgxTypeMap()
}
