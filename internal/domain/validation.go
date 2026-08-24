package domain

import (
	"fmt"
	"strings"
)

func ValidateSearch(query RouteSearch) error {
	if len(strings.TrimSpace(query.District)) > 80 {
		return fmt.Errorf("district filter is too long")
	}
	if len(strings.TrimSpace(query.Term)) > 100 {
		return fmt.Errorf("search term is too long")
	}
	return nil
}

func NormalizeStatus(status string) string {
	status = NormalizeTerm(status)
	if status == "" {
		return "pending"
	}
	return status
}

func ValidBookingStatus(status string) bool {
	switch NormalizeStatus(status) {
	case "pending", "confirmed", "cancelled":
		return true
	default:
		return false
	}
}

func NoticeFor(route NightTourRoute) string {
	if route.MeetingPoint.Instructions != "" {
		return route.MeetingPoint.Instructions
	}
	return "Arrive ten minutes before departure and show your confirmation."
}
