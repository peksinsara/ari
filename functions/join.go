package functions

import (
	"fmt"

	"github.com/abourget/ari"
)

// JoinCall joins a call by adding endpoints to an existing bridge
func JoinCall(client *ari.Client, callID string, endpointIDs []string) error {
	bridge, err := client.Bridges.Get(callID)
	if err != nil {
		return fmt.Errorf("failed to get bridge %s: %s", callID, err)
	}

	for _, endpointID := range endpointIDs {
		channel, err := CreateChannel(client, endpointID)
		if err != nil {
			return fmt.Errorf("failed to create channel for endpoint %s: %s", endpointID, err)
		}

		if err := bridge.AddChannel(channel.ID, ari.Participant); err != nil {
			return fmt.Errorf("failed to add channel %s to bridge %s: %s", channel.ID, callID, err)
		}

	}
	if len(bridge.Channels) >= 2 {
		fmt.Println("Chaning 'call' to 'conference'. ")
		bridgeType[bridge.ID] = "conference"

	}

	return nil
}
