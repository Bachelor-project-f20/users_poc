package test

import (
	"fmt"
	"testing"

	stan "github.com/Bachelor-project-f20/eventToGo/nats"
	"github.com/Bachelor-project-f20/go-outbox"
	"github.com/grammeaway/users_poc/users/lib/configure"
	"github.com/nats-io/go-nats"
)

func TestOutboxSetup(t *testing.T) {
	configFile := "testConfig"
	config, err := configure.ExtractConfiguration(configFile)

	if err != nil {
		fmt.Printf("Error extracting config file: %v \n", err)
		fmt.Println("Using default configuration")

	}

	encodedConn, err := setupNatsConn()

	if err != nil {
		fmt.Printf("Error connecting to Nats: %v \n", err)
		t.Error(err)
	}

	exchange := "test"
	queueType := "queue"

	eventEmitter, err := stan.NewNatsEventEmitter(encodedConn, exchange, queueType)

	_, obErr := outbox.NewOutbox(config.DatabaseType, config.DatabaseConnection, eventEmitter)

	if obErr != nil {
		fmt.Printf("Error creating Outbox: %v \n", err)
		t.Error(err)
	}
}

func setupNatsConn() (*nats.EncodedConn, error) {

	natsConn, err := nats.Connect("localhost:4222")

	if err != nil {
		fmt.Println("Connection to Nats failed")
		return nil, err
	}

	encodedConn, err := nats.NewEncodedConn(natsConn, "json")

	if err != nil {
		fmt.Println("Creation of encoded connection failed")
		return nil, err
	}

	return encodedConn, nil

}
