package store

import (
	"context"
	"database/sql"
	"example.com/nightguide/internal/domain"
)

func (s *Store) SaveRoute(route domain.NightTourRoute) error {
	args, err := routeArgs(route)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(context.Background(), `INSERT INTO routes (id,name,district,summary,meeting_name,meeting_address,meeting_instructions,stops,duration_minutes,capacity,active) VALUES (?,?,?,?,?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET name=excluded.name,district=excluded.district,summary=excluded.summary,meeting_name=excluded.meeting_name,meeting_address=excluded.meeting_address,meeting_instructions=excluded.meeting_instructions,stops=excluded.stops,duration_minutes=excluded.duration_minutes,capacity=excluded.capacity,active=excluded.active`, args...)
	return err
}

func (s *Store) ListRoutes() ([]domain.NightTourRoute, error) {
	rows, err := s.db.QueryContext(context.Background(), `SELECT id,name,district,summary,meeting_name,meeting_address,meeting_instructions,stops,duration_minutes,capacity,active FROM routes ORDER BY district,name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []domain.NightTourRoute
	for rows.Next() {
		route, scanErr := scanRoute(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		result = append(result, route)
	}
	return result, rows.Err()
}

func (s *Store) FindRoute(id string) (domain.NightTourRoute, error) {
	row := s.db.QueryRowContext(context.Background(), `SELECT id,name,district,summary,meeting_name,meeting_address,meeting_instructions,stops,duration_minutes,capacity,active FROM routes WHERE id=?`, id)
	route, err := scanRoute(row)
	if err == sql.ErrNoRows {
		return domain.NightTourRoute{}, domain.ErrRouteNotFound
	}
	return route, err
}

type rowScanner interface{ Scan(dest ...any) error }

func scanRoute(row rowScanner) (domain.NightTourRoute, error) {
	var route domain.NightTourRoute
	var stops string
	var active int
	err := row.Scan(&route.ID, &route.Name, &route.District, &route.Summary, &route.MeetingPoint.Name, &route.MeetingPoint.Address, &route.MeetingPoint.Instructions, &stops, &route.DurationMinutes, &route.Capacity, &active)
	if err != nil {
		return route, err
	}
	route.Stops, err = decodeStops(stops)
	route.Active = active == 1
	return route, err
}
