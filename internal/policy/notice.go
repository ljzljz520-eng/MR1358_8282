package policy

import (
	"example.com/nightguide/internal/domain"
	"fmt"
	"strings"
)

type Notice struct {
	Title    string
	Lines    []string
	Severity string
}

func BuildNotice(route domain.NightTourRoute, policy BookingPolicy) Notice {
	lines := []string{"Arrive ten minutes before departure.", "Bring the confirmation shown after booking."}
	if route.DurationMinutes >= 120 {
		lines = append(lines, "This route lasts two hours or more; plan a rest break.")
	}
	if route.MeetingPoint.Instructions != "" {
		lines = append(lines, route.MeetingPoint.Instructions)
	}
	if policy.LeadMinutes > 0 {
		lines = append(lines, fmt.Sprintf("Changes are accepted before the %d-minute arrival lead time.", policy.LeadMinutes))
	}
	return Notice{Title: route.Name, Lines: lines, Severity: "info"}
}

func (n Notice) Text() string { return strings.Join(n.Lines, " ") }

func (n Notice) HasLine(term string) bool {
	for _, line := range n.Lines {
		if strings.Contains(strings.ToLower(line), strings.ToLower(term)) {
			return true
		}
	}
	return false
}

func MeetingInstruction(route domain.NightTourRoute) string {
	if route.MeetingPoint.Instructions != "" {
		return route.MeetingPoint.Instructions
	}
	return fmt.Sprintf("Meet at %s, %s.", route.MeetingPoint.Name, route.MeetingPoint.Address)
}

func GuestReminder(route domain.NightTourRoute, party int) string {
	if party == 1 {
		return fmt.Sprintf("One guest is booked for %s.", route.Name)
	}
	return fmt.Sprintf("%d guests are booked for %s.", party, route.Name)
}

func CancellationText(route domain.NightTourRoute) string {
	return fmt.Sprintf("Your reservation for %s was cancelled. Please contact the guide before selecting another route.", route.Name)
}
