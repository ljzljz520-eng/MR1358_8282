package ledger

import (
	"fmt"
	"sort"
)

type Summary struct {
	Submitted int
	Changed   int
	Confirmed int
	Cancelled int
}

func Summarize(events []Event) Summary {
	result := Summary{}
	for _, event := range events {
		switch event.Kind {
		case "submitted":
			result.Submitted++
		case "changed":
			result.Changed++
		case "confirmed":
			result.Confirmed++
		case "cancelled":
			result.Cancelled++
		}
	}
	return result
}

func (s Summary) Active() int {
	value := s.Submitted - s.Cancelled
	if value < 0 {
		return 0
	}
	return value
}

func (s Summary) ConfirmationRate() float64 {
	if s.Submitted == 0 {
		return 0
	}
	return float64(s.Confirmed) / float64(s.Submitted)
}

func GroupByRoute(events []Event) map[string][]Event {
	result := make(map[string][]Event)
	for _, event := range events {
		result[event.RouteID] = append(result[event.RouteID], event)
	}
	return result
}

func RouteIDs(events []Event) []string {
	groups := GroupByRoute(events)
	result := make([]string, 0, len(groups))
	for routeID := range groups {
		result = append(result, routeID)
	}
	sort.Strings(result)
	return result
}

func RenderSummary(summary Summary) string {
	return fmt.Sprintf("submitted=%d changed=%d confirmed=%d cancelled=%d active=%d", summary.Submitted, summary.Changed, summary.Confirmed, summary.Cancelled, summary.Active())
}
