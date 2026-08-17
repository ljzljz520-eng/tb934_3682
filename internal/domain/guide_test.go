package domain

import "testing"

func TestWeddingGuideValidation(t *testing.T) {
	guide := WeddingGuide{ID: "g", Title: "A day", Couple: "A & B", Schedule: []ScheduleItem{{ID: "s", Title: "Ceremony", StartsAt: "10:00"}}, Venue: VenueAddress{Name: "Hall", Line1: "1 Road", City: "City", Country: "Country"}, Attire: AttireTip{Summary: "Formal"}}
	if err := guide.Validate(); err != nil {
		t.Fatalf("valid guide rejected: %v", err)
	}
	guide.Title = ""
	if err := guide.Validate(); err != ErrMissingTitle {
		t.Fatalf("expected missing title, got %v", err)
	}
}

func TestVisitorLimitAndClone(t *testing.T) {
	visitor := VisitorRecord{ID: "v", GuideID: "g", VisitorKey: "key"}
	visitor = visitor.Increment("view", "now")
	if visitor.VisitCount != 1 || visitor.LastAction != "view" {
		t.Fatalf("unexpected visitor: %#v", visitor)
	}
	if !visitor.Reached(1) {
		t.Fatal("visitor should be at limit")
	}
	guide := WeddingGuide{ID: "g", Schedule: []ScheduleItem{{ID: "s"}}, Links: []ActionLink{{ID: "l"}}}
	clone := guide.Clone()
	clone.Schedule[0].ID = "changed"
	if guide.Schedule[0].ID == clone.Schedule[0].ID {
		t.Fatal("clone shares schedule storage")
	}
}
