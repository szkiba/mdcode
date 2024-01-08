package cmd

import "github.com/szkiba/mdcode/internal/mdcode"

func walk(source []byte, walker mdcode.Walker, filter filterFunc) (bool, []byte, error) {
	return mdcode.Walk(source, func(block *mdcode.Block) error {
		if filter(block.Lang, block.Meta) {
			return walker(block)
		}

		return nil
	})
}
