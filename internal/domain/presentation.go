package domain

import (
	"sort"
	"strings"
)

type GuideView struct {
	ID             string
	Title          string
	Couple         string
	Welcome        string
	Schedule       []ScheduleView
	Venue          VenueView
	Attire         AttireView
	Actions        []ActionSummary
	PublishedLabel string
	RevisionLabel  string
}

type ScheduleView struct {
	ID          string
	Title       string
	Details     string
	TimeLabel   string
	Location    string
	Rank        int
	IsMilestone bool
}

type VenueView struct {
	Name        string
	AddressLine string
	Locality    string
	MapLabel    string
}

type AttireView struct {
	Summary     string
	Description string
	Palette     string
	Weather     string
}

type ActionSummary struct {
	ID         string
	Label      string
	Kind       string
	URL        string
	Emphasis   string
	Accessible string
}

func BuildGuideView(guide WeddingGuide) GuideView {
	view := GuideView{ID: guide.ID, Title: guide.Title, Couple: guide.Couple, Welcome: guide.Welcome, Venue: buildVenueView(guide.Venue), Attire: buildAttireView(guide.Attire), PublishedLabel: publicationLabel(guide.Published), RevisionLabel: revisionLabel(guide.Revision)}
	view.Schedule = buildScheduleViews(guide.Schedule)
	view.Actions = buildActionSummaries(guide.Links)
	return view
}

func buildVenueView(venue VenueAddress) VenueView {
	lines := []string{venue.Line1}
	if strings.TrimSpace(venue.Line2) != "" {
		lines = append(lines, venue.Line2)
	}
	locality := strings.TrimSpace(strings.Join([]string{venue.City, venue.Region, venue.PostalCode}, " "))
	return VenueView{Name: venue.Name, AddressLine: strings.Join(lines, ", "), Locality: locality, MapLabel: venue.Latitude + ", " + venue.Longitude}
}

func buildAttireView(attire AttireTip) AttireView {
	palette := strings.TrimSpace(attire.ColorHint)
	if palette == "" {
		palette = "Choose something that feels like you"
	}
	weather := strings.TrimSpace(attire.WeatherNote)
	if weather == "" {
		weather = "A light layer is a thoughtful choice"
	}
	return AttireView{Summary: attire.Summary, Description: attire.Description, Palette: palette, Weather: weather}
}

func buildScheduleViews(items []ScheduleItem) []ScheduleView {
	ordered := append([]ScheduleItem(nil), items...)
	sort.SliceStable(ordered, func(i, j int) bool {
		if ordered[i].DisplayRank == ordered[j].DisplayRank {
			return ordered[i].StartsAt < ordered[j].StartsAt
		}
		return ordered[i].DisplayRank < ordered[j].DisplayRank
	})
	views := make([]ScheduleView, 0, len(ordered))
	for index, item := range ordered {
		views = append(views, ScheduleView{ID: item.ID, Title: item.Title, Details: item.Details, TimeLabel: timeLabel(item), Location: item.Location, Rank: index + 1, IsMilestone: index == 0 || strings.Contains(strings.ToLower(item.Title), "ceremony")})
	}
	return views
}

func buildActionSummaries(links []ActionLink) []ActionSummary {
	ordered := append([]ActionLink(nil), links...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].DisplayRank < ordered[j].DisplayRank })
	summaries := make([]ActionSummary, 0, len(ordered))
	for _, link := range ordered {
		emphasis := "secondary"
		if link.Kind == "navigation" {
			emphasis = "primary"
		}
		accessible := link.Label
		if link.Kind != "" {
			accessible = link.Label + ", " + link.Kind
		}
		summaries = append(summaries, ActionSummary{ID: link.ID, Label: link.Label, Kind: link.Kind, URL: link.URL, Emphasis: emphasis, Accessible: accessible})
	}
	return summaries
}

func timeLabel(item ScheduleItem) string {
	start := strings.TrimSpace(item.StartsAt)
	end := strings.TrimSpace(item.EndsAt)
	if end == "" {
		return start
	}
	return start + " - " + end
}

func publicationLabel(published bool) string {
	if published {
		return "Published for guests"
	}
	return "Private draft"
}

func revisionLabel(revision int) string {
	if revision < 1 {
		revision = 1
	}
	return "Revision " + formatInt(revision)
}

func formatInt(value int) string {
	if value == 0 {
		return "0"
	}
	negative := value < 0
	if negative {
		value = -value
	}
	digits := make([]byte, 0, 12)
	for value > 0 {
		digits = append(digits, byte('0'+value%10))
		value /= 10
	}
	for left, right := 0, len(digits)-1; left < right; left, right = left+1, right-1 {
		digits[left], digits[right] = digits[right], digits[left]
	}
	if negative {
		return "-" + string(digits)
	}
	return string(digits)
}

func ValidateScheduleOrder(items []ScheduleItem) bool {
	lastRank := 0
	for _, item := range items {
		if item.DisplayRank < lastRank {
			return false
		}
		lastRank = item.DisplayRank
	}
	return true
}

func VisibleActions(guide WeddingGuide, published bool) []ActionLink {
	if !published && guide.Published {
		published = true
	}
	if !published {
		return nil
	}
	visible := make([]ActionLink, 0, len(guide.Links))
	for _, link := range guide.Links {
		if strings.TrimSpace(link.URL) == "" || strings.TrimSpace(link.Label) == "" {
			continue
		}
		visible = append(visible, link)
	}
	return visible
}

func ActionKinds(links []ActionLink) map[string]int {
	counts := make(map[string]int)
	for _, link := range links {
		kind := strings.TrimSpace(link.Kind)
		if kind == "" {
			kind = "other"
		}
		counts[kind]++
	}
	return counts
}

func HasRequiredActions(links []ActionLink) bool {
	seen := make(map[string]bool)
	for _, link := range links {
		seen[link.Kind] = true
	}
	return seen["navigation"] && seen["seating"] && seen["blessing"]
}
