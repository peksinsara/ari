package functions

import (
	"fmt"
	"log"
	"time"

	"github.com/abourget/ari"
)

func DialEndpoint(client *ari.Client, callerEndpoint string, calleeEndpoint string) error {
	callerParams := ari.OriginateParams{
		Endpoint:  callerEndpoint,
		Context:   "public",
		Extension: "s",
		App:       "myari",
	}

	callerChannel, err := client.Channels.Create(callerParams)
	if err != nil {
		return fmt.Errorf("error creating channel for caller %s: %s", callerEndpoint, err)
	}
	log.Printf("Created channel %s for caller %s", callerChannel.ID, callerEndpoint)

	for callerChannel.State != "Up" {
		time.Sleep(time.Millisecond * 100)
		callerChannel, err = client.Channels.Get(callerChannel.ID)
		if err != nil {
			return fmt.Errorf("error getting channel %s: %s", callerChannel.ID, err)
		}
	}

	calleeParams := ari.OriginateParams{
		Endpoint:  calleeEndpoint,
		Context:   "public",
		Extension: "s",
		App:       "myari",
	}

	calleeChannel, err := client.Channels.Create(calleeParams)
	if err != nil {
		return fmt.Errorf("error creating channel for callee %s: %s", calleeEndpoint, err)
	}
	log.Printf("Created channel %s for callee %s", calleeChannel.ID, calleeEndpoint)

	for calleeChannel.State != "Up" {
		time.Sleep(time.Millisecond * 100)
		calleeChannel, err = client.Channels.Get(calleeChannel.ID)
		if err != nil {
			return fmt.Errorf("error getting channel %s: %s", calleeChannel.ID, err)
		}
	}

	bridgeParams := ari.CreateBridgeParams{
		Type: "mixing",
		Name: "myBridge",
	}

	bridge, err := client.Bridges.Create(bridgeParams)
	if err != nil {
		return fmt.Errorf("error creating bridge: %s", err)
	}
	log.Printf("Created bridge %s", bridge.ID)

	if err := bridge.AddChannel(callerChannel.ID, ari.Participant); err != nil {
		return fmt.Errorf("error adding caller channel %s to bridge: %s", callerChannel.ID, err)
	}
	log.Printf("Added caller channel %s to bridge %s", callerChannel.ID, bridge.ID)

	if err := bridge.AddChannel(calleeChannel.ID, ari.Participant); err != nil {
		return fmt.Errorf("error adding callee channel %s to bridge: %s", calleeChannel.ID, err)
	}
	log.Printf("Added callee channel %s to bridge %s", calleeChannel.ID, bridge.ID)

	return nil
}
