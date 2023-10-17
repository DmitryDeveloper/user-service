package main

import (
	"fmt"
	"net/http"
	"os"

	"user-service/controllers"
	u "user-service/utils"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		u.Respond(w, u.Message(true, "OK"))
	}).Methods("GET")

	apiRouter.HandleFunc("/register", controllers.CreateUser).Methods("POST")
	apiRouter.HandleFunc("/login", controllers.Authenticate).Methods("POST")

	apiRouter.HandleFunc("/users/{user_id}", controllers.GetUser).Methods("GET")
	apiRouter.HandleFunc("/users", controllers.GetAll).Methods("GET")
	apiRouter.HandleFunc("/users/{user_id}", controllers.UpdateUser).Methods("PUT")

	// router.Use(app.JwtAuthentication) //attach JWT auth middleware

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
