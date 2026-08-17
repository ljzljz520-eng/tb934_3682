package web

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestWorkflowTwo(t *testing.T) {
	app, repository := testApp(t)
	view := httptest.NewRecorder()
	app.ServeHTTP(view, httptest.NewRequest(http.MethodGet, "/guide/demo-guide?visitor=workflow-two", nil))
	if view.Code != http.StatusOK {
		t.Fatalf("view status %d", view.Code)
	}
	action := httptest.NewRecorder()
	app.ServeHTTP(action, httptest.NewRequest(http.MethodGet, "/action/demo-guide/seats?visitor=workflow-two", nil))
	if action.Code != http.StatusFound || action.Header().Get("Location") != "/seats/demo-guide" {
		t.Fatalf("seat action mismatch: %d %s", action.Code, action.Header().Get("Location"))
	}
	form := url.Values{"visitor": {"workflow-two"}, "name": {"Rae"}, "message": {"May your home be full of laughter"}}
	blessing := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/bless/demo-guide", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	app.ServeHTTP(blessing, request)
	if blessing.Code != http.StatusCreated || !strings.Contains(blessing.Body.String(), "Rae") {
		t.Fatalf("blessing status %d body %s", blessing.Code, blessing.Body.String())
	}
	audits, err := repository.ListAudits("demo-guide")
	if err != nil || len(audits) < 5 {
		t.Fatalf("workflow two audit trail incomplete: %d %v", len(audits), err)
	}
}
