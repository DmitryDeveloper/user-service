package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/DmitryDeveloper/user-service/models"
	u "github.com/DmitryDeveloper/user-service/utils"
	"github.com/gorilla/mux"
)

var GetUser = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := models.GetUser(vars["user_id"])

	if user == nil {
		u.Respond(w, u.Message(false, "User not found"))
		return
	}

	response := u.Message(true, "Success")
	response["data"] = user

	u.Respond(w, response)
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	users := models.GetAll()

	response := u.Message(true, "Success")
	response["data"] = users

	u.Respond(w, response)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var updatedUserData models.Account
	err := json.NewDecoder(r.Body).Decode(&updatedUserData)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	user := models.GetUser(vars["user_id"])
	if user == nil {
		u.Respond(w, u.Message(false, "User not found"))
		return
	}

	user.FirstName = updatedUserData.FirstName
	user.LastName = updatedUserData.LastName
	user.BirthDate = updatedUserData.BirthDate
	user.Save()

	response := u.Message(true, "Success")
	response["data"] = user
	u.Respond(w, response)
}
