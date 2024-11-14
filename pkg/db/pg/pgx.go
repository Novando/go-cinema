package pg

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InitPGXv5
//
// Initialize database connection
func InitPGXv5(
	user string,
	pass string,
	host string,
	port int,
	name string,
	schema string,
	maxPool int,
) (
	pool *pgxpool.Pool,
	query *Queries,
	err error,
) {
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?pool_max_conns=%d&search_path=%s&sslmode=disable",
		user,
		pass,
		host,
		port,
		name,
		maxPool,
		schema,
	)
	c, err := pgxpool.ParseConfig(url)
	if err != nil {
		return
	}

	pool, err = pgxpool.NewWithConfig(context.Background(), c)
	if err != nil {
		return
	}

	query = pgx.NewQuery(pool)
	return
}
