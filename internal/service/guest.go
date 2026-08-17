package service

import (
	"errors"
	"fmt"
	"strings"

	"weddingguide/internal/clock"
	"weddingguide/internal/domain"
	"weddingguide/internal/ids"
	"weddingguide/internal/store"
)

var (
	ErrGuideUnavailable = errors.New("guide is not published")
	ErrUnknownAction    = errors.New("action link not found")
	ErrBlessingTooShort = errors.New("blessing message is too short")
)

type GuestService struct {
	store *store.Store
	clock clock.Clock
	ids   *ids.Generator
}

type VisitResult struct {
	Guide      domain.WeddingGuide
	Visitor    domain.VisitorRecord
	GentleHint string
	FirstVisit bool
	AtLimit    bool
}

func NewGuestService(repository *store.Store, now clock.Clock, generator *ids.Generator) *GuestService {
	return &GuestService{store: repository, clock: now, ids: generator}
}

func (g *GuestService) Visit(guideID, visitorKey, action string) (VisitResult, error) {
	guide, err := g.store.ReadGuide(guideID)
	if err != nil {
		return VisitResult{}, err
	}
	if !guide.Published {
		return VisitResult{}, ErrGuideUnavailable
	}
	visitor, err := g.store.ReadVisitor(guideID, visitorKey)
	first := false
	if err != nil {
		if !errors.Is(err, store.ErrNotFound) {
			return VisitResult{}, err
		}
		visitor = domain.VisitorRecord{ID: fmt.Sprintf("visitor-%s-%s", guideID, visitorKey), GuideID: guideID, VisitorKey: visitorKey}
		first = true
	}
	visitor = visitor.Increment(action, g.clock.Now())
	if err := g.store.SaveVisitor(visitor); err != nil {
		return VisitResult{}, err
	}
	if err := g.record(guideID, visitorKey, domain.EventGuestVisited, "VisitorRecord", visitor.ID, action); err != nil {
		return VisitResult{}, err
	}
	result := VisitResult{Guide: guide, Visitor: visitor, FirstVisit: first}
	if visitor.Reached(guide.VisitLimit) {
		result.AtLimit = true
		result.GentleHint = "You have already seen this guide a few times; the details will still be here when you need them."
	}
	return result, nil
}

func (g *GuestService) ResolveAction(guideID, linkID string) (domain.ActionLink, error) {
	guide, err := g.store.ReadGuide(guideID)
	if err != nil {
		return domain.ActionLink{}, err
	}
	link, ok := guide.FindLink(linkID)
	if !ok {
		return domain.ActionLink{}, ErrUnknownAction
	}
	return link, nil
}

func (g *GuestService) SubmitBlessing(guideID, visitorKey, name, message string) (domain.Blessing, error) {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(message) == "" {
		return domain.Blessing{}, domain.ErrInvalidGuide
	}
	if len([]rune(strings.TrimSpace(message))) < 8 {
		return domain.Blessing{}, ErrBlessingTooShort
	}
	blessing := domain.Blessing{ID: g.ids.NextFor("blessing"), GuideID: guideID, VisitorKey: visitorKey, Name: strings.TrimSpace(name), Message: strings.TrimSpace(message), CreatedAt: g.clock.Now()}
	if err := g.store.SaveBlessing(blessing); err != nil {
		return domain.Blessing{}, err
	}
	if err := g.record(guideID, visitorKey, domain.EventBlessingSubmitted, "Blessing", blessing.ID, blessing.Name); err != nil {
		return domain.Blessing{}, err
	}
	return blessing, nil
}

func (g *GuestService) record(guideID, actor string, event domain.EventName, entity, entityID, detail string) error {
	entry := (domain.DomainEvent{Name: event, GuideID: guideID, Actor: actor, Entity: entity, EntityID: entityID, Detail: detail, At: g.clock.Now()}).Audit()
	entry.ID = g.ids.NextFor("audit")
	return g.store.AppendAudit(entry)
}
