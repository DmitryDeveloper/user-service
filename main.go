package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/DmitryDeveloper/user-service/app"
	"github.com/DmitryDeveloper/user-service/controllers"
	u "github.com/DmitryDeveloper/user-service/utils"
	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		u.Respond(w, u.Message(true, "OK"))
	}).Methods("GET")

	router.HandleFunc("/api/register", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/login", controllers.Authenticate).Methods("POST")

	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	//router.NotFoundHandler = app.NotFoundHandler
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Print(err)
	}
}
