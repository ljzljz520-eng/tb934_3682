package audit

import (
	"path/filepath"
	"testing"

	"weddingguide/internal/domain"
	"weddingguide/internal/service"
	"weddingguide/internal/store"
)

func TestSearchAndTimeline(t *testing.T) {
	repository, err := store.New(filepath.Join(t.TempDir(), "audit.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer repository.Close()
	for _, entry := range []domain.AuditEntry{{ID: "1", GuideID: "g", Actor: "planner", Action: "guide.updated", Entity: "WeddingGuide", EntityID: "g", Detail: "title", CreatedAt: "02"}, {ID: "2", GuideID: "g", Actor: "sam", Action: "guest.visited", Entity: "VisitorRecord", EntityID: "v", Detail: "view", CreatedAt: "01"}} {
		if err := repository.AppendAudit(entry); err != nil {
			t.Fatal(err)
		}
	}
	query := NewQuery(service.NewAuditService(repository))
	entries, err := query.Search("g", "guest")
	if err != nil || len(entries) != 1 || entries[0].Actor != "sam" {
		t.Fatalf("search mismatch: %#v %v", entries, err)
	}
	lines, err := query.Timeline("g")
	if err != nil || len(lines) != 2 || lines[0][:2] != "01" {
		t.Fatalf("timeline mismatch: %#v %v", lines, err)
	}
}
