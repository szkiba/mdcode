package mdcode

import (
	"bytes"
	"regexp"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

var reInfo = regexp.MustCompile(`\s*(\w+)\s*(.*)\s*`)

type Walker func(block *Block) error

type change struct {
	fcb   *ast.FencedCodeBlock
	block *Block
}

func (c *change) bounds() (int, int) {
	lines := c.fcb.Lines()
	if lines.Len() == 0 {
		return c.fcb.Info.Segment.Stop + 1, c.fcb.Info.Segment.Stop + 1
	}

	return lines.At(0).Start, lines.At(lines.Len() - 1).Stop
}

func (c *change) sizeIncrement() int {
	start, stop := c.bounds()

	return len(c.block.Code) - (stop - start)
}

func Walk(source []byte, walker Walker) (bool, []byte, error) {
	parser := goldmark.DefaultParser()
	reader := text.NewReader(source)
	root := parser.Parse(reader).OwnerDocument()

	var changes []*change

	err := ast.Walk(root, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		node = transformCommentedCodeBlock(node, entering, source)

		fcb := asFencedCodeBlock(node, entering)
		if fcb == nil {
			return ast.WalkContinue, nil
		}

		block, berr := extractBlock(fcb, source)
		if berr != nil {
			return ast.WalkContinue, berr
		}

		code := block.Code

		berr = walker(block)
		if berr != nil {
			return ast.WalkContinue, berr
		}

		if !bytes.Equal(code, block.Code) {
			changes = append(changes, &change{fcb: fcb, block: block})
		}

		return ast.WalkContinue, nil
	})
	if err != nil {
		return false, nil, err
	}

	if len(changes) == 0 {
		return false, nil, nil
	}

	return true, applyChanges(changes, source), nil
}

func asFencedCodeBlock(node ast.Node, entering bool) *ast.FencedCodeBlock {
	if entering || node.Kind() != ast.KindFencedCodeBlock {
		return nil
	}

	if fcb, ok := node.(*ast.FencedCodeBlock); ok {
		return fcb
	}

	return nil
}

func extractBlock(fcb *ast.FencedCodeBlock, source []byte) (*Block, error) {
	lang, meta, err := extractInfo(fcb, source)
	if err != nil {
		return nil, err
	}

	return &Block{Lang: lang, Meta: meta, Code: extractCode(fcb, source)}, nil
}

func extractCode(fcb *ast.FencedCodeBlock, source []byte) []byte {
	var buff bytes.Buffer

	lines := fcb.Lines()
	for i := 0; i < lines.Len(); i++ {
		seg := lines.At(i)

		buff.Write(seg.Value(source))
	}

	return buff.Bytes()
}

func extractInfo(fcb *ast.FencedCodeBlock, source []byte) (string, Meta, error) {
	if fcb.Info == nil {
		return "", nil, nil
	}

	return parseInfo(fcb.Info.Text(source))
}

func parseInfo(text []byte) (string, Meta, error) {
	all := reInfo.FindSubmatch(text)
	if all == nil {
		return "", nil, nil
	}

	var (
		lang string
		meta Meta
		err  error
	)

	if len(all) > 1 {
		lang = string(all[1])
	}

	if len(all) <= 2 { //nolint:gomnd
		return lang, meta, nil
	}

	meta, err = parseMeta(all[2])

	return lang, meta, err
}

func applyChanges(changes []*change, source []byte) []byte {
	resSize := len(source)

	for _, change := range changes {
		resSize += change.sizeIncrement()
	}

	result := make([]byte, resSize)

	var srcIdx, resIdx int

	for _, change := range changes {
		start, stop := change.bounds()

		copy(result[resIdx:], source[srcIdx:start])
		resIdx += (start - srcIdx)

		copy(result[resIdx:], change.block.Code)
		resIdx += len(change.block.Code)

		srcIdx = stop
	}

	copy(result[resIdx:], source[srcIdx:])

	return result
}

var (
	reCommentedCodeBlock = regexp.MustCompile(`^\s*(<!--)?\s*<script\s*type=["']text/markdown["']\s*>\s*$`)
	reFences             = regexp.MustCompile("^\\s*```")
)

func transformCommentedCodeBlock(node ast.Node, entering bool, source []byte) ast.Node { //nolint:ireturn
	if entering || node.Kind() != ast.KindHTMLBlock {
		return node
	}

	html, ok := node.(*ast.HTMLBlock)
	if !ok {
		return node
	}

	const minLines = 2

	lines := html.Lines()
	if lines.Len() < minLines {
		return node
	}

	seg := lines.At(0)
	line := seg.Value(source)

	if !reCommentedCodeBlock.Match(line) {
		return node
	}

	seg = lines.At(1)
	line = seg.Value(source)

	loc := reFences.FindIndex(line)
	if loc == nil {
		return node
	}

	info := ast.NewTextSegment(text.NewSegment(seg.Start+loc[1], seg.Stop-1))
	fcb := ast.NewFencedCodeBlock(info)

	seg = lines.At(lines.Len() - 1)
	line = seg.Value(source)

	if !reFences.Match(line) {
		return node
	}

	segs := text.NewSegments()

	for i := 2; i < lines.Len()-1; i++ {
		segs.Append(lines.At(i))
	}

	fcb.SetLines(segs)

	return fcb
}
