package catalog

import (
	"example.com/nightguide/internal/domain"
	"sync"
)

type Catalog struct {
	mu     sync.RWMutex
	routes map[string]domain.NightTourRoute
	cached *domain.NightTourRoute
}

func New(routes []domain.NightTourRoute) (*Catalog, error) {
	c := &Catalog{routes: make(map[string]domain.NightTourRoute)}
	for _, route := range routes {
		if err := route.Validate(); err != nil {
			return nil, err
		}
		c.routes[route.ID] = domain.CopyRoute(route)
	}
	return c, nil
}

func (c *Catalog) List(query domain.RouteSearch) ([]domain.NightTourRoute, error) {
	if err := domain.ValidateSearch(query); err != nil {
		return nil, err
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make([]domain.NightTourRoute, 0, len(c.routes))
	for _, route := range c.routes {
		if domain.RouteMatches(route, query) {
			result = append(result, domain.CopyRoute(route))
		}
	}
	return domain.SortRoutes(result), nil
}

func (c *Catalog) GetRouteDetail(routeID string) (domain.NightTourRoute, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cached != nil {
		return domain.CopyRoute(*c.cached), nil
	}
	route, ok := c.routes[routeID]
	if !ok {
		return domain.NightTourRoute{}, domain.ErrRouteNotFound
	}
	copy := domain.CopyRoute(route)
	c.cached = &copy
	return copy, nil
}

func (c *Catalog) UpdateRoute(route domain.NightTourRoute) error {
	if err := route.Validate(); err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.routes[route.ID] = domain.CopyRoute(route)
	if c.cached != nil && c.cached.ID == route.ID {
		c.cached = nil
	}
	return nil
}

func (c *Catalog) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.routes)
}
