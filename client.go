package rpc

import (
	// RPC
	"github.com/klyed/hiverpc-go/apis/database"
	//"github.com/klyed/hiverpc-go/apis/networkbroadcast"
	"github.com/klyed/hiverpc-go/interfaces"
	"time"
)

// Client can be used to access Steem remote APIs.
//
// There is a public field for every Steem API available,
// e.g. Client.Database corresponds to database_api.
type Client struct {
	cc       interfaces.CallCloser
	Database *database.API
}

func (Client *Client) Deadline() (deadline time.Time, ok bool) {
	return deadline, true
}

func (Client *Client) Done() interfaces.CallCloser {
	return Client.cc
}

func (Client *Client) Err() interfaces.CallCloser {
	return Client.Done()
}

// NewClient creates a new RPC client that use the given CallCloser internally.
func NewClient(cc interfaces.CallCloser) (*Client, error) {
	Client := &Client{cc: cc}
	//blank := [][]string{}
	paraminner := []interface{}{"database_api", "get_dynamic_global_properties"}
	params := []interface{}{paraminner}

	Client.Database = database.NewAPI("call", params, Client.cc)

	//networkBroadcastAPI, err := networkbroadcast.NewAPI("network_broadcast_api", ",", client.cc)
	//if err != nil {
	//	return nil, err
	//}
	//client.NetworkBroadcast = networkBroadcastAPI

	return Client, nil
}

// Close should be used to close the client when no longer needed.
// It simply calls Close() on the underlying CallCloser.
func (Client *Client) Close() error {
	return Client.cc.Close()
}
