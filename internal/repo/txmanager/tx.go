package txmanager

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"test_news/pkg/postgres"
)

type Executor interface {
	Exec(sql string, args ...any) (pgconn.CommandTag, error)
	QueryRow(sql string, args ...any) pgx.Row
	Query(sql string, args ...any) (pgx.Rows, error)
}

type TX interface {
	Executor
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type Manager interface {
	DB(ctx context.Context) Executor
	TX(ctx context.Context) (TX, error)
	TxFunc(ctx context.Context, f func(tx TX) error) (err error)
}

type manager struct {
	pg postgres.Postgres
}

func NewManager(pg postgres.Postgres) Manager {
	return &manager{
		pg: pg,
	}
}

func (m *manager) DB(ctx context.Context) Executor {
	return &poolExec{
		ctx:  ctx,
		pool: m.pg.GetPool(),
	}
}

func (m *manager) TX(ctx context.Context) (TX, error) {
	tx, err := m.pg.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &txExec{
		ctx: ctx,
		tx:  tx,
	}, nil
}

type poolExec struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

func (p *poolExec) Exec(sql string, args ...any) (pgconn.CommandTag, error) {
	return p.pool.Exec(p.ctx, sql, args...)
}

func (p *poolExec) QueryRow(sql string, args ...any) pgx.Row {
	return p.pool.QueryRow(p.ctx, sql, args...)
}

func (p *poolExec) Query(sql string, args ...any) (pgx.Rows, error) {
	return p.pool.Query(p.ctx, sql, args...)
}

type txExec struct {
	ctx context.Context
	tx  pgx.Tx
}

func (t *txExec) Exec(sql string, args ...any) (pgconn.CommandTag, error) {
	return t.tx.Exec(t.ctx, sql, args...)
}

func (t *txExec) QueryRow(sql string, args ...any) pgx.Row {
	return t.tx.QueryRow(t.ctx, sql, args...)
}

func (t *txExec) Query(sql string, args ...any) (pgx.Rows, error) {
	return t.tx.Query(t.ctx, sql, args...)
}

func (t *txExec) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *txExec) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

func (m *manager) TxFunc(ctx context.Context, f func(tx TX) error) (err error) {
	const op = "txmanager.tx.TxFunc"

	tx, err := m.TX(ctx)
	if err != nil {
		return fmt.Errorf("%s init TX error: %w", op, err)
	}
	defer func() {
		if tx == nil {
			return
		}
		if e := tx.Rollback(ctx); e != nil {
			err = fmt.Errorf("%s rollback TX error: %w", op, e)
		}
	}()

	if err = f(tx); err != nil {
		return fmt.Errorf("%s exec func error: %w", op, err)
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s commit TX error: %w", op, err)
	}
	tx = nil
	return nil
}
