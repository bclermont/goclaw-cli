package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func workstationBodyFromFlags(cmd *cobra.Command, changedOnly bool) (map[string]any, error) {
	body := map[string]any{}
	for _, name := range []string{"workstation-key", "name", "backend-type", "default-cwd"} {
		if changedOnly && !cmd.Flags().Changed(name) {
			continue
		}
		v, _ := cmd.Flags().GetString(name)
		if v == "" {
			continue
		}
		key := map[string]string{
			"workstation-key": "workstationKey",
			"backend-type":    "backendType",
			"default-cwd":     "defaultCwd",
		}[name]
		if key == "" {
			key = name
		}
		body[key] = v
	}
	for _, spec := range []struct{ flag, key string }{{"metadata", "metadata"}, {"default-env", "defaultEnv"}} {
		if changedOnly && !cmd.Flags().Changed(spec.flag) {
			continue
		}
		raw, _ := cmd.Flags().GetString(spec.flag)
		if raw == "" {
			continue
		}
		var parsed any
		if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
			return nil, fmt.Errorf("invalid --%s JSON: %w", spec.flag, err)
		}
		body[spec.key] = parsed
	}
	return body, nil
}

func callWorkstationLink(cmd *cobra.Command, method string) error {
	ws, err := newWS("cli")
	if err != nil {
		return err
	}
	if _, err := ws.Connect(); err != nil {
		return err
	}
	defer ws.Close()
	agentID, _ := cmd.Flags().GetString("agent")
	workstationID, _ := cmd.Flags().GetString("workstation")
	params := map[string]any{"agentId": agentID, "workstationId": workstationID}
	if cmd.Flags().Changed("default") {
		isDefault, _ := cmd.Flags().GetBool("default")
		params["isDefault"] = isDefault
	}
	data, err := ws.Call(method, params)
	if err != nil {
		return err
	}
	printer.Print(unmarshalMap(data))
	return nil
}
