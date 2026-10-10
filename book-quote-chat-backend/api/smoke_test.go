package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSmokePublicQuoteEndpoints(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })

	req := httptest.NewRequest(http.MethodGet, "/api/quotes", nil)
	rec := httptest.NewRecorder()
	HandleGetQuotes(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/quotes status %d body %s", rec.Code, rec.Body.String())
	}
	var quotes []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &quotes); err != nil {
		t.Fatalf("GET /api/quotes json: %v", err)
	}
	if len(quotes) == 0 {
		t.Fatal("GET /api/quotes returned no quotes")
	}

	req = httptest.NewRequest(http.MethodGet, "/api/friend-links", nil)
	rec = httptest.NewRecorder()
	HandleGetLinks(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/friend-links status %d body %s", rec.Code, rec.Body.String())
	}
	var links []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &links); err != nil {
		t.Fatalf("GET /api/friend-links json: %v", err)
	}
}
