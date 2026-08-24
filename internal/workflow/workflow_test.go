package workflow

import (
	"example.com/nightguide/internal/booking"
	"example.com/nightguide/internal/catalog"
	"example.com/nightguide/internal/domain"
	"example.com/nightguide/internal/store"
	"path/filepath"
	"testing"
)

func TestWorkflowSearchThenReserve(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "workflow.db"))
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
	service := booking.NewService(c, s)
	query := booking.NewQuery(c, service)
	flow := NewReservationFlow(service, query)
	items, err := query.Search(domain.RouteSearch{Term: "lantern", ActiveOnly: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("items=%d", len(items))
	}
	preview, err := flow.Preview(items[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if preview.MeetingPoint.Name != "Old Street Gate" {
		t.Fatal("preview should show route meeting point")
	}
	confirmation, err := flow.ReserveAndConfirm(items[0].ID, "Mia", "mia@example.test", 2, "")
	if err != nil {
		t.Fatal(err)
	}
	if confirmation.Status != "confirmed" {
		t.Fatalf("status=%s", confirmation.Status)
	}
}

func TestWorkflowChangePending(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "change-workflow.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	route := catalog.DefaultRoutes()[1]
	if err := s.SaveRoute(route); err != nil {
		t.Fatal(err)
	}
	c, err := catalog.New([]domain.NightTourRoute{route})
	if err != nil {
		t.Fatal(err)
	}
	service := booking.NewService(c, s)
	query := booking.NewQuery(c, service)
	flow := NewReservationFlow(service, query)
	b, err := flow.Submit(route.ID, "Kai", "kai@example.test", 1, "")
	if err != nil {
		t.Fatal(err)
	}
	changed, err := flow.ModifyPending(b.ID, "kai+changed@example.test", 3, "near the front")
	if err != nil {
		t.Fatal(err)
	}
	if changed.PartySize != 3 || changed.Notes != "near the front" {
		t.Fatalf("changed=%+v", changed)
	}
}

func TestWorkflowConfirmNightRiver(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "river-workflow.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	route := catalog.DefaultRoutes()[1]
	if err := s.SaveRoute(route); err != nil {
		t.Fatal(err)
	}
	c, err := catalog.New([]domain.NightTourRoute{route})
	if err != nil {
		t.Fatal(err)
	}
	service := booking.NewService(c, s)
	query := booking.NewQuery(c, service)
	flow := NewReservationFlow(service, query)
	confirmation, err := flow.ReserveAndConfirm("night-river", "River Guest", "river@example.test", 2, "arrive early")
	if err != nil {
		t.Fatal(err)
	}
	if confirmation.RouteName != "Night River Story Route" || confirmation.MeetingPoint.Name != "Moon Gate Pier" {
		t.Fatalf("confirmation=%+v", confirmation)
	}
}
