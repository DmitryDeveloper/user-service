package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"user-service/models"
	u "user-service/utils"
)

type UserRegisteredEvent struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	EventType string `json:"event_type"`
}

var CreateUser = func(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := user.Create()

	if resp["status"] == true {
		fmt.Println("User created")

		userRegisteredEventData := UserRegisteredEvent{
			ID:        int(user.ID),
			Email:     user.Email,
			EventType: "user_registered",
		}
		jsonData, _ := json.Marshal(userRegisteredEventData)

		u.SendToQueue(jsonData)
	}

	u.Respond(w, resp)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := models.Login(user.Email, user.Password)
	u.Respond(w, resp)
}
