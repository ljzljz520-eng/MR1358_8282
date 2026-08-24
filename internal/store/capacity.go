package store

import (
	"context"
	"example.com/nightguide/internal/domain"
)

func (s *Store) ReservedParty(routeID string) (int, error) {
	var reserved int
	err := s.db.QueryRowContext(context.Background(), `SELECT reserved_party FROM route_capacity WHERE route_id=?`, routeID).Scan(&reserved)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return 0, nil
	}
	return reserved, err
}

func (s *Store) Reserve(route domain.NightTourRoute, party int) error {
	reserved, err := s.ReservedParty(route.ID)
	if err != nil {
		return err
	}
	if reserved+party > route.Capacity {
		return domain.ErrCapacityReached
	}
	_, err = s.db.ExecContext(context.Background(), `INSERT INTO route_capacity(route_id,reserved_party) VALUES (?,?) ON CONFLICT(route_id) DO UPDATE SET reserved_party=excluded.reserved_party`, route.ID, reserved+party)
	return err
}
