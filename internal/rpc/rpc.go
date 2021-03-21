package rpc

import (
	// Stdlib
	"encoding/json"

	// RPC
	"github.com/klyed/hiverpc-go/interfaces"

	// Vendor
	"github.com/pkg/errors"
)

func GetNumericAPIID(caller interfaces.Caller, apiName string) (int, error) {
	params := []interface{}{apiName}

	var resp json.RawMessage
	if err := caller.Call("call", []interface{}{params}, &resp); err != nil {
		return 0, err
	}

	if string(resp) == "null" {
		return 0, errors.Errorf("API not available: %v", apiName)
	}

	var id int
	if err := json.Unmarshal([]byte(resp), &id); err != nil {
		return 0, err
	}
	return id, nil
}
