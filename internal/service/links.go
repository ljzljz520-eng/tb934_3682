package service

import (
	"net/url"
	"sort"
	"strings"

	"weddingguide/internal/domain"
)

type LinkService struct{}

type LinkCheck struct {
	Link       domain.ActionLink
	Valid      bool
	Internal   bool
	Reason     string
	TargetPath string
}

func NewLinkService() *LinkService {
	return &LinkService{}
}

func (l *LinkService) Check(link domain.ActionLink) LinkCheck {
	check := LinkCheck{Link: link}
	if strings.TrimSpace(link.ID) == "" {
		check.Reason = "missing id"
		return check
	}
	if strings.TrimSpace(link.Label) == "" {
		check.Reason = "missing label"
		return check
	}
	parsed, err := url.Parse(link.URL)
	if err != nil || parsed.Scheme == "" && !strings.HasPrefix(parsed.Path, "/") {
		check.Reason = "invalid target"
		return check
	}
	check.Valid = true
	check.Internal = parsed.IsAbs() == false
	check.TargetPath = parsed.Path
	if check.TargetPath == "" {
		check.TargetPath = "/"
	}
	return check
}

func (l *LinkService) Validate(links []domain.ActionLink) []LinkCheck {
	checks := make([]LinkCheck, 0, len(links))
	for _, link := range links {
		checks = append(checks, l.Check(link))
	}
	return checks
}

func (l *LinkService) AllValid(links []domain.ActionLink) bool {
	checks := l.Validate(links)
	if len(checks) == 0 {
		return false
	}
	for _, check := range checks {
		if !check.Valid {
			return false
		}
	}
	return true
}

func (l *LinkService) Ordered(links []domain.ActionLink) []domain.ActionLink {
	ordered := append([]domain.ActionLink(nil), links...)
	sort.SliceStable(ordered, func(left, right int) bool {
		if ordered[left].DisplayRank == ordered[right].DisplayRank {
			return ordered[left].Label < ordered[right].Label
		}
		return ordered[left].DisplayRank < ordered[right].DisplayRank
	})
	return ordered
}

func (l *LinkService) ByKind(links []domain.ActionLink, kind string) []domain.ActionLink {
	kind = strings.ToLower(strings.TrimSpace(kind))
	result := make([]domain.ActionLink, 0)
	for _, link := range links {
		if strings.ToLower(strings.TrimSpace(link.Kind)) == kind {
			result = append(result, link)
		}
	}
	return l.Ordered(result)
}

func (l *LinkService) DestinationMap(links []domain.ActionLink) map[string]string {
	result := make(map[string]string)
	for _, link := range links {
		check := l.Check(link)
		if check.Valid {
			result[link.ID] = check.TargetPath
		}
	}
	return result
}

func (l *LinkService) MissingKinds(links []domain.ActionLink) []string {
	seen := make(map[string]bool)
	for _, link := range links {
		seen[strings.ToLower(strings.TrimSpace(link.Kind))] = true
	}
	missing := make([]string, 0, 3)
	for _, required := range []string{"navigation", "seating", "blessing"} {
		if !seen[required] {
			missing = append(missing, required)
		}
	}
	return missing
}
