package service

import (
	"sort"
	"strings"

	"weddingguide/internal/domain"
	"weddingguide/internal/store"
)

type InsightsService struct {
	store *store.Store
	audit *AuditService
}

type GuideInsights struct {
	GuideID         string
	Published       bool
	Revision        int
	ScheduleItems   int
	ActionItems     int
	UniqueVisitors  int
	TotalVisits     int
	Blessings       int
	AuditEvents     int
	ImportEvents    int
	LastActivity    string
	RecommendedStep string
	ActionBreakdown map[string]int
}

func NewInsightsService(repository *store.Store, audits *AuditService) *InsightsService {
	return &InsightsService{store: repository, audit: audits}
}

func (i *InsightsService) Build(guideID string) (GuideInsights, error) {
	guide, err := i.store.ReadGuide(guideID)
	if err != nil {
		return GuideInsights{}, err
	}
	visitors, err := i.store.ListVisitors(guideID)
	if err != nil {
		return GuideInsights{}, err
	}
	blessings, err := i.store.ListBlessings(guideID)
	if err != nil {
		return GuideInsights{}, err
	}
	audits, err := i.audit.ForGuide(guideID)
	if err != nil {
		return GuideInsights{}, err
	}
	result := GuideInsights{GuideID: guideID, Published: guide.Published, Revision: guide.Revision, ScheduleItems: len(guide.Schedule), ActionItems: len(guide.Links), UniqueVisitors: len(visitors), Blessings: len(blessings), AuditEvents: len(audits), ActionBreakdown: make(map[string]int)}
	for _, visitor := range visitors {
		result.TotalVisits += visitor.VisitCount
		if visitor.LastAction != "" {
			result.ActionBreakdown[visitor.LastAction]++
		}
	}
	for _, entry := range audits {
		if entry.Action == string(domain.EventImportCompleted) {
			result.ImportEvents++
		}
		if entry.CreatedAt > result.LastActivity {
			result.LastActivity = entry.CreatedAt
		}
	}
	result.RecommendedStep = recommendStep(result)
	return result, nil
}

func recommendStep(insights GuideInsights) string {
	if !insights.Published {
		return "Publish the guide when the details are ready"
	}
	if insights.ScheduleItems == 0 {
		return "Add the ceremony schedule"
	}
	if insights.ActionItems < 3 {
		return "Add navigation, seating and blessing buttons"
	}
	if insights.UniqueVisitors == 0 {
		return "Share the guest link"
	}
	if insights.Blessings == 0 {
		return "Invite guests to leave a blessing"
	}
	return "Keep an eye on the activity timeline"
}

func (i *InsightsService) VisitorSummary(guideID string) (map[string]int, error) {
	visitors, err := i.store.ListVisitors(guideID)
	if err != nil {
		return nil, err
	}
	summary := map[string]int{"one": 0, "two_to_three": 0, "four_plus": 0}
	for _, visitor := range visitors {
		switch {
		case visitor.VisitCount <= 1:
			summary["one"]++
		case visitor.VisitCount <= 3:
			summary["two_to_three"]++
		default:
			summary["four_plus"]++
		}
	}
	return summary, nil
}

func (i *InsightsService) AuditByActor(guideID string) (map[string][]domain.AuditEntry, error) {
	entries, err := i.audit.ForGuide(guideID)
	if err != nil {
		return nil, err
	}
	grouped := make(map[string][]domain.AuditEntry)
	for _, entry := range entries {
		actor := strings.TrimSpace(entry.Actor)
		if actor == "" {
			actor = "system"
		}
		grouped[actor] = append(grouped[actor], entry)
	}
	return grouped, nil
}

func (i *InsightsService) Recent(guideID string, limit int) ([]domain.AuditEntry, error) {
	entries, err := i.audit.ForGuide(guideID)
	if err != nil {
		return nil, err
	}
	if limit < 1 {
		return []domain.AuditEntry{}, nil
	}
	if len(entries) > limit {
		entries = entries[len(entries)-limit:]
	}
	sort.SliceStable(entries, func(left, right int) bool { return entries[left].CreatedAt > entries[right].CreatedAt })
	return entries, nil
}
