package guide

import (
	"example.com/nightguide/internal/domain"
	"example.com/nightguide/internal/itinerary"
	"fmt"
	"strings"
)

type StoryBeat struct {
	Stop      string
	Title     string
	Narration string
	Minutes   int
	Prompt    string
}

type Script struct {
	RouteID   string
	RouteName string
	Beats     []StoryBeat
	Opening   string
	Closing   string
}

func BuildScript(route domain.NightTourRoute) (Script, error) {
	plan, err := itinerary.BuildPlan(route)
	if err != nil {
		return Script{}, err
	}
	beats := make([]StoryBeat, 0, len(plan.Stops))
	for index, stop := range plan.Stops {
		minutes := 8
		if index < len(plan.Segments) {
			minutes = plan.Segments[index].Minutes / 2
		}
		if minutes < 4 {
			minutes = 4
		}
		beats = append(beats, StoryBeat{Stop: stop, Title: beatTitle(stop, index), Narration: narrationFor(route, stop, index), Minutes: minutes, Prompt: promptFor(stop)})
	}
	return Script{RouteID: route.ID, RouteName: route.Name, Beats: beats, Opening: openingFor(route), Closing: closingFor(route)}, nil
}

func beatTitle(stop string, index int) string {
	if index == 0 {
		return "Welcome and orientation"
	}
	if strings.Contains(strings.ToLower(stop), "bridge") {
		return "The bridge after dark"
	}
	if strings.Contains(strings.ToLower(stop), "gate") {
		return "The city gate"
	}
	return fmt.Sprintf("Story stop %d", index+1)
}

func narrationFor(route domain.NightTourRoute, stop string, index int) string {
	if index == 0 {
		return fmt.Sprintf("Welcome to %s. We begin at %s and keep the group together.", route.Name, stop)
	}
	if index == len(route.Stops)-1 {
		return fmt.Sprintf("Our final story is here at %s; thank you for exploring the city tonight.", stop)
	}
	return fmt.Sprintf("At %s, notice how the night changes the character of %s.", stop, route.District)
}

func promptFor(stop string) string {
	if strings.Contains(strings.ToLower(stop), "river") || strings.Contains(strings.ToLower(stop), "pier") {
		return "Invite guests to listen for the water before speaking."
	}
	if strings.Contains(strings.ToLower(stop), "hill") || strings.Contains(strings.ToLower(stop), "lookout") {
		return "Pause for a wide view and allow photographs."
	}
	return "Ask guests what detail they noticed first."
}

func openingFor(route domain.NightTourRoute) string {
	return fmt.Sprintf("Good evening. This is %s through %s. Please stay within sight of the guide.", route.Name, route.District)
}

func closingFor(route domain.NightTourRoute) string {
	return fmt.Sprintf("Thank you for joining %s. The route ends at %s.", route.Name, route.Stops[len(route.Stops)-1])
}

func (s Script) Validate() error {
	if s.RouteID == "" || s.RouteName == "" {
		return fmt.Errorf("script route identity is required")
	}
	if len(s.Beats) < 2 {
		return fmt.Errorf("script needs at least two beats")
	}
	if strings.TrimSpace(s.Opening) == "" || strings.TrimSpace(s.Closing) == "" {
		return fmt.Errorf("script opening and closing are required")
	}
	for _, beat := range s.Beats {
		if beat.Stop == "" || beat.Narration == "" || beat.Minutes <= 0 {
			return fmt.Errorf("script beat is incomplete")
		}
	}
	return nil
}

func (s Script) Beat(stop string) (StoryBeat, bool) {
	for _, beat := range s.Beats {
		if strings.EqualFold(beat.Stop, stop) {
			return beat, true
		}
	}
	return StoryBeat{}, false
}

func (s Script) TotalMinutes() int {
	total := 0
	for _, beat := range s.Beats {
		total += beat.Minutes
	}
	return total
}

func (s Script) Text() string {
	lines := []string{s.Opening}
	for index, beat := range s.Beats {
		lines = append(lines, fmt.Sprintf("%d. %s: %s", index+1, beat.Title, beat.Narration))
	}
	lines = append(lines, s.Closing)
	return strings.Join(lines, "\n")
}

func FilterBeats(script Script, term string) []StoryBeat {
	result := make([]StoryBeat, 0)
	term = strings.ToLower(strings.TrimSpace(term))
	for _, beat := range script.Beats {
		if term == "" || strings.Contains(strings.ToLower(beat.Title), term) || strings.Contains(strings.ToLower(beat.Narration), term) {
			result = append(result, beat)
		}
	}
	return result
}
