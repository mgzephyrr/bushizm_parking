package pgstorage

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"subscription/internal/storage/txmanager"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

const maxDequeueTimes = 20

type PostgresStorage struct {
	tx               *txmanager.TxManager
	lastDequeueTimes []time.Time
	mu               sync.Mutex
}

func NewPostgresStore(connStr string) (*PostgresStorage, *txmanager.TxManager, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, nil, err
	}

	slog.Info("Connected to PostgreSQL")

	tx := txmanager.NewTxManager(db)

	return &PostgresStorage{
		tx:               tx,
		lastDequeueTimes: []time.Time{},
	}, tx, nil
}

func (s *PostgresStorage) Shutdown(ctx context.Context) error {
	slog.Info("Closing database connection")
	if err := s.tx.Shutdown(); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		slog.Info("Postgres connection closed")
		return nil
	}
}

func (s *PostgresStorage) Migrate(ctx context.Context, path string) error {
	if err := goose.Up(s.tx.GetQueryEngine(ctx).(*sql.DB), path); err != nil {
		slog.Error("Error while applying migrations", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (s *PostgresStorage) GetByID(ctx context.Context, tableName string, ID int) (*sql.Rows, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", tableName)
	rows, err := s.tx.GetQueryEngine(ctx).Query(query, ID)
	if err != nil {
		slog.Error("Error while selecting rows", slog.String("table", tableName), slog.String("error", err.Error()))
		return nil, fmt.Errorf("unexpected error while selecting rows")
	}

	return rows, nil
}

func (s *PostgresStorage) GetAll(ctx context.Context, tableName string) (*sql.Rows, error) {
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := s.tx.GetQueryEngine(ctx).Query(query)
	if err != nil {
		slog.Error("Error while selecting rows", slog.String("table", tableName), slog.String("error", err.Error()))
		return nil, fmt.Errorf("unexpected error while selecting rows")
	}

	return rows, nil
}

func (s *PostgresStorage) GetByQuery(ctx context.Context, tableName string, query string, args ...any) (*sql.Rows, error) {
	rows, err := s.tx.GetQueryEngine(ctx).Query(query, args...)
	if err != nil {
		slog.Error("Error while selecting rows", slog.String("table", tableName), slog.String("error", err.Error()))
		return nil, fmt.Errorf("unexpected error while selecting rows")
	}

	return rows, nil
}

func (s *PostgresStorage) GetPage(ctx context.Context, tableName string, limit, lastID int) (*sql.Rows, error) {
	query := fmt.Sprintf(`
		SELECT * FROM %s
		WHERE id > $1
		ORDER BY id
		LIMIT $2
	`, tableName)

	rows, err := s.tx.GetQueryEngine(ctx).Query(query, lastID, limit)
	if err != nil {
		slog.Error("Error while selecting rows", slog.String("table", tableName), slog.String("error", err.Error()))
		return nil, fmt.Errorf("unexpected error while selecting rows")
	}

	return rows, nil
}

func (s *PostgresStorage) Create(ctx context.Context, query, tableName string, args ...any) (int, error) {
	var id int

	err := s.tx.GetQueryEngine(ctx).QueryRow(query, args...).Scan(&id)
	if err == nil {
		return id, nil
	}

	slog.Error("Error while inserting values", slog.String("table", tableName), slog.String("error", err.Error()))
	return 0, err
}

func (s *PostgresStorage) Update(ctx context.Context, query, msgNotFound, tableName string, args ...any) error {
	res, err := s.tx.GetQueryEngine(ctx).Exec(query, args...)
	if err != nil {
		slog.Error("Error while updating rows", slog.String("table", tableName), slog.String("error", err.Error()))
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		slog.Error("Error while receiving number of affected rows", slog.String("table", tableName), slog.String("error", err.Error()))
		return err
	}

	if rowsAffected == 0 {
		return &ErrorDoesNotExist{
			What:  msgNotFound,
			Inner: nil,
		}
	}

	return nil
}

func (s *PostgresStorage) Delete(ctx context.Context, msgNotFound, tableName string, ID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", tableName)
	res, err := s.tx.GetQueryEngine(ctx).Exec(query, ID)
	if err != nil {
		slog.Error("Error while deleting rows", slog.String("table", tableName), slog.String("error", err.Error()))
		return fmt.Errorf("unexpected error while deleting rows")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		slog.Error("Error while receiving number of affected rows", slog.String("table", tableName), slog.String("error", err.Error()))
		return fmt.Errorf("unexpected error while receiving number of affected rows")
	}

	if rowsAffected == 0 {
		return &ErrorDoesNotExist{
			What:  msgNotFound,
			Inner: nil,
		}
	}

	return nil
}

func (s *PostgresStorage) EstimateWaitTime(position int) time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()

	times := s.lastDequeueTimes
	if len(times) < 2 {
		return 0
	}

	var total time.Duration
	for i := 1; i < len(times); i++ {
		total += times[i].Sub(times[i-1])
	}

	avgInterval := total / time.Duration(len(times)-1)
	return avgInterval * time.Duration(position-1)
}

func (s *PostgresStorage) GetUserPosition(ctx context.Context, userID int) (int, error) {
	query := `SELECT COUNT(*) + 1 AS position
	FROM parking_queue
	AND created_at < (
		SELECT created_at
		FROM parking_queue
		WHERE user_id = $1
		ORDER BY created_at
		LIMIT 1
	);`

	var position int

	err := s.tx.GetQueryEngine(ctx).QueryRow(query).Scan(&position)
	if err != nil {
		slog.Error("Error while selecting position", slog.String("table", "parking_queue"), slog.String("error", err.Error()))
		return 0, fmt.Errorf("unexpected error while checking position")
	}

	return position, nil
}
