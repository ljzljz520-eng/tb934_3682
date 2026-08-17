package web

import (
	"html/template"
	"strings"

	"weddingguide/internal/domain"
)

type displayHelper struct{}

func displayTitle(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "Wedding guide"
	}
	return value
}

func displayStatus(published bool) string {
	if published {
		return "Published"
	}
	return "Draft"
}

func displayVisitCount(count int) string {
	if count == 0 {
		return "Not visited yet"
	}
	if count == 1 {
		return "Visited once"
	}
	return "Visited " + integerText(count) + " times"
}

func integerText(value int) string {
	if value < 0 {
		return "0"
	}
	if value == 0 {
		return "0"
	}
	digits := make([]byte, 0, 10)
	for value > 0 {
		digits = append(digits, byte('0'+value%10))
		value /= 10
	}
	for left, right := 0, len(digits)-1; left < right; left, right = left+1, right-1 {
		digits[left], digits[right] = digits[right], digits[left]
	}
	return string(digits)
}

func actionClass(kind string) string {
	switch strings.ToLower(strings.TrimSpace(kind)) {
	case "navigation":
		return "button primary"
	case "seating":
		return "button secondary"
	case "blessing":
		return "button secondary"
	default:
		return "button secondary"
	}
}

func actionLabel(link domain.ActionLink) string {
	label := strings.TrimSpace(link.Label)
	if label == "" {
		label = "Open details"
	}
	return label
}

func addressLines(address domain.VenueAddress) []string {
	lines := make([]string, 0, 4)
	for _, value := range []string{address.Line1, address.Line2, address.City, address.Region} {
		if strings.TrimSpace(value) != "" {
			lines = append(lines, strings.TrimSpace(value))
		}
	}
	if address.PostalCode != "" {
		lines = append(lines, address.PostalCode)
	}
	if address.Country != "" {
		lines = append(lines, address.Country)
	}
	return lines
}

func linkTarget(link domain.ActionLink, visitor string) string {
	if link.Kind == "blessing" {
		return link.URL + "?visitor=" + template.URLQueryEscaper(visitor)
	}
	return link.URL
}

func auditDescription(entry domain.AuditEntry) string {
	parts := []string{entry.Actor, entry.Action}
	if entry.EntityID != "" {
		parts = append(parts, entry.EntityID)
	}
	return strings.Join(parts, " · ")
}

func truncate(value string, max int) string {
	value = strings.TrimSpace(value)
	if max < 1 || len([]rune(value)) <= max {
		return value
	}
	runes := []rune(value)
	return string(runes[:max-1]) + "..."
}
