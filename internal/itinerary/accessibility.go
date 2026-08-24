package itinerary

import (
	"fmt"
	"strings"
)

type AccessCheck struct {
	Name   string
	Passed bool
	Detail string
}

func CheckAccessibility(plan TourPlan) []AccessCheck {
	checks := []AccessCheck{
		{Name: "meeting point", Passed: strings.TrimSpace(plan.Meeting.Name) != "", Detail: "The guide starts at a named meeting point."},
		{Name: "route continuity", Passed: len(plan.Segments) == len(plan.Stops)-1, Detail: "Every stop has an onward segment except the final stop."},
		{Name: "duration", Passed: plan.DurationMinutes > 0 && SumSegmentMinutes(plan.Segments) == plan.DurationMinutes, Detail: "Segment durations add up to the advertised duration."},
	}
	if UncoveredMinutes(plan) > plan.DurationMinutes/2 {
		checks = append(checks, AccessCheck{Name: "outdoor share", Passed: false, Detail: "More than half the route is outdoors."})
	} else {
		checks = append(checks, AccessCheck{Name: "outdoor share", Passed: true, Detail: "Outdoor sections are limited."})
	}
	return checks
}

func AccessibilityNote(plan TourPlan) string {
	checks := CheckAccessibility(plan)
	failed := make([]string, 0)
	for _, check := range checks {
		if !check.Passed {
			failed = append(failed, check.Name)
		}
	}
	if len(failed) == 0 {
		return "No route access warnings."
	}
	return fmt.Sprintf("Check before booking: %s.", strings.Join(failed, ", "))
}

func HasAccessibleMeeting(plan TourPlan) bool {
	return strings.Contains(strings.ToLower(plan.Meeting.Instructions), "accessible") || strings.Contains(strings.ToLower(plan.Meeting.Instructions), "step-free")
}

func RequiresComfortableShoes(plan TourPlan) bool {
	return UncoveredMinutes(plan) > 0 || len(plan.Stops) > 3
}

func WalkingLoad(plan TourPlan) string {
	if plan.DurationMinutes >= 130 {
		return "long"
	}
	if plan.DurationMinutes >= 90 {
		return "moderate"
	}
	return "short"
}

func ExplainLoad(plan TourPlan) string {
	load := WalkingLoad(plan)
	if load == "long" {
		return "Plan for a long evening walk with rest points."
	}
	if load == "moderate" {
		return "A moderate walk with time for stories and photos."
	}
	return "A short route suitable for a compact evening outing."
}
