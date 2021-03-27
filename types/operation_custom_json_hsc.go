package types

type HscOperation struct {
	ref_block_num        int         `json:"ref_block_num"`
	ref_block_prefix     int         `json:"ref_block_prefix"`
	expiration           string      `json:"expiration"`
	operations           interface{} `json:"op[0]"`
	required_auths       []string    `json:"op[0].required_auths"`
	id                   OpType      `json:"op[0].id"`
	RequiredPostingAuths []string    `json:"op[0].required_posting_auths"`
	extensions           string      `json:"extensions"`
	signatures           string      `json:"signatures"`
	transaction_id       string      `json:"transaction_id"`
	block_num            int         `json:"block_num"`
	transaction_num      int         `json:"transaction_num"`
}

func (op *HscOperation) Type() OpType {
	return op.id
}

func (op *HscOperation) Data() interface{} {
	return op.operations
}
