# Contract Lock Report

## Source Inventory

- Server trace list: `GET /v1/traces`
- Server trace follow: `GET /v1/traces/follow`
- Server trace detail: `GET /v1/traces/{traceID}`
- Server trace export: `GET /v1/traces/{traceID}/export`
- Server run timeline: `GET /v1/runs/{runID}/timeline`

## Locked Response Shapes

- Trace list returns `{traces,total,limit,offset}`.
- Trace detail returns `{trace,spans}`.
- Trace follow returns `{traces,spans_by_trace_id,server_time,next_since,limit}`.
- Run timeline returns `{run_id,session_key,items,limit,offset}`.

## Locked Field Names

- Trace rows use `id`, `agent_id`, `session_key`, `status`, `duration_ms`, `total_input_tokens`, `total_output_tokens`, `total_cost`.
- Span rows use `id`, `parent_span_id`, `span_type`, `name`, `duration_ms`, `status`.
- Timeline items use `seq`, `item_type`, `status`, `title`, `tool_name`, `trace_id`, `span_id`, `created_at`.

## Unsupported

- `POST /v1/traces/{id}/replay` is absent on server `dev`; no CLI replay command added.

## Fixture Safety

- Fixtures are synthetic.
- No bearer tokens, prompts, API keys, live user IDs, or live tenant IDs.

## Red Gate Captured

- Focused test gate initially failed at compile because `tracesTimelineCmd` was missing.
- Existing list/get tests also encoded stale array/flat payload assumptions before update.

## Unresolved Questions

None.
