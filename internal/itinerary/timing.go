package itinerary

import (
	"fmt"
	"strconv"
	"strings"
)

type DepartureWindow struct {
	StartMinute int
	EndMinute   int
	Label       string
}

func ParseClock(value string) (int, error) {
	parts := strings.Split(strings.TrimSpace(value), ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("clock must use HH:MM")
	}
	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid hour")
	}
	minute, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid minute")
	}
	if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return 0, fmt.Errorf("clock is out of range")
	}
	return hour*60 + minute, nil
}

func FormatClock(minutes int) string {
	minutes = ((minutes % 1440) + 1440) % 1440
	return fmt.Sprintf("%02d:%02d", minutes/60, minutes%60)
}

func BuildWindow(start, end string) (DepartureWindow, error) {
	startMinute, err := ParseClock(start)
	if err != nil {
		return DepartureWindow{}, err
	}
	endMinute, err := ParseClock(end)
	if err != nil {
		return DepartureWindow{}, err
	}
	if endMinute <= startMinute {
		return DepartureWindow{}, fmt.Errorf("departure window must end after it starts")
	}
	return DepartureWindow{StartMinute: startMinute, EndMinute: endMinute, Label: fmt.Sprintf("%s-%s", FormatClock(startMinute), FormatClock(endMinute))}, nil
}

func (w DepartureWindow) Contains(minutes int) bool {
	return minutes >= w.StartMinute && minutes <= w.EndMinute
}

func (w DepartureWindow) Duration() int {
	if w.EndMinute < w.StartMinute {
		return 0
	}
	return w.EndMinute - w.StartMinute
}

func (w DepartureWindow) Slots(step int) []string {
	if step <= 0 {
		return nil
	}
	result := make([]string, 0)
	for minute := w.StartMinute; minute <= w.EndMinute; minute += step {
		result = append(result, FormatClock(minute))
	}
	return result
}

func WeekdayName(day int) string {
	names := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	if day < 0 || day >= len(names) {
		return "Unknown"
	}
	return names[day]
}

func EveningLabel(window DepartureWindow) string {
	if window.StartMinute >= 21*60 {
		return "late evening"
	}
	if window.StartMinute >= 18*60 {
		return "evening"
	}
	return "early night"
}

func ValidateDeparture(start string, route TourPlan) error {
	minute, err := ParseClock(start)
	if err != nil {
		return err
	}
	if minute < 17*60 {
		return fmt.Errorf("night routes begin after 17:00")
	}
	if route.DurationMinutes > 150 && minute > 21*60 {
		return fmt.Errorf("long routes cannot start after 21:00")
	}
	return nil
}
