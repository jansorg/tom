package report

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeArgs(t *testing.T) {
	assert.Equal(t, []string{"report", "--month=-1"}, NormalizeArgs([]string{"report", "--month", "-1"}))
	assert.Equal(t, []string{"report", "--month=-12", "--json"}, NormalizeArgs([]string{"report", "-m", "-12", "--json"}))
	assert.Equal(t, []string{"report", "--month", "--json"}, NormalizeArgs([]string{"report", "--month", "--json"}))
	assert.Equal(t, []string{"report", "--month"}, NormalizeArgs([]string{"report", "--month"}))
	assert.Equal(t, []string{"report", "--", "--month", "-1"}, NormalizeArgs([]string{"report", "--", "--month", "-1"}))
}
