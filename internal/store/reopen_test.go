package store

import (
	"path/filepath"
	"testing"

	"weddingguide/internal/domain"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "guide.db")
	first, err := New(path)
	if err != nil {
		t.Fatal(err)
	}
	guide := domain.WeddingGuide{ID: "reopen", Title: "Reopen day", Couple: "A & B", Schedule: []domain.ScheduleItem{{ID: "ceremony", Title: "Ceremony", StartsAt: "10:00"}}, Venue: domain.VenueAddress{Name: "Hall", Line1: "1 Road", City: "City", Country: "Country"}, Attire: domain.AttireTip{Summary: "Formal"}}
	if err := first.SaveGuide(guide); err != nil {
		t.Fatal(err)
	}
	if err := first.AppendAudit(domain.AuditEntry{ID: "a1", GuideID: "reopen", Action: "created", Entity: "WeddingGuide", CreatedAt: "fixed"}); err != nil {
		t.Fatal(err)
	}
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}
	second, err := New(path)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()
	loaded, err := second.ReadGuide("reopen")
	if err != nil || loaded.Title != guide.Title {
		t.Fatalf("reopen lost guide: %#v %v", loaded, err)
	}
	audits, err := second.ListAudits("reopen")
	if err != nil || len(audits) != 1 {
		t.Fatalf("reopen lost audit: %#v %v", audits, err)
	}
}
