package service

import (
	"path/filepath"
	"testing"

	"weddingguide/internal/clock"
	"weddingguide/internal/domain"
	"weddingguide/internal/ids"
	"weddingguide/internal/store"
)

func TestGuideEditorPublishesAndAudits(t *testing.T) {
	repository, err := store.New(filepath.Join(t.TempDir(), "editor.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer repository.Close()
	editor := NewGuideEditor(repository, clock.NewFixed("t1"), ids.New("test"))
	guide := domain.WeddingGuide{Title: "Draft", Couple: "A & B", Schedule: []domain.ScheduleItem{{ID: "s", Title: "Ceremony", StartsAt: "10:00"}}, Venue: domain.VenueAddress{Name: "Hall", Line1: "Road", City: "City", Country: "Country"}, Attire: domain.AttireTip{Summary: "Formal"}}
	created, err := editor.CreateDraft(guide, "planner")
	if err != nil {
		t.Fatal(err)
	}
	created.Title = "Updated"
	if _, err := editor.UpdateDraft(created, "planner"); err != nil {
		t.Fatal(err)
	}
	published, err := editor.Publish(created.ID, "planner")
	if err != nil || !published.Published {
		t.Fatalf("publish failed: %#v %v", published, err)
	}
	audits, err := repository.ListAudits(created.ID)
	if err != nil || len(audits) != 3 {
		t.Fatalf("expected three editor audits: %#v %v", audits, err)
	}
}
