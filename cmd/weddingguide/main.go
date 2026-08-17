package main

import (
	"log"
	"net/http"
	"os"

	"weddingguide/internal/audit"
	"weddingguide/internal/clock"
	"weddingguide/internal/config"
	"weddingguide/internal/ids"
	"weddingguide/internal/importer"
	"weddingguide/internal/service"
	"weddingguide/internal/store"
	"weddingguide/internal/web"
)

func main() {
	cfg, err := config.FromArgs(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	repository, err := store.NewWithReaderLimit(cfg.DatabasePath, cfg.ReaderLimit)
	if err != nil {
		log.Fatal(err)
	}
	defer repository.Close()
	now := clock.NewFixed("2026-01-02T03:04:05Z")
	generator := ids.New("wg")
	editor := service.NewGuideEditor(repository, now, generator)
	guest := service.NewGuestService(repository, now, generator)
	preview := service.NewPreviewService(repository)
	audits := audit.NewQuery(service.NewAuditService(repository))
	imports := importer.New(repository, now, generator)
	app := web.New(cfg, repository, editor, guest, preview, audits, imports)
	if _, err := editor.Get(cfg.DefaultGuide); err != nil {
		if err := web.SeedDemo(app); err != nil {
			log.Fatal(err)
		}
	}
	server := &http.Server{Addr: cfg.Address, Handler: app}
	log.Printf("wedding guide listening on %s", cfg.Address)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
