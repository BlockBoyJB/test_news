package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"runtime"
	"time"
)

const (
	defaultConnAttempts = 10
	defaultConnTimeout  = 2 * time.Second
)

type Postgres interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Begin(ctx context.Context) (pgx.Tx, error)
	Close()
	GetPool() *pgxpool.Pool
}

type postgres struct {
	*pgxpool.Pool
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration
}

func NewPG(url string) (Postgres, error) {
	pg := &postgres{
		maxPoolSize:  runtime.NumCPU(),
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
	}

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	poolConfig.MaxConns = int32(pg.maxPoolSize)
	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}
		log.Printf("Postgres trying to connect, attemps left: %d", pg.connAttempts)
		time.Sleep(pg.connTimeout)
		pg.connAttempts--
	}
	if err != nil {
		return nil, fmt.Errorf("error connect to postgres, %w", err)
	}
	return pg, err
}

func (p *postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

func (p *postgres) GetPool() *pgxpool.Pool {
	return p.Pool
}
