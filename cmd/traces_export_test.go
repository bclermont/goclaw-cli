package cmd

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func resetTracesExportFlags(t *testing.T) {
	t.Helper()
	resetTestFlag(tracesExportCmd, "output", "")
}

func TestTracesExport_RejectsMalformedIDBeforeHTTP(t *testing.T) {
	t.Cleanup(func() { resetTracesExportFlags(t) })
	var calls int64
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&calls, 1)
	}))
	defer srv.Close()
	t.Setenv("GOCLAW_SERVER", srv.URL)
	t.Setenv("GOCLAW_TOKEN", "test-token")

	err := runCmd(t, "traces", "export", "../bad")
	if err == nil {
		t.Fatal("expected validation error")
	}
	if atomic.LoadInt64(&calls) != 0 {
		t.Fatalf("malformed trace id made %d HTTP calls", calls)
	}
	if !strings.Contains(err.Error(), "trace id") {
		t.Fatalf("error should mention trace id: %v", err)
	}
}

func TestTracesExport_UsesSafePathAndOutputFile(t *testing.T) {
	t.Cleanup(func() { resetTracesExportFlags(t) })
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.EscapedPath()
		w.Header().Set("Content-Type", "application/gzip")
		_, _ = w.Write([]byte("gzip-bytes"))
	}))
	defer srv.Close()
	t.Setenv("GOCLAW_SERVER", srv.URL)
	t.Setenv("GOCLAW_TOKEN", "test-token")

	outFile := t.TempDir() + "/trace.json.gz"
	if err := runCmd(t, "traces", "export", "trace_FIXTURE.001", "--output", outFile); err != nil {
		t.Fatalf("traces export: %v", err)
	}
	if gotPath != "/v1/traces/trace_FIXTURE.001/export" {
		t.Fatalf("path = %q", gotPath)
	}
}
