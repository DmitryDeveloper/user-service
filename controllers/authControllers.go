package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/DmitryDeveloper/user-service/models"
	u "github.com/DmitryDeveloper/user-service/utils"
	"github.com/streadway/amqp"
)

type UserRegisteredEvent struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := account.Create()

	if resp["status"] == true {
		fmt.Println("Account created")

		userRegisteredEventData := UserRegisteredEvent{
			ID:    int(account.ID),
			Email: account.Email,
		}
		jsonData, _ := json.Marshal(userRegisteredEventData)

		fmt.Println("User data prepared")

		sendToQueue(jsonData)

		fmt.Println("User data has sent")
	}

	u.Respond(w, resp)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := models.Login(account.Email, account.Password)
	u.Respond(w, resp)
}

func sendToQueue(data []byte) {
	conn, err := amqp.Dial(os.Getenv("RABBIT_MQ_CONNECTION_STRING"))
	if err != nil {
		fmt.Println("Cannot create connection to RabbitMQ")
		fmt.Println(err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("Cannot create channel in RabbitMQ")
		fmt.Println(err)
		return
	}
	defer ch.Close()

	fmt.Println("SENT DATA TO RABBITMQ")

	err = ch.Publish(
		os.Getenv("RABBIT_MQ_USERS_QUEUE_NAME"),
		os.Getenv("RABBIT_MQ_USERS_QUEUE_KEY"),
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,               // pass data
		})
	if err != nil {
		fmt.Println("Cannot publsh message in RabbitMQ")
		fmt.Println(err)
		return
	}
}
