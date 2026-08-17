package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"weddingguide/internal/domain"
	"weddingguide/internal/store"
)

func (a *App) handleGuide(writer http.ResponseWriter, request *http.Request, segments []string) {
	if len(segments) < 2 {
		http.NotFound(writer, request)
		return
	}
	guideID := segments[1]
	visitor := request.URL.Query().Get("visitor")
	if visitor == "" {
		visitor = request.Header.Get(a.config.VisitorHeader)
	}
	if visitor == "" {
		visitor = "guest"
	}
	result, err := a.guest.Visit(guideID, visitor, "view")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}
	data := map[string]any{"Guide": result.Guide, "Visitor": result.Visitor, "Hint": result.GentleHint, "AtLimit": result.AtLimit, "VisitorKey": visitor}
	a.render(writer, "guide", data)
}

func (a *App) handleAdmin(writer http.ResponseWriter, request *http.Request, segments []string) {
	if len(segments) < 2 {
		http.NotFound(writer, request)
		return
	}
	guideID := segments[1]
	if request.Method == http.MethodPost {
		if err := request.ParseForm(); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		guide, err := a.preview.Draft(guideID)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		}
		guide.Title = request.FormValue("title")
		guide.Welcome = request.FormValue("welcome")
		guide.Attire.Summary = request.FormValue("attire")
		if _, err := a.editor.UpdateDraft(guide, request.FormValue("actor")); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if request.FormValue("publish") == "yes" {
			if _, err := a.editor.Publish(guideID, request.FormValue("actor")); err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
		}
		http.Redirect(writer, request, "/admin/"+guideID, http.StatusSeeOther)
		return
	}
	model, err := a.preview.RenderModel(guideID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}
	a.render(writer, "admin", model)
}

func (a *App) handleAction(writer http.ResponseWriter, request *http.Request, segments []string) {
	if len(segments) < 3 {
		http.NotFound(writer, request)
		return
	}
	guideID, linkID := segments[1], segments[2]
	visitor := request.URL.Query().Get("visitor")
	if visitor == "" {
		visitor = "guest"
	}
	if _, err := a.guest.Visit(guideID, visitor, "action:"+linkID); err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}
	link, err := a.guest.ResolveAction(guideID, linkID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}
	http.Redirect(writer, request, link.URL, http.StatusFound)
}

func (a *App) handleBlessing(writer http.ResponseWriter, request *http.Request, segments []string) {
	if len(segments) < 2 {
		http.NotFound(writer, request)
		return
	}
	guideID := segments[1]
	if request.Method == http.MethodPost {
		if err := request.ParseForm(); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		blessing, err := a.guest.SubmitBlessing(guideID, request.FormValue("visitor"), request.FormValue("name"), request.FormValue("message"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(writer, http.StatusCreated, blessing)
		return
	}
	a.render(writer, "blessing", map[string]any{"GuideID": guideID})
}

func (a *App) handleAudits(writer http.ResponseWriter, request *http.Request, segments []string) {
	if len(segments) < 2 {
		http.NotFound(writer, request)
		return
	}
	term := request.URL.Query().Get("q")
	entries, err := a.audits.Search(segments[1], term)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}
	a.render(writer, "audits", map[string]any{"GuideID": segments[1], "Entries": entries, "Term": term})
}

func (a *App) handleImport(writer http.ResponseWriter, request *http.Request, segments []string) {
	if request.Method != http.MethodPost || len(segments) < 2 {
		http.Error(writer, "POST /import/{guide}", http.StatusMethodNotAllowed)
		return
	}
	if err := request.ParseForm(); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	count, err := strconv.Atoi(request.FormValue("count"))
	if err != nil || count < 1 || count > 100 {
		http.Error(writer, "count must be between 1 and 100", http.StatusBadRequest)
		return
	}
	rows := make([]store.ImportRow, 0, count)
	for index := 0; index < count; index++ {
		rows = append(rows, store.ImportRow{VisitorKey: fmt.Sprintf("guest-%02d", index), Action: "import", SeenAt: "2026-01-02T03:04:05Z"})
	}
	report, importErr := a.imports.ImportRows(segments[1], request.FormValue("actor"), rows)
	if importErr != nil {
		writeJSON(writer, http.StatusInternalServerError, map[string]any{"report": report, "error": importErr.Error()})
		return
	}
	writeJSON(writer, http.StatusCreated, report)
}

func (a *App) render(writer http.ResponseWriter, name string, data any) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := a.templates.ExecuteTemplate(writer, name, data); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(value)
}

func linkURL(link domain.ActionLink, visitor string) string {
	if link.Kind == "blessing" {
		return link.URL + "?visitor=" + template.URLQueryEscaper(visitor)
	}
	return link.URL
}
