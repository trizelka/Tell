package main

import (
	// Standard library packages
	"net/http"

	// Third party packages
	log "Tell/Sirupsen/logrus"
	"Tell/julienschmidt/httprouter"
	"Tell/trizelka/Messaging/controllers"
	"gopkg.in/mgo.v2"
)

func main() {
	// Instantiate a new router
	r := httprouter.New()

	// Get a UserController instance
	uc := controllers.NewUserController(getSession())

	// Get a user resource
	r.GET("/api/user/:from/:phone", uc.GetUser)

	// Send Messages a user resource
        r.POST("/api/messages/:from", uc.SendMessages)

	// Create or Update a new user
	r.POST("/api/user/:methode", uc.CreateUser)

	// Remove an existing user
	r.DELETE("/api/user/:phone", uc.RemoveUser)

	// Fire up the server
	log.Fatal(http.ListenAndServe("128.199.73.20:8000", r))
}

// getSession creates a new mongo session and panics if connection error occurs
func getSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://localhost")

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}

	// Deliver session
	return s
}
