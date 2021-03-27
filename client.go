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

	cc interfaces.CallCloser
	Database *database.API

}

func (client *Client) Deadline() (deadline time.Time, ok bool) {
	panic("implement me")
	return deadline, true
}

func (client *Client) Done() *Client{
	panic("implement me")
	return client.cc.Close()
}

func (client *Client) Err() error {
	panic("implement me")
	return client.cc.Close()
}


// NewClient creates a new RPC client that use the given CallCloser internally.
func NewClient(cc interfaces.CallCloser) (*Client, error) {
	client := &Client{cc: cc}
	//blank := [][]string{}
	paraminner := []interface{}{"database_api", "get_dynamic_global_properties"}
	params := []interface{}{paraminner}

  client.Database = database.NewAPI("call", params, client.cc)

	//networkBroadcastAPI, err := networkbroadcast.NewAPI("network_broadcast_api", ",", client.cc)
	//if err != nil {
	//	return nil, err
	//}
	//client.NetworkBroadcast = networkBroadcastAPI

	return client, nil
}

// Close should be used to close the client when no longer needed.
// It simply calls Close() on the underlying CallCloser.
func (client *Client) Close() error {
	return client.cc.Close()
}
