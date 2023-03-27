package functions

import (
	"fmt"

	"github.com/abourget/ari"
)

func ListOngoingCalls(client *ari.Client) error {
	calls, err := client.Channels.List()
	if err != nil {
		return err
	}

	fmt.Println("Ongoing Calls:")
	for _, call := range calls {
		fmt.Printf("Call ID: %s, Participants: %s\n", call.ID, call.Name)
	}
	return nil
}
