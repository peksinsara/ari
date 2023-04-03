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

func DialConference(client *ari.Client, endpoints []string) error {
	bridgeParams := ari.CreateBridgeParams{
		Type: "mixing",
		Name: "myBridge",
	}

	bridge, err := client.Bridges.Create(bridgeParams)
	if err != nil {
		return fmt.Errorf("error creating bridge: %s", err)
	}
	log.Printf("Created bridge %s", bridge.ID)

	var channels []*ari.Channel
	for _, endpoint := range endpoints {
		channelParams := ari.OriginateParams{
			Endpoint:  endpoint,
			Context:   "public",
			Extension: "s",
			App:       "myari",
		}

		channel, err := client.Channels.Create(channelParams)
		if err != nil {
			return fmt.Errorf("error creating channel for endpoint %s: %s", endpoint, err)
		}
		log.Printf("Created channel %s for endpoint %s", channel.ID, endpoint)

		for channel.State != "Up" {
			time.Sleep(time.Millisecond * 100)
			channel, err = client.Channels.Get(channel.ID)
			if err != nil {
				return fmt.Errorf("error getting channel %s: %s", channel.ID, err)
			}
		}

		channels = append(channels, channel)

		if err := bridge.AddChannel(channel.ID, ari.Participant); err != nil {
			return fmt.Errorf("error adding channel %s to bridge: %s", channel.ID, err)
		}
		log.Printf("Added channel %s to bridge %s", channel.ID, bridge.ID)
	}

	log.Printf("Started conference %s with %d participant(s)", bridge.ID, len(channels))

	for {
		bridge, err = client.Bridges.Get(bridge.ID)
		if err != nil {
			return fmt.Errorf("error getting bridge %s: %s", bridge.ID, err)
		}

		if bridge.BridgeClass == "destroyed" {
			log.Printf("Conference %s ended", bridge.ID)
			break
		}

		if len(bridge.Channels) == 0 {
			log.Printf("All participants have left conference %s, ending conference", bridge.ID)
			if err := bridge.Destroy(); err != nil {
				return fmt.Errorf("error destroying bridge %s: %s", bridge.ID, err)
			}
			break
		}

		time.Sleep(time.Second)
	}

	return nil
}
