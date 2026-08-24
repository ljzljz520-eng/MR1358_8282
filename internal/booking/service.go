package booking

import (
	"example.com/nightguide/internal/catalog"
	"example.com/nightguide/internal/domain"
	"example.com/nightguide/internal/ledger"
	"example.com/nightguide/internal/policy"
	"example.com/nightguide/internal/store"
	"fmt"
)

type Service struct {
	catalog *catalog.Catalog
	store   *store.Store
	policy  policy.BookingPolicy
	ledger  *ledger.Ledger
}

func NewService(c *catalog.Catalog, s *store.Store) *Service {
	return &Service{catalog: c, store: s, policy: policy.DefaultPolicy(), ledger: ledger.New()}
}

func (s *Service) Submit(routeID, guestName, guestEmail string, partySize int, notes string) (domain.TourBooking, error) {
	route, err := s.catalog.GetRouteDetail(routeID)
	if err != nil {
		return domain.TourBooking{}, err
	}
	booking := domain.TourBooking{ID: BookingID(routeID, guestName), RouteID: routeID, GuestName: guestName, GuestEmail: guestEmail, PartySize: partySize, Status: "pending", Notes: notes}
	if err := booking.Validate(); err != nil {
		return booking, err
	}
	if decision := s.policy.CheckBooking(booking, route); !decision.Allowed {
		return booking, decision.Error()
	}
	if err := s.store.Reserve(route, partySize); err != nil {
		return booking, err
	}
	if err := s.store.SaveBooking(booking); err != nil {
		return booking, fmt.Errorf("save booking: %w", err)
	}
	if _, err := s.ledger.Append(booking.ID, booking.RouteID, "submitted", "guest reservation recorded"); err != nil {
		return booking, err
	}
	return booking, nil
}

func (s *Service) Confirm(bookingID string) (domain.BookingConfirmation, error) {
	booking, err := s.store.FindBooking(bookingID)
	if err != nil {
		return domain.BookingConfirmation{}, err
	}
	if booking.IsConfirmed() {
		return s.store.FindConfirmation(bookingID)
	}
	if !booking.IsPending() {
		return domain.BookingConfirmation{}, fmt.Errorf("booking status %s cannot be confirmed", booking.Status)
	}
	route, err := s.catalog.GetRouteDetail(booking.RouteID)
	if err != nil {
		return domain.BookingConfirmation{}, err
	}
	confirmation := domain.BookingConfirmation{BookingID: booking.ID, RouteID: route.ID, RouteName: route.Name, MeetingPoint: route.MeetingPoint, Stops: route.Stops, PartySize: booking.PartySize, Notice: domain.NoticeFor(route), Status: "confirmed"}
	booking.Status = "confirmed"
	if err := s.store.SaveBooking(booking); err != nil {
		return domain.BookingConfirmation{}, err
	}
	if err := s.store.SaveConfirmation(confirmation); err != nil {
		return domain.BookingConfirmation{}, err
	}
	if _, err := s.ledger.Append(confirmation.BookingID, confirmation.RouteID, "confirmed", "route-specific confirmation recorded"); err != nil {
		return domain.BookingConfirmation{}, err
	}
	return confirmation, nil
}

func (s *Service) Change(bookingID, guestEmail string, partySize int, notes string) (domain.TourBooking, error) {
	booking, err := s.store.FindBooking(bookingID)
	if err != nil {
		return booking, err
	}
	if decision := s.policy.CheckChange(booking); !decision.Allowed {
		if booking.IsConfirmed() {
			return booking, domain.ErrAlreadyConfirmed
		}
		return booking, decision.Error()
	}
	if partySize < 1 || partySize > 12 {
		return booking, domain.ErrInvalidBooking
	}
	booking.GuestEmail = guestEmail
	booking.PartySize = partySize
	booking.Notes = notes
	if err := booking.Validate(); err != nil {
		return booking, err
	}
	if err := s.store.SaveBooking(booking); err != nil {
		return booking, err
	}
	if _, err := s.ledger.Append(booking.ID, booking.RouteID, "changed", "pending guest details updated"); err != nil {
		return booking, err
	}
	return booking, nil
}

func (s *Service) Lookup(bookingID string) (domain.TourBooking, error) {
	return s.store.FindBooking(bookingID)
}

func (s *Service) Audit(bookingID string) []ledger.Event {
	return s.ledger.ForBooking(bookingID)
}
