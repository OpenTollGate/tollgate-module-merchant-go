package main

import (
	"flag"
	"context"
	"fmt"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"strconv"
	"time"
	"runtime/debug"
)

var (
	Version    string
	CommitHash string
	BuildTime  string
)

var relayUrl = "ws://localhost:3334"

// var relayUrl = "wss://relay.damus.io"
var purchaseKind = 21000
var sessionKind = 22000
var valvePubkey = "714161f55b3198b6f95f1d23ca9ee8132052574f7785fcc859cb1f3cf2a2cf5f"

var privateKey = "f4be433e9648024b8d3ce6ab4798f0b8bfd87c3344a633a72af0fbdc6c352ac5"
var pubkey, _ = nostr.GetPublicKey(privateKey)
var nsec, _ = nip19.EncodePrivateKey(privateKey)
var npub, _ = nip19.EncodePublicKey(pubkey)


func getVersionInfo() string {
    if info, ok := debug.ReadBuildInfo(); ok {
        for _, setting := range info.Settings {
            switch setting.Key {
            case "vcs.revision":
                CommitHash = setting.Value[:7]
            case "vcs.time":
                BuildTime = setting.Value
            }
        }
    }
    return fmt.Sprintf("Version: %s\nCommit: %s\nBuild Time: %s", 
        Version, CommitHash, BuildTime)
}

func main() {
	// Add a version flag
	versionFlag := flag.Bool("version", false, "Print version information")
	flag.Parse()

	if *versionFlag {
		fmt.Println(getVersionInfo())
		return
	}

	fmt.Println("Starting Tollgate - merchant")
	fmt.Println("privateKey:", privateKey, "/", nsec)
	fmt.Println("pk:", pubkey, "/", npub)

	listenForPayments()

	fmt.Println("Shutting down Tollgate - merchant")
}

func listenForPayments() {

	backgroundContext := context.Background()
	relay, err := nostr.RelayConnect(backgroundContext, relayUrl)

	if err != nil {
		panic(err)
	}

	var filters nostr.Filters
	if _, _, err := nip19.Decode(npub); err == nil {
		filters = []nostr.Filter{{
			Kinds: []int{purchaseKind},
			//Limit: 10,
		}}
	} else {
		panic(err)
	}

	//backgroundContext, cancel := context.WithTimeout(backgroundContext, 3*time.Second)
	//defer cancel()

	subscription, err := relay.Subscribe(backgroundContext, filters)
	if err != nil {
		panic(err)
	}

	for evnt := range subscription.Events {
		handlePurchaseEvent(backgroundContext, evnt)
	}
}

func handlePurchaseEvent(backgroundContext context.Context, evnt *nostr.Event) {

	var customerPubKey = evnt.PubKey
	//var payment = evnt.Content

	var purchasedTimeSeconds = int64(60)

	var now = time.Now().Unix()
	sessionEndUnix := now + purchasedTimeSeconds

	var macAddress = evnt.Tags.GetFirst([]string{"mac"})

	if macAddress == nil {
		println("macAddress: nil")
		return
	}

	println("macAddress: " + macAddress.Value())
	println("Customer " + customerPubKey + " purchased " + strconv.FormatInt(purchasedTimeSeconds/int64(60), 10) + " min of access. Authenticating...")

	relay, err := nostr.RelayConnect(backgroundContext, relayUrl)

	if err != nil {
		fmt.Println(err)
	}

	sessionEvent := nostr.Event{
		PubKey:    pubkey,
		CreatedAt: nostr.Now(),
		Kind:      sessionKind,
		Tags: []nostr.Tag{
			{"p", valvePubkey},
			{"mac", macAddress.Value()},
			{"session-end", strconv.FormatInt(sessionEndUnix, 10)},
		},
		Content: "",
	}

	sessionEvent.Sign(privateKey)

	if err := relay.Publish(backgroundContext, sessionEvent); err != nil {
		fmt.Println(err)
	}

	println("Sent out session event to valve for MAC " + macAddress.Value())
}
