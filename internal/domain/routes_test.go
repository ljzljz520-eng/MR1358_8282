package domain

import "testing"

func TestRouteMatchingAndSorting(t *testing.T) {
	routes := []NightTourRoute{{ID: "b", Name: "River Walk", District: "Riverside"}, {ID: "a", Name: "Old Walk", District: "Old Town"}}
	if !RouteMatches(routes[0], RouteSearch{Term: "river", ActiveOnly: false}) {
		t.Fatal("term should match route")
	}
	result := SortRoutes(routes)
	if result[0].ID != "a" {
		t.Fatalf("first route=%s", result[0].ID)
	}
	if RouteMatches(routes[0], RouteSearch{District: "Old Town", ActiveOnly: false}) {
		t.Fatal("district should exclude route")
	}
}
