package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/nextlevelbuilder/goclaw-cli/internal/tui"
	"github.com/spf13/cobra"
)

var mcpUserCredentialsCmd = &cobra.Command{Use: "user-credentials", Short: "Manage per-user MCP server credentials"}

var mcpUserCredentialsGetCmd = &cobra.Command{
	Use:   "get <serverID>",
	Short: "Show credential presence for an MCP server",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newHTTP()
		if err != nil {
			return err
		}
		data, err := c.Get(mcpUserCredentialsPath(cmd, args[0]))
		if err != nil {
			return err
		}
		printer.Print(unmarshalMap(data))
		return nil
	},
}

var mcpUserCredentialsSetCmd = &cobra.Command{
	Use:   "set <serverID>",
	Short: "Set per-user MCP credentials",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		raw, _ := cmd.Flags().GetString("body")
		if raw == "" {
			return fmt.Errorf("--body is required (JSON object)")
		}
		var body map[string]any
		if err := json.Unmarshal([]byte(raw), &body); err != nil {
			return fmt.Errorf("invalid --body JSON: %w", err)
		}
		c, err := newHTTP()
		if err != nil {
			return err
		}
		data, err := c.Put(mcpUserCredentialsPath(cmd, args[0]), body)
		if err != nil {
			return err
		}
		printer.Print(unmarshalMap(data))
		return nil
	},
}

var mcpUserCredentialsDeleteCmd = &cobra.Command{
	Use:   "delete <serverID>",
	Short: "Delete per-user MCP credentials",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !tui.Confirm("Delete MCP user credentials?", cfg.Yes) {
			return nil
		}
		c, err := newHTTP()
		if err != nil {
			return err
		}
		data, err := c.Delete(mcpUserCredentialsPath(cmd, args[0]))
		if err != nil {
			return err
		}
		printer.Print(unmarshalMap(data))
		return nil
	},
}

func mcpUserCredentialsPath(cmd *cobra.Command, serverID string) string {
	path := "/v1/mcp/servers/" + url.PathEscape(serverID) + "/user-credentials"
	if userID, _ := cmd.Flags().GetString("user"); userID != "" {
		path += "?user_id=" + url.QueryEscape(userID)
	}
	return path
}

func init() {
	for _, c := range []*cobra.Command{mcpUserCredentialsGetCmd, mcpUserCredentialsSetCmd, mcpUserCredentialsDeleteCmd} {
		c.Flags().String("user", "", "Target user ID (admin/tenant admin only)")
	}
	mcpUserCredentialsSetCmd.Flags().String("body", "", "Credentials JSON object")
	_ = mcpUserCredentialsSetCmd.MarkFlagRequired("body")
	mcpUserCredentialsCmd.AddCommand(mcpUserCredentialsGetCmd, mcpUserCredentialsSetCmd, mcpUserCredentialsDeleteCmd)
	mcpServersCmd.AddCommand(mcpUserCredentialsCmd)
}
