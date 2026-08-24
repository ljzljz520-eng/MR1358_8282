package app

import (
	"example.com/nightguide/internal/booking"
	"example.com/nightguide/internal/catalog"
	"example.com/nightguide/internal/httpapi"
	"example.com/nightguide/internal/quality"
	"example.com/nightguide/internal/store"
	"example.com/nightguide/internal/workflow"
)

type Application struct {
	Store    *store.Store
	Catalog  *catalog.Catalog
	Bookings *booking.Service
	Query    *booking.Query
	Flow     *workflow.ReservationFlow
	HTTP     *httpapi.Server
}

func Open(path string) (*Application, error) {
	s, err := store.Open(path)
	if err != nil {
		return nil, err
	}
	routes := catalog.DefaultRoutes()
	for _, route := range routes {
		if err := quality.Gate(route); err != nil {
			s.Close()
			return nil, err
		}
		if err := s.SaveRoute(route); err != nil {
			s.Close()
			return nil, err
		}
	}
	c, err := catalog.New(routes)
	if err != nil {
		s.Close()
		return nil, err
	}
	bookings := booking.NewService(c, s)
	query := booking.NewQuery(c, bookings)
	flow := workflow.NewReservationFlow(bookings, query)
	return &Application{Store: s, Catalog: c, Bookings: bookings, Query: query, Flow: flow, HTTP: httpapi.NewServer(query, flow)}, nil
}

func (a *Application) Close() error { return a.Store.Close() }
