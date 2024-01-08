package mdcode

import (
	"bytes"
	"embed"
	"io/fs"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed testdata/testdoc.md
var testdoc []byte

//go:embed testdata/testdocmod.md
var testdocmod []byte

//go:embed testdata
var testfs embed.FS

func testBlocks(t *testing.T, filter func(*Block) bool) Blocks {
	t.Helper()

	var blocks Blocks

	mod, data, err := Walk(testdoc, func(block *Block) error {
		if filter(block) {
			blocks = append(blocks, block)
		}

		return nil
	})

	require.NoError(t, err)
	require.False(t, mod)
	require.Nil(t, data)

	return blocks
}

func Test_Walk_entire(t *testing.T) {
	t.Parallel()

	blocks := testBlocks(t, func(b *Block) bool { return strings.HasPrefix(b.Meta.Get("file"), "entire") })

	require.Len(t, blocks, 3)

	for _, block := range blocks {
		file := block.Meta.Get("file")
		require.NotEmpty(t, file)

		require.Len(t, block.Meta, 1)

		want, err := fs.ReadFile(testfs, path.Join("testdata", file))
		require.NoError(t, err)

		require.Equal(t, want, block.Code)
	}
}

func Test_Walk_partial(t *testing.T) {
	t.Parallel()

	blocks := testBlocks(t, func(b *Block) bool { return strings.HasPrefix(b.Meta.Get("file"), "partial") })

	require.Len(t, blocks, 2)

	for _, block := range blocks {
		file := block.Meta.Get("file")
		require.NotEmpty(t, file)

		require.Len(t, block.Meta, 2)

		region := block.Meta.Get("region")
		outline := block.Meta.Get("outline")

		if len(region) == 0 {
			require.Equal(t, "true", outline)
		} else {
			require.Equal(t, "function", region)
		}
	}
}

func Test_Walk_mod(t *testing.T) {
	t.Parallel()

	eol := "\n"
	if runtime.GOOS == "windows" {
		eol = "\r\n"
	}

	mod, got, err := Walk(testdoc, func(block *Block) error {
		if !strings.HasPrefix(block.Meta.Get("file"), "entire") {
			return nil
		}

		var buff bytes.Buffer

		buff.WriteString("/*" + eol)
		buff.Write(block.Code)
		buff.WriteString("*/" + eol)

		block.Code = buff.Bytes()

		return nil
	})

	require.NoError(t, err)
	require.True(t, mod)
	require.NotNil(t, got)

	require.Equal(t, testdocmod, got)
}
