package blockchain

type Transaction struct {
	ID      []byte
	Inputs  []TInput
	Outputs []TOutput
}

type TOutput struct {
	Value  int
	PubKey string
}

type TInput struct {
	ID  []byte
	Out int
	Sig string
}
