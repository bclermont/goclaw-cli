package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nextlevelbuilder/goclaw-cli/internal/config"
	"github.com/nextlevelbuilder/goclaw-cli/internal/output"
	"github.com/spf13/cobra"
)

func setupWebhookTest(serverURL string) {
	cfg = &config.Config{Server: serverURL, Token: "test-token", OutputFormat: "json"}
	printer = output.NewPrinter("json")
	resetWebhookFlags(webhooksCreateCmd)
	resetWebhookFlags(webhooksUpdateCmd)
}

func resetWebhookFlags(cmd *cobra.Command) {
	for _, name := range []string{
		"body", "name", "kind", "agent-id", "channel-id", "scopes",
		"rate-limit-per-min", "ip-allowlist", "require-hmac", "localhost-only",
	} {
		flag := cmd.Flags().Lookup(name)
		if flag == nil {
			continue
		}
		_ = flag.Value.Set(flag.DefValue)
		flag.Changed = false
	}
}

func TestWebhooksCreateUsesServerContractPayload(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/webhooks" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		okJSON(t, w, map[string]any{"id": "webhook-1", "secret": "secret-once"})
	}))
	defer srv.Close()
	setupWebhookTest(srv.URL)

	_ = webhooksCreateCmd.Flags().Set("name", "agent hook")
	_ = webhooksCreateCmd.Flags().Set("kind", "llm")
	_ = webhooksCreateCmd.Flags().Set("agent-id", "11111111-1111-1111-1111-111111111111")
	_ = webhooksCreateCmd.Flags().Set("scopes", "chat.run,agent.read")
	_ = webhooksCreateCmd.Flags().Set("rate-limit-per-min", "30")
	_ = webhooksCreateCmd.Flags().Set("ip-allowlist", "127.0.0.1,10.0.0.0/8")
	_ = webhooksCreateCmd.Flags().Set("require-hmac", "true")

	if err := webhooksCreateCmd.RunE(webhooksCreateCmd, nil); err != nil {
		t.Fatalf("webhooks create: %v", err)
	}

	if gotBody["kind"] != "llm" || gotBody["agent_id"] != "11111111-1111-1111-1111-111111111111" {
		t.Fatalf("wrong server contract payload: %#v", gotBody)
	}
	if _, exists := gotBody["target_type"]; exists {
		t.Fatalf("legacy target_type must not be sent: %#v", gotBody)
	}
	if _, exists := gotBody["enabled"]; exists {
		t.Fatalf("legacy enabled must not be sent: %#v", gotBody)
	}
	if gotBody["require_hmac"] != true || gotBody["rate_limit_per_min"] != float64(30) {
		t.Fatalf("missing typed fields: %#v", gotBody)
	}
}

func TestWebhooksUpdateOnlySendsChangedFlags(t *testing.T) {
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch || r.URL.Path != "/v1/webhooks/hook-1" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		okJSON(t, w, map[string]any{"id": "hook-1", "name": "renamed"})
	}))
	defer srv.Close()
	setupWebhookTest(srv.URL)

	_ = webhooksUpdateCmd.Flags().Set("name", "renamed")

	if err := webhooksUpdateCmd.RunE(webhooksUpdateCmd, []string{"hook-1"}); err != nil {
		t.Fatalf("webhooks update: %v", err)
	}
	if len(gotBody) != 1 || gotBody["name"] != "renamed" {
		t.Fatalf("update must send only changed fields, got %#v", gotBody)
	}
}
