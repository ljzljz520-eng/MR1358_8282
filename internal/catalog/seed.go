package catalog

import "example.com/nightguide/internal/domain"

func DefaultRoutes() []domain.NightTourRoute {
	return []domain.NightTourRoute{
		{ID: "old-street", Name: "Old Street Lantern Walk", District: "Old Town", Summary: "Courtyards, tea houses, and lantern alleys.", MeetingPoint: domain.MeetingPoint{Name: "Old Street Gate", Address: "18 Lantern Lane", Instructions: "Meet beside the red lantern at Old Street Gate."}, Stops: []string{"Old Street Gate", "Bell Courtyard", "River Steps"}, DurationMinutes: 105, Capacity: 14, Active: true},
		{ID: "night-river", Name: "Night River Story Route", District: "Riverside", Summary: "A quiet story walk beneath the illuminated bridges.", MeetingPoint: domain.MeetingPoint{Name: "Moon Gate Pier", Address: "7 Moon Gate Road", Instructions: "Meet at the blue route marker by Moon Gate Pier."}, Stops: []string{"Moon Gate Pier", "Glass Bridge", "South Embankment"}, DurationMinutes: 120, Capacity: 10, Active: true},
		{ID: "hill-lights", Name: "Hill Lights Observatory", District: "North Hill", Summary: "Night views and neighborhood astronomy tales.", MeetingPoint: domain.MeetingPoint{Name: "Observatory Steps", Address: "2 North Hill Rise", Instructions: "Meet under the observatory clock."}, Stops: []string{"Observatory Steps", "Cedar Lookout", "Wind Bell Terrace"}, DurationMinutes: 135, Capacity: 8, Active: true},
	}
}
