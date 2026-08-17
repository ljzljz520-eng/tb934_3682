package service

import (
	"path/filepath"
	"testing"

	"weddingguide/internal/clock"
	"weddingguide/internal/domain"
	"weddingguide/internal/ids"
	"weddingguide/internal/store"
)

func TestGuestVisitShowsGentleLimit(t *testing.T) {
	repository, err := store.New(filepath.Join(t.TempDir(), "guest.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer repository.Close()
	guide := domain.WeddingGuide{ID: "g", Title: "Day", Couple: "A & B", Published: true, VisitLimit: 2, Schedule: []domain.ScheduleItem{{ID: "s", Title: "Ceremony", StartsAt: "10:00"}}, Venue: domain.VenueAddress{Name: "Hall", Line1: "Road", City: "City", Country: "Country"}, Attire: domain.AttireTip{Summary: "Formal"}}
	if err := repository.SaveGuide(guide); err != nil {
		t.Fatal(err)
	}
	guest := NewGuestService(repository, clock.NewFixed("t1"), ids.New("guest"))
	if _, err := guest.Visit("g", "v", "view"); err != nil {
		t.Fatal(err)
	}
	result, err := guest.Visit("g", "v", "view")
	if err != nil {
		t.Fatal(err)
	}
	if !result.AtLimit || result.GentleHint == "" || result.Visitor.VisitCount != 2 {
		t.Fatalf("unexpected limit result: %#v", result)
	}
	if _, err := guest.SubmitBlessing("g", "v", "Sam", "Wishing you a lifetime of joy"); err != nil {
		t.Fatal(err)
	}
}
