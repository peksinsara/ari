package functions

import (
	"fmt"
	"log"
	"time"

	"github.com/abourget/ari"
)

func DialEndpoint(client *ari.Client, endpoints []string, extensions []string, direction string) error {
	var channels []*ari.Channel

	for i, endpoint := range endpoints {
		params := ari.OriginateParams{
			Endpoint:  endpoint,
			Context:   "public",
			Extension: extensions[i],
			App:       "myari",
		}

		channel, err := client.Channels.Create(params)
		if err != nil {
			return fmt.Errorf("error creating channel for endpoint %s: %s", endpoint, err)
		}
		log.Printf("Created channel %s for endpoint %s with extension %s", channel.ID, endpoint, extensions[i])

		for channel.State != "Up" {
			time.Sleep(time.Millisecond * 100)
			channel, err = client.Channels.Get(channel.ID)
			if err != nil {
				return fmt.Errorf("error getting channel %s: %s", channel.ID, err)
			}
		}

		channels = append(channels, channel)
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

	for _, channel := range channels {
		if err := bridge.AddChannel(channel.ID, ari.Participant); err != nil {
			return fmt.Errorf("error adding channel %s to bridge: %s", channel.ID, err)
		}
		log.Printf("Added channel %s to bridge %s", channel.ID, bridge.ID)
	}

	return nil
}
