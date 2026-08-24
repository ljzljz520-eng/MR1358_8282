package domain

import "testing"

func TestRouteAndBookingValidation(t *testing.T) {
	route := NightTourRoute{ID: "route", Name: "Night Route", MeetingPoint: MeetingPoint{Name: "Gate", Address: "1 Road"}, Stops: []string{"Gate", "Pier"}, DurationMinutes: 60, Capacity: 8}
	if err := route.Validate(); err != nil {
		t.Fatal(err)
	}
	booking := TourBooking{ID: "booking", RouteID: "route", GuestName: "Lin", GuestEmail: "lin@example.test", PartySize: 2, Status: "pending"}
	if err := booking.Validate(); err != nil {
		t.Fatal(err)
	}
	if !booking.IsPending() {
		t.Fatal("new booking should default to pending semantics")
	}
}

func TestNoticeUsesRouteInstructions(t *testing.T) {
	route := NightTourRoute{MeetingPoint: MeetingPoint{Instructions: "Use the blue marker."}}
	if got := NoticeFor(route); got != "Use the blue marker." {
		t.Fatalf("notice=%q", got)
	}
}
