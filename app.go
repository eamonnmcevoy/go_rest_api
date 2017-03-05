package main

import (
  "os"
  "fmt"
  "log"
  "net/http"
  "github.com/gorilla/mux"
  "github.com/gorilla/handlers"
  "go_rest_api/user"
  "go_rest_api/mongo"
)

type App struct {
  Router *mux.Router
}

func(a *App) Initialize() {
  a.Router = mux.NewRouter()
  a.SetupRoutes()
}

func(a *App) SetupRoutes() {
  a.Router.Handle("/user", user.NewUserRouter(a.getSubrouter("/user")))
}

func(a *App) getSubrouter(path string) *mux.Router {
  return a.Router.PathPrefix(path).Subrouter()
}

func(a *App) Run() {
  fmt.Println("Run")
  defer mongo.Close()

  log.Println("Listening on port 1337...")
  if err := http.ListenAndServe(":1337", handlers.LoggingHandler(os.Stdout, a.Router)); err != nil {
      log.Fatal("http.ListenAndServe: ", err)
  }
}