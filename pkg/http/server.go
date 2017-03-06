package http

import (
  "os"
  "log"
  "net/http"
  "encoding/json"
  "github.com/gorilla/mux"
  "github.com/gorilla/handlers"
  "go_rest_api/pkg"
)

type Server struct {
  userService root.UserService
  router *mux.Router
}

func NewServer(u root.UserService) *Server {
  s := Server { 
    userService: u,
    router: mux.NewRouter() }
    
  s.router.HandleFunc("/user", s.createUserHandler).Methods("PUT")
  s.router.HandleFunc("/user/profile", validate(s.profileHandler)).Methods("GET")
  s.router.HandleFunc("/user/{username}", s.getUserHandler).Methods("GET")
  s.router.HandleFunc("/user/login", s.loginHandler).Methods("POST")
  
  return &s
}

func(s *Server) Start() {
  log.Println("Listening on port 1337...")
  if err := http.ListenAndServe(":1337", handlers.LoggingHandler(os.Stdout, s.router)); err != nil {
      log.Fatal("http.ListenAndServe: ", err)
  }
}

func(s* Server) createUserHandler(w http.ResponseWriter, r *http.Request) {
  log.Println("createUserHandler")
  err, user := decodeUser(r)
  if err != nil {
    Error(w, http.StatusBadRequest, "Invalid request payload")
    return
  }

  err = s.userService.CreateUser(&user)
  if err != nil {
    Error(w, http.StatusInternalServerError, err.Error())
    return
  }

  Json(w, http.StatusOK, err)
}

func(s* Server) profileHandler(w http.ResponseWriter, r *http.Request) {
  _claims := r.Context().Value(contextKeyAuthtoken).(claims)
  username := _claims.Username

  err, user := s.userService.GetUserByUsername(username)
  if err != nil {
    Error(w, http.StatusInternalServerError, err.Error())
    return
  }

  Json(w, http.StatusOK, user)
}

func(s *Server) getUserHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  username := vars["username"]
  
  err, user := s.userService.GetUserByUsername(username)
  if err != nil {
    Error(w, http.StatusInternalServerError, err.Error())
    return
  }

  Json(w, http.StatusOK, user)
}

func(s* Server) loginHandler(w http.ResponseWriter, r *http.Request) {
  log.Println("loginHandler")
  err, credentials := decodeCredentials(r)
  if err != nil {
    Error(w, http.StatusBadRequest, "Invalid request payload")
    return
  }

  var user root.User
  err, user = s.userService.Login(credentials)
  if err == nil {
    cookie := newAuthCookie(user)
    JsonWithCookie(w, http.StatusOK, user, cookie)
  } else {
    Error(w, http.StatusInternalServerError, "Incorrect password")
  }
}

func decodeUser(r *http.Request) (error,root.User) {
  var c root.User
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&c)
  return err, c
}

func decodeCredentials(r *http.Request) (error,root.Credentials) {
  var c root.Credentials
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&c)
  return err, c
}