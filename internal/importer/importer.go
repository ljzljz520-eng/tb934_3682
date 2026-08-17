package importer

import (
	"errors"
	"fmt"

	"weddingguide/internal/clock"
	"weddingguide/internal/domain"
	"weddingguide/internal/ids"
	"weddingguide/internal/store"
)

var ErrImportStopped = errors.New("visitor import stopped")

type Service struct {
	store *store.Store
	clock clock.Clock
	ids   *ids.Generator
}

type Report struct {
	Processed int
	Created   int
	Updated   int
	Rejected  int
	Errors    []error
}

func New(repository *store.Store, now clock.Clock, generator *ids.Generator) *Service {
	return &Service{store: repository, clock: now, ids: generator}
}

func (s *Service) ImportRows(guideID, actor string, rows []store.ImportRow) (Report, error) {
	report := Report{Errors: make([]error, 0)}
	for index, row := range rows {
		if err := s.processRow(guideID, actor, index, row, &report); err != nil {
			return report, err
		}
	}
	return report, nil
}

// processRow handles a single import row. Opening and closing the reader is
// scoped to this call so the reader slot is released as soon as the row is
// done — never held across rows. This keeps a batch from exhausting the
// store's active-reader limit before every record is imported and audited.
func (s *Service) processRow(guideID, actor string, index int, row store.ImportRow, report *Report) error {
	reader, err := s.store.OpenReader()
	if err != nil {
		report.Errors = append(report.Errors, fmt.Errorf("row %d: %w", index, err))
		return ErrImportStopped
	}
	defer reader.Close()
	if err := s.validateRow(reader, guideID, row); err != nil {
		report.Rejected++
		report.Errors = append(report.Errors, fmt.Errorf("row %d: %w", index, err))
		return nil
	}
	_, created, err := s.store.ImportVisitor(row, guideID)
	if err != nil {
		report.Errors = append(report.Errors, fmt.Errorf("row %d: %w", index, err))
		return ErrImportStopped
	}
	report.Processed++
	if created {
		report.Created++
	} else {
		report.Updated++
	}
	if err := s.audit(guideID, actor, row, report.Processed); err != nil {
		report.Errors = append(report.Errors, err)
		return ErrImportStopped
	}
	return nil
}

func (s *Service) validateRow(reader *store.Reader, guideID string, row store.ImportRow) error {
	if reader == nil {
		return store.ErrReaderLimit
	}
	if guideID == "" || row.VisitorKey == "" {
		return domain.ErrMissingVisitor
	}
	if row.Action == "" {
		return domain.ErrInvalidGuide
	}
	return nil
}

func (s *Service) audit(guideID, actor string, row store.ImportRow, count int) error {
	entry := domain.DomainEvent{Name: domain.EventImportCompleted, GuideID: guideID, Actor: actor, Entity: "VisitorRecord", EntityID: row.VisitorKey, Detail: fmt.Sprintf("processed=%d action=%s", count, row.Action), At: s.clock.Now()}.Audit()
	entry.ID = s.ids.NextFor("audit")
	return s.store.AppendAudit(entry)
}

func (s *Service) ImportSingle(guideID, actor string, row store.ImportRow) (Report, error) {
	return s.ImportRows(guideID, actor, []store.ImportRow{row})
}

func (s *Service) Summarize(report Report) string {
	return fmt.Sprintf("processed=%d created=%d updated=%d rejected=%d errors=%d", report.Processed, report.Created, report.Updated, report.Rejected, len(report.Errors))
}
