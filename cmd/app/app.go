package main

import (
  "fmt"
  "log"
  "go_rest_api/pkg/mongo"
  "go_rest_api/pkg/server"
)

type App struct {
  server *server.Server
  session *mongo.Session
}

func(a *App) Initialize() {
  err := a.session.Open()
  if(err != nil) {
    log.Fatalln("unable to connect to mongodb")
  }
  u := mongo.NewUserService(a.session.Copy())
  a.server = server.NewServer(u)
}

func(a *App) Run() {
  fmt.Println("Run")
  defer a.session.Close()
  a.server.Start()
}