package main

import (
	"fmt"
	"go_rest_api/pkg"
	"go_rest_api/pkg/config"
	"go_rest_api/pkg/mongo"
	"go_rest_api/pkg/server"
	"log"
)

type App struct {
	server  *server.Server
	session *mongo.Session
	config  *root.Config
}

func (a *App) Initialize() {
	a.config = config.GetConfig()
	var err error
	a.session, err = mongo.NewSession(a.config.Mongo)
	if err != nil {
		log.Fatalln("unable to connect to mongodb")
	}

	u := mongo.NewUserService(a.session.Copy(), a.config.Mongo)
	a.server = server.NewServer(u, a.config)
}

func (a *App) Run() {
	fmt.Println("Run")
	defer a.session.Close()
	a.server.Start()
}
