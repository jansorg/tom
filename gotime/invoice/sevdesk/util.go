package sevdesk

import "time"

func iso8601(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

func boolToString(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

func stringOr0(v string) string {
	if v == "" {
		return "0"
	}
	return v
}
