package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

const gatewayUpgradeTokenHeader = "X-GoClaw-Upgrade-Token"

var gatewayUpgradeTagPattern = regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+(-(beta|rc)\.[0-9]+)?$`)

var systemCmd = &cobra.Command{Use: "system", Aliases: []string{"gateway"}, Short: "Manage gateway system operations"}
var gatewayUpgradeCmd = &cobra.Command{Use: "upgrade", Short: "Manage host gateway release upgrades"}

var gatewayUpgradeStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show gateway upgrade status",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := gatewayUpgradeToken(cmd)
		if err != nil {
			return err
		}
		c, err := newHTTP()
		if err != nil {
			return err
		}
		resp, err := c.GetRawWithHeaders("/v1/system/gateway/upgrade/status",
			map[string]string{gatewayUpgradeTokenHeader: token})
		if err != nil {
			return err
		}
		result, err := decodeRawResponse(resp)
		if err != nil {
			return err
		}
		printer.Print(result)
		return nil
	},
}

var gatewayUpgradeStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a gateway release upgrade",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := gatewayUpgradeToken(cmd)
		if err != nil {
			return err
		}
		tag, _ := cmd.Flags().GetString("tag")
		if !validGatewayUpgradeTag(tag) {
			return fmt.Errorf("--tag must be latest or vMAJOR.MINOR.PATCH[-beta.N|-rc.N]")
		}
		c, err := newHTTP()
		if err != nil {
			return err
		}
		body, _ := json.Marshal(map[string]any{"tag": tag})
		resp, err := c.PostRawWithHeaders("/v1/system/gateway/upgrade", "application/json", bytes.NewReader(body),
			map[string]string{gatewayUpgradeTokenHeader: token})
		if err != nil {
			return err
		}
		result, err := decodeRawResponse(resp)
		if err != nil {
			return err
		}
		printer.Print(result)
		return nil
	},
}

func gatewayUpgradeToken(cmd *cobra.Command) (string, error) {
	token, _ := cmd.Flags().GetString("upgrade-token")
	if token == "" {
		token = os.Getenv("GOCLAW_UPGRADE_TRIGGER_TOKEN")
	}
	if token == "" {
		return "", fmt.Errorf("--upgrade-token or GOCLAW_UPGRADE_TRIGGER_TOKEN is required")
	}
	return token, nil
}

func validGatewayUpgradeTag(tag string) bool {
	return tag == "latest" || gatewayUpgradeTagPattern.MatchString(tag)
}

func init() {
	for _, c := range []*cobra.Command{gatewayUpgradeStatusCmd, gatewayUpgradeStartCmd} {
		c.Flags().String("upgrade-token", "", "Gateway upgrade trigger token (env: GOCLAW_UPGRADE_TRIGGER_TOKEN)")
	}
	gatewayUpgradeStartCmd.Flags().String("tag", "", "Release tag: latest or vMAJOR.MINOR.PATCH[-beta.N|-rc.N]")
	_ = gatewayUpgradeStartCmd.MarkFlagRequired("tag")
	gatewayUpgradeCmd.AddCommand(gatewayUpgradeStatusCmd, gatewayUpgradeStartCmd)
	systemCmd.AddCommand(gatewayUpgradeCmd)
	rootCmd.AddCommand(systemCmd)
}
