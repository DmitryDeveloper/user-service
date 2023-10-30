package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/streadway/amqp"
)

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func SendToQueue(data []byte) {
	fmt.Println("Setup connection with RABBIT_MQ")
	conn, err := amqp.Dial(os.Getenv("RABBIT_MQ_CONNECTION_STRING"))
	if err != nil {
		fmt.Println("Cannot create connection to RabbitMQ")
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Should I keep channel in memory?
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("Cannot create channel in RabbitMQ")
		fmt.Println(err)
		return
	}
	defer ch.Close()

	fmt.Println("SENT DATA TO RABBITMQ")

	err = ch.Publish(
		os.Getenv("RABBIT_MQ_USERS_EXCHANGE_NAME"),
		os.Getenv("RABBIT_MQ_USERS_ROUTE_KEY"),
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		})
	if err != nil {
		fmt.Println("Cannot publsh message in RabbitMQ")
		fmt.Println(err)
		return
	}

	fmt.Println("DATA is published to RABBITMQ")
}
