package pg

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/novando/go-cinema/pkg/logger"
	"github.com/spf13/viper"
)

type PGTX struct {
	db pgx.Tx
}

// BeginTx start transaction mode
func (q *PG) BeginTx() (pgx.Tx, error) {
	if viper.GetBool("db.pg.logging") {
		logger.Call().Infof("Transaction start")
	}
	return q.db.BeginTx(context.Background(), pgx.TxOptions{})
}

// Exec execute transaction without returning any rows
func (q *PGTX) Exec(sql string, arg ...any) (pgconn.CommandTag, error) {
	if viper.GetBool("db.pg.logging") {
		logger.Call().Infof("Exec: %v Arguments: %v", sql, arg)
	}
	return q.db.Exec(context.Background(), sql, arg)
}

// Query execute transaction, returning single row or an error
func (q *PGTX) Query(sql string, arg ...any) (pgx.Rows, error) {
	if viper.GetBool("db.pg.logging") {
		logger.Call().Infof("Query): %v Arguments: %v", sql, arg)
	}
	return q.db.Query(context.Background(), sql, arg)
}

// QueryRow execute transaction, returning 0 or multiple rows
func (q *PGTX) QueryRow(sql string, arg ...any) pgx.Row {
	if viper.GetBool("db.pg.logging") {
		logger.Call().Infof("QueryRow): %v Arguments: %v", sql, arg)
	}
	return q.db.QueryRow(context.Background(), sql, arg)
}

// Rollback cancel the transaction
func (q *PGTX) Rollback() error {
	if viper.GetBool("db.pg.logging") {
		logger.Call().Infof("Rollingback transaction")
	}
	return q.db.Rollback(context.Background())
}

// Commit proceed the transaction
func (q *PGTX) Commit() error {
	if viper.GetBool("db.pg.logging") {
		logger.Call().Infof("Commiting transaction")
	}
	return q.db.Commit(context.Background())
}
