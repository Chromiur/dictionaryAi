package main

import (
	"dictionaryAi/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

// Create a new ServeMux using Gorilla
var rMux = mux.NewRouter()

// PORT is where the web server listens to
var PORT = ":1234"

func main() {
	arguments := os.Args
	if len(arguments) >= 2 {
		PORT = ":" + arguments[1]
	}

	s := http.Server{
		Addr:         PORT,
		Handler:      rMux,
		ErrorLog:     nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	rMux.NotFoundHandler = http.HandlerFunc(handlers.DefaultHandler)

	notAllowed := handlers.NotAllowedHandler{}
	rMux.MethodNotAllowedHandler = notAllowed

	rMux.HandleFunc("/", handlers.MainPageHandler)

	// Define Handler Functions
	// Register GET
	getMux := rMux.Methods(http.MethodGet).Subrouter()
	//Get list of the words
	getMux.HandleFunc("/list", handlers.WordsListHandler)
	//Get list of Users
	//getMux.HandleFunc("/getid/{username}", GetIDHandler)
	//getMux.HandleFunc("/logged", LoggedUsersHandler)
	//getMux.HandleFunc("/username/{id:[0-9]+}", GetUserDataHandler)

	// Register PUT
	// Update User
	putMux := rMux.Methods(http.MethodPut).Subrouter()
	putMux.HandleFunc("/update", handlers.UpdateUserHandler)

	// Register POST
	postMux := rMux.Methods(http.MethodPost).Subrouter()
	// Add Words
	postMux.HandleFunc("/add", handlers.AddHandler)
	// Generate Sentence from words
	postMux.HandleFunc("/generateSentence", handlers.GenerateSentenceHandler)
	// Add User + Login + Logout
	postMux.HandleFunc("/addUser", handlers.AddUserHandler)
	//postMux.HandleFunc("/login", LoginHandler)
	//postMux.HandleFunc("/logout", LogoutHandler)

	// Register DELETE
	// Delete Words
	deleteMux := rMux.Methods(http.MethodDelete).Subrouter()
	deleteMux.HandleFunc("/deleteWords", handlers.DeleteHandler)
	// Delete User
	deleteMux.HandleFunc("/deleteUser", handlers.DeleteUserHandler)

	go func() {
		log.Println("Listening to", PORT)
		err := s.ListenAndServe()
		if err != nil {
			log.Printf("Error starting server: %s\n", err)
			return
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	sig := <-sigs
	log.Println("Quitting after signal:", sig)
	time.Sleep(5 * time.Second)
	s.Shutdown(nil)
}
