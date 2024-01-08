package mdcode

type Block struct {
	Lang string
	Meta Meta
	Code []byte
}

type Blocks []*Block
