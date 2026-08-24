package store

import (
	"context"
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
)

type Store struct{ db *sql.DB }

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) migrate() error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS routes (id TEXT PRIMARY KEY, name TEXT NOT NULL, district TEXT NOT NULL, summary TEXT NOT NULL, meeting_name TEXT NOT NULL, meeting_address TEXT NOT NULL, meeting_instructions TEXT NOT NULL, stops TEXT NOT NULL, duration_minutes INTEGER NOT NULL, capacity INTEGER NOT NULL, active INTEGER NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS bookings (id TEXT PRIMARY KEY, route_id TEXT NOT NULL, guest_name TEXT NOT NULL, guest_email TEXT NOT NULL, party_size INTEGER NOT NULL, status TEXT NOT NULL, notes TEXT NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS booking_confirmations (booking_id TEXT PRIMARY KEY, route_id TEXT NOT NULL, route_name TEXT NOT NULL, meeting_name TEXT NOT NULL, meeting_address TEXT NOT NULL, meeting_instructions TEXT NOT NULL, stops TEXT NOT NULL, party_size INTEGER NOT NULL, notice TEXT NOT NULL, status TEXT NOT NULL)`,
		`CREATE TABLE IF NOT EXISTS route_capacity (route_id TEXT PRIMARY KEY, reserved_party INTEGER NOT NULL)`,
	}
	ctx := context.Background()
	for _, statement := range statements {
		if _, err := s.db.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
	}
	return nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) Ping() error { return s.db.Ping() }
