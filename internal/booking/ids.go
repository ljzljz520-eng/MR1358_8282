package booking

import "fmt"

func BookingID(routeID, guestName string) string {
	clean := ""
	for _, ch := range guestName {
		if ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' {
			clean += string(ch)
		}
	}
	if clean == "" {
		clean = "guest"
	}
	return fmt.Sprintf("%s-%s", routeID, clean)
}
