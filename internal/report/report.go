package report

import (
	"example.com/nightguide/internal/domain"
	"fmt"
	"strings"
)

func ConfirmationText(confirm domain.BookingConfirmation) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Booking %s\n", confirm.BookingID)
	fmt.Fprintf(&b, "%s for %d guests\n", confirm.RouteName, confirm.PartySize)
	fmt.Fprintf(&b, "Meet at %s, %s\n", confirm.MeetingPoint.Name, confirm.MeetingPoint.Address)
	fmt.Fprintf(&b, "Note: %s\n", confirm.Notice)
	if len(confirm.Stops) > 0 {
		fmt.Fprintf(&b, "Stops: %s\n", strings.Join(confirm.Stops, " -> "))
	}
	return b.String()
}

func RouteSummary(route domain.NightTourRoute) string {
	return fmt.Sprintf("%s (%d minutes, %d places)", route.Name, route.DurationMinutes, route.Capacity)
}

func CompareRoutes(routes []domain.NightTourRoute) string {
	lines := make([]string, 0, len(routes))
	for _, route := range routes {
		lines = append(lines, RouteSummary(route))
	}
	return strings.Join(lines, "\n")
}
