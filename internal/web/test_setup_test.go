package web

import (
	"path/filepath"
	"testing"

	"weddingguide/internal/audit"
	"weddingguide/internal/clock"
	"weddingguide/internal/config"
	"weddingguide/internal/fixture"
	"weddingguide/internal/ids"
	"weddingguide/internal/importer"
	"weddingguide/internal/service"
	"weddingguide/internal/store"
)

func testApp(t *testing.T) (*App, *store.Store) {
	t.Helper()
	repository, err := store.New(filepath.Join(t.TempDir(), "web.db"))
	if err != nil {
		t.Fatal(err)
	}
	now := clock.NewFixed("2026-01-02T03:04:05Z")
	generator := ids.New("web")
	editor := service.NewGuideEditor(repository, now, generator)
	guest := service.NewGuestService(repository, now, generator)
	preview := service.NewPreviewService(repository)
	audits := audit.NewQuery(service.NewAuditService(repository))
	imports := importer.New(repository, now, generator)
	app := New(config.Default(), repository, editor, guest, preview, audits, imports)
	if _, err := editor.CreateDraft(fixture.Wedding(), "seed"); err != nil {
		t.Fatal(err)
	}
	if _, err := editor.Publish("demo-guide", "seed"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { repository.Close() })
	return app, repository
}
