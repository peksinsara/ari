package main

import (
	"fmt"
	"log"

	"github.com/abourget/ari"
)

func main() {
	client := ari.NewClient("adminari", "1234", "localhost", 8088, "myari")

	fmt.Println("Choose an option:")
	fmt.Println("1. Dial")
	fmt.Println("2. Join")
	fmt.Println("3. List")

	var choice int
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		fmt.Println("Functionality not implemented")
	case 2:
		fmt.Println("Functionality not implemented")
	case 3:
		calls, err := client.Channels.List()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Ongoing Calls:")
		for _, call := range calls {
			fmt.Printf("Call ID: %s, Participants: %s\n", call.ID, call.Name)
		}
	default:
		fmt.Println("Invalid choice")
	}

}
