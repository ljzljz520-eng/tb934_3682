package fixture

import "weddingguide/internal/domain"

func Wedding() domain.WeddingGuide {
	return domain.WeddingGuide{
		ID:      "demo-guide",
		Title:   "A day with Lin and Morgan",
		Couple:  "Lin & Morgan",
		Welcome: "We are so glad you can join us. Everything you need for the day is gathered here.",
		Schedule: []domain.ScheduleItem{
			{ID: "arrival", Title: "Guest arrival", Details: "Welcome drinks in the courtyard", StartsAt: "2026-10-18 15:30", EndsAt: "2026-10-18 16:00", Location: "Courtyard", DisplayRank: 1},
			{ID: "ceremony", Title: "Ceremony", Details: "Please take your seat ten minutes early", StartsAt: "2026-10-18 16:00", EndsAt: "2026-10-18 16:45", Location: "Garden", DisplayRank: 2},
			{ID: "dinner", Title: "Dinner and toasts", Details: "Dinner service begins after the family toast", StartsAt: "2026-10-18 18:00", EndsAt: "2026-10-18 20:00", Location: "Glass Hall", DisplayRank: 3},
			{ID: "dance", Title: "Dancing", Details: "Music continues until late", StartsAt: "2026-10-18 20:00", EndsAt: "2026-10-18 23:00", Location: "Glass Hall", DisplayRank: 4},
		},
		Venue:  domain.VenueAddress{Name: "Willow House", Line1: "18 River Lane", Line2: "South Gate", City: "Hangzhou", Region: "Zhejiang", PostalCode: "310000", Country: "China", Latitude: "30.2741", Longitude: "120.1551"},
		Attire: domain.AttireTip{Summary: "Garden formal", Description: "Light tailoring, midi dresses and shoes comfortable on grass are all welcome.", ColorHint: "Jade, coral and soft neutrals", WeatherNote: "The evening can be cool; bring a light layer."},
		Links: []domain.ActionLink{
			{ID: "navigate", Label: "Open navigation", URL: "https://maps.example.test/willow-house", Kind: "navigation", DisplayRank: 1},
			{ID: "seats", Label: "Find my seat", URL: "/seats/demo-guide", Kind: "seating", DisplayRank: 2},
			{ID: "blessing", Label: "Leave a blessing", URL: "/bless/demo-guide", Kind: "blessing", DisplayRank: 3},
		},
		VisitLimit: 4,
	}
}

func PublishedWedding() domain.WeddingGuide {
	guide := Wedding()
	guide.Published = true
	guide.Revision = 2
	return guide
}

func ImportRows(count int) []domain.VisitorRecord {
	rows := make([]domain.VisitorRecord, 0, count)
	for i := 0; i < count; i++ {
		rows = append(rows, domain.VisitorRecord{ID: key(i), GuideID: "demo-guide", VisitorKey: key(i), VisitCount: 1, LastAction: "import", LastSeenAt: "2026-01-02T03:04:05Z"})
	}
	return rows
}

func key(index int) string {
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return "guest-" + string(alphabet[index%len(alphabet)])
}
