package tours

import (
	"example.com/nightguide/internal/app"
	"path/filepath"
	"testing"
)

func TestNightTourConfirmationUsesOwnMeetingPoint(t *testing.T) {
	application, err := app.Open(filepath.Join(t.TempDir(), "night-tour.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer application.Close()
	if _, err := application.Flow.Preview("old-street"); err != nil {
		t.Fatal(err)
	}
	confirmation, err := application.Flow.ReserveAndConfirm("night-river", "Nora", "nora@example.test", 2, "quiet pace")
	if err != nil {
		t.Fatal(err)
	}
	if confirmation.MeetingPoint.Name != "Moon Gate Pier" {
		t.Fatalf("meeting=%s", confirmation.MeetingPoint.Name)
	}
	if confirmation.RouteID != "night-river" {
		t.Fatalf("route=%s", confirmation.RouteID)
	}
	if confirmation.Notice != "Meet at the blue route marker by Moon Gate Pier." {
		t.Fatalf("notice=%s", confirmation.Notice)
	}
}
