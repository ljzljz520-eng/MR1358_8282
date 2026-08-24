package domain

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrRouteNotFound    = errors.New("route not found")
	ErrBookingNotFound  = errors.New("booking not found")
	ErrInvalidBooking   = errors.New("invalid booking")
	ErrCapacityReached  = errors.New("route capacity reached")
	ErrAlreadyConfirmed = errors.New("booking already confirmed")
)

type MeetingPoint struct {
	Name         string
	Address      string
	Instructions string
}

type NightTourRoute struct {
	ID              string
	Name            string
	District        string
	Summary         string
	MeetingPoint    MeetingPoint
	Stops           []string
	DurationMinutes int
	Capacity        int
	Active          bool
}

type TourBooking struct {
	ID         string
	RouteID    string
	GuestName  string
	GuestEmail string
	PartySize  int
	Status     string
	Notes      string
}

type BookingConfirmation struct {
	BookingID    string
	RouteID      string
	RouteName    string
	MeetingPoint MeetingPoint
	Stops        []string
	PartySize    int
	Notice       string
	Status       string
}

type RouteSearch struct {
	District   string
	ActiveOnly bool
	Term       string
}

func (r NightTourRoute) Validate() error {
	if strings.TrimSpace(r.ID) == "" || strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("route id and name are required")
	}
	if r.DurationMinutes <= 0 || r.Capacity <= 0 {
		return fmt.Errorf("route duration and capacity must be positive")
	}
	if r.MeetingPoint.Name == "" || r.MeetingPoint.Address == "" {
		return fmt.Errorf("route meeting point is incomplete")
	}
	if len(r.Stops) < 2 {
		return fmt.Errorf("route must have at least two stops")
	}
	return nil
}

func (b TourBooking) Validate() error {
	if strings.TrimSpace(b.ID) == "" || strings.TrimSpace(b.RouteID) == "" {
		return ErrInvalidBooking
	}
	if strings.TrimSpace(b.GuestName) == "" || !strings.Contains(b.GuestEmail, "@") {
		return ErrInvalidBooking
	}
	if b.PartySize < 1 || b.PartySize > 12 {
		return ErrInvalidBooking
	}
	return nil
}

func (b TourBooking) IsPending() bool { return b.Status == "pending" }

func (b TourBooking) IsConfirmed() bool { return b.Status == "confirmed" }

func NormalizeTerm(value string) string { return strings.ToLower(strings.TrimSpace(value)) }
