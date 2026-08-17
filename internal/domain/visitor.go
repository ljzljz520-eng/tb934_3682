package domain

import (
	"errors"
	"strings"
)

var (
	ErrMissingVisitor = errors.New("visitor id is required")
	ErrMissingGuide   = errors.New("guide id is required")
	ErrInvalidVisit   = errors.New("visit count must be positive")
)

type VisitorRecord struct {
	ID         string
	GuideID    string
	VisitorKey string
	VisitCount int
	LastAction string
	LastSeenAt string
	EditorView bool
}

type Blessing struct {
	ID         string
	GuideID    string
	VisitorKey string
	Name       string
	Message    string
	CreatedAt  string
	Approved   bool
}

type AuditEntry struct {
	ID        string
	GuideID   string
	Actor     string
	Action    string
	Entity    string
	EntityID  string
	Detail    string
	CreatedAt string
}

func (v VisitorRecord) Validate() error {
	if strings.TrimSpace(v.ID) == "" {
		return ErrMissingVisitor
	}
	if strings.TrimSpace(v.GuideID) == "" {
		return ErrMissingGuide
	}
	if strings.TrimSpace(v.VisitorKey) == "" {
		return ErrMissingVisitor
	}
	if v.VisitCount < 0 {
		return ErrInvalidVisit
	}
	return nil
}

func (b Blessing) Validate() error {
	if strings.TrimSpace(b.ID) == "" || strings.TrimSpace(b.GuideID) == "" {
		return ErrInvalidGuide
	}
	if strings.TrimSpace(b.Name) == "" || strings.TrimSpace(b.Message) == "" {
		return ErrInvalidGuide
	}
	return nil
}

func (a AuditEntry) Validate() error {
	if strings.TrimSpace(a.ID) == "" || strings.TrimSpace(a.GuideID) == "" {
		return ErrInvalidGuide
	}
	if strings.TrimSpace(a.Action) == "" || strings.TrimSpace(a.Entity) == "" {
		return ErrInvalidGuide
	}
	if strings.TrimSpace(a.CreatedAt) == "" {
		return ErrInvalidGuide
	}
	return nil
}

func (v VisitorRecord) Reached(limit int) bool {
	return limit > 0 && v.VisitCount >= limit
}

func (v VisitorRecord) Increment(action, seenAt string) VisitorRecord {
	v.VisitCount++
	v.LastAction = action
	v.LastSeenAt = seenAt
	return v
}

func (b Blessing) Approve() Blessing {
	b.Approved = true
	return b
}
