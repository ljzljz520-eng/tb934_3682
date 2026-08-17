package service

import (
	"path/filepath"
	"testing"

	"weddingguide/internal/domain"
	"weddingguide/internal/store"
)

func TestPreviewDoesNotConsumeVisit(t *testing.T) {
	repository, err := store.New(filepath.Join(t.TempDir(), "preview.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer repository.Close()
	guide := domain.WeddingGuide{ID: "g", Title: "Day", Couple: "A & B", Schedule: []domain.ScheduleItem{{ID: "s", Title: "Ceremony", StartsAt: "10:00"}}, Venue: domain.VenueAddress{Name: "Hall", Line1: "Road", City: "City", Country: "Country"}, Attire: domain.AttireTip{Summary: "Formal"}}
	if err := repository.SaveGuide(guide); err != nil {
		t.Fatal(err)
	}
	preview := NewPreviewService(repository)
	model, err := preview.RenderModel("g")
	if err != nil || !model.PreviewOnly {
		t.Fatalf("preview failed: %#v %v", model, err)
	}
	if _, err := repository.ReadVisitor("g", "preview"); err == nil {
		t.Fatal("preview created visitor record")
	}
}
