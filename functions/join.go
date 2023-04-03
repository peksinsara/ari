package functions

import (
	"fmt"
	"log"
	"time"

	"github.com/abourget/ari"
)

// JoinCall joins a call by adding endpoints to an existing bridge
func JoinCall(client *ari.Client, callID string, endpointIDs []string) error {
	bridge, err := client.Bridges.Get(callID)
	if err != nil {
		return fmt.Errorf("failed to get bridge %s: %s", callID, err)
	}

	for _, endpointID := range endpointIDs {
		channelParams := ari.OriginateParams{
			Endpoint:  endpointID,
			Context:   "public",
			Extension: "s",
			App:       "myari",
		}

		channel, err := client.Channels.Create(channelParams)
		if err != nil {
			return fmt.Errorf("error creating channel for endpoint %s: %s", endpointID, err)
		}
		log.Printf("Created channel %s for endpoint %s", channel.ID, endpointID)

		for channel.State != "Up" {
			time.Sleep(time.Millisecond * 100)
			channel, err = client.Channels.Get(channel.ID)
			if err != nil {
				return fmt.Errorf("error getting channel %s: %s", channel.ID, err)
			}
		}

		if err := bridge.AddChannel(channel.ID, ari.Participant); err != nil {
			return fmt.Errorf("error adding channel %s to bridge %s: %s", channel.ID, callID, err)
		}
	}

	return nil
}
