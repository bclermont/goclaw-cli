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
		body, err := webhookBody(cmd)
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
		body, err := webhookBody(cmd)
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

func webhookBody(cmd *cobra.Command) (map[string]any, error) {
	if raw, _ := cmd.Flags().GetString("body"); raw != "" {
		var body map[string]any
		if err := json.Unmarshal([]byte(raw), &body); err != nil {
			return nil, fmt.Errorf("invalid --body JSON: %w", err)
		}
		return body, nil
	}
	return buildBody(
		"name", mustString(cmd, "name"),
		"target_type", mustString(cmd, "target-type"),
		"target_id", mustString(cmd, "target-id"),
		"enabled", mustBool(cmd, "enabled"),
	), nil
}

func mustBool(cmd *cobra.Command, name string) bool {
	v, _ := cmd.Flags().GetBool(name)
	return v
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
		c.Flags().String("target-type", "", "Target type")
		c.Flags().String("target-id", "", "Target ID")
		c.Flags().Bool("enabled", true, "Whether webhook is enabled")
	}
	webhooksCmd.AddCommand(webhooksListCmd, webhooksGetCmd, webhooksCreateCmd,
		webhooksUpdateCmd, webhooksRotateCmd, webhooksDeleteCmd)
	rootCmd.AddCommand(webhooksCmd)
}
