package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"handler"
	"images"
	"log"
	"middlewaree"
	"mysql"
	"net/http"
	"session"
	"user"
)

func init() {
	//Assign a user store
	store, err := user.NewFileUserStore("./data/users.json")
	fmt.Println(store)
	if err != nil {
		panic(fmt.Errorf("Error creating users store: %s", err))
	}
	user.GlobalUserStore = store

	//Assign a session store
	sessionStore, err := session.NewFileSessionStore("./data/sessions.json")
	if err != nil {
		panic(fmt.Errorf("Error creating session store: %s", err))
	}
	session.GlobalSessionStore = sessionStore

	//Assign a sql database
	db, err := mysql.NewMySQLDB("root:1234@tcp(127.0.0.1:3306)/gophr")
	if err != nil {
		panic(err)
	}
	mysql.GlobalMYSQLDB = db

	//Assign a image store
	images.GlobalImageStore = images.NewDBImageStore()
}

func main() {
	router := NewRouter()

	router.Handle("GET", "/", handler.HandleHome)
	router.Handle("GET", "/register", handler.HandlerUserNew)
	router.Handle("POST", "/register", handler.HandleUserCreate)
	router.Handle("GET", "/login", handler.HandleSessionNew)
	router.Handle("POST", "/login", handler.HandleSessionCreate)
	router.Handle("GET", "/image/:imageID", handler.HandleImageShow)
	router.Handle("GET", "/user/:userID", handler.HandleUserShow)

	router.ServeFiles(
		"/assets/*filepath",
		http.Dir("assets/"),
	)
	router.ServeFiles(
		"/im/*filepath",
		http.Dir("data/images/"),
	)

	secureRouter := NewRouter()
	secureRouter.Handle("GET", "/sign-out", handler.HandleSessionDestroy)
	secureRouter.Handle("GET", "/account", handler.HandleUserEdit)
	secureRouter.Handle("POST", "/account", handler.HandleUserUpdate)
	secureRouter.Handle("GET", "/images/new", handler.HandleImageNew)
	secureRouter.Handle("POST", "/images/new", handler.HandleImageCreate)

	middleware := middlewaree.Middleware{}
	middleware.Add(router)
	middleware.Add(http.HandlerFunc(session.RequireLogin))
	middleware.Add(secureRouter)

	fmt.Println("Server running......")
	log.Fatal(http.ListenAndServe(":3000", middleware))
}

// Creates a new router
func NewRouter() *httprouter.Router {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	return router
}
