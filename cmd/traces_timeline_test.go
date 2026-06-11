package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
)

func resetTracesTimelineFlags(t *testing.T) {
	t.Helper()
	for _, name := range []string{"session-key"} {
		resetTestFlag(tracesTimelineCmd, name, "")
	}
	resetTestFlag(tracesTimelineCmd, "limit", "0")
	resetTestFlag(tracesTimelineCmd, "offset", "0")
}

func TestTracesTimeline_BuildsPathAndQuery(t *testing.T) {
	t.Cleanup(func() { resetTracesTimelineFlags(t) })
	var gotPath string
	var query url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.EscapedPath()
		query = r.URL.Query()
		rawJSON(t, w, timelineFixture())
	}))
	defer srv.Close()
	t.Setenv("GOCLAW_SERVER", srv.URL)
	t.Setenv("GOCLAW_TOKEN", "test-token")
	t.Setenv("GOCLAW_OUTPUT", "json")

	if err := runCmd(t, "traces", "timeline", "run_FIXTURE.001", "--session-key=session-1", "--limit=25", "--offset=5"); err != nil {
		t.Fatalf("traces timeline: %v", err)
	}
	if gotPath != "/v1/runs/run_FIXTURE.001/timeline" {
		t.Fatalf("path = %q", gotPath)
	}
	if query.Get("session_key") != "session-1" || query.Get("limit") != "25" || query.Get("offset") != "5" {
		t.Fatalf("query = %s", query.Encode())
	}
}

func TestTracesTimeline_JSONPreservesEnvelope(t *testing.T) {
	t.Cleanup(func() { resetTracesTimelineFlags(t) })
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawJSON(t, w, timelineFixture())
	}))
	defer srv.Close()
	t.Setenv("GOCLAW_SERVER", srv.URL)
	t.Setenv("GOCLAW_TOKEN", "test-token")
	t.Setenv("GOCLAW_OUTPUT", "json")

	out, err := captureStdout(t, func() error {
		return runCmd(t, "traces", "timeline", "run_FIXTURE.001", "--output", "json")
	})
	if err != nil {
		t.Fatalf("traces timeline: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("stdout is not JSON: %v\n%s", err, out)
	}
	if got["run_id"] != "run_FIXTURE.001" || got["limit"] != float64(25) {
		t.Fatalf("timeline envelope not preserved: %#v", got)
	}
	if _, ok := got["items"].([]any); !ok {
		t.Fatalf("items array not preserved: %#v", got["items"])
	}
}

func TestTracesTimeline_TableRendersItems(t *testing.T) {
	t.Cleanup(func() { resetTracesTimelineFlags(t) })
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawJSON(t, w, timelineFixture())
	}))
	defer srv.Close()
	t.Setenv("GOCLAW_SERVER", srv.URL)
	t.Setenv("GOCLAW_TOKEN", "test-token")
	t.Setenv("GOCLAW_OUTPUT", "table")

	out, err := captureStdout(t, func() error {
		return runCmd(t, "traces", "timeline", "run_FIXTURE.001", "--output", "table")
	})
	if err != nil {
		t.Fatalf("traces timeline: %v", err)
	}
	headerRE := regexp.MustCompile(`SEQ.*TYPE.*STATUS.*TITLE.*TOOL.*TRACE_ID.*SPAN_ID.*CREATED_AT`)
	if !headerRE.MatchString(out) {
		t.Fatalf("table headers missing in:\n%s", out)
	}
	for _, want := range []string{"2", "tool.call", "running", "web_fetch", "trace-1", "span-1"} {
		if !strings.Contains(out, want) {
			t.Fatalf("timeline table missing %q:\n%s", want, out)
		}
	}
}

func TestTracesTimeline_RejectsMalformedRunIDBeforeHTTP(t *testing.T) {
	t.Cleanup(func() { resetTracesTimelineFlags(t) })
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("malformed run id should not make an HTTP request")
	}))
	defer srv.Close()
	t.Setenv("GOCLAW_SERVER", srv.URL)
	t.Setenv("GOCLAW_TOKEN", "test-token")

	err := runCmd(t, "traces", "timeline", "../run")
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestTracesTimeline_CommandRegistered(t *testing.T) {
	for _, c := range tracesCmd.Commands() {
		if c.Name() == "timeline" {
			return
		}
	}
	t.Fatal("traces timeline command is not registered")
}

func timelineFixture() map[string]any {
	return map[string]any{
		"run_id":      "run_FIXTURE.001",
		"session_key": "session-1",
		"items": []map[string]any{{
			"id":           "item-1",
			"run_id":       "run_FIXTURE.001",
			"session_key":  "session-1",
			"seq":          2,
			"item_type":    "tool.call",
			"status":       "running",
			"title":        "web_fetch",
			"tool_name":    "web_fetch",
			"tool_call_id": "call-1",
			"trace_id":     "trace-1",
			"span_id":      "span-1",
			"created_at":   "2026-05-29T10:00:00Z",
		}},
		"limit":  25,
		"offset": 0,
	}
}
