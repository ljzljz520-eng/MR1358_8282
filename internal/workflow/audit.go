package workflow

import (
	"example.com/nightguide/internal/domain"
	"fmt"
)

type AuditEntry struct {
	Event     string
	BookingID string
	RouteID   string
}

func RecordSubmission(booking domain.TourBooking) AuditEntry {
	return AuditEntry{Event: "submitted", BookingID: booking.ID, RouteID: booking.RouteID}
}

func RecordConfirmation(confirm domain.BookingConfirmation) AuditEntry {
	return AuditEntry{Event: "confirmed", BookingID: confirm.BookingID, RouteID: confirm.RouteID}
}

func FormatAudit(entry AuditEntry) string {
	return fmt.Sprintf("%s booking=%s route=%s", entry.Event, entry.BookingID, entry.RouteID)
}
