package pgstorage

import (
	"database/sql"
	"fmt"
	"log/slog"
	"subscription/internal/models"
)

func scanIntoManInQueue(rows *sql.Rows) (*models.ManInQueue, error) {
	miq := &models.ManInQueue{}
	err := rows.Scan(&miq.ID, &miq.UserID)
	if err != nil {
		return nil, fmt.Errorf("error while scanning data to struct")
	}

	return miq, nil
}

func getManInQueueSlice(rows *sql.Rows) ([]*models.ManInQueue, error) {
	defer func() {
		err := rows.Close()
		if err != nil {
			slog.Error("Error while trying to close rows", slog.String("error", err.Error()))
		}
	}()
	miqs := []*models.ManInQueue{}
	for rows.Next() {
		miq, err := scanIntoManInQueue(rows)
		if err != nil {
			return nil, err
		}
		miqs = append(miqs, miq)
	}

	return miqs, nil
}
