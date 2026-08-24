package store

import (
	"example.com/nightguide/internal/catalog"
	"example.com/nightguide/internal/domain"
	"path/filepath"
	"testing"
)

func TestStoreRoundTrip(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "roundtrip.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	route := catalog.DefaultRoutes()[0]
	if err := s.SaveRoute(route); err != nil {
		t.Fatal(err)
	}
	booking := domain.TourBooking{ID: "old-street-Mina", RouteID: route.ID, GuestName: "Mina", GuestEmail: "mina@example.test", PartySize: 2, Status: "pending"}
	if err := s.SaveBooking(booking); err != nil {
		t.Fatal(err)
	}
	got, err := s.FindBooking(booking.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.RouteID != route.ID || got.PartySize != 2 {
		t.Fatalf("booking=%+v", got)
	}
}

func TestStoreCapacity(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "capacity.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	route := catalog.DefaultRoutes()[2]
	if err := s.SaveRoute(route); err != nil {
		t.Fatal(err)
	}
	if err := s.Reserve(route, route.Capacity); err != nil {
		t.Fatal(err)
	}
	if err := s.Reserve(route, 1); err == nil {
		t.Fatal("capacity should reject overflow")
	}
}
