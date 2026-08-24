package store

import (
	"context"
	"database/sql"
	"example.com/nightguide/internal/domain"
)

func (s *Store) SaveBooking(booking domain.TourBooking) error {
	_, err := s.db.ExecContext(context.Background(), `INSERT INTO bookings (id,route_id,guest_name,guest_email,party_size,status,notes) VALUES (?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET route_id=excluded.route_id,guest_name=excluded.guest_name,guest_email=excluded.guest_email,party_size=excluded.party_size,status=excluded.status,notes=excluded.notes`, booking.ID, booking.RouteID, booking.GuestName, booking.GuestEmail, booking.PartySize, booking.Status, booking.Notes)
	return err
}

func (s *Store) FindBooking(id string) (domain.TourBooking, error) {
	var booking domain.TourBooking
	err := s.db.QueryRowContext(context.Background(), `SELECT id,route_id,guest_name,guest_email,party_size,status,notes FROM bookings WHERE id=?`, id).Scan(&booking.ID, &booking.RouteID, &booking.GuestName, &booking.GuestEmail, &booking.PartySize, &booking.Status, &booking.Notes)
	if err == sql.ErrNoRows {
		return booking, domain.ErrBookingNotFound
	}
	return booking, err
}

func (s *Store) SaveConfirmation(confirm domain.BookingConfirmation) error {
	stops, err := encodeStops(confirm.Stops)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(context.Background(), `INSERT INTO booking_confirmations (booking_id,route_id,route_name,meeting_name,meeting_address,meeting_instructions,stops,party_size,notice,status) VALUES (?,?,?,?,?,?,?,?,?,?) ON CONFLICT(booking_id) DO UPDATE SET route_id=excluded.route_id,route_name=excluded.route_name,meeting_name=excluded.meeting_name,meeting_address=excluded.meeting_address,meeting_instructions=excluded.meeting_instructions,stops=excluded.stops,party_size=excluded.party_size,notice=excluded.notice,status=excluded.status`, confirm.BookingID, confirm.RouteID, confirm.RouteName, confirm.MeetingPoint.Name, confirm.MeetingPoint.Address, confirm.MeetingPoint.Instructions, stops, confirm.PartySize, confirm.Notice, confirm.Status)
	return err
}

func (s *Store) FindConfirmation(id string) (domain.BookingConfirmation, error) {
	var result domain.BookingConfirmation
	var stops string
	err := s.db.QueryRowContext(context.Background(), `SELECT booking_id,route_id,route_name,meeting_name,meeting_address,meeting_instructions,stops,party_size,notice,status FROM booking_confirmations WHERE booking_id=?`, id).Scan(&result.BookingID, &result.RouteID, &result.RouteName, &result.MeetingPoint.Name, &result.MeetingPoint.Address, &result.MeetingPoint.Instructions, &stops, &result.PartySize, &result.Notice, &result.Status)
	if err != nil {
		return result, err
	}
	result.Stops, err = decodeStops(stops)
	return result, err
}
