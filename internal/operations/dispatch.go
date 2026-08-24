package operations

import (
	"example.com/nightguide/internal/domain"
	"fmt"
	"sort"
	"strings"
)

type Guide struct {
	ID        string
	Name      string
	Languages []string
	MaxGuests int
	Active    bool
}

type DispatchPlan struct {
	RouteID    string
	GuideID    string
	GuestCount int
	Language   string
	Status     string
	Reason     string
}

func ChooseGuide(route domain.NightTourRoute, guides []Guide, language string, guests int) (Guide, error) {
	if guests < 1 {
		return Guide{}, fmt.Errorf("guest count must be positive")
	}
	candidates := make([]Guide, 0)
	for _, guide := range guides {
		if !guide.Active || guide.MaxGuests < guests {
			continue
		}
		if language != "" && !hasLanguage(guide, language) {
			continue
		}
		candidates = append(candidates, guide)
	}
	if len(candidates) == 0 {
		return Guide{}, fmt.Errorf("no guide can cover %s", route.Name)
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].MaxGuests == candidates[j].MaxGuests {
			return candidates[i].Name < candidates[j].Name
		}
		return candidates[i].MaxGuests < candidates[j].MaxGuests
	})
	return candidates[0], nil
}

func hasLanguage(guide Guide, language string) bool {
	language = strings.ToLower(strings.TrimSpace(language))
	for _, item := range guide.Languages {
		if strings.ToLower(item) == language {
			return true
		}
	}
	return false
}

func MakeDispatch(route domain.NightTourRoute, guides []Guide, language string, guests int) (DispatchPlan, error) {
	guide, err := ChooseGuide(route, guides, language, guests)
	if err != nil {
		return DispatchPlan{RouteID: route.ID, GuestCount: guests, Language: language, Status: "unassigned", Reason: err.Error()}, err
	}
	return DispatchPlan{RouteID: route.ID, GuideID: guide.ID, GuestCount: guests, Language: language, Status: "assigned", Reason: fmt.Sprintf("%s covers %s", guide.Name, route.Name)}, nil
}

func (p DispatchPlan) Assigned() bool { return p.Status == "assigned" && p.GuideID != "" }

func (p DispatchPlan) Summary() string {
	if p.Assigned() {
		return fmt.Sprintf("%s assigned to %s for %d guests", p.RouteID, p.GuideID, p.GuestCount)
	}
	return fmt.Sprintf("%s is unassigned: %s", p.RouteID, p.Reason)
}

func (g Guide) Validate() error {
	if g.ID == "" || g.Name == "" {
		return fmt.Errorf("guide identity is required")
	}
	if g.MaxGuests < 1 {
		return fmt.Errorf("guide capacity is required")
	}
	return nil
}

func ActiveGuides(guides []Guide) []Guide {
	result := make([]Guide, 0)
	for _, guide := range guides {
		if guide.Active {
			result = append(result, guide)
		}
	}
	return result
}
