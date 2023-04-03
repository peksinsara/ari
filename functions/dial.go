package functions

import (
	"fmt"
	"log"
	"time"

	"github.com/abourget/ari"
)

func DialEndpoints(client *ari.Client, endpoints []string) error {
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

	if len(channels) == 2 {
		log.Printf("Started call %s, Participants: %s, %s", bridge.ID, endpoints[0], endpoints[1])
	} else {
		log.Printf("Started conference %s with %d participant(s)", bridge.ID, len(channels))
	}

	go func() {
		for {
			bridge, err = client.Bridges.Get(bridge.ID)
			if err != nil {
				log.Printf("error getting bridge %s: %s", bridge.ID, err)
				break
			}

			if bridge.BridgeClass == "destroyed" {
				log.Printf("Conference %s ended", bridge.ID)
				break
			}

			if len(bridge.Channels) == 0 {
				log.Printf("All participants have left %s, ending", bridge.ID)
				if err := bridge.Destroy(); err != nil {
					log.Printf("error destroying bridge %s: %s", bridge.ID, err)
				}
				break
			}

			time.Sleep(time.Second)
		}
	}()

	return nil
}
