package i18n

import (
	"testing"
	"time"

	"github.com/go-playground/locales/en_US"
	"github.com/stretchr/testify/require"
)

func Test_isoTime(t *testing.T) {
	date := time.Date(2020, time.January, 12, 9, 0, 0, 0, time.UTC)

	f := isoDelegate{locateDelegate{en_US.New()}}
	require.EqualValues(t, "09:00", f.FmtTimeShort(date))
	require.EqualValues(t, "09:00:00", f.FmtTimeMedium(date))
	require.EqualValues(t, "09:00:00", f.FmtTimeLong(date))
}
