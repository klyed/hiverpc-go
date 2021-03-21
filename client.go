package rpc

import (
	// RPC
	"github.com/klyed/hiverpc-go/apis/database"
	"github.com/klyed/hiverpc-go/apis/networkbroadcast"
	"github.com/klyed/hiverpc-go/interfaces"
)

// Client can be used to access Steem remote APIs.
//
// There is a public field for every Steem API available,
// e.g. Client.Database corresponds to database_api.
type Client struct {
	cc interfaces.CallCloser

	// Login represents login_api.
	//Login *login.API

	// Database represents database_api.
	Database *database.API

	// NetworkBroadcast represents network_broadcast_api.
	NetworkBroadcast *networkbroadcast.API
}

// NewClient creates a new RPC client that use the given CallCloser internally.
func NewClient(cc interfaces.CallCloser) (*Client, error) {
	client := &Client{cc: cc}

	client.Database = database.NewAPI(client.cc)

	if err != nil {
		return nil, err
	}

	networkBroadcastAPI, err := networkbroadcast.NewAPI(client.cc)
	if err != nil {
		return nil, err
	}
	client.NetworkBroadcast = networkBroadcastAPI

	return client, nil
}

// Close should be used to close the client when no longer needed.
// It simply calls Close() on the underlying CallCloser.
func (client *Client) Close() error {
	return client.cc.Close()
}
