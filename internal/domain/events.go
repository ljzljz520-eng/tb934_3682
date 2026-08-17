package domain

type EventName string

const (
	EventGuideCreated      EventName = "guide.created"
	EventGuideUpdated      EventName = "guide.updated"
	EventGuidePublished    EventName = "guide.published"
	EventGuestVisited      EventName = "guest.visited"
	EventBlessingSubmitted EventName = "blessing.submitted"
	EventImportCompleted   EventName = "import.completed"
	EventImportRejected    EventName = "import.rejected"
)

type DomainEvent struct {
	Name     EventName
	GuideID  string
	Actor    string
	Entity   string
	EntityID string
	Detail   string
	At       string
}

func (e DomainEvent) Audit() AuditEntry {
	return AuditEntry{
		ID:        e.At + ":" + string(e.Name) + ":" + e.EntityID,
		GuideID:   e.GuideID,
		Actor:     e.Actor,
		Action:    string(e.Name),
		Entity:    e.Entity,
		EntityID:  e.EntityID,
		Detail:    e.Detail,
		CreatedAt: e.At,
	}
}
