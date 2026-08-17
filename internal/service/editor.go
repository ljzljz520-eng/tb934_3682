package service

import (
	"fmt"

	"weddingguide/internal/clock"
	"weddingguide/internal/domain"
	"weddingguide/internal/ids"
	"weddingguide/internal/store"
)

type GuideEditor struct {
	store *store.Store
	clock clock.Clock
	ids   *ids.Generator
}

func NewGuideEditor(repository *store.Store, now clock.Clock, generator *ids.Generator) *GuideEditor {
	return &GuideEditor{store: repository, clock: now, ids: generator}
}

func (e *GuideEditor) CreateDraft(input domain.WeddingGuide, actor string) (domain.WeddingGuide, error) {
	if input.ID == "" {
		input.ID = e.ids.NextFor("guide")
	}
	input.Published = false
	input.Revision = 1
	if err := e.store.SaveGuide(input); err != nil {
		return domain.WeddingGuide{}, err
	}
	if err := e.record(input.ID, actor, domain.EventGuideCreated, "WeddingGuide", input.ID, input.Title); err != nil {
		return domain.WeddingGuide{}, err
	}
	return input, nil
}

func (e *GuideEditor) UpdateDraft(input domain.WeddingGuide, actor string) (domain.WeddingGuide, error) {
	existing, err := e.store.ReadGuide(input.ID)
	if err != nil {
		return domain.WeddingGuide{}, err
	}
	input.Revision = existing.Revision + 1
	input.Published = existing.Published
	if err := e.store.SaveGuide(input); err != nil {
		return domain.WeddingGuide{}, err
	}
	if err := e.record(input.ID, actor, domain.EventGuideUpdated, "WeddingGuide", input.ID, fmt.Sprintf("revision=%d", input.Revision)); err != nil {
		return domain.WeddingGuide{}, err
	}
	return input, nil
}

func (e *GuideEditor) Publish(guideID, actor string) (domain.WeddingGuide, error) {
	guide, err := e.store.ReadGuide(guideID)
	if err != nil {
		return domain.WeddingGuide{}, err
	}
	guide.Published = true
	guide.Revision++
	if err := e.store.SaveGuide(guide); err != nil {
		return domain.WeddingGuide{}, err
	}
	if err := e.record(guide.ID, actor, domain.EventGuidePublished, "WeddingGuide", guide.ID, "published"); err != nil {
		return domain.WeddingGuide{}, err
	}
	return guide, nil
}

func (e *GuideEditor) Get(guideID string) (domain.WeddingGuide, error) {
	return e.store.ReadGuide(guideID)
}

func (e *GuideEditor) record(guideID, actor string, event domain.EventName, entity, entityID, detail string) error {
	entry := (domain.DomainEvent{Name: event, GuideID: guideID, Actor: actor, Entity: entity, EntityID: entityID, Detail: detail, At: e.clock.Now()}).Audit()
	entry.ID = e.ids.NextFor("audit")
	return e.store.AppendAudit(entry)
}
