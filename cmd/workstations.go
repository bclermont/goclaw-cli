package cmd

import (
	"net/url"

	"github.com/nextlevelbuilder/goclaw-cli/internal/output"
	"github.com/nextlevelbuilder/goclaw-cli/internal/tui"
	"github.com/spf13/cobra"
)

var workstationsCmd = &cobra.Command{Use: "workstations", Short: "Manage coding-agent workstations"}

var workstationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List workstations",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newHTTP()
		if err != nil {
			return err
		}
		data, err := c.Get("/v1/workstations")
		if err != nil {
			return err
		}
		m := unmarshalMap(data)
		items := toList(m["workstations"])
		if cfg.OutputFormat != "table" {
			printer.Print(items)
			return nil
		}
		tbl := output.NewTable("ID", "KEY", "NAME", "BACKEND", "ACTIVE")
		for _, ws := range items {
			tbl.AddRow(str(ws, "id"), str(ws, "workstationKey"), str(ws, "name"), str(ws, "backendType"), str(ws, "active"))
		}
		printer.Print(tbl)
		return nil
	},
}

var workstationsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get workstation details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newHTTP()
		if err != nil {
			return err
		}
		data, err := c.Get("/v1/workstations/" + url.PathEscape(args[0]))
		if err != nil {
			return err
		}
		printer.Print(unmarshalMap(data))
		return nil
	},
}

var workstationsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a workstation",
	RunE: func(cmd *cobra.Command, args []string) error {
		body, err := workstationBodyFromFlags(cmd, false)
		if err != nil {
			return err
		}
		c, err := newHTTP()
		if err != nil {
			return err
		}
		data, err := c.Post("/v1/workstations", body)
		if err != nil {
			return err
		}
		printer.Print(unmarshalMap(data))
		return nil
	},
}

var workstationsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a workstation",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		body, err := workstationBodyFromFlags(cmd, true)
		if err != nil {
			return err
		}
		c, err := newHTTP()
		if err != nil {
			return err
		}
		data, err := c.Put("/v1/workstations/"+url.PathEscape(args[0]), body)
		if err != nil {
			return err
		}
		printer.Print(unmarshalMap(data))
		return nil
	},
}

var workstationsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a workstation",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !tui.Confirm("Delete this workstation?", cfg.Yes) {
			return nil
		}
		c, err := newHTTP()
		if err != nil {
			return err
		}
		_, err = c.Delete("/v1/workstations/" + url.PathEscape(args[0]))
		if err != nil {
			return err
		}
		printer.Success("Workstation deleted")
		return nil
	},
}

var workstationsLinkCmd = &cobra.Command{
	Use:   "link-agent",
	Short: "Link an agent to a workstation over WS RPC",
	RunE: func(cmd *cobra.Command, args []string) error {
		return callWorkstationLink(cmd, "workstations.linkAgent")
	},
}

var workstationsUnlinkCmd = &cobra.Command{
	Use:   "unlink-agent",
	Short: "Unlink an agent from a workstation over WS RPC",
	RunE: func(cmd *cobra.Command, args []string) error {
		return callWorkstationLink(cmd, "workstations.unlinkAgent")
	},
}

func init() {
	for _, c := range []*cobra.Command{workstationsCreateCmd, workstationsUpdateCmd} {
		c.Flags().String("workstation-key", "", "Stable workstation key")
		c.Flags().String("name", "", "Display name")
		c.Flags().String("backend-type", "", "Backend type")
		c.Flags().String("metadata", "", "Metadata JSON object")
		c.Flags().String("default-cwd", "", "Default working directory")
		c.Flags().String("default-env", "", "Default environment JSON object")
	}
	_ = workstationsCreateCmd.MarkFlagRequired("workstation-key")
	_ = workstationsCreateCmd.MarkFlagRequired("backend-type")
	for _, c := range []*cobra.Command{workstationsLinkCmd, workstationsUnlinkCmd} {
		c.Flags().String("agent", "", "Agent ID")
		c.Flags().String("workstation", "", "Workstation ID")
		_ = c.MarkFlagRequired("agent")
		_ = c.MarkFlagRequired("workstation")
	}
	workstationsLinkCmd.Flags().Bool("default", false, "Make this the agent default workstation")
	workstationsCmd.AddCommand(workstationsListCmd, workstationsGetCmd, workstationsCreateCmd,
		workstationsUpdateCmd, workstationsDeleteCmd, workstationsLinkCmd, workstationsUnlinkCmd)
	rootCmd.AddCommand(workstationsCmd)
}
