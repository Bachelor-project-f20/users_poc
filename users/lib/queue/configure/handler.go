package configure

import (
	"encoding/json"
	"log"
	"os"

	"github.com/dueruen/go-outbox"
)

var (
	DBTypeDefault                  = outbox.MySQL
	DBConnectionDefault            = "root:root@/root?charset=utf8&parseTime=True&loc=Local"
	MessageBrokerTypeDefault       = "nats"
	MessageBrokerDefaultConnection = "localhost:4222"
)

type ServiceConfig struct {
	DatabaseType            outbox.DbType
	DatabaseConnection      string
	MessageBrokerType       string
	MessageBrokerConnection string
}

func ExtractConfiguration(filename string) (ServiceConfig, error) {
	config := ServiceConfig{
		DBTypeDefault,
		DBConnectionDefault,
		MessageBrokerTypeDefault,
		MessageBrokerDefaultConnection,
	}
	config.extractFromFile(filename)
	config.extractFromEnv()

	log.Println("!!!ALERT!!! Use of default values")
	return config, nil
}

func (config *ServiceConfig) extractFromFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Config file not found")
	}
	json.NewDecoder(file).Decode(config)
}

func (config *ServiceConfig) extractFromEnv() {

	if v := os.Getenv("MYSQL_URL"); v != "" {
		config.DatabaseType = outbox.MySQL
		config.DatabaseConnection = v
	} else if v := os.Getenv("POSTGRES_URL"); v != "" {
		config.DatabaseType = outbox.Postgres
		config.DatabaseConnection = v
	}

	if v := os.Getenv("NATS_BROKER_URL"); v != "" {
		config.MessageBrokerType = "nats"
		config.MessageBrokerConnection = v
	}
}
