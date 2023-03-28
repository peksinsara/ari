package functions

import (
	"github.com/abourget/ari"
)

func DialEndpoint(client *ari.Client, endpoint1 string, ext1 string, endpoint2 string, ext2 string, direction string) error {
	params1 := ari.OriginateParams{
		Endpoint:  endpoint1,
		Context:   "public",
		Extension: ext1,
		App:       "myari",
	}

	params2 := ari.OriginateParams{
		Endpoint:  endpoint2,
		Context:   "public",
		Extension: ext2,
		App:       "myari",
	}

	channel1, err := client.Channels.Create(params1)
	if err != nil {
		return err
	}

	channel2, err := client.Channels.Create(params2)
	if err != nil {
		return err
	}

	bridgeParams := ari.CreateBridgeParams{
		Type: "mixing",
		Name: "myBridge",
	}
	bridge, err := client.Bridges.Create(bridgeParams)

	if err := bridge.AddChannel(channel1.ID, ari.Participant); err != nil {
		return err
	}

	if err := bridge.AddChannel(channel2.ID, ari.Participant); err != nil {
		return err
	}

	return nil
}
