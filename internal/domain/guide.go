package domain

import (
	"errors"
	"strings"
)

var (
	ErrInvalidGuide  = errors.New("invalid wedding guide")
	ErrMissingTitle  = errors.New("guide title is required")
	ErrMissingCouple = errors.New("couple name is required")
	ErrNoSchedule    = errors.New("at least one schedule item is required")
)

type WeddingGuide struct {
	ID         string
	Title      string
	Couple     string
	Welcome    string
	Schedule   []ScheduleItem
	Venue      VenueAddress
	Attire     AttireTip
	Links      []ActionLink
	VisitLimit int
	Published  bool
	Revision   int
}

type ScheduleItem struct {
	ID          string
	Title       string
	Details     string
	StartsAt    string
	EndsAt      string
	Location    string
	DisplayRank int
}

type VenueAddress struct {
	Name       string
	Line1      string
	Line2      string
	City       string
	Region     string
	PostalCode string
	Country    string
	Latitude   string
	Longitude  string
}

type AttireTip struct {
	Summary     string
	Description string
	ColorHint   string
	WeatherNote string
}

type ActionLink struct {
	ID          string
	Label       string
	URL         string
	Kind        string
	DisplayRank int
}

func (g WeddingGuide) Validate() error {
	if strings.TrimSpace(g.ID) == "" {
		return ErrInvalidGuide
	}
	if strings.TrimSpace(g.Title) == "" {
		return ErrMissingTitle
	}
	if strings.TrimSpace(g.Couple) == "" {
		return ErrMissingCouple
	}
	if len(g.Schedule) == 0 {
		return ErrNoSchedule
	}
	if g.VisitLimit < 0 {
		return ErrInvalidGuide
	}
	if err := g.Venue.Validate(); err != nil {
		return err
	}
	if err := g.Attire.Validate(); err != nil {
		return err
	}
	for _, item := range g.Schedule {
		if err := item.Validate(); err != nil {
			return err
		}
	}
	for _, link := range g.Links {
		if err := link.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (s ScheduleItem) Validate() error {
	if strings.TrimSpace(s.ID) == "" || strings.TrimSpace(s.Title) == "" {
		return ErrInvalidGuide
	}
	if strings.TrimSpace(s.StartsAt) == "" {
		return ErrInvalidGuide
	}
	return nil
}

func (v VenueAddress) Validate() error {
	if strings.TrimSpace(v.Name) == "" || strings.TrimSpace(v.Line1) == "" {
		return ErrInvalidGuide
	}
	if strings.TrimSpace(v.City) == "" || strings.TrimSpace(v.Country) == "" {
		return ErrInvalidGuide
	}
	return nil
}

func (a AttireTip) Validate() error {
	if strings.TrimSpace(a.Summary) == "" {
		return ErrInvalidGuide
	}
	return nil
}

func (l ActionLink) Validate() error {
	if strings.TrimSpace(l.ID) == "" || strings.TrimSpace(l.Label) == "" {
		return ErrInvalidGuide
	}
	if strings.TrimSpace(l.URL) == "" || strings.TrimSpace(l.Kind) == "" {
		return ErrInvalidGuide
	}
	return nil
}

func (g WeddingGuide) FindLink(id string) (ActionLink, bool) {
	for _, link := range g.Links {
		if link.ID == id {
			return link, true
		}
	}
	return ActionLink{}, false
}

func (g WeddingGuide) FindSchedule(id string) (ScheduleItem, bool) {
	for _, item := range g.Schedule {
		if item.ID == id {
			return item, true
		}
	}
	return ScheduleItem{}, false
}

func (g WeddingGuide) Clone() WeddingGuide {
	copyGuide := g
	copyGuide.Schedule = append([]ScheduleItem(nil), g.Schedule...)
	copyGuide.Links = append([]ActionLink(nil), g.Links...)
	return copyGuide
}
