package http

import (
  "os"
  "log"
  "net/http"
  "github.com/gorilla/mux"
  "github.com/gorilla/handlers"
  "go_rest_api/pkg"
)

type Server struct {
  router *mux.Router
}

func NewServer(u root.UserService) *Server {
  s := Server { router: mux.NewRouter() }
  
  s.router.Handle("/user", NewUserRouter(u, s.getSubrouter("/user")))
  
  return &s
}

func(s *Server) Start() {
  log.Println("Listening on port 1337...")
  if err := http.ListenAndServe(":1337", handlers.LoggingHandler(os.Stdout, s.router)); err != nil {
      log.Fatal("http.ListenAndServe: ", err)
  }
}

func(s *Server) getSubrouter(path string) *mux.Router {
  return s.router.PathPrefix(path).Subrouter()
}