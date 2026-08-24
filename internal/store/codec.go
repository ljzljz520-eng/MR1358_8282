package store

import (
	"encoding/json"
	"example.com/nightguide/internal/domain"
	"fmt"
)

func encodeStops(stops []string) (string, error) {
	b, err := json.Marshal(stops)
	if err != nil {
		return "", fmt.Errorf("encode stops: %w", err)
	}
	return string(b), nil
}

func decodeStops(value string) ([]string, error) {
	var stops []string
	if err := json.Unmarshal([]byte(value), &stops); err != nil {
		return nil, fmt.Errorf("decode stops: %w", err)
	}
	return stops, nil
}

func routeArgs(route domain.NightTourRoute) ([]any, error) {
	stops, err := encodeStops(route.Stops)
	if err != nil {
		return nil, err
	}
	active := 0
	if route.Active {
		active = 1
	}
	return []any{route.ID, route.Name, route.District, route.Summary, route.MeetingPoint.Name, route.MeetingPoint.Address, route.MeetingPoint.Instructions, stops, route.DurationMinutes, route.Capacity, active}, nil
}
