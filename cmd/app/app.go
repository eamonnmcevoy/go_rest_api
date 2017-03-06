package main

import (
  "fmt"
  "go_rest_api/pkg/mongo"
  "go_rest_api/pkg/http"
)

type App struct {
  server *http.Server
}

func(a *App) Initialize() {
  u := mongo.NewUserService()
  a.server = http.NewServer(u)
}

func(a *App) Run() {
  fmt.Println("Run")
  defer mongo.Close()
  a.server.Start()
}