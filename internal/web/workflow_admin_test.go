package web

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestWorkflowOne(t *testing.T) {
	app, repository := testApp(t)
	form := url.Values{"title": {"A brighter day"}, "welcome": {"Welcome to the river garden"}, "attire": {"Garden formal"}, "actor": {"planner"}, "publish": {"yes"}}
	request := httptest.NewRequest(http.MethodPost, "/admin/demo-guide", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)
	if response.Code != http.StatusSeeOther {
		t.Fatalf("save status %d", response.Code)
	}
	page := httptest.NewRecorder()
	app.ServeHTTP(page, httptest.NewRequest(http.MethodGet, "/guide/demo-guide?visitor=workflow-one", nil))
	body, _ := io.ReadAll(page.Result().Body)
	if page.Code != http.StatusOK || !strings.Contains(string(body), "A brighter day") {
		t.Fatalf("guest page missing updated title: %d", page.Code)
	}
	audits, err := repository.ListAudits("demo-guide")
	if err != nil || len(audits) < 4 {
		t.Fatalf("workflow one audit trail incomplete: %d %v", len(audits), err)
	}
}
