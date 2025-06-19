package txmanager

import (
	"context"
	"database/sql"
	"subscription/internal/storage"

	_ "github.com/lib/pq"
)

type txManagerKey struct{}

type TxManager struct {
	db storage.QueryEngine
}

func NewTxManager(db storage.QueryEngine) *TxManager {
	return &TxManager{db: db}
}

func (m *TxManager) RunTx(
	ctx context.Context,
	fn func(ctxTx context.Context) error,
	isoLevel sql.IsolationLevel,
	readOnly bool,
) error {
	opts := &sql.TxOptions{
		Isolation: isoLevel,
		ReadOnly:  readOnly,
	}

	return m.beginFunc(ctx, opts, fn)
}

func (m *TxManager) beginFunc(ctx context.Context, opts *sql.TxOptions, fn func(ctxTx context.Context) error) error {
	tx, err := m.db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	ctx = context.WithValue(ctx, txManagerKey{}, tx)
	if apiErr := fn(ctx); apiErr != nil {
		return apiErr
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m *TxManager) GetQueryEngine(ctx context.Context) storage.QueryEngine {
	v, ok := ctx.Value(txManagerKey{}).(storage.QueryEngine)
	if ok && v != nil {
		return v
	}

	return m.db
}

func (m *TxManager) Shutdown() error {
	return m.db.Close()
}
