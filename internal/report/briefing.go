package report

import (
	"example.com/nightguide/internal/domain"
	"example.com/nightguide/internal/guide"
	"example.com/nightguide/internal/itinerary"
	"fmt"
	"strings"
)

type Briefing struct {
	RouteName    string
	Arrival      string
	RouteShape   string
	Access       string
	Safety       string
	StoryOpening string
}

func BuildBriefing(route domain.NightTourRoute) (Briefing, error) {
	kit, err := guide.BuildKit(route)
	if err != nil {
		return Briefing{}, err
	}
	brief := Briefing{RouteName: route.Name, Arrival: kit.ArrivalCard, RouteShape: itinerary.FormatSegments(kit.Plan), Access: kit.Plan.Accessibility, Safety: guide.EmergencyCard(route), StoryOpening: kit.Script.Opening}
	if err := brief.Validate(); err != nil {
		return Briefing{}, err
	}
	return brief, nil
}

func (b Briefing) Validate() error {
	if b.RouteName == "" || b.Arrival == "" || b.RouteShape == "" || b.Safety == "" {
		return fmt.Errorf("briefing is incomplete")
	}
	return nil
}

func (b Briefing) Lines() []string {
	return []string{b.RouteName, b.Arrival, b.RouteShape, b.Access, b.Safety, b.StoryOpening}
}

func (b Briefing) Text() string { return strings.Join(b.Lines(), "\n") }

func (b Briefing) Contains(term string) bool {
	term = strings.ToLower(strings.TrimSpace(term))
	if term == "" {
		return true
	}
	for _, line := range b.Lines() {
		if strings.Contains(strings.ToLower(line), term) {
			return true
		}
	}
	return false
}

func GuestFacing(route domain.NightTourRoute) (string, error) {
	brief, err := BuildBriefing(route)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s\nMeet: %s\nRoute: %s\n%s", brief.RouteName, brief.Arrival, brief.RouteShape, brief.Access), nil
}

func Compact(route domain.NightTourRoute) string {
	return fmt.Sprintf("%s at %s", route.Name, route.MeetingPoint.Name)
}
