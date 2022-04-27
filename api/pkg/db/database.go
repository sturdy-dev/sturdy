package db

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func setup(dbSourceURL string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dbSourceURL)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

func SetupWithTimeout(dbSourceURL string, timeout time.Duration) (*sqlx.DB, error) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			db, err := setup(dbSourceURL)
			if err != nil {
				log.Println("error opening db, retrying...", zap.Error(err))
				continue
			}

			if _, err := db.Exec("SELECT 1"); err != nil {
				log.Println("error opening db, retrying...", zap.Error(err))
				continue
			}

			return db, nil
		case <-time.After(timeout):
			return nil, fmt.Errorf("timeout waiting for db to start")
		}
	}
}
