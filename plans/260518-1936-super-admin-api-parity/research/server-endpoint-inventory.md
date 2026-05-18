# Server Endpoint Inventory

## Source

Route scan was based on server checkout `D:\www\digitop\goclaw` after `git pull` to commit `99460a7f`.

## P0 Endpoints

| Domain | Method | Path | Server file | CLI status |
|---|---|---|---|---|
| API keys | POST | `/v1/api-keys/{id}/revoke` | `internal/http/api_keys.go` | Broken: CLI uses DELETE. |
| Gateway upgrade | GET | `/v1/system/gateway/upgrade/status` | `internal/http/gateway_upgrade.go` | Missing. |
| Gateway upgrade | POST | `/v1/system/gateway/upgrade` | `internal/http/gateway_upgrade.go` | Missing. |
| Packages | GET | `/v1/packages/updates` | `internal/http/packages.go` | Missing. |
| Packages | POST | `/v1/packages/updates/refresh` | `internal/http/packages.go` | Missing. |
| Packages | POST | `/v1/packages/update` | `internal/http/packages.go` | Missing. |
| Packages | POST | `/v1/packages/updates/apply-all` | `internal/http/packages.go` | Missing. |
| Workstations | GET | `/v1/workstations` | `internal/http/workstations.go` | Missing. |
| Workstations | POST | `/v1/workstations` | `internal/http/workstations.go` | Missing. |
| Workstations | GET | `/v1/workstations/{id}` | `internal/http/workstations.go` | Missing. |
| Workstations | PUT | `/v1/workstations/{id}` | `internal/http/workstations.go` | Missing. |
| Workstations | DELETE | `/v1/workstations/{id}` | `internal/http/workstations.go` | Missing. |
| Workstations | GET | `/v1/workstations/{id}/permissions` | `internal/http/workstations.go` | Missing. |
| Workstations | POST | `/v1/workstations/{id}/permissions` | `internal/http/workstations.go` | Missing. |
| Workstations | DELETE | `/v1/workstations/{id}/permissions/{permId}` | `internal/http/workstations.go` | Missing. |
| Workstations | PUT | `/v1/workstations/{id}/permissions/{permId}/toggle` | `internal/http/workstations.go` | Missing. |
| Workstations | GET | `/v1/workstations/{id}/activity` | `internal/http/workstations.go` | Missing. |

## P1 Endpoints

| Domain | Method | Path | Server file | CLI status |
|---|---|---|---|---|
| Webhooks | POST | `/v1/webhooks` | `internal/http/webhooks_admin.go` | Missing. |
| Webhooks | GET | `/v1/webhooks` | `internal/http/webhooks_admin.go` | Missing. |
| Webhooks | GET | `/v1/webhooks/{id}` | `internal/http/webhooks_admin.go` | Missing. |
| Webhooks | PATCH | `/v1/webhooks/{id}` | `internal/http/webhooks_admin.go` | Missing. |
| Webhooks | POST | `/v1/webhooks/{id}/rotate` | `internal/http/webhooks_admin.go` | Missing. |
| Webhooks | DELETE | `/v1/webhooks/{id}` | `internal/http/webhooks_admin.go` | Missing. |
| CLI credentials | POST | `/v1/cli-credentials/{id}/agent-grants/{grantId}/env:reveal` | `internal/http/secure_cli_agent_grants.go` | Missing. |
| MCP credentials | GET | `/v1/mcp/servers/{id}/user-credentials` | `internal/http/mcp_user_credentials.go` | Missing. |
| MCP credentials | PUT | `/v1/mcp/servers/{id}/user-credentials` | `internal/http/mcp_user_credentials.go` | Missing. |
| MCP credentials | DELETE | `/v1/mcp/servers/{id}/user-credentials` | `internal/http/mcp_user_credentials.go` | Missing. |
| Media | POST | `/v1/media/upload` | `internal/http/media_upload.go` | Stub in CLI. |
| TTS | GET | `/v1/tts/capabilities` | `internal/http/tts_capabilities.go` | Missing. |
| TTS | GET | `/v1/tts/config` | `internal/http/tts_config.go` | Missing. |
| TTS | POST | `/v1/tts/config` | `internal/http/tts_config.go` | Missing. |
| TTS | POST | `/v1/tts/synthesize` | `internal/http/tts.go` | Missing. |

## P2 Endpoints

| Domain | Method | Path | Server file | CLI status |
|---|---|---|---|---|
| Storage | POST | `/v1/storage/files` | `internal/http/storage.go` | Missing. |
| Storage | PUT | `/v1/storage/move` | `internal/http/storage.go` | Missing. |
| Channels | POST | `/v1/contacts/unmerge` | `internal/http/channel_instances.go` | Missing. |
| Channels | GET | `/v1/tenant-users` | `internal/http/channel_instances.go` | Missing. |
| Channels | GET | `/v1/channels/instances/{id}/writers/groups` | `internal/http/channel_instances.go` | Missing. |
| Teams | GET | `/v1/teams/{teamId}/attachments/{attachmentId}/download` | `internal/http/team_attachments.go` | Missing. |

## Notes

- The server OpenAPI JSON is incomplete relative to source-registered routes. Source route scan is more authoritative for this plan.
- Some WS methods are also uncovered, but HTTP is preferred when both exist because command tests are simpler and automation output is easier to reason about.
