package web

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestWorkflowThree(t *testing.T) {
	app, repository := testApp(t)
	form := url.Values{"count": {"3"}, "actor": {"planner"}}
	request := httptest.NewRequest(http.MethodPost, "/import/demo-guide", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)
	if response.Code != http.StatusCreated || !strings.Contains(response.Body.String(), "Processed") {
		t.Fatalf("import status %d body %s", response.Code, response.Body.String())
	}
	for index := 0; index < 3; index++ {
		key := "guest-0" + string(rune('0'+index))
		if _, err := repository.ReadVisitor("demo-guide", key); err != nil {
			t.Fatalf("import did not persist %s: %v", key, err)
		}
	}
	audits, err := repository.ListAudits("demo-guide")
	if err != nil || len(audits) < 5 {
		t.Fatalf("workflow three audit trail incomplete: %d %v", len(audits), err)
	}
}
