package cmd

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nextlevelbuilder/goclaw-cli/internal/config"
	"github.com/nextlevelbuilder/goclaw-cli/internal/output"
)

func TestStorageListSubpathUsesQueryRoute(t *testing.T) {
	var gotPath, gotQueryPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQueryPath = r.URL.Query().Get("path")
		if r.Method != http.MethodGet || gotPath != "/v1/storage/files" {
			t.Errorf("unexpected request: %s %s", r.Method, gotPath)
		}
		okJSON(t, w, []map[string]any{{"name": "child.txt", "size": 12}})
	}))
	defer srv.Close()

	cfg = &config.Config{Server: srv.URL, Token: "test-token", OutputFormat: "json"}
	printer = output.NewPrinter("json")
	flag := storageListCmd.Flags().Lookup("path")
	_ = flag.Value.Set(flag.DefValue)
	flag.Changed = false
	_ = storageListCmd.Flags().Set("path", "tenants/team alpha")

	if err := storageListCmd.RunE(storageListCmd, nil); err != nil {
		t.Fatalf("storage list: %v", err)
	}
	if gotQueryPath != "tenants/team alpha" {
		t.Fatalf("path query = %q", gotQueryPath)
	}
}
