package guide

import (
	"example.com/nightguide/internal/domain"
	"example.com/nightguide/internal/itinerary"
	"fmt"
	"strings"
)

type GuideKit struct {
	Route       domain.NightTourRoute
	Plan        itinerary.TourPlan
	Script      Script
	QuickFacts  []string
	ArrivalCard string
}

func BuildKit(route domain.NightTourRoute) (GuideKit, error) {
	plan, err := itinerary.BuildPlan(route)
	if err != nil {
		return GuideKit{}, err
	}
	script, err := BuildScript(route)
	if err != nil {
		return GuideKit{}, err
	}
	kit := GuideKit{Route: domain.CopyRoute(route), Plan: plan, Script: script, QuickFacts: quickFacts(route), ArrivalCard: arrivalCard(route)}
	if err := kit.Validate(); err != nil {
		return GuideKit{}, err
	}
	return kit, nil
}

func quickFacts(route domain.NightTourRoute) []string {
	facts := []string{fmt.Sprintf("%d minutes", route.DurationMinutes), fmt.Sprintf("%d stops", len(route.Stops)), fmt.Sprintf("up to %d guests", route.Capacity)}
	if route.Active {
		facts = append(facts, "currently accepting reservations")
	} else {
		facts = append(facts, "not currently accepting reservations")
	}
	return facts
}

func arrivalCard(route domain.NightTourRoute) string {
	return fmt.Sprintf("Meet at %s, %s. %s", route.MeetingPoint.Name, route.MeetingPoint.Address, route.MeetingPoint.Instructions)
}

func (k GuideKit) Validate() error {
	if err := k.Route.Validate(); err != nil {
		return err
	}
	if err := k.Plan.Validate(); err != nil {
		return err
	}
	if err := k.Script.Validate(); err != nil {
		return err
	}
	if len(k.QuickFacts) < 3 {
		return fmt.Errorf("guide kit needs quick facts")
	}
	return nil
}

func (k GuideKit) Search(text string) []string {
	term := strings.ToLower(strings.TrimSpace(text))
	result := make([]string, 0)
	for _, fact := range k.QuickFacts {
		if term == "" || strings.Contains(strings.ToLower(fact), term) {
			result = append(result, fact)
		}
	}
	for _, beat := range FilterBeats(k.Script, term) {
		result = append(result, beat.Title)
	}
	return result
}

func (k GuideKit) Summary() string {
	return fmt.Sprintf("%s\n%s\n%s", k.Route.Name, k.ArrivalCard, itinerary.FormatStops(k.Plan))
}

func (k GuideKit) StopCount() int { return len(k.Plan.Stops) }

func (k GuideKit) HasStop(name string) bool { _, ok := itinerary.FindStop(k.Plan, name); return ok }

func (k GuideKit) NextAfter(name string) string { return itinerary.NextStop(k.Plan, name) }
