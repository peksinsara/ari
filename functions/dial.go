package functions

import (
	"fmt"
	"log"
	"time"

	"github.com/abourget/ari"
)

var bridgeType = make(map[string]string)

func DialEndpoints(client *ari.Client, endpoints []string) error {

	var bridgeParams ari.CreateBridgeParams

	if len(endpoints) == 2 {
		bridgeParams = ari.CreateBridgeParams{
			Type: "mixing",
			Name: "myBridge",
		}
		bridgeType[bridgeParams.Name] = "call"
	} else {
		bridgeParams = ari.CreateBridgeParams{
			Type: "mixing",
			Name: "myBridge",
		}
		bridgeType[bridgeParams.Name] = "conference"
	}

	bridge, err := client.Bridges.Create(bridgeParams)
	if err != nil {
		return fmt.Errorf("error creating bridge: %s", err)
	}

	log.Printf("Created %s bridge %s", bridgeType[bridgeParams.Name], bridge.ID)

	bridgeType[bridge.ID] = bridgeType[bridgeParams.Name]

	var channels []*ari.Channel
	for _, endpoint := range endpoints {
		channel, err := CreateChannel(client, endpoint)
		if err != nil {
			return fmt.Errorf("error creating channel for endpoint %s: %s", endpoint, err)
		}

		channels = append(channels, channel)

		if err := bridge.AddChannel(channel.ID, ari.Participant); err != nil {
			return fmt.Errorf("error adding channel %s to bridge: %s", channel.ID, err)
		}
		log.Printf("Added channel %s to bridge %s", channel.ID, bridge.ID)

	}

	go func() {
		for {
			bridge, err = client.Bridges.Get(bridge.ID)
			if err != nil {
				log.Printf("error getting bridge %s: %s", bridge.ID, err)
				break
			}

			if bridge.BridgeClass == "destroyed" {
				log.Printf("%s %s ended", bridgeType[bridge.ID], bridge.ID)
				break
			}

			if len(bridge.Channels) == 0 {
				log.Printf("All participants have left %s %s", bridgeType[bridge.ID], bridge.ID)
				if err := bridge.Destroy(); err != nil {
					log.Printf("error destroying bridge %s: %s", bridge.ID, err)
				}
				break
			}

			if bridgeType[bridge.ID] == "call" && len(bridge.Channels) == 1 {
				log.Printf("Only one channel left in %s %s. Destroying bridge.", bridgeType[bridge.ID], bridge.ID)
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

func CreateChannel(client *ari.Client, endpoint string) (*ari.Channel, error) {
	channelParams := ari.OriginateParams{
		Endpoint:  endpoint,
		Context:   "public",
		Extension: "s",
		App:       "myari",
	}

	channel, err := client.Channels.Create(channelParams)
	if err != nil {
		return nil, fmt.Errorf("error creating channel for endpoint %s: %s", endpoint, err)
	}
	log.Printf("Created channel %s for endpoint %s", channel.ID, endpoint)

	for channel.State != "Up" {
		time.Sleep(time.Millisecond * 100)
		channel, err = client.Channels.Get(channel.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting channel %s: %s", channel.ID, err)
		}
	}

	return channel, nil
}
