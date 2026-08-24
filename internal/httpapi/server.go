package httpapi

import (
	"encoding/json"
	"example.com/nightguide/internal/booking"
	"example.com/nightguide/internal/domain"
	"example.com/nightguide/internal/workflow"
	"net/http"
)

type Server struct {
	query *booking.Query
	flow  *workflow.ReservationFlow
	mux   *http.ServeMux
}

func NewServer(query *booking.Query, flow *workflow.ReservationFlow) *Server {
	s := &Server{query: query, flow: flow, mux: http.NewServeMux()}
	s.mux.HandleFunc("/routes", s.routes)
	s.mux.HandleFunc("/routes/detail", s.routeDetail)
	s.mux.HandleFunc("/routes/briefing", s.routeBriefing)
	s.mux.HandleFunc("/routes/checklist", s.routeChecklist)
	s.mux.HandleFunc("/reservations", s.reservations)
	s.mux.HandleFunc("/reservations/confirm", s.confirm)
	s.mux.HandleFunc("/reservations/change", s.change)
	return s
}

func (s *Server) Handler() http.Handler { return s.mux }

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func decodeBody(r *http.Request, target any) error { return json.NewDecoder(r.Body).Decode(target) }

func routeSearch(r *http.Request) domain.RouteSearch {
	return domain.RouteSearch{District: r.URL.Query().Get("district"), Term: r.URL.Query().Get("q"), ActiveOnly: r.URL.Query().Get("active") != "false"}
}
