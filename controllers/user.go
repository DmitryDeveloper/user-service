package controllers

import (
	"encoding/json"
	"net/http"

	"user-service/models"
	u "user-service/utils"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
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

	var updatedUserData models.User
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

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var newPassData map[string]string

	err := json.NewDecoder(r.Body).Decode(&newPassData)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	if newPassData["password"] != newPassData["confirmed"] {
		u.Respond(w, u.Message(false, "Password and confirmation doesn't match"))
		return
	}

	user := models.GetUser(vars["user_id"])
	if user == nil {
		u.Respond(w, u.Message(false, "User not found"))
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassData["password"]), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	user.Save()

	response := u.Message(true, "Success")
	user.Password = ""
	response["data"] = user
	u.Respond(w, response)
}
