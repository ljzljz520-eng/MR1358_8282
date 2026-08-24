package booking

import (
	"example.com/nightguide/internal/catalog"
	"example.com/nightguide/internal/domain"
	"example.com/nightguide/internal/store"
	"path/filepath"
	"testing"
)

func TestSubmitAndConfirm(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "booking.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	routes := catalog.DefaultRoutes()
	for _, route := range routes {
		if err := s.SaveRoute(route); err != nil {
			t.Fatal(err)
		}
	}
	c, err := catalog.New(routes)
	if err != nil {
		t.Fatal(err)
	}
	service := NewService(c, s)
	booking, err := service.Submit("night-river", "Rui", "rui@example.test", 2, "window seat")
	if err != nil {
		t.Fatal(err)
	}
	confirmation, err := service.Confirm(booking.ID)
	if err != nil {
		t.Fatal(err)
	}
	if confirmation.Status != "confirmed" || confirmation.MeetingPoint.Name != "Moon Gate Pier" {
		t.Fatalf("confirmation=%+v", confirmation)
	}
}

func TestChangePending(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "change.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	route := catalog.DefaultRoutes()[0]
	if err := s.SaveRoute(route); err != nil {
		t.Fatal(err)
	}
	c, err := catalog.New([]domain.NightTourRoute{route})
	if err != nil {
		t.Fatal(err)
	}
	service := NewService(c, s)
	booking, err := service.Submit(route.ID, "Bo", "bo@example.test", 1, "")
	if err != nil {
		t.Fatal(err)
	}
	changed, err := service.Change(booking.ID, "new@example.test", 2, "after work")
	if err != nil {
		t.Fatal(err)
	}
	if changed.GuestEmail != "new@example.test" || changed.PartySize != 2 {
		t.Fatalf("changed=%+v", changed)
	}
}
