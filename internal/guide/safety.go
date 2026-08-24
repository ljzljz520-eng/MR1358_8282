package guide

import (
	"example.com/nightguide/internal/domain"
	"example.com/nightguide/internal/itinerary"
	"fmt"
	"strings"
)

type SafetyBrief struct {
	RouteID       string
	MeetingRule   string
	GroupRule     string
	WeatherRule   string
	EmergencyRule string
	Warnings      []string
}

func BuildSafetyBrief(route domain.NightTourRoute) (SafetyBrief, error) {
	plan, err := itinerary.BuildPlan(route)
	if err != nil {
		return SafetyBrief{}, err
	}
	warnings := make([]string, 0)
	if itinerary.UncoveredMinutes(plan) > plan.DurationMinutes/2 {
		warnings = append(warnings, "more than half of the route is outdoors")
	}
	if route.DurationMinutes >= 120 {
		warnings = append(warnings, "schedule a water break")
	}
	if len(route.Stops) > 4 {
		warnings = append(warnings, "count the group at every stop")
	}
	brief := SafetyBrief{RouteID: route.ID, MeetingRule: fmt.Sprintf("Gather at %s before departure.", route.MeetingPoint.Name), GroupRule: "Keep the guide in sight and cross only at marked crossings.", WeatherRule: "Pause the walk when conditions make the path unsafe.", EmergencyRule: "Return to the nearest staffed public place and contact local emergency services.", Warnings: warnings}
	if err := brief.Validate(); err != nil {
		return SafetyBrief{}, err
	}
	return brief, nil
}

func (b SafetyBrief) Validate() error {
	if b.RouteID == "" || b.MeetingRule == "" || b.GroupRule == "" || b.EmergencyRule == "" {
		return fmt.Errorf("safety brief is incomplete")
	}
	return nil
}

func (b SafetyBrief) Text() string {
	lines := []string{b.MeetingRule, b.GroupRule, b.WeatherRule, b.EmergencyRule}
	lines = append(lines, b.Warnings...)
	return strings.Join(lines, " ")
}

func (b SafetyBrief) HasWarning(term string) bool {
	term = strings.ToLower(term)
	for _, warning := range b.Warnings {
		if strings.Contains(strings.ToLower(warning), term) {
			return true
		}
	}
	return false
}

func (b SafetyBrief) WarningCount() int { return len(b.Warnings) }

func EmergencyCard(route domain.NightTourRoute) string {
	return fmt.Sprintf("If the group separates on %s, return to %s.", route.Name, route.MeetingPoint.Name)
}
