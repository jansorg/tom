package dataImport_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/jansorg/tom/go-tom/dataImport"
	"github.com/jansorg/tom/go-tom/dataImport/fanurio"
	"github.com/jansorg/tom/go-tom/dataImport/macTimeTracker"
	"github.com/jansorg/tom/go-tom/dataImport/watson"
	"github.com/jansorg/tom/go-tom/test_setup"
)

// malformed input must result in an error, not in a panic
func TestImportMalformedInput(t *testing.T) {
	tests := []struct {
		name     string
		importer dataImport.Handler
		content  string
	}{
		{"watson short frame", watson.NewImporter(), `[[1, 2]]`},
		{"watson start not a number", watson.NewImporter(), `[["a", 2, "p", "id", [], 3]]`},
		{"watson project not a string", watson.NewImporter(), `[[1, 2, 42, "id", [], 3]]`},
		{"watson tags not a list", watson.NewImporter(), `[[1, 2, "p", "id", "tag", 3]]`},
		{"watson tag not a string", watson.NewImporter(), `[[1, 2, "p", "id", [1], 3]]`},
		{"watson updated not a number", watson.NewImporter(), `[[1, 2, "p", "id", [], "x"]]`},
		{"fanurio too few columns", fanurio.NewCSVImporter(), "a,b,c\n1,2,3\n"},
		{"mac time tracker short duration", macTimeTracker.NewImporter(), "a;b;c;d;e;f;g\np;t;x;2015-12-18 17:00:37;y;23:17;notes\n"},
		{"mac time tracker empty duration", macTimeTracker.NewImporter(), "a;b;c;d;e;f;g\np;t;x;2015-12-18 17:00:37;y;;notes\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, err := test_setup.CreateTestContext(language.English)
			require.NoError(t, err)
			defer test_setup.CleanupTestContext(ctx)

			file := filepath.Join(t.TempDir(), "import")
			require.NoError(t, os.WriteFile(file, []byte(test.content), 0600))

			assert.NotPanics(t, func() {
				_, err = test.importer.Import(file, ctx)
			})
			assert.Error(t, err)
		})
	}
}

func TestImportWatsonOptionalValues(t *testing.T) {
	ctx, err := test_setup.CreateTestContext(language.English)
	require.NoError(t, err)
	defer test_setup.CleanupTestContext(ctx)

	file := filepath.Join(t.TempDir(), "frames")
	require.NoError(t, os.WriteFile(file, []byte(`[[1, 2, "p", "id", null], [3, 4, "p", "id", ["a"], 5]]`), 0600))

	result, err := watson.NewImporter().Import(file, ctx)
	require.NoError(t, err)
	assert.EqualValues(t, 2, result.CreatedFrames)
	assert.EqualValues(t, 1, result.CreatedTags)
}
