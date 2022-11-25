package extract_test

import (
	_ "embed"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bastjan/declextract/extract"
)

//go:embed testdata/basic.go
var basicDeclarations string

func Test(t *testing.T) {
	testCases := []struct {
		desc        string
		src         string
		declaration string
		want        any
	}{
		{
			desc:        "Extract string",
			src:         basicDeclarations,
			declaration: "str",
			want:        "str",
		},
		{
			desc:        "Extract string multiline",
			src:         basicDeclarations,
			declaration: "strMultiline",
			want:        "\na:\n  b: c\n",
		},
		{
			desc:        "Extract string UTF8",
			src:         basicDeclarations,
			declaration: "strUTF8",
			want:        "\t\n\r 🚵‍♀️",
		},
		{
			desc:        "Extract float",
			src:         basicDeclarations,
			declaration: "f",
			want:        1.7,
		},
		{
			desc:        "Extract int",
			src:         basicDeclarations,
			declaration: "i",
			want:        int64(1),
		},
		{
			desc:        "Extract int octal",
			src:         basicDeclarations,
			declaration: "iOctal",
			want:        int64(0o644),
		},
		{
			desc:        "Extract complex",
			src:         basicDeclarations,
			declaration: "c",
			want:        1i,
		},
		{
			desc:        "Extract rune",
			src:         basicDeclarations,
			declaration: "r",
			want:        'r',
		},
		{
			desc:        "Extract variable string",
			src:         basicDeclarations,
			declaration: "strVar",
			want:        "str",
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			td := t.TempDir()
			require.NoError(t, os.WriteFile(td+"/src.go", []byte(tC.src), 0o644))

			res, err := extract.ExtractDeclarationFromFile(td+"/src.go", tC.declaration)
			require.NoError(t, err)
			assert.Equal(t, tC.want, res)
		})
	}
}

func TestNotFoundError(t *testing.T) {
	t.Parallel()

	td := t.TempDir()
	require.NoError(t, os.WriteFile(td+"/src.go", []byte("package blub\n"), 0o644))

	_, err := extract.ExtractDeclarationFromFile(td+"/src.go", "notFound")

	assert.ErrorIs(t, err, &extract.NotFoundError{})
	assert.ErrorIs(t, err, extract.NotFoundError{})
	assert.ErrorContains(t, err, "notFound")
}

func TestMangledFile(t *testing.T) {
	t.Parallel()

	td := t.TempDir()
	require.NoError(t, os.WriteFile(td+"/src.go", []byte("asdasfg\n"), 0o644))

	_, err := extract.ExtractDeclarationFromFile(td+"/src.go", "notFound")
	assert.Error(t, err)
}

func TestUnsupported(t *testing.T) {
	t.Parallel()

	const src = `
package blub
var function = func() {}
`

	td := t.TempDir()
	require.NoError(t, os.WriteFile(td+"/src.go", []byte(src), 0o644))

	_, err := extract.ExtractDeclarationFromFile(td+"/src.go", "function")
	assert.Error(t, err)
}
