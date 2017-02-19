package user

import (
  "net/http"
  "encoding/json"
  "github.com/gorilla/mux"
  "gopkg.in/mgo.v2"
  "go_rest_api/util")

type userService struct {
  provider *provider
}

func NewUserRouter(session *mgo.Session, router *mux.Router) *mux.Router {
  s := userService{NewUserProvider(session)}

  router.HandleFunc("/create", s.createUserHandler)
  router.HandleFunc("/get", s.getUserHandler)
  router.HandleFunc("/authenticate", s.authenticateHandler)
  return router
}

func(u* userService) createUserHandler(w http.ResponseWriter, r *http.Request) {
  err, credentials := decodeCredentials(r)
  if err != nil {
    util.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
    return
  }

  var user User
  user, err = NewUser(credentials)
  if err != nil {
    util.RespondWithError(w, http.StatusInternalServerError, err.Error())
    return
  }

  err = u.provider.InsertUser(user)
  if err != nil {
    util.RespondWithError(w, http.StatusInternalServerError, err.Error())
    return
  }

  util.RespondWithJSON(w, http.StatusOK, err)
}

func(u* userService) getUserHandler(w http.ResponseWriter, r *http.Request) {
  err, findUser := decodeUser(r)
  if err != nil {
    util.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
    return
  }

  var user User
  err, user = u.provider.GetUser(findUser.Username)
  if err != nil {
    util.RespondWithError(w, http.StatusInternalServerError, err.Error())
    return
  }

  util.RespondWithJSON(w, http.StatusOK, user)
}

func(u *userService) authenticateHandler(w http.ResponseWriter, r *http.Request) {
  err, credentials := decodeCredentials(r)
  if err != nil {
    util.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
    return
  }

  var user User
  err, user = u.provider.GetUser(credentials.Username)
  if err != nil {
    util.RespondWithError(w, http.StatusInternalServerError, "No such user")
    return
  }

  if user.comparePassword(credentials.Password) {
    util.RespondWithJSON(w, http.StatusOK, map[string]string{"Success: ":credentials.Username})
  } else {
    util.RespondWithError(w, http.StatusInternalServerError, "Incorrect password")
  }
  
}

func decodeUser(r *http.Request) (error,User) {
  var u User
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&u)
  return err, u
}

func decodeCredentials(r *http.Request) (error,Credentials) {
  var c Credentials
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&c)
  return err, c
}  