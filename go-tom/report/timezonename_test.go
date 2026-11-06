package report

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimezoneNameUnmarshal(t *testing.T) {
	var tz TimezoneName
	require.NoError(t, json.Unmarshal([]byte(`"UTC+2"`), &tz))
	assert.Equal(t, 2*60*60, func() int { _, offset := tzNow(tz); return offset }())

	require.NoError(t, json.Unmarshal([]byte(`"Europe/Berlin"`), &tz))
	assert.Error(t, json.Unmarshal([]byte(`"No/Such_Zone"`), &tz))
}

func tzNow(tz TimezoneName) (string, int) {
	return time.Date(2024, 1, 1, 0, 0, 0, 0, tz.AsTimezone()).Zone()
}
