---
phase: 5
title: "Media TTS Storage Channel Fillers"
status: completed
effort: "4h"
---

# Phase 5: Media TTS Storage Channel Fillers

## Overview

Priority: P2
Current status: completed

Close the remaining super-admin operational gaps that are visible to coding agents but less central than workstation/package/credential control.

## Context Links

- `D:\www\digitop\goclaw\internal\http\media_upload.go`
- `D:\www\digitop\goclaw\internal\http\tts.go`
- `D:\www\digitop\goclaw\internal\http\tts_config.go`
- `D:\www\digitop\goclaw\internal\http\tts_capabilities.go`
- `D:\www\digitop\goclaw\internal\http\storage.go`
- `D:\www\digitop\goclaw\internal\http\channel_instances.go`
- `D:\www\nextlevelbuilder\goclaw-cli\cmd\admin_tts_media.go`
- `D:\www\nextlevelbuilder\goclaw-cli\cmd\storage.go`
- `D:\www\nextlevelbuilder\goclaw-cli\cmd\channels_contacts.go`
- `D:\www\nextlevelbuilder\goclaw-cli\cmd\channels_writers.go`

## Key Insights

- `media upload` currently is a stub and should become real multipart upload.
- TTS has HTTP config/capabilities/synthesize routes in addition to current WS status/provider commands.
- Storage lacks upload/move despite server routes.
- Channels lack `contacts unmerge`, `tenant-users`, and writer group listing.

## Requirements

- Replace `media upload` stub with streaming multipart upload.
- Add `tts config get/set`, `tts capabilities`, `tts synthesize`.
- Add `storage upload`, `storage move`.
- Add `channels contacts unmerge`, `channels tenant-users`, `channels writers groups`.
- `tts synthesize` must require `--file` because server returns raw audio bytes.

## Architecture

```
multipart helper -> media upload + storage upload
tts HTTP routes -> config/capabilities/synthesize
channel fillers -> small REST GET/POST commands
```

## Related Code Files

| File | Action |
|---|---|
| `cmd/admin_tts_media.go` | modify or split |
| `cmd/media_test.go` | create |
| `cmd/tts_http_test.go` | create |
| `cmd/storage.go` | modify or split |
| `cmd/storage_test.go` | create |
| `cmd/channels_contacts.go` | modify |
| `cmd/channels_writers.go` | modify |

## Implementation Steps

1. Add failing test proving media upload performs multipart POST.
2. Reuse `newMultipartWriter` from vault/team workspace helper if compatible; otherwise extract a small shared helper without breaking existing tests.
3. Add TTS HTTP commands:
   - `tts capabilities`
   - `tts config get`
   - `tts config set --body @file-or-json`
   - `tts synthesize --text|--text-file --provider --voice --file`
   - reject synthesize without `--file`; do not pass audio bytes to `printer.Print`
4. Add storage upload/move using existing endpoint shapes.
5. Add channel filler commands with small tests.

## Tests Before

- Media upload test should fail because current command is stub.
- TTS config/capabilities/synthesize tests fail due missing commands.
- Storage upload/move tests fail due missing commands.
- TTS synthesize test must assert file write path and no binary stdout.

## Refactor

- Split `admin_tts_media.go` if it grows past 200 lines after additions.
- Prefer `@file` support for large TTS config bodies.
- Binary/audio output must write to file, not table printer.
- Check response status and content type before writing success output.

## Tests After

- `go test ./cmd -run "TestMedia|TestTTS|TestStorage|TestChannels"`
- `go build ./...`

## Todo List

- [x] Real media upload.
- [x] TTS HTTP commands.
- [x] Storage upload/move.
- [x] Channel filler commands.
- [x] Add TTS raw-audio output tests.

## Success Criteria

- [x] Media upload returns server JSON, not a placeholder message.
- [x] TTS config and synthesize can be managed by automation.
- [x] Storage upload/move support workspace maintenance.
- [x] Channel helpers remove remaining admin dashboard dependency for listed routes.

## Risk Assessment

- Risk: audio response handling corrupts terminal output. Mitigation: require `--file` for synthesize unless server returns JSON.
- Risk: multipart helper regressions. Mitigation: reuse tests from vault upload patterns.

## Security Considerations

- Avoid printing media binary to stdout accidentally.
- Treat TTS config as potentially secret-bearing.
- Destructive storage delete already gated; move may overwrite server-side, document behavior.
- Do not include TTS provider secrets in docs or fixtures.

## Regression Gate

```powershell
go test ./cmd -run "TestMedia|TestTTS|TestStorage|TestChannels"
go build ./...
```

## Next Steps

- Phase 6 validates docs and route coverage.
