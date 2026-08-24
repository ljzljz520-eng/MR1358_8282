package httpapi

import (
	"encoding/json"
	"example.com/nightguide/internal/booking"
	"example.com/nightguide/internal/catalog"
	"example.com/nightguide/internal/domain"
	"example.com/nightguide/internal/store"
	"example.com/nightguide/internal/workflow"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestRoutesEndpoint(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "http.db"))
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
	server := NewServer(query, workflow.NewReservationFlow(service, query))
	req := httptest.NewRequest(http.MethodGet, "/routes?district=Riverside", nil)
	res := httptest.NewRecorder()
	server.Handler().ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("status=%d", res.Code)
	}
	var items []map[string]any
	if err := json.NewDecoder(res.Body).Decode(&items); err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("items=%d", len(items))
	}
}

func TestReservationEndpoint(t *testing.T) {
	s, err := store.Open(filepath.Join(t.TempDir(), "http-reservation.db"))
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
	server := NewServer(query, workflow.NewReservationFlow(service, query))
	body := strings.NewReader(`{"route_id":"night-river","guest_name":"Ada","guest_email":"ada@example.test","party_size":2}`)
	req := httptest.NewRequest(http.MethodPost, "/reservations", body)
	res := httptest.NewRecorder()
	server.Handler().ServeHTTP(res, req)
	if res.Code != http.StatusCreated {
		t.Fatalf("status=%d body=%s", res.Code, res.Body.String())
	}
}
