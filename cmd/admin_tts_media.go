package cmd

import (
	"github.com/nextlevelbuilder/goclaw-cli/internal/tui"
	"github.com/spf13/cobra"
)

// admin_tts_media.go holds TTS and Media subcommands, extracted from admin.go
// to keep that file under 200 LoC.

// --- TTS ---

var ttsCmd = &cobra.Command{Use: "tts", Short: "Text-to-speech operations"}

var ttsStatusCmd = &cobra.Command{
	Use: "status", Short: "TTS status",
	RunE: func(cmd *cobra.Command, args []string) error {
		ws, err := newWS("cli")
		if err != nil {
			return err
		}
		if _, err := ws.Connect(); err != nil {
			return err
		}
		defer ws.Close()
		data, err := ws.Call("tts.status", nil)
		if err != nil {
			return err
		}
		printer.Print(unmarshalMap(data))
		return nil
	},
}

var ttsEnableCmd = &cobra.Command{
	Use: "enable", Short: "Enable TTS",
	RunE: func(cmd *cobra.Command, args []string) error {
		ws, err := newWS("cli")
		if err != nil {
			return err
		}
		if _, err := ws.Connect(); err != nil {
			return err
		}
		defer ws.Close()
		_, err = ws.Call("tts.enable", nil)
		if err != nil {
			return err
		}
		printer.Success("TTS enabled")
		return nil
	},
}

var ttsDisableCmd = &cobra.Command{
	Use: "disable", Short: "Disable TTS",
	RunE: func(cmd *cobra.Command, args []string) error {
		ws, err := newWS("cli")
		if err != nil {
			return err
		}
		if _, err := ws.Connect(); err != nil {
			return err
		}
		defer ws.Close()
		_, err = ws.Call("tts.disable", nil)
		if err != nil {
			return err
		}
		printer.Success("TTS disabled")
		return nil
	},
}

var ttsProvidersCmd = &cobra.Command{
	Use: "providers", Short: "List TTS providers",
	RunE: func(cmd *cobra.Command, args []string) error {
		ws, err := newWS("cli")
		if err != nil {
			return err
		}
		if _, err := ws.Connect(); err != nil {
			return err
		}
		defer ws.Close()
		data, err := ws.Call("tts.providers", nil)
		if err != nil {
			return err
		}
		printer.Print(unmarshalList(data))
		return nil
	},
}

var ttsSetProviderCmd = &cobra.Command{
	Use: "set-provider", Short: "Set TTS provider",
	RunE: func(cmd *cobra.Command, args []string) error {
		ws, err := newWS("cli")
		if err != nil {
			return err
		}
		if _, err := ws.Connect(); err != nil {
			return err
		}
		defer ws.Close()
		name, _ := cmd.Flags().GetString("name")
		_, err = ws.Call("tts.setProvider", map[string]any{"provider": name})
		if err != nil {
			return err
		}
		printer.Success("TTS provider set")
		return nil
	},
}

var ttsTestConnectionCmd = &cobra.Command{
	Use: "test-connection", Short: "Test TTS provider connection",
	Long: `POST /v1/tts/test-connection

Flags:
  --provider=<name>   Provider to test (e.g. elevenlabs, azure)
  --voice=<id>        Optional voice ID to test`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newHTTP()
		if err != nil {
			return err
		}
		body := buildBody(
			"provider", mustString(cmd, "provider"),
			"voice", mustString(cmd, "voice"),
		)
		data, err := c.Post("/v1/tts/test-connection", body)
		if err != nil {
			return err
		}
		printer.Print(unmarshalMap(data))
		return nil
	},
}

// --- Voices ---

var voicesCmd = &cobra.Command{Use: "voices", Short: "Manage voice catalog"}

var voicesListCmd = &cobra.Command{
	Use: "list", Short: "List available voices (GET /v1/voices)",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newHTTP()
		if err != nil {
			return err
		}
		data, err := c.Get("/v1/voices")
		if err != nil {
			return err
		}
		printer.Print(unmarshalList(data))
		return nil
	},
}

var voicesRefreshCmd = &cobra.Command{
	Use: "refresh", Short: "Refresh voice catalog from providers (admin)",
	Long: `POST /v1/voices/refresh

Forces a re-fetch of the voice catalog from configured TTS providers.
Requires --yes to confirm.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !tui.Confirm("Refresh voice catalog from upstream providers?", cfg.Yes) {
			return nil
		}
		c, err := newHTTP()
		if err != nil {
			return err
		}
		_, err = c.Post("/v1/voices/refresh", nil)
		if err != nil {
			return err
		}
		printer.Success("Voice catalog refresh triggered")
		return nil
	},
}

// mustString reads a string flag (returns empty on error/missing) — used for buildBody pairs.
func mustString(cmd *cobra.Command, name string) string {
	v, _ := cmd.Flags().GetString(name)
	return v
}

func init() {
	ttsSetProviderCmd.Flags().String("name", "", "Provider name")
	_ = ttsSetProviderCmd.MarkFlagRequired("name")
	ttsTestConnectionCmd.Flags().String("provider", "", "Provider to test")
	ttsTestConnectionCmd.Flags().String("voice", "", "Voice ID to test")
	_ = ttsTestConnectionCmd.MarkFlagRequired("provider")
	ttsCmd.AddCommand(ttsStatusCmd, ttsEnableCmd, ttsDisableCmd,
		ttsProvidersCmd, ttsSetProviderCmd, ttsTestConnectionCmd)

	voicesCmd.AddCommand(voicesListCmd, voicesRefreshCmd)
}
