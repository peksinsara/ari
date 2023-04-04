package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/abourget/ari"
	"github.com/peksinsara/ari/functions"
	"github.com/peksinsara/ari/server"
)

func main() {
	client, err := server.ConnectToARI("adminari", "1234", "localhost", 8088, "myari")
	if err != nil {
		log.Fatal(err)
	}

	eventCh := client.LaunchListener()
	go func() {
		for {
			event := <-eventCh
			switch event.GetType() {
			case "ChannelLeftBridge":
				if bridgeEvent, ok := event.(*ari.ChannelLeftBridge); ok {
					go functions.DestroyBridge(client, bridgeEvent.Bridge)
				} else {
					fmt.Println("error in event")
				}
			default:
			}
		}
	}()

	// eventCh := client.LaunchListener()
	// go func() {
	// 	for {
	// 		event := <-eventCh
	// 		switch event := event.(type) {
	// 		//Gettype, returns string
	// 		case *ari.ChannelLeftBridge:
	// 			//log.Printf("Received ChannelLeftBridge event: %+v\n", event)
	// 			go functions.DestroyBridge(client, event.Bridge, functions.BridgeType[event.Bridge.ID])
	// 		default:
	// 		}
	// 	}
	// }()

	for {
		fmt.Println("Choose an option:")
		fmt.Println("1. Dial")
		fmt.Println("2. Join")
		fmt.Println("3. List")
		fmt.Println("4. Exit")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:

			reader := bufio.NewReader(os.Stdin)
			fmt.Println("Enter extensions separated by space (example: dial 1000 1001 or dial 1000 1002 1003):")
			input, _ := reader.ReadString('\n')
			args := strings.Fields(input)

			if len(args) < 2 {
				fmt.Println("Invalid input, please enter 'dial' and at least two endpoints separated by a space.")
				return
			}

			if args[0] == "dial" {
				endpoints := args[1:]

				if len(endpoints) < 2 {
					fmt.Println("Invalid input, please enter at least two endpoints separated by a space.")
					return
				}
				if len(endpoints) >= 2 {
					var channels []string
					for _, endpoint := range endpoints {
						channels = append(channels, "SIP/"+strings.TrimSpace(endpoint))
					}

					err = functions.DialEndpoints(client, channels)
					if err != nil {
						log.Fatal(err)
					}
				}
			} else {
				fmt.Println("Invalid command, please enter 'dial' followed by extensions separated by space.")
				return
			}

		case 2:
			if err := functions.ListOngoingCalls(client); err != nil {
				log.Fatal(err)
			}

			reader := bufio.NewReader(os.Stdin)

			fmt.Println("Enter command in the format of (example: join CALLID 1001 ): ")
			command, _ := reader.ReadString('\n')
			command = strings.TrimSpace(command)

			commandArgs := strings.Split(command, " ")
			if len(commandArgs) < 3 {
				fmt.Println("Invalid command. Usage: join CALLID endpointID1 endpointID2 ...")
				continue
			}

			callID := commandArgs[1]
			endpointIDs := commandArgs[2:]

			for i, endpointID := range endpointIDs {
				endpointIDs[i] = "SIP/" + endpointID
			}

			if err := functions.JoinCall(client, callID, endpointIDs); err != nil {
				log.Printf("Error joining call: %s", err)
			}

		case 3:
			if err := functions.ListOngoingCalls(client); err != nil {
				log.Fatal(err)
			}
		case 4:
			return
		default:
			fmt.Println("Invalid choice")

		}

	}
}
