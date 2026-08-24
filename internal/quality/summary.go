package quality

import (
	"fmt"
	"sort"
)

type ReviewSummary struct {
	RouteID   string
	Published int
	Average   float64
	Positive  int
	Negative  int
	Mixed     int
}

func SummarizeReviews(routeID string, reviews []Review) ReviewSummary {
	published := FilterPublished(reviews)
	summary := ReviewSummary{RouteID: routeID, Published: len(published), Average: AverageStars(published)}
	for _, review := range published {
		switch review.Sentiment() {
		case "positive":
			summary.Positive++
		case "negative":
			summary.Negative++
		default:
			summary.Mixed++
		}
	}
	return summary
}

func (s ReviewSummary) Trusted() bool {
	if s.Published < 3 {
		return false
	}
	return s.Average >= 3
}

func (s ReviewSummary) Label() string {
	if s.Published == 0 {
		return "No guest reviews yet."
	}
	if s.Trusted() {
		return fmt.Sprintf("%.1f/5 from %d guests", s.Average, s.Published)
	}
	return fmt.Sprintf("%.1f/5 from %d early guests", s.Average, s.Published)
}

func SortReviews(reviews []Review) []Review {
	result := append([]Review(nil), reviews...)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Stars == result[j].Stars {
			return result[i].BookingID < result[j].BookingID
		}
		return result[i].Stars > result[j].Stars
	})
	return result
}

func HighlightedReviews(reviews []Review, limit int) []Review {
	if limit <= 0 {
		return nil
	}
	sorted := SortReviews(FilterPublished(reviews))
	if len(sorted) > limit {
		sorted = sorted[:limit]
	}
	return sorted
}

func TagCounts(reviews []Review) map[string]int {
	counts := make(map[string]int)
	for _, review := range reviews {
		if !review.Published {
			continue
		}
		for _, tag := range review.Tags {
			if tag != "" {
				counts[tag]++
			}
		}
	}
	return counts
}
