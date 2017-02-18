package user

import (
  "fmt"
  "net/http"
  "encoding/json"
  "github.com/gorilla/mux"
  "gopkg.in/mgo.v2"
  "go_rest_api/util")

type userService struct {
  provider *provider
}

func NewUserRouter(session *mgo.Session, router *mux.Router) *mux.Router {
  fmt.Println("set up routes - user")

  s := userService{NewUserProvider(session)}

  router.HandleFunc("/create", s.createUserHandler)
  router.HandleFunc("/get", s.getUserHandler)
  return router
}

func(u* userService) createUserHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Println("createUserHandler")
  var user User
  decoder := json.NewDecoder(r.Body)
  if err := decoder.Decode(&user); err != nil {
    util.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
    return
  }
  err := u.provider.InsertUser(user)
  if err != nil {
    util.RespondWithError(w, http.StatusInternalServerError, err.Error())
    return
  }

  util.RespondWithJSON(w, http.StatusOK, err)
}

func(u* userService) getUserHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Println("getUserHandler")
  var user User
  decoder := json.NewDecoder(r.Body)
  if err := decoder.Decode(&user); err != nil {
    util.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
    return
  }
  err, user := u.provider.GetUser(user.Username)
  if err != nil {
    util.RespondWithError(w, http.StatusInternalServerError, err.Error())
    return
  }

  util.RespondWithJSON(w, http.StatusOK, user)
}