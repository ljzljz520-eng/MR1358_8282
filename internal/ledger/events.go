package ledger

import (
	"fmt"
	"sort"
)

type Event struct {
	Sequence  int
	BookingID string
	RouteID   string
	Kind      string
	Detail    string
}

type Ledger struct{ events []Event }

func New() *Ledger { return &Ledger{events: make([]Event, 0)} }

func (l *Ledger) Append(bookingID, routeID, kind, detail string) (Event, error) {
	if bookingID == "" || routeID == "" || kind == "" {
		return Event{}, fmt.Errorf("event identity is required")
	}
	if !validKind(kind) {
		return Event{}, fmt.Errorf("unknown event kind %s", kind)
	}
	event := Event{Sequence: len(l.events) + 1, BookingID: bookingID, RouteID: routeID, Kind: kind, Detail: detail}
	l.events = append(l.events, event)
	return event, nil
}

func validKind(kind string) bool {
	switch kind {
	case "submitted", "changed", "confirmed", "cancelled":
		return true
	default:
		return false
	}
}

func (l *Ledger) Events() []Event { return append([]Event(nil), l.events...) }

func (l *Ledger) ForBooking(bookingID string) []Event {
	result := make([]Event, 0)
	for _, event := range l.events {
		if event.BookingID == bookingID {
			result = append(result, event)
		}
	}
	return result
}

func (l *Ledger) Latest(bookingID string) (Event, bool) {
	events := l.ForBooking(bookingID)
	if len(events) == 0 {
		return Event{}, false
	}
	return events[len(events)-1], true
}

func ValidateSequence(events []Event) error {
	if len(events) == 0 {
		return nil
	}
	copyEvents := append([]Event(nil), events...)
	sort.Slice(copyEvents, func(i, j int) bool { return copyEvents[i].Sequence < copyEvents[j].Sequence })
	for index, event := range copyEvents {
		if event.Sequence != index+1 {
			return fmt.Errorf("event sequence has a gap at %d", index+1)
		}
	}
	return nil
}

func Describe(event Event) string {
	return fmt.Sprintf("%d %s booking=%s route=%s %s", event.Sequence, event.Kind, event.BookingID, event.RouteID, event.Detail)
}
