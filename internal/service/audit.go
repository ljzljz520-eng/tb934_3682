package service

import (
	"sort"

	"weddingguide/internal/domain"
	"weddingguide/internal/store"
)

type AuditService struct {
	store *store.Store
}

func NewAuditService(repository *store.Store) *AuditService {
	return &AuditService{store: repository}
}

func (a *AuditService) ForGuide(guideID string) ([]domain.AuditEntry, error) {
	entries, err := a.store.ListAudits(guideID)
	if err != nil {
		return nil, err
	}
	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].CreatedAt == entries[j].CreatedAt {
			return entries[i].ID < entries[j].ID
		}
		return entries[i].CreatedAt < entries[j].CreatedAt
	})
	return entries, nil
}

func (a *AuditService) CountAction(guideID string, action domain.EventName) (int, error) {
	entries, err := a.ForGuide(guideID)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, entry := range entries {
		if entry.Action == string(action) {
			count++
		}
	}
	return count, nil
}

func (a *AuditService) Latest(guideID string) (domain.AuditEntry, error) {
	entries, err := a.ForGuide(guideID)
	if err != nil {
		return domain.AuditEntry{}, err
	}
	if len(entries) == 0 {
		return domain.AuditEntry{}, store.ErrNotFound
	}
	return entries[len(entries)-1], nil
}
