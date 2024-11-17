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
func (q *PG) BeginTx(c context.Context) (PGTX, error) {
	if viper.GetBool("db.pg.logging") {
		logger.Call().Infof("Transaction start")
	}
	tx, err := q.db.BeginTx(c, pgx.TxOptions{})
	return PGTX{tx}, err
}

// Exec execute transaction without returning any rows
func (q *PGTX) Exec(c context.Context, sql string, arg ...any) (pgconn.CommandTag, error) {
	if viper.GetBool("db.pg.logging") {
		logger.Call().Infof("Exec: %v Arguments: %v", sql, arg)
	}
	return q.db.Exec(c, sql, arg...)
}

// Query execute transaction, returning single row or an error
func (q *PGTX) Query(c context.Context, sql string, arg ...any) (pgx.Rows, error) {
	if viper.GetBool("db.pg.logging") {
		logger.Call().Infof("Query: %v Arguments: %v", sql, arg)
	}
	return q.db.Query(c, sql, arg...)
}

// QueryRow execute transaction, returning 0 or multiple rows
func (q *PGTX) QueryRow(c context.Context, sql string, arg ...any) pgx.Row {
	if viper.GetBool("db.pg.logging") {
		logger.Call().Infof("QueryRow: %v Arguments: %v", sql, arg)
	}
	return q.db.QueryRow(c, sql, arg...)
}

// Rollback cancel the transaction
func (q *PGTX) Rollback(c context.Context) error {
	if viper.GetBool("db.pg.logging") {
		logger.Call().Infof("Rollingback transaction")
	}
	return q.db.Rollback(c)
}

// Commit proceed the transaction
func (q *PGTX) Commit(c context.Context) error {
	if viper.GetBool("db.pg.logging") {
		logger.Call().Infof("Commiting transaction")
	}
	return q.db.Commit(c)
}
