package itinerary

import (
	"example.com/nightguide/internal/domain"
	"fmt"
	"strings"
)

type Segment struct {
	From        string
	To          string
	Minutes     int
	Description string
	Covered     bool
}

type TourPlan struct {
	RouteID         string
	RouteName       string
	Meeting         domain.MeetingPoint
	Stops           []string
	Segments        []Segment
	DurationMinutes int
	Accessibility   string
}

func BuildPlan(route domain.NightTourRoute) (TourPlan, error) {
	if err := route.Validate(); err != nil {
		return TourPlan{}, err
	}
	segments, err := buildSegments(route.Stops, route.DurationMinutes)
	if err != nil {
		return TourPlan{}, err
	}
	plan := TourPlan{RouteID: route.ID, RouteName: route.Name, Meeting: route.MeetingPoint, Stops: append([]string(nil), route.Stops...), Segments: segments, DurationMinutes: route.DurationMinutes}
	plan.Accessibility = AccessibilitySummary(plan)
	return plan, nil
}

func buildSegments(stops []string, duration int) ([]Segment, error) {
	if len(stops) < 2 {
		return nil, fmt.Errorf("itinerary needs two stops")
	}
	parts := len(stops) - 1
	base := duration / parts
	extra := duration % parts
	segments := make([]Segment, 0, parts)
	for i := 0; i < parts; i++ {
		minutes := base
		if i < extra {
			minutes++
		}
		segments = append(segments, Segment{From: stops[i], To: stops[i+1], Minutes: minutes, Description: segmentDescription(stops[i], stops[i+1]), Covered: i%2 == 0})
	}
	return segments, nil
}

func segmentDescription(from, to string) string {
	if strings.Contains(strings.ToLower(to), "bridge") {
		return fmt.Sprintf("Cross from %s toward the lit bridge at %s.", from, to)
	}
	if strings.Contains(strings.ToLower(from), "gate") {
		return fmt.Sprintf("Gather at %s, then follow the guide toward %s.", from, to)
	}
	return fmt.Sprintf("Walk from %s to %s while the guide shares local stories.", from, to)
}

func (p TourPlan) Validate() error {
	if p.RouteID == "" || p.RouteName == "" {
		return fmt.Errorf("plan route identity is required")
	}
	if len(p.Stops) != len(p.Segments)+1 {
		return fmt.Errorf("plan segments do not match stops")
	}
	if p.DurationMinutes <= 0 {
		return fmt.Errorf("plan duration must be positive")
	}
	if SumSegmentMinutes(p.Segments) != p.DurationMinutes {
		return fmt.Errorf("plan duration does not match segments")
	}
	return nil
}

func SumSegmentMinutes(segments []Segment) int {
	total := 0
	for _, segment := range segments {
		total += segment.Minutes
	}
	return total
}

func FindStop(plan TourPlan, name string) (int, bool) {
	term := strings.ToLower(strings.TrimSpace(name))
	for index, stop := range plan.Stops {
		if strings.ToLower(stop) == term {
			return index, true
		}
	}
	return -1, false
}

func NextStop(plan TourPlan, current string) string {
	index, ok := FindStop(plan, current)
	if !ok || index+1 >= len(plan.Stops) {
		return ""
	}
	return plan.Stops[index+1]
}

func StopsBefore(plan TourPlan, destination string) []string {
	index, ok := FindStop(plan, destination)
	if !ok {
		return append([]string(nil), plan.Stops...)
	}
	return append([]string(nil), plan.Stops[:index+1]...)
}

func CoveredSegments(plan TourPlan) []Segment {
	result := make([]Segment, 0)
	for _, segment := range plan.Segments {
		if segment.Covered {
			result = append(result, segment)
		}
	}
	return result
}

func UncoveredMinutes(plan TourPlan) int {
	total := 0
	for _, segment := range plan.Segments {
		if !segment.Covered {
			total += segment.Minutes
		}
	}
	return total
}

func AccessibilitySummary(plan TourPlan) string {
	if len(plan.Stops) == 0 {
		return "route details pending"
	}
	if UncoveredMinutes(plan) == 0 {
		return "fully sheltered route"
	}
	if UncoveredMinutes(plan) <= plan.DurationMinutes/3 {
		return "mostly sheltered with short outdoor walks"
	}
	return "outdoor walking route; comfortable shoes recommended"
}

func FormatStops(plan TourPlan) string {
	parts := make([]string, 0, len(plan.Stops))
	for index, stop := range plan.Stops {
		parts = append(parts, fmt.Sprintf("%d. %s", index+1, stop))
	}
	return strings.Join(parts, " | ")
}

func FormatSegments(plan TourPlan) string {
	parts := make([]string, 0, len(plan.Segments))
	for _, segment := range plan.Segments {
		marker := "outdoor"
		if segment.Covered {
			marker = "sheltered"
		}
		parts = append(parts, fmt.Sprintf("%s -> %s (%d min, %s)", segment.From, segment.To, segment.Minutes, marker))
	}
	return strings.Join(parts, "; ")
}
