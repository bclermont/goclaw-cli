package cmd

import (
	"fmt"
	"net/url"

	"github.com/nextlevelbuilder/goclaw-cli/internal/tui"
	"github.com/spf13/cobra"
)

var workstationsPermissionsCmd = &cobra.Command{Use: "permissions", Short: "Manage workstation path permissions"}

var workstationsPermissionsListCmd = &cobra.Command{
	Use:   "list <workstationID>",
	Short: "List workstation permissions",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newHTTP()
		if err != nil {
			return err
		}
		data, err := c.Get("/v1/workstations/" + url.PathEscape(args[0]) + "/permissions")
		if err != nil {
			return err
		}
		printer.Print(toList(unmarshalMap(data)["permissions"]))
		return nil
	},
}

var workstationsPermissionsAddCmd = &cobra.Command{
	Use:   "add <workstationID>",
	Short: "Add a workstation permission pattern",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pattern, _ := cmd.Flags().GetString("pattern")
		c, err := newHTTP()
		if err != nil {
			return err
		}
		data, err := c.Post("/v1/workstations/"+url.PathEscape(args[0])+"/permissions", map[string]any{"pattern": pattern})
		if err != nil {
			return err
		}
		printer.Print(unmarshalMap(data))
		return nil
	},
}

var workstationsPermissionsRemoveCmd = &cobra.Command{
	Use:   "remove <workstationID> <permissionID>",
	Short: "Remove a workstation permission",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !tui.Confirm("Remove this workstation permission?", cfg.Yes) {
			return nil
		}
		c, err := newHTTP()
		if err != nil {
			return err
		}
		_, err = c.Delete(fmt.Sprintf("/v1/workstations/%s/permissions/%s", url.PathEscape(args[0]), url.PathEscape(args[1])))
		if err != nil {
			return err
		}
		printer.Success("Workstation permission removed")
		return nil
	},
}

var workstationsPermissionsToggleCmd = &cobra.Command{
	Use:   "toggle <workstationID> <permissionID>",
	Short: "Enable or disable a workstation permission",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		enabled, _ := cmd.Flags().GetBool("enabled")
		c, err := newHTTP()
		if err != nil {
			return err
		}
		data, err := c.Put(fmt.Sprintf("/v1/workstations/%s/permissions/%s/toggle", url.PathEscape(args[0]), url.PathEscape(args[1])),
			map[string]any{"enabled": enabled})
		if err != nil {
			return err
		}
		printer.Print(unmarshalMap(data))
		return nil
	},
}

var workstationsActivityCmd = &cobra.Command{
	Use:   "activity <workstationID>",
	Short: "List workstation activity",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		q := url.Values{}
		if v, _ := cmd.Flags().GetInt("limit"); v > 0 {
			q.Set("limit", fmt.Sprintf("%d", v))
		}
		if v, _ := cmd.Flags().GetString("cursor"); v != "" {
			q.Set("cursor", v)
		}
		path := "/v1/workstations/" + url.PathEscape(args[0]) + "/activity"
		if len(q) > 0 {
			path += "?" + q.Encode()
		}
		c, err := newHTTP()
		if err != nil {
			return err
		}
		data, err := c.Get(path)
		if err != nil {
			return err
		}
		printer.Print(unmarshalMap(data))
		return nil
	},
}

func init() {
	workstationsPermissionsAddCmd.Flags().String("pattern", "", "Allowed path pattern")
	_ = workstationsPermissionsAddCmd.MarkFlagRequired("pattern")
	workstationsPermissionsToggleCmd.Flags().Bool("enabled", true, "Whether the permission is enabled")
	workstationsActivityCmd.Flags().Int("limit", 50, "Max activity rows")
	workstationsActivityCmd.Flags().String("cursor", "", "Pagination cursor")
	workstationsPermissionsCmd.AddCommand(workstationsPermissionsListCmd, workstationsPermissionsAddCmd,
		workstationsPermissionsRemoveCmd, workstationsPermissionsToggleCmd)
	workstationsCmd.AddCommand(workstationsPermissionsCmd, workstationsActivityCmd)
}
