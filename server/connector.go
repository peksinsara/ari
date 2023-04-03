package server

import (
	"log"

	"github.com/abourget/ari"
)

// ConnectToARI creates a connection to the ARI server and returns an ari.Client instance
func ConnectToARI(username, password, hostname string, port int, appName string) (*ari.Client, error) { // Create a new ARI client instance
	client := ari.NewClient(username, password, hostname, port, appName)

	// List the available ARI applications
	applications, err := client.Applications.List()
	if err != nil {
		return nil, err
	}

	log.Println("Available ARI applications:")
	for _, app := range applications {
		log.Printf("- %s\n", app.Name)
	}

	return client, nil
}
