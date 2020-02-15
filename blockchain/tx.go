package blockchain

// TxOutput is the output part of a transaction, containing a value and a public key
type TxOutput struct {
	Value  int
	PubKey string
}

// TxInput is the input part of a transaction, containing an ID, an out value and a signature
type TxInput struct {
	ID  []byte
	Out int
	Sig string
}

// CanUnlock checks whether an input can unlock some given data
func (i *TxInput) CanUnlock(d string) bool {
	return i.Sig == d
}

// CanBeUnlocked checks whether an output can be unlocked
func (o *TxOutput) CanBeUnlocked(d string) bool {
	return o.PubKey == d
}
