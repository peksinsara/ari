package functions

import (
	"fmt"

	"github.com/abourget/ari"
)

// ListOngoingCalls prints out the list of ongoing calls on the ARI server
func ListOngoingCalls(client *ari.Client) error {
	bridges, err := client.Bridges.List()
	if err != nil {
		return fmt.Errorf("failed to get the list of bridges: %s", err)
	}

	fmt.Println("Ongoing Calls:")

	for _, bridge := range bridges {
		if bridge.BridgeType == "mixing" && len(bridge.Channels) > 0 {
			fmt.Printf("Call ID: %s, Participants: ", bridge.ID)

			for _, channelID := range bridge.Channels {
				channel, err := client.Channels.Get(channelID)
				if err != nil {
					return fmt.Errorf("failed to get channel %s: %s", channelID, err)
				}
				fmt.Printf("%s ", channel.Name)
			}
			fmt.Println()
		}
	}
	return nil
}
