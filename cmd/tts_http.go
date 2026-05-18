package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var ttsConfigCmd = &cobra.Command{Use: "config", Short: "Manage tenant TTS config"}

var ttsConfigGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get tenant TTS config",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newHTTP()
		if err != nil {
			return err
		}
		data, err := c.Get("/v1/tts/config")
		if err != nil {
			return err
		}
		printer.Print(unmarshalMap(data))
		return nil
	},
}

var ttsConfigSaveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save tenant TTS config",
	RunE: func(cmd *cobra.Command, args []string) error {
		body, err := ttsConfigBody(cmd)
		if err != nil {
			return err
		}
		c, err := newHTTP()
		if err != nil {
			return err
		}
		data, err := c.Post("/v1/tts/config", body)
		if err != nil {
			return err
		}
		printer.Print(unmarshalMap(data))
		return nil
	},
}

var ttsCapabilitiesCmd = &cobra.Command{
	Use:   "capabilities",
	Short: "List TTS provider capabilities",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newHTTP()
		if err != nil {
			return err
		}
		data, err := c.Get("/v1/tts/capabilities")
		if err != nil {
			return err
		}
		printer.Print(unmarshalMap(data))
		return nil
	},
}

var ttsSynthesizeCmd = &cobra.Command{
	Use:   "synthesize",
	Short: "Synthesize text to an audio file",
	RunE: func(cmd *cobra.Command, args []string) error {
		outFile, _ := cmd.Flags().GetString("file")
		if outFile == "" {
			return fmt.Errorf("--file is required for raw audio output")
		}
		body := buildBody(
			"text", mustString(cmd, "text"),
			"provider", mustString(cmd, "provider"),
			"voice_id", mustString(cmd, "voice-id"),
			"model_id", mustString(cmd, "model-id"),
		)
		raw, _ := json.Marshal(body)
		c, err := newHTTP()
		if err != nil {
			return err
		}
		resp, err := c.PostRaw("/v1/tts/synthesize", "application/json", bytes.NewReader(raw))
		if err != nil {
			return err
		}
		if resp.StatusCode >= 400 {
			return rawResponseError(resp)
		}
		defer resp.Body.Close()
		f, err := os.Create(outFile)
		if err != nil {
			return err
		}
		defer f.Close()
		n, err := io.Copy(f, resp.Body)
		if err != nil {
			return err
		}
		printer.Success(fmt.Sprintf("Saved %d bytes to %s", n, outFile))
		return nil
	},
}

func ttsConfigBody(cmd *cobra.Command) (map[string]any, error) {
	if raw, _ := cmd.Flags().GetString("body"); raw != "" {
		var body map[string]any
		if err := json.Unmarshal([]byte(raw), &body); err != nil {
			return nil, fmt.Errorf("invalid --body JSON: %w", err)
		}
		return body, nil
	}
	return buildBody(
		"provider", mustString(cmd, "provider"),
		"auto", mustString(cmd, "auto"),
		"mode", mustString(cmd, "mode"),
		"max_length", mustInt(cmd, "max-length"),
		"timeout_ms", mustInt(cmd, "timeout-ms"),
	), nil
}

func mustInt(cmd *cobra.Command, name string) int {
	v, _ := cmd.Flags().GetInt(name)
	return v
}

func init() {
	ttsConfigSaveCmd.Flags().String("body", "", "Full TTS config JSON object")
	ttsConfigSaveCmd.Flags().String("provider", "", "Provider name")
	ttsConfigSaveCmd.Flags().String("auto", "", "Auto mode")
	ttsConfigSaveCmd.Flags().String("mode", "", "TTS mode")
	ttsConfigSaveCmd.Flags().Int("max-length", 0, "Max text length")
	ttsConfigSaveCmd.Flags().Int("timeout-ms", 0, "Synthesis timeout in milliseconds")
	ttsSynthesizeCmd.Flags().String("text", "", "Text to synthesize")
	ttsSynthesizeCmd.Flags().String("provider", "", "Provider override")
	ttsSynthesizeCmd.Flags().String("voice-id", "", "Voice ID override")
	ttsSynthesizeCmd.Flags().String("model-id", "", "Model ID override")
	ttsSynthesizeCmd.Flags().StringP("file", "f", "", "Output audio file")
	_ = ttsSynthesizeCmd.MarkFlagRequired("text")
	ttsConfigCmd.AddCommand(ttsConfigGetCmd, ttsConfigSaveCmd)
	ttsCmd.AddCommand(ttsConfigCmd, ttsCapabilitiesCmd, ttsSynthesizeCmd)
}
