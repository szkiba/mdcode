package mdcode

func Unfence(source []byte) (Blocks, error) {
	var blocks Blocks

	_, _, err := Walk(source, func(block *Block) error {
		blocks = append(blocks, block)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return blocks, nil
}
