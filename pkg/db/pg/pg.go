package pg

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/novando/go-cinema/pkg/logger"
)

type Config struct {
	User      string
	Pass      string
	Host      string
	Port      uint
	Name      string
	Schema    string
	MaxPool   uint
	SSL       bool
	PgBouncer bool
}

type PG struct {
	db *pgxpool.Pool
}

// InitPGXv5
// Initialize database connection
func InitPGXv5(cfg Config, l ...*logger.Logger) (query *PG, err error) {
	log := logger.Call()
	if len(l) > 0 {
		log = l[0]
	}
	if cfg.Name == "" {
		return nil, errors.New("DB_NAME_REQUIRED")
	}
	if cfg.MaxPool == 0 {
		cfg.MaxPool = 5
	}
	if cfg.Host == "" {
		cfg.Host = "localhost"
	}
	if cfg.Port == 0 {
		cfg.Port = 5432
	}
	if cfg.User == "" {
		cfg.User = "postgres"
	}
	if cfg.Schema == "" {
		cfg.Schema = "public"
	}
	sslMode := "disable"
	if cfg.SSL {
		sslMode = "require"
	}
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?pool_max_conns=%d&search_path=%s&sslmode=%s&pgbouncer=%v",
		cfg.User,
		cfg.Pass,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.MaxPool,
		cfg.Schema,
		sslMode,
		cfg.PgBouncer,
	)
	c, err := pgxpool.ParseConfig(url)
	if err != nil {
		if log != nil {
			log.Errorf("PG_CONFIG_ERR:%v", err.Error())
		}
		return
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), c)
	if err != nil {
		if log != nil {
			log.Errorf("PG_CONN_ERR:%v", err.Error())
		}
		return
	}
	if err = pool.Ping(context.Background()); err != nil {
		if log != nil {
			log.Errorf("PG_CONN_ERR:%v", err.Error())
		}
		return
	}
	if log != nil {
		log.Infof("PG connected with %s@%s", cfg.Name, cfg.Host)
	}
	return &PG{pool}, err
}
