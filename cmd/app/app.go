package main

import (
  "fmt"
  "go_rest_api/pkg/mongo"
  "go_rest_api/pkg/server"
)

type App struct {
  s *server.Server
}

func(a *App) Initialize() {
  u := mongo.NewUserService()
  a.s = server.NewServer(u)
}

func(a *App) Run() {
  fmt.Println("Run")
  defer mongo.Close()
  a.s.Start()
}