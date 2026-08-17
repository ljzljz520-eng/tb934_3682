package web

import (
	"html/template"
	"net/http"

	"weddingguide/internal/audit"
	"weddingguide/internal/config"
	"weddingguide/internal/fixture"
	"weddingguide/internal/importer"
	"weddingguide/internal/service"
	"weddingguide/internal/store"
)

type App struct {
	config    config.Config
	store     *store.Store
	editor    *service.GuideEditor
	guest     *service.GuestService
	preview   *service.PreviewService
	audits    *audit.Query
	imports   *importer.Service
	templates *template.Template
}

func New(cfg config.Config, repository *store.Store, editor *service.GuideEditor, guest *service.GuestService, preview *service.PreviewService, auditQuery *audit.Query, importService *importer.Service) *App {
	tmpl := template.Must(template.New("pages").Funcs(template.FuncMap{"inc": func(value int) int { return value + 1 }, "displayTitle": displayTitle, "displayStatus": displayStatus, "displayVisits": displayVisitCount, "actionClass": actionClass, "actionLabel": actionLabel, "addressLines": addressLines, "linkTarget": linkTarget, "auditDescription": auditDescription, "truncate": truncate}).Parse(templates))
	return &App{config: cfg, store: repository, editor: editor, guest: guest, preview: preview, audits: auditQuery, imports: importService, templates: tmpl}
}

func SeedDemo(app *App) error {
	guide := fixture.Wedding()
	if _, err := app.editor.CreateDraft(guide, "seed"); err != nil {
		if guide, readErr := app.editor.Get(guide.ID); readErr == nil && guide.ID != "" {
			return nil
		}
		return err
	}
	_, err := app.editor.Publish(guide.ID, "seed")
	return err
}

func (a *App) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path == "/health" {
		writeText(writer, http.StatusOK, "ok")
		return
	}
	if request.URL.Path == "/" {
		http.Redirect(writer, request, "/guide/"+a.config.DefaultGuide, http.StatusFound)
		return
	}
	segments := splitPath(request.URL.Path)
	if len(segments) == 0 {
		http.NotFound(writer, request)
		return
	}
	switch segments[0] {
	case "guide":
		a.handleGuide(writer, request, segments)
	case "admin":
		a.handleAdmin(writer, request, segments)
	case "action":
		a.handleAction(writer, request, segments)
	case "bless":
		a.handleBlessing(writer, request, segments)
	case "audits":
		a.handleAudits(writer, request, segments)
	case "import":
		a.handleImport(writer, request, segments)
	default:
		http.NotFound(writer, request)
	}
}

func splitPath(path string) []string {
	trimmed := path
	for len(trimmed) > 0 && trimmed[0] == '/' {
		trimmed = trimmed[1:]
	}
	for len(trimmed) > 0 && trimmed[len(trimmed)-1] == '/' {
		trimmed = trimmed[:len(trimmed)-1]
	}
	if trimmed == "" {
		return nil
	}
	parts := make([]string, 0, 4)
	for _, part := range stringsSplit(trimmed, '/') {
		if part != "" {
			parts = append(parts, part)
		}
	}
	return parts
}

func stringsSplit(value string, separator rune) []string {
	result := make([]string, 0, 4)
	start := 0
	for index, character := range value {
		if character == separator {
			result = append(result, value[start:index])
			start = index + 1
		}
	}
	return append(result, value[start:])
}

func writeText(writer http.ResponseWriter, status int, value string) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(status)
	_, _ = writer.Write([]byte(value))
}
