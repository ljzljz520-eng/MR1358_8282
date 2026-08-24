package httpapi

import (
	"example.com/nightguide/internal/domain"
	"example.com/nightguide/internal/operations"
	"example.com/nightguide/internal/report"
	"net/http"
)

func (s *Server) routes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, domain.ErrInvalidBooking)
		return
	}
	items, err := s.query.Search(routeSearch(r))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) routeDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, domain.ErrInvalidBooking)
		return
	}
	id := r.URL.Query().Get("id")
	item, err := s.query.Detail(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (s *Server) routeBriefing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, domain.ErrInvalidBooking)
		return
	}
	route, err := s.query.Detail(r.URL.Query().Get("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	briefing, err := report.BuildBriefing(route)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, briefing)
}

func (s *Server) routeChecklist(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, domain.ErrInvalidBooking)
		return
	}
	route, err := s.query.Detail(r.URL.Query().Get("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	checklist := operations.NewOpeningChecklist(route)
	writeJSON(w, http.StatusOK, map[string]any{"summary": checklist.Summary(), "complete": checklist.Complete(), "pending": checklist.Pending()})
}

type reservationRequest struct {
	RouteID    string `json:"route_id"`
	GuestName  string `json:"guest_name"`
	GuestEmail string `json:"guest_email"`
	PartySize  int    `json:"party_size"`
	Notes      string `json:"notes"`
}

func (s *Server) reservations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, domain.ErrInvalidBooking)
		return
	}
	var req reservationRequest
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	b, err := s.flow.Submit(req.RouteID, req.GuestName, req.GuestEmail, req.PartySize, req.Notes)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, b)
}

func (s *Server) confirm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, domain.ErrInvalidBooking)
		return
	}
	id := r.URL.Query().Get("id")
	result, err := s.flow.Confirm(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

type changeRequest struct {
	BookingID  string `json:"booking_id"`
	GuestEmail string `json:"guest_email"`
	PartySize  int    `json:"party_size"`
	Notes      string `json:"notes"`
}

func (s *Server) change(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		writeError(w, http.StatusMethodNotAllowed, domain.ErrInvalidBooking)
		return
	}
	var req changeRequest
	if err := decodeBody(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	result, err := s.flow.Change(req.BookingID, req.GuestEmail, req.PartySize, req.Notes)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}
