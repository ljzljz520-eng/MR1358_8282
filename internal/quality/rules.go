package quality

import (
	"example.com/nightguide/internal/domain"
	"fmt"
	"strings"
)

type Result struct {
	Rule   string
	Passed bool
	Detail string
}

func ReviewRoute(route domain.NightTourRoute) []Result {
	results := []Result{{Rule: "identity", Passed: route.ID != "" && route.Name != "", Detail: "route has an identifier and name"}, {Rule: "meeting", Passed: route.MeetingPoint.Name != "" && route.MeetingPoint.Address != "", Detail: "meeting point has a name and address"}, {Rule: "stops", Passed: len(route.Stops) >= 2, Detail: "route has at least two stops"}, {Rule: "capacity", Passed: route.Capacity > 0, Detail: "route capacity is positive"}}
	if strings.TrimSpace(route.MeetingPoint.Instructions) == "" {
		results = append(results, Result{Rule: "arrival instructions", Passed: false, Detail: "guests have no arrival guidance"})
	} else {
		results = append(results, Result{Rule: "arrival instructions", Passed: true, Detail: "arrival guidance is present"})
	}
	return results
}

func Passed(results []Result) bool {
	for _, result := range results {
		if !result.Passed {
			return false
		}
	}
	return true
}

func FailedRules(results []Result) []Result {
	failed := make([]Result, 0)
	for _, result := range results {
		if !result.Passed {
			failed = append(failed, result)
		}
	}
	return failed
}

func Report(results []Result) string {
	if Passed(results) {
		return "route quality passed"
	}
	parts := make([]string, 0)
	for _, result := range FailedRules(results) {
		parts = append(parts, result.Rule)
	}
	return fmt.Sprintf("route quality failed: %s", strings.Join(parts, ", "))
}

func Score(results []Result) int {
	if len(results) == 0 {
		return 0
	}
	passed := 0
	for _, result := range results {
		if result.Passed {
			passed++
		}
	}
	return passed * 100 / len(results)
}

func Gate(route domain.NightTourRoute) error {
	results := ReviewRoute(route)
	if !Passed(results) {
		return fmt.Errorf("quality gate: %s", Report(results))
	}
	return nil
}
