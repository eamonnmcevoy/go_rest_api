package user

import (
  "net/http"
  "encoding/json"
  "github.com/gorilla/mux"
  "gopkg.in/mgo.v2"
  "go_rest_api/util/response"
)

type userRouter struct {
  service IuserService
}

func NewUserRouter(session *mgo.Session, router *mux.Router) *mux.Router {
  u := userRouter{NewUserService(session)}

  router.HandleFunc("/create", u.createUserHandler)
  router.HandleFunc("/get", u.getUserHandler)
  router.HandleFunc("/authenticate", u.authenticateHandler)
  return router
}

func(u* userRouter) createUserHandler(w http.ResponseWriter, r *http.Request) {
  err, credentials := decodeCredentials(r)
  if err != nil {
    response.Error(w, http.StatusBadRequest, "Invalid request payload")
    return
  }

  err = u.service.createUser(credentials)
  if err != nil {
    response.Error(w, http.StatusInternalServerError, err.Error())
    return
  }

  response.Json(w, http.StatusOK, err)
}

func(u* userRouter) getUserHandler(w http.ResponseWriter, r *http.Request) {
  err, findUser := decodeUser(r)
  if err != nil {
    response.Error(w, http.StatusBadRequest, "Invalid request payload")
    return
  }

  var user User
  err, user = u.service.getUserByUsername(findUser.Username)
  if err != nil {
    response.Error(w, http.StatusInternalServerError, err.Error())
    return
  }

  response.Json(w, http.StatusOK, user)
}

func(u *userRouter) authenticateHandler(w http.ResponseWriter, r *http.Request) {
  err, credentials := decodeCredentials(r)
  if err != nil {
    response.Error(w, http.StatusBadRequest, "Invalid request payload")
    return
  }

  var user User
  err, user = u.service.authenticate(credentials)
  if err == nil {
    response.Json(w, http.StatusOK, user)
  } else {
    response.Error(w, http.StatusInternalServerError, "Incorrect password")
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