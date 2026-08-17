package service

import (
	"weddingguide/internal/domain"
	"weddingguide/internal/store"
)

type PreviewService struct {
	store *store.Store
}

func NewPreviewService(repository *store.Store) *PreviewService {
	return &PreviewService{store: repository}
}

func (p *PreviewService) Draft(guideID string) (domain.WeddingGuide, error) {
	return p.store.ReadGuide(guideID)
}

func (p *PreviewService) RenderModel(guideID string) (PreviewModel, error) {
	guide, err := p.Draft(guideID)
	if err != nil {
		return PreviewModel{}, err
	}
	return PreviewModel{Guide: guide, ScheduleCount: len(guide.Schedule), LinkCount: len(guide.Links), PreviewOnly: true}, nil
}

type PreviewModel struct {
	Guide         domain.WeddingGuide
	ScheduleCount int
	LinkCount     int
	PreviewOnly   bool
}
