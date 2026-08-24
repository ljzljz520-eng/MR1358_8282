package workflow

import (
	"example.com/nightguide/internal/booking"
	"example.com/nightguide/internal/domain"
	"fmt"
)

type ReservationFlow struct {
	bookings *booking.Service
	query    *booking.Query
}

func NewReservationFlow(bookings *booking.Service, query *booking.Query) *ReservationFlow {
	return &ReservationFlow{bookings: bookings, query: query}
}

func (f *ReservationFlow) Preview(routeID string) (domain.NightTourRoute, error) {
	return f.query.Detail(routeID)
}

func (f *ReservationFlow) Submit(routeID, guestName, email string, party int, notes string) (domain.TourBooking, error) {
	return f.bookings.Submit(routeID, guestName, email, party, notes)
}

func (f *ReservationFlow) Confirm(bookingID string) (domain.BookingConfirmation, error) {
	return f.bookings.Confirm(bookingID)
}

func (f *ReservationFlow) Change(bookingID, email string, party int, notes string) (domain.TourBooking, error) {
	return f.bookings.Change(bookingID, email, party, notes)
}

func (f *ReservationFlow) ReserveAndConfirm(routeID, guestName, email string, party int, notes string) (domain.BookingConfirmation, error) {
	bookingRecord, err := f.bookings.Submit(routeID, guestName, email, party, notes)
	if err != nil {
		return domain.BookingConfirmation{}, err
	}
	confirmation, err := f.bookings.Confirm(bookingRecord.ID)
	if err != nil {
		return domain.BookingConfirmation{}, fmt.Errorf("confirm reservation: %w", err)
	}
	return confirmation, nil
}

func (f *ReservationFlow) ModifyPending(bookingID, email string, party int, notes string) (domain.TourBooking, error) {
	return f.bookings.Change(bookingID, email, party, notes)
}
