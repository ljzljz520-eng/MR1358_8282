package app

import (
	"example.com/nightguide/internal/booking"
	"example.com/nightguide/internal/catalog"
	"example.com/nightguide/internal/domain"
	"example.com/nightguide/internal/store"
	"path/filepath"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "persistent.db")
	s, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	route := catalog.DefaultRoutes()[1]
	if err := s.SaveRoute(route); err != nil {
		t.Fatal(err)
	}
	c, err := catalog.New([]domain.NightTourRoute{route})
	if err != nil {
		t.Fatal(err)
	}
	service := booking.NewService(c, s)
	record, err := service.Submit(route.ID, "Reopen Guest", "reopen@example.test", 2, "")
	if err != nil {
		t.Fatal(err)
	}
	confirmation, err := service.Confirm(record.ID)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	gotRoute, err := reopened.FindRoute(route.ID)
	if err != nil {
		t.Fatal(err)
	}
	gotConfirmation, err := reopened.FindConfirmation(confirmation.BookingID)
	if err != nil {
		t.Fatal(err)
	}
	if gotRoute.MeetingPoint.Name != "Moon Gate Pier" || gotConfirmation.MeetingPoint.Name != "Moon Gate Pier" {
		t.Fatalf("route=%+v confirmation=%+v", gotRoute, gotConfirmation)
	}
}
