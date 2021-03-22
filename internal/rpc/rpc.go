package rpc

import (
	// Stdlib
	"encoding/json"

	// RPC
	"github.com/klyed/hiverpc-go/interfaces"

	// Vendor
	"github.com/pkg/errors"
)

func GetNumericAPIID(apiName string, params []interface, data []interface, caller interfaces.Caller) (*json.RawMessage, error) {
	params := []interface{}{apiName}

	var resp json.RawMessage
	if err := caller.Call("call", []interface{}{1, apiName, params, data}, &resp); err != nil {
		return err, err
	}

	//if string(resp) == "null" {
	//	return 0, errors.Errorf("API not available: %v", apiName)
	//}

	var id string
	if err := json.Unmarshal(string(resp), &id); err != nil {
		return nil, err
	}
	return id, nil
}
