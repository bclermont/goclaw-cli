package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/nextlevelbuilder/goclaw-cli/internal/tui"
	"github.com/spf13/cobra"
)

var webhooksCmd = &cobra.Command{Use: "webhooks", Short: "Manage inbound webhooks"}

var webhooksListCmd = &cobra.Command{
	Use:   "list",
	Short: "List webhooks",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newHTTP()
		if err != nil {
			return err
		}
		data, err := c.Get("/v1/webhooks")
		if err != nil {
			return err
		}
		printer.Print(unmarshalList(data))
		return nil
	},
}

var webhooksGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get webhook details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newHTTP()
		if err != nil {
			return err
		}
		data, err := c.Get("/v1/webhooks/" + url.PathEscape(args[0]))
		if err != nil {
			return err
		}
		printer.Print(unmarshalMap(data))
		return nil
	},
}

var webhooksCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a webhook",
	RunE: func(cmd *cobra.Command, args []string) error {
		body, err := webhookBody(cmd, false)
		if err != nil {
			return err
		}
		c, err := newHTTP()
		if err != nil {
			return err
		}
		data, err := c.Post("/v1/webhooks", body)
		if err != nil {
			return err
		}
		printWebhookSecretResult(unmarshalMap(data))
		return nil
	},
}

var webhooksUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a webhook",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		body, err := webhookBody(cmd, true)
		if err != nil {
			return err
		}
		c, err := newHTTP()
		if err != nil {
			return err
		}
		data, err := c.Patch("/v1/webhooks/"+url.PathEscape(args[0]), body)
		if err != nil {
			return err
		}
		printer.Print(unmarshalMap(data))
		return nil
	},
}

var webhooksRotateCmd = &cobra.Command{
	Use:   "rotate <id>",
	Short: "Rotate webhook secret",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !tui.Confirm("Rotate this webhook secret?", cfg.Yes) {
			return nil
		}
		c, err := newHTTP()
		if err != nil {
			return err
		}
		data, err := c.Post("/v1/webhooks/"+url.PathEscape(args[0])+"/rotate", nil)
		if err != nil {
			return err
		}
		printWebhookSecretResult(unmarshalMap(data))
		return nil
	},
}

var webhooksDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a webhook",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !tui.Confirm("Delete this webhook?", cfg.Yes) {
			return nil
		}
		c, err := newHTTP()
		if err != nil {
			return err
		}
		_, err = c.Delete("/v1/webhooks/" + url.PathEscape(args[0]))
		if err != nil {
			return err
		}
		printer.Success("Webhook deleted")
		return nil
	},
}

func webhookBody(cmd *cobra.Command, changedOnly bool) (map[string]any, error) {
	if raw, _ := cmd.Flags().GetString("body"); raw != "" {
		var body map[string]any
		if err := json.Unmarshal([]byte(raw), &body); err != nil {
			return nil, fmt.Errorf("invalid --body JSON: %w", err)
		}
		return body, nil
	}
	body := map[string]any{}
	addStringFlag(body, cmd, changedOnly, "name", "name")
	addStringFlag(body, cmd, changedOnly, "kind", "kind")
	addStringFlag(body, cmd, changedOnly, "agent-id", "agent_id")
	addStringFlag(body, cmd, changedOnly, "channel-id", "channel_id")
	addIntFlag(body, cmd, changedOnly, "rate-limit-per-min", "rate_limit_per_min")
	addBoolFlag(body, cmd, changedOnly, "require-hmac", "require_hmac")
	addBoolFlag(body, cmd, changedOnly, "localhost-only", "localhost_only")
	addCSVFlag(body, cmd, changedOnly, "scopes", "scopes")
	addCSVFlag(body, cmd, changedOnly, "ip-allowlist", "ip_allowlist")
	return body, nil
}

func addStringFlag(body map[string]any, cmd *cobra.Command, changedOnly bool, flag, key string) {
	if changedOnly && !cmd.Flags().Changed(flag) {
		return
	}
	v := mustString(cmd, flag)
	if v != "" || cmd.Flags().Changed(flag) {
		body[key] = v
	}
}

func addIntFlag(body map[string]any, cmd *cobra.Command, changedOnly bool, flag, key string) {
	if changedOnly && !cmd.Flags().Changed(flag) {
		return
	}
	v, _ := cmd.Flags().GetInt(flag)
	if v != 0 || cmd.Flags().Changed(flag) {
		body[key] = v
	}
}

func addBoolFlag(body map[string]any, cmd *cobra.Command, changedOnly bool, flag, key string) {
	if changedOnly && !cmd.Flags().Changed(flag) {
		return
	}
	v, _ := cmd.Flags().GetBool(flag)
	if v || cmd.Flags().Changed(flag) {
		body[key] = v
	}
}

func addCSVFlag(body map[string]any, cmd *cobra.Command, changedOnly bool, flag, key string) {
	if changedOnly && !cmd.Flags().Changed(flag) {
		return
	}
	raw := mustString(cmd, flag)
	if raw != "" || cmd.Flags().Changed(flag) {
		body[key] = splitCSV(raw)
	}
}

func printWebhookSecretResult(result map[string]any) {
	if cfg.OutputFormat == "table" && str(result, "secret") != "" {
		fmt.Printf("Webhook secret: %s\n", str(result, "secret"))
		return
	}
	printer.Print(result)
}

func init() {
	for _, c := range []*cobra.Command{webhooksCreateCmd, webhooksUpdateCmd} {
		c.Flags().String("body", "", "Webhook payload JSON object")
		c.Flags().String("name", "", "Webhook name")
		c.Flags().String("kind", "", "Webhook kind: llm or message")
		c.Flags().String("agent-id", "", "Agent UUID for LLM webhooks")
		c.Flags().String("channel-id", "", "Channel UUID for message webhooks")
		c.Flags().String("scopes", "", "Comma-separated webhook scopes")
		c.Flags().Int("rate-limit-per-min", 0, "Per-webhook rate limit per minute")
		c.Flags().String("ip-allowlist", "", "Comma-separated allowed IP/CIDR entries")
		c.Flags().Bool("require-hmac", false, "Require HMAC signatures")
		c.Flags().Bool("localhost-only", false, "Restrict webhook to localhost callers")
	}
	_ = webhooksCreateCmd.MarkFlagRequired("name")
	_ = webhooksCreateCmd.MarkFlagRequired("kind")
	webhooksCmd.AddCommand(webhooksListCmd, webhooksGetCmd, webhooksCreateCmd,
		webhooksUpdateCmd, webhooksRotateCmd, webhooksDeleteCmd)
	rootCmd.AddCommand(webhooksCmd)
}
