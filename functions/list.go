package functions

import (
	"fmt"

	"github.com/abourget/ari"
)

func ListOngoingCalls(client *ari.Client) error {
	bridges, err := client.Bridges.List()
	if err != nil {
		return err
	}

	fmt.Println("Ongoing Calls:")
	for _, bridge := range bridges {
		if bridge.BridgeType == "mixing" && len(bridge.Channels) > 0 {
			fmt.Printf("Call ID: %s, Participants: ", bridge.ID)
			for _, channelID := range bridge.Channels {
				channel, err := client.Channels.Get(channelID)
				if err != nil {
					return fmt.Errorf("error getting channel %s: %s", channelID, err)
				}
				fmt.Printf("%s ", channel.Name)
			}
			fmt.Println()
		}
	}
	return nil
}
