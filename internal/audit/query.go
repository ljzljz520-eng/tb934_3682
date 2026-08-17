package audit

import (
	"strings"

	"weddingguide/internal/domain"
	"weddingguide/internal/service"
)

type Query struct {
	service *service.AuditService
}

func NewQuery(audits *service.AuditService) *Query {
	return &Query{service: audits}
}

func (q *Query) Search(guideID, term string) ([]domain.AuditEntry, error) {
	entries, err := q.service.ForGuide(guideID)
	if err != nil {
		return nil, err
	}
	term = strings.ToLower(strings.TrimSpace(term))
	if term == "" {
		return entries, nil
	}
	filtered := make([]domain.AuditEntry, 0, len(entries))
	for _, entry := range entries {
		if strings.Contains(strings.ToLower(entry.Action), term) || strings.Contains(strings.ToLower(entry.Detail), term) || strings.Contains(strings.ToLower(entry.Actor), term) {
			filtered = append(filtered, entry)
		}
	}
	return filtered, nil
}

func (q *Query) Timeline(guideID string) ([]string, error) {
	entries, err := q.service.ForGuide(guideID)
	if err != nil {
		return nil, err
	}
	lines := make([]string, 0, len(entries))
	for _, entry := range entries {
		lines = append(lines, entry.CreatedAt+" "+entry.Action+" "+entry.EntityID)
	}
	return lines, nil
}
