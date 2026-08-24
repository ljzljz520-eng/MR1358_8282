package operations

import (
	"example.com/nightguide/internal/domain"
	"fmt"
	"strings"
)

type CheckItem struct {
	Key      string
	Label    string
	Required bool
	Done     bool
	Note     string
}

type OpeningChecklist struct {
	RouteID string
	Items   []CheckItem
}

func NewOpeningChecklist(route domain.NightTourRoute) OpeningChecklist {
	return OpeningChecklist{RouteID: route.ID, Items: []CheckItem{{Key: "route-active", Label: "route accepts reservations", Required: true, Done: route.Active}, {Key: "meeting-point", Label: "meeting point is posted", Required: true, Done: route.MeetingPoint.Name != ""}, {Key: "stops", Label: "at least two stops are ready", Required: true, Done: len(route.Stops) >= 2}, {Key: "capacity", Label: "capacity is positive", Required: true, Done: route.Capacity > 0}, {Key: "instructions", Label: "arrival instructions are present", Required: false, Done: route.MeetingPoint.Instructions != ""}}}
}

func (c OpeningChecklist) Complete() bool {
	for _, item := range c.Items {
		if item.Required && !item.Done {
			return false
		}
	}
	return true
}

func (c OpeningChecklist) Pending() []CheckItem {
	result := make([]CheckItem, 0)
	for _, item := range c.Items {
		if !item.Done {
			result = append(result, item)
		}
	}
	return result
}

func (c *OpeningChecklist) Mark(key string, done bool, note string) bool {
	for index := range c.Items {
		if c.Items[index].Key == key {
			c.Items[index].Done = done
			c.Items[index].Note = note
			return true
		}
	}
	return false
}

func (c OpeningChecklist) Summary() string {
	state := "ready"
	if !c.Complete() {
		state = "blocked"
	}
	return fmt.Sprintf("%s: %s (%d pending)", c.RouteID, state, len(c.Pending()))
}

func (c OpeningChecklist) MissingLabels() string {
	labels := make([]string, 0)
	for _, item := range c.Pending() {
		if item.Required {
			labels = append(labels, item.Label)
		}
	}
	return strings.Join(labels, ", ")
}

func ValidateChecklist(c OpeningChecklist) error {
	if c.RouteID == "" {
		return fmt.Errorf("checklist route is required")
	}
	if len(c.Items) == 0 {
		return fmt.Errorf("checklist has no items")
	}
	return nil
}
