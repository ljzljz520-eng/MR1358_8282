package booking

import (
	"example.com/nightguide/internal/catalog"
	"example.com/nightguide/internal/domain"
)

type Query struct {
	catalog *catalog.Catalog
	service *Service
}

func NewQuery(c *catalog.Catalog, service *Service) *Query {
	return &Query{catalog: c, service: service}
}

func (q *Query) Search(query domain.RouteSearch) ([]domain.NightTourRoute, error) {
	return q.catalog.List(query)
}

func (q *Query) Detail(routeID string) (domain.NightTourRoute, error) {
	return q.catalog.GetRouteDetail(routeID)
}

func (q *Query) Confirmation(bookingID string) (domain.BookingConfirmation, error) {
	return q.service.store.FindConfirmation(bookingID)
}
