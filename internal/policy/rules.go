package policy

import (
	"example.com/nightguide/internal/domain"
	"fmt"
	"strings"
)

type BookingPolicy struct {
	MinimumParty        int
	MaximumParty        int
	LeadMinutes         int
	RequireEmail        bool
	AllowPendingChanges bool
}

func DefaultPolicy() BookingPolicy {
	return BookingPolicy{MinimumParty: 1, MaximumParty: 12, LeadMinutes: 10, RequireEmail: true, AllowPendingChanges: true}
}

type Decision struct {
	Allowed bool
	Code    string
	Message string
}

func (p BookingPolicy) Validate() error {
	if p.MinimumParty < 1 || p.MaximumParty < p.MinimumParty {
		return fmt.Errorf("party limits are invalid")
	}
	if p.LeadMinutes < 0 {
		return fmt.Errorf("lead time cannot be negative")
	}
	return nil
}

func (p BookingPolicy) CheckBooking(booking domain.TourBooking, route domain.NightTourRoute) Decision {
	if err := p.Validate(); err != nil {
		return Decision{Code: "policy", Message: err.Error()}
	}
	if booking.PartySize < p.MinimumParty || booking.PartySize > p.MaximumParty {
		return Decision{Code: "party_size", Message: "party size is outside the route policy"}
	}
	if p.RequireEmail && !strings.Contains(booking.GuestEmail, "@") {
		return Decision{Code: "email", Message: "a reachable email is required"}
	}
	if route.ID == "" {
		return Decision{Code: "route", Message: "selected route is not identified"}
	}
	return Decision{Allowed: true, Code: "ok", Message: "booking meets the route policy"}
}

func (p BookingPolicy) CheckChange(booking domain.TourBooking) Decision {
	if !p.AllowPendingChanges {
		return Decision{Code: "changes_disabled", Message: "pending changes are closed"}
	}
	if booking.Status != "pending" {
		return Decision{Code: "status", Message: "only pending bookings can change"}
	}
	return Decision{Allowed: true, Code: "ok", Message: "pending booking can change"}
}

func (d Decision) Error() error {
	if d.Allowed {
		return nil
	}
	return fmt.Errorf("%s: %s", d.Code, d.Message)
}

func CapacityDecision(route domain.NightTourRoute, reserved, requested int) Decision {
	if requested <= 0 {
		return Decision{Code: "party_size", Message: "requested party must be positive"}
	}
	if reserved < 0 {
		return Decision{Code: "capacity_data", Message: "reserved capacity cannot be negative"}
	}
	if reserved+requested > route.Capacity {
		return Decision{Code: "capacity", Message: "route capacity is already reached"}
	}
	return Decision{Allowed: true, Code: "ok", Message: fmt.Sprintf("%d places remain after this request", route.Capacity-reserved-requested)}
}

func StatusTransition(from, to string) Decision {
	valid := map[string]map[string]bool{"pending": {"confirmed": true, "cancelled": true}, "confirmed": {"cancelled": true}, "cancelled": {}}
	if valid[from][to] {
		return Decision{Allowed: true, Code: "ok", Message: "status transition is allowed"}
	}
	return Decision{Code: "transition", Message: fmt.Sprintf("cannot change status from %s to %s", from, to)}
}
