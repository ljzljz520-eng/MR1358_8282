package quality

import (
	"fmt"
	"strings"
)

type Review struct {
	BookingID string
	RouteID   string
	Stars     int
	Text      string
	Tags      []string
	Published bool
}

func (r Review) Validate() error {
	if r.BookingID == "" || r.RouteID == "" {
		return fmt.Errorf("review identity is required")
	}
	if r.Stars < 1 || r.Stars > 5 {
		return fmt.Errorf("stars must be from one to five")
	}
	if strings.TrimSpace(r.Text) == "" {
		return fmt.Errorf("review text is required")
	}
	return nil
}

func (r Review) Sentiment() string {
	if r.Stars >= 4 {
		return "positive"
	}
	if r.Stars <= 2 {
		return "negative"
	}
	return "mixed"
}

func (r Review) HasTag(tag string) bool {
	tag = strings.ToLower(strings.TrimSpace(tag))
	for _, item := range r.Tags {
		if strings.ToLower(item) == tag {
			return true
		}
	}
	return false
}

func (r *Review) Publish() error {
	if err := r.Validate(); err != nil {
		return err
	}
	r.Published = true
	return nil
}

func (r Review) DisplayName() string {
	if r.Published {
		return fmt.Sprintf("%d/5: %s", r.Stars, r.Text)
	}
	return "pending moderation"
}

func FilterPublished(reviews []Review) []Review {
	result := make([]Review, 0)
	for _, review := range reviews {
		if review.Published {
			result = append(result, review)
		}
	}
	return result
}

func AverageStars(reviews []Review) float64 {
	if len(reviews) == 0 {
		return 0
	}
	total := 0
	for _, review := range reviews {
		total += review.Stars
	}
	return float64(total) / float64(len(reviews))
}
