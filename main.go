package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/peksinsara/ari/functions"
	"github.com/peksinsara/ari/server"
)

func main() {
	client, err := server.ConnectToARI("adminari", "1234", "localhost", 8088, "myari")
	if err != nil {
		log.Fatal(err)
	}

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
			fmt.Println("Enter extensions separated by space (example 1000 1001 or 1000 1002 1003):")
			input, _ := reader.ReadString('\n')
			endpoints := strings.Fields(input)

			if len(endpoints) < 2 {
				fmt.Println("Invalid input, please enter at least two endpoints separated by a space.")
				return
			}

			if len(endpoints) == 2 {
				err = functions.DialEndpoint(client, "SIP/"+strings.TrimSpace(endpoints[0]), "SIP/"+strings.TrimSpace(endpoints[1]))
				if err != nil {
					log.Fatal(err)
				}
			} else {
				var channels []string
				for _, endpoint := range endpoints {
					channels = append(channels, "SIP/"+strings.TrimSpace(endpoint))
				}

				err = functions.DialConference(client, channels)
				if err != nil {
					log.Fatal(err)
				}
			}

		case 2:
			fmt.Println("Functionality not implemented")
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
