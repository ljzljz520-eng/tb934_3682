package store

import (
	"fmt"

	"weddingguide/internal/domain"
)

type ImportRow struct {
	VisitorKey string
	Action     string
	SeenAt     string
}

type ImportResult struct {
	Processed int
	Created   int
	Updated   int
	Rejected  int
	Errors    []error
}

func (s *Store) ImportVisitor(row ImportRow, guideID string) (domain.VisitorRecord, bool, error) {
	visitor, err := s.ReadVisitor(guideID, row.VisitorKey)
	created := false
	if err != nil {
		if err != ErrNotFound {
			return domain.VisitorRecord{}, false, err
		}
		visitor = domain.VisitorRecord{ID: fmt.Sprintf("visitor-%s-%s", guideID, row.VisitorKey), GuideID: guideID, VisitorKey: row.VisitorKey}
		created = true
	}
	visitor = visitor.Increment(row.Action, row.SeenAt)
	if err := s.SaveVisitor(visitor); err != nil {
		return domain.VisitorRecord{}, created, err
	}
	return visitor, created, nil
}
