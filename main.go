package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/abourget/ari"
	"github.com/peksinsara/ari/functions"
)

func main() {
	client := ari.NewClient("adminari", "1234", "localhost", 8088, "myari")

	applications, err := client.Applications.List()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Available ARI applications:")
	for _, app := range applications {
		fmt.Println("-", app.Name)
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
			fmt.Println("Enter two extensions separated by space (example 1000 1001):")
			input, _ := reader.ReadString('\n')
			extensions := strings.Fields(input)

			if len(extensions) != 2 {
				fmt.Println("Invalid input, please enter two extensions separated by space.")
				continue
			}

			endpoint1 := "SIP/" + extensions[0]
			endpoint2 := "SIP/" + extensions[1]

			err := functions.DialEndpoint(client, endpoint1, extensions[0], endpoint2, extensions[1], extensions[0], extensions[1], "outgoing")
			if err != nil {
				log.Fatal(err)
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
