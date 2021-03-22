package types

import (
	// Stdlib
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"strings"

	// Vendor
	"github.com/pkg/errors"
)

var customJSONDataObjects = map[OpType]HscOperation {
	//map[OpType]: HscOperation{}{}
}




// FC_REFLECT( hive::chain::custom_json_operation,
//             (required_auths)
//             (required_posting_auths)
//             (id)
//             (json) )

// CustomJSONOperation represents custom_json operation data.
type CustomJSONOperation struct {
	ref_block_num int `json:"ref_block_num"`
	ref_block_prefix int `json:"ref_block_prefix"`
	expiration string `json:"expiration"`
	operations []string `json:"op[0]"`
	required_auths       []string `json:"op[0].required_auths"`
	id OpType `json:"op[0].id"`
	json string `json:"op[0].json"`
	RequiredPostingAuths []string `json:"op[0].required_posting_auths"`
	extensions string `json:"extensions"`
	signatures string `json:"signatures"`
	transaction_id	string	`json:"transaction_id"`
	block_num	int	`json:"block_num"`
	transaction_num	int	`json:"transaction_num"`
}

func (op *CustomJSONOperation) Type() OpType {
	return op.id
}

func (op *CustomJSONOperation) Data() interface{} {
	return op.operations
}

func (op *CustomJSONOperation) UnmarshalData() (interface{}, error) {
	// Get the corresponding data object template.
	template, ok := customJSONDataObjects[op.id]
	if !ok {
		// In case there is no corresponding template, return nil.
		return nil, nil
	}

	// Clone the template.
	opData := reflect.New(reflect.Indirect(reflect.ValueOf(template)).Type()).Interface()

	// Prepare the whole operation tuple.
	var bodyReader io.Reader
	if op.json[0] == '[' {
		rawTuple := make([]json.RawMessage, 2)
		if err := json.NewDecoder(strings.NewReader(op.json)).Decode(&rawTuple); err != nil {
			return nil, errors.Wrapf(err,
				"failed to unmarshal CustomJSONOperation.JSON: \n%v", op.json)
		}
		if rawTuple[1] == nil {
			return nil, errors.Errorf("invalid CustomJSONOperation.JSON: \n%v", op.json)
		}
		bodyReader = bytes.NewReader([]byte(rawTuple[1]))
	} else {
		bodyReader = strings.NewReader(op.json)
	}

	// Unmarshal into the new object instance.
	if err := json.NewDecoder(bodyReader).Decode(opData); err != nil {
		return nil, errors.Wrapf(err,
			"failed to unmarshal CustomJSONOperation.JSON: \n%v", op.json)
	}

	return opData, nil
}
