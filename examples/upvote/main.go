package main

import (
	// Stdlib
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	// RPC
	"github.com/klyed/hiverpc-go"
	"github.com/klyed/hiverpc-go/encoding/wif"
	"github.com/klyed/hiverpc-go/transactions"
	"github.com/klyed/hiverpc-go/transports/websocket"
	"github.com/klyed/hiverpc-go/types"

	// Vendor
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln("Error:", err)
	}
}

func run() (err error) {
	// Process flags.
	flagAddress := flag.String("rpc_endpoint", "ws://localhost:8090", "hived RPC endpoint address")
	flag.Parse()

	url := *flagAddress

	// Process args.
	args := flag.Args()
	if len(args) != 3 {
		return errors.New("3 arguments required")
	}
	author, permlink, voter := args[0], args[1], args[2]

	// Prompt for WIF.
	wifKey, err := promptWIF(voter)
	if err != nil {
		return err
	}

	// Start catching signals.
	var interrupted bool
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// Drop the error in case it is a request being interrupted.
	defer func() {
		if err == websocket.ErrClosing && interrupted {
			err = nil
		}
	}()

	// Instantiate the WebSocket transport.
	t, err := websocket.NewTransport(url)
	if err != nil {
		return err
	}

	// Use the transport to get an RPC client.
	client, err := rpc.NewClient(t)
	if err != nil {
		return err
	}
	defer func() {
		if !interrupted {
			client.Close()
		}
	}()

	// Start processing signals.
	go func() {
		<-signalCh
		fmt.Println()
		log.Println("Signal received, exiting...")
		signal.Stop(signalCh)
		interrupted = true
		client.Close()
	}()

	// Get the props to get the head block number and ID
	// so that we can use that for the transaction.
	props, err := client.Database.GetDynamicGlobalProperties()
	if err != nil {
		return err
	}

	// Prepare the transaction.
	refBlockPrefix, err := transactions.RefBlockPrefix(props.HeadBlockID)
	if err != nil {
		return err
	}

	tx := transactions.NewSignedTransaction(&types.Transaction{
		RefBlockNum:    transactions.RefBlockNum(props.HeadBlockNumber),
		RefBlockPrefix: refBlockPrefix,
	})

	tx.PushOperation(&types.VoteOperation{
		Voter:    voter,
		Author:   author,
		Permlink: permlink,
		Weight:   10000,
	})

	// Sign.
	privKey, err := wif.Decode(wifKey)
	if err != nil {
		return err
	}
	privKeys := [][]byte{privKey}

	if err := tx.Sign(privKeys, transactions.HiveChain); err != nil {
		return err
	}

	// Broadcast.
	resp, err := client.NetworkBroadcast.BroadcastTransactionSynchronous(tx.Transaction)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", *resp)

	// Success!
	return nil
}

func promptWIF(accountName string) (string, error) {
	fmt.Printf("Please insert WIF for account @%v: ", accountName)
	passwd, err := terminal.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", errors.Wrap(err, "failed to read WIF from the terminal")
	}
	fmt.Println()
	return string(passwd), nil
}
