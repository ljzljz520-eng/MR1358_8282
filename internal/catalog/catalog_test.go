package catalog

import (
	"example.com/nightguide/internal/domain"
	"testing"
)

func TestCatalogSearchAndUpdate(t *testing.T) {
	c, err := New(DefaultRoutes())
	if err != nil {
		t.Fatal(err)
	}
	items, err := c.List(domain.RouteSearch{District: "Riverside", ActiveOnly: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].ID != "night-river" {
		t.Fatalf("items=%v", items)
	}
	route, err := c.GetRouteDetail("night-river")
	if err != nil {
		t.Fatal(err)
	}
	route.Name = "Night River Stories"
	if err := c.UpdateRoute(route); err != nil {
		t.Fatal(err)
	}
	updated, err := c.GetRouteDetail("night-river")
	if err != nil {
		t.Fatal(err)
	}
	if updated.Name != "Night River Stories" {
		t.Fatalf("name=%s", updated.Name)
	}
}

func TestCatalogReportsMissingRoute(t *testing.T) {
	c, err := New(DefaultRoutes())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.GetRouteDetail("missing"); err == nil {
		t.Fatal("missing route should fail")
	}
}
