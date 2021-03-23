package types

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"strings"

	"github.com/klyed/hiverpc-go/encoding/transaction"
	"github.com/pkg/errors"
)

var (
	TypeFollow         = "follow"
	TypeReblog         = "reblog"
	TypeLogin          = "login"
	TypePrivateMessage = "private_message"
)

var customJSONDataObjects = map[string]interface{}{
	TypeFollow:         &FollowOperation{},
	TypeReblog:         &ReblogOperation{},
	TypeLogin:          &LoginOperation{},
	TypePrivateMessage: &PrivateMessageOperation{},
}

// FC_REFLECT( steemit::chain::custom_json_operation,
//             (required_auths)
//             (required_posting_auths)
//             (id)
//             (json) )

//CustomJSONOperation represents custom_json operation data.
type CustomJSONOperation struct {
	RequiredAuths        []string `json:"required_auths"`
	RequiredPostingAuths []string `json:"required_posting_auths"`
	ID                   string   `json:"id"`
	JSON                 string   `json:"json"`
}

//Type function that defines the type of operation.
func (op *CustomJSONOperation) Type() OpType {
	return TypeCustomJSON
}

//Data returns the operation data.
func (op *CustomJSONOperation) Data() interface{} {
	return op
}

//UnmarshalData unpacking the JSON parameter in the CustomJSONOperation type.
func (op *CustomJSONOperation) UnmarshalData() (interface{}, error) {
	// Get the corresponding data object template.
	template, ok := customJSONDataObjects[op.ID]
	if !ok {
		// In case there is no corresponding template, return nil.
		return nil, nil
	}

	// Clone the template.
	opData := reflect.New(reflect.Indirect(reflect.ValueOf(template)).Type()).Interface()

	// Prepare the whole operation tuple.
	var bodyReader io.Reader
	if op.JSON[0] == '[' {
		rawTuple := make([]json.RawMessage, 2)
		if err := json.NewDecoder(strings.NewReader(op.JSON)).Decode(&rawTuple); err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal CustomJSONOperation.JSON: \n%v", op.JSON)
		}
		if len(rawTuple) < 2 || rawTuple[1] == nil {
			return nil, errors.Errorf("invalid CustomJSONOperation.JSON: \n%v", op.JSON)
		}
		bodyReader = bytes.NewReader([]byte(rawTuple[1]))
	} else {
		bodyReader = strings.NewReader(op.JSON)
	}

	// Unmarshal into the new object instance.
	if err := json.NewDecoder(bodyReader).Decode(opData); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal CustomJSONOperation.JSON: \n%v", op.JSON)
	}

	return opData, nil
}

//MarshalTransaction is a function of converting type CustomJSONOperation to bytes.
func (op *CustomJSONOperation) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.EncodeUVarint(uint64(TypeCustomJSON.Code()))
	enc.EncodeArrString(op.RequiredAuths)
	enc.EncodeArrString(op.RequiredPostingAuths)
	enc.Encode(op.ID)
	enc.Encode(op.JSON)
	return enc.Err()
}

//FollowOperation the structure for the operation CustomJSONOperation.
type FollowOperation struct {
	Follower  string   `json:"follower"`
	Following string   `json:"following"`
	What      []string `json:"what"`
}

//ReblogOperation the structure for the operation CustomJSONOperation.
type ReblogOperation struct {
	Account  string `json:"account"`
	Author   string `json:"author"`
	Permlink string `json:"permlink"`
}

//LoginOperation the structure for the operation CustomJSONOperation.
type LoginOperation struct {
	Account string `json:"account"`
}

//PrivateMessageOperation the structure for the operation CustomJSONOperation.
type PrivateMessageOperation struct {
	From             string `json:"from"`
	To               string `json:"to"`
	FromMemoKey      string `json:"from_memo_key"`
	ToMemoKey        string `json:"to_memo_key"`
	SentTime         uint64 `json:"sent_time"`
	Checksum         uint32 `json:"checksum"`
	EncryptedMessage string `json:"encrypted_message"`
}

//MarshalCustomJSON generate a row from the structure fields.
func MarshalCustomJSON(v interface{}) (string, error) {
	var tmp []interface{}

	typeInterface := reflect.TypeOf(v).Name()
	switch typeInterface {
	case "FollowOperation":
		tmp = append(tmp, TypeFollow)
	case "ReblogOperation":
		tmp = append(tmp, TypeReblog)
	case "LoginOperation":
		tmp = append(tmp, TypeLogin)
	case "PrivateMessageOperation":
		tmp = append(tmp, TypePrivateMessage)
	default:
		return "", errors.New("Unknown type")
	}

	tmp = append(tmp, v)

	b, err := json.Marshal(tmp)
	if err != nil {
		return "", err
	}

	return string(b), nil //strings.Replace(string(b), "\"", "\\\"", -1), nil
}

/*
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

var customJSONDataObjects = map[OpType][]string {
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
	operations []string `json:"operations"`
	required_auths       []string `json:"op.required_auths"`
	id OpType `json:"op.id"`
	json string `json:"op.json"`
	RequiredPostingAuths []string `json:"op.required_posting_auths"`
	extensions string `json:"extensions"`
	signatures string `json:"signatures"`
	transaction_id	string	`json:"transaction_id"`
	block_num	int	`json:"block_num"`
	transaction_num	int	`json:"transaction_num"`
}

type CustomJSONOperations []interface {
	Type() OpType
	Data() interface{}
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
*/
