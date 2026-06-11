package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/nextlevelbuilder/goclaw-cli/internal/output"
)

func decodeTraceListPayload(data json.RawMessage) (map[string]any, []any, error) {
	var envelope map[string]any
	if err := json.Unmarshal(data, &envelope); err == nil {
		rows, _ := envelope["traces"].([]any)
		if rows == nil {
			rows = []any{}
		}
		return envelope, rows, nil
	}

	var legacy []any
	if err := json.Unmarshal(data, &legacy); err != nil {
		return nil, nil, fmt.Errorf("decode traces payload: %w", err)
	}
	return map[string]any{
		"traces": legacy,
		"total":  len(legacy),
		"limit":  len(legacy),
		"offset": 0,
	}, legacy, nil
}

func printTraceRowsTable(rows []any) {
	tbl := output.NewTable("ID", "AGENT", "STATUS", "DURATION_MS", "TOTAL_INPUT_TOKENS", "TOTAL_OUTPUT_TOKENS", "TOTAL_COST")
	for _, raw := range rows {
		t, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		tbl.AddRow(traceField(t, "id", "trace_id"), str(t, "agent_id"), str(t, "status"),
			str(t, "duration_ms"), traceField(t, "total_input_tokens", "input_tokens"),
			traceField(t, "total_output_tokens", "output_tokens"), traceField(t, "total_cost", "cost"))
	}
	printer.Print(tbl)
}

// renderTraceTable prints a human-readable summary and span tree.
func renderTraceTable(payload map[string]any, w io.Writer) {
	trace := payload
	if nested, ok := payload["trace"].(map[string]any); ok {
		trace = nested
	}
	for _, row := range [][2]string{
		{"TRACE_ID", traceField(trace, "id", "trace_id")},
		{"AGENT_ID", str(trace, "agent_id")},
		{"SESSION_KEY", str(trace, "session_key")},
		{"RUN_ID", str(trace, "run_id")},
		{"STATUS", str(trace, "status")},
		{"DURATION_MS", str(trace, "duration_ms")},
		{"TOTAL_INPUT_TOKENS", traceField(trace, "total_input_tokens", "input_tokens")},
		{"TOTAL_OUTPUT_TOKENS", traceField(trace, "total_output_tokens", "output_tokens")},
		{"TOTAL_COST", traceField(trace, "total_cost", "cost")},
	} {
		if row[1] != "" {
			fmt.Fprintf(w, "%-20s %s\n", row[0]+":", row[1])
		}
	}
	spans, _ := payload["spans"].([]any)
	if len(spans) == 0 {
		fmt.Fprintln(w, "\nSPANS: (none)")
		return
	}
	fmt.Fprintln(w, "\nSPANS:")
	output.PrintTreeRoot(buildSpanTree(spans), w)
}

// buildSpanTree links spans via parent_span_id. Children are kept in insertion order.
func buildSpanTree(spans []any) output.TreeNode {
	order := make([]string, 0, len(spans))
	labels := make(map[string]string, len(spans))
	children := make(map[string][]string, len(spans))
	parentOf := make(map[string]string, len(spans))
	for _, s := range spans {
		m, ok := s.(map[string]any)
		if !ok {
			continue
		}
		id := traceField(m, "id", "span_id")
		if id == "" {
			continue
		}
		label := id
		if name := str(m, "name"); name != "" {
			label = name + " [" + id + "]"
		}
		if kind := traceField(m, "span_type", "kind"); kind != "" {
			label += " span_type=" + kind
		}
		if dur := str(m, "duration_ms"); dur != "" {
			label += " " + dur + "ms"
		}
		if status := str(m, "status"); status != "" {
			label += " status=" + status
		}
		labels[id] = label
		order = append(order, id)
		parentOf[id], _ = m["parent_span_id"].(string)
	}
	for _, id := range order {
		if p := parentOf[id]; p != "" {
			if _, ok := labels[p]; ok {
				children[p] = append(children[p], id)
				continue
			}
		}
		children[""] = append(children[""], id)
	}
	var build func(id string) output.TreeNode
	build = func(id string) output.TreeNode {
		n := output.TreeNode{Name: labels[id]}
		for _, c := range children[id] {
			n.Children = append(n.Children, build(c))
		}
		return n
	}
	root := output.TreeNode{Name: "trace"}
	for _, id := range children[""] {
		root.Children = append(root.Children, build(id))
	}
	return root
}

func traceField(m map[string]any, primary, fallback string) string {
	if v := str(m, primary); v != "" && v != "<nil>" {
		return v
	}
	if fallback == "" {
		return ""
	}
	if v := str(m, fallback); v != "" && v != "<nil>" {
		return v
	}
	return ""
}
