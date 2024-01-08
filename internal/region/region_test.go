package region_test

import (
	"embed"
	"io/fs"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/szkiba/mdcode/internal/region"
)

//go:embed testdata/testdoc.js
var testdoc []byte

//go:embed testdata/testdocmod.js
var testdocmod []byte

//go:embed testdata/testdocoutline.js
var testdocoutline []byte

//go:embed testdata
var testfs embed.FS

func Test_Outline(t *testing.T) {
	t.Parallel()

	got, ok, err := region.Outline(testdoc)

	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, testdocoutline, got)
}

func Test_Read(t *testing.T) {
	t.Parallel()

	got, found, err := region.Read(testdoc, "empty")

	require.NoError(t, err)
	require.True(t, found)
	require.NotNil(t, got)
	require.Empty(t, got)

	got, found, err = region.Read(testdoc, "nonempty")

	require.NoError(t, err)
	require.True(t, found)

	want, err := fs.ReadFile(testfs, "testdata/nonempty.js")
	require.NoError(t, err)

	require.Equal(t, want, got)

	got, found, err = region.Read(testdoc, "block")

	require.NoError(t, err)
	require.True(t, found)

	want, err = fs.ReadFile(testfs, "testdata/block.js")
	require.NoError(t, err)

	require.Equal(t, want, got)
}

func Test_Replace(t *testing.T) {
	t.Parallel()

	eol := "\n"
	if runtime.GOOS == "windows" {
		eol = "\r\n"
	}

	begin := "/* begin */" + eol
	end := "/* end */" + eol

	data, found, err := region.Replace(testdoc, "empty", []byte(begin+end))
	require.NoError(t, err)
	require.True(t, found)

	body, _, _ := region.Read(data, "nonempty")

	data, found, err = region.Replace(data, "nonempty", []byte(begin+string(body)+end))
	require.NoError(t, err)
	require.True(t, found)

	body, _, _ = region.Read(data, "block")

	data, found, err = region.Replace(data, "block", []byte(begin+string(body)+end))
	require.NoError(t, err)
	require.True(t, found)

	require.Equal(t, string(testdocmod), string(data))
}
