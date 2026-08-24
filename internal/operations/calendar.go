package operations

import (
	"fmt"
	"sort"
	"strings"
)

type Slot struct {
	RouteID  string
	Date     string
	Start    string
	Capacity int
	Reserved int
	Closed   bool
}

type Calendar struct{ slots []Slot }

func NewCalendar() *Calendar { return &Calendar{slots: make([]Slot, 0)} }

func (c *Calendar) Add(slot Slot) error {
	if slot.RouteID == "" || slot.Date == "" || slot.Start == "" {
		return fmt.Errorf("slot identity is required")
	}
	if slot.Capacity < 1 || slot.Reserved < 0 || slot.Reserved > slot.Capacity {
		return fmt.Errorf("slot capacity is invalid")
	}
	c.slots = append(c.slots, slot)
	return nil
}

func (c *Calendar) Slots() []Slot {
	result := append([]Slot(nil), c.slots...)
	sort.Slice(result, func(i, j int) bool {
		if result[i].Date == result[j].Date {
			return result[i].Start < result[j].Start
		}
		return result[i].Date < result[j].Date
	})
	return result
}

func (c *Calendar) ForRoute(routeID string) []Slot {
	result := make([]Slot, 0)
	for _, slot := range c.slots {
		if slot.RouteID == routeID {
			result = append(result, slot)
		}
	}
	return result
}

func (c *Calendar) OpenOn(date string) []Slot {
	result := make([]Slot, 0)
	for _, slot := range c.slots {
		if slot.Date == date && !slot.Closed && slot.Reserved < slot.Capacity {
			result = append(result, slot)
		}
	}
	return result
}

func (s Slot) Remaining() int {
	remaining := s.Capacity - s.Reserved
	if remaining < 0 {
		return 0
	}
	return remaining
}

func (s Slot) Label() string {
	state := "open"
	if s.Closed {
		state = "closed"
	} else if s.Remaining() == 0 {
		state = "full"
	}
	return fmt.Sprintf("%s %s (%s)", s.Date, s.Start, state)
}

func NormalizeDate(value string) string { return strings.TrimSpace(value) }

func (c *Calendar) CloseRoute(routeID string) int {
	count := 0
	for index := range c.slots {
		if c.slots[index].RouteID == routeID && !c.slots[index].Closed {
			c.slots[index].Closed = true
			count++
		}
	}
	return count
}

func (c *Calendar) Reserve(routeID, date, start string, party int) error {
	for index := range c.slots {
		slot := &c.slots[index]
		if slot.RouteID == routeID && slot.Date == date && slot.Start == start {
			if slot.Closed {
				return fmt.Errorf("slot is closed")
			}
			if slot.Remaining() < party {
				return fmt.Errorf("slot has insufficient capacity")
			}
			slot.Reserved += party
			return nil
		}
	}
	return fmt.Errorf("slot not found")
}
