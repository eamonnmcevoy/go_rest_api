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

  router.HandleFunc("/", u.createUserHandler).Methods("PUT")
  router.HandleFunc("/", validate(u.profileHandler)).Methods("GET")
  router.HandleFunc("/{username}", u.getUserHandler).Methods("GET")
  router.HandleFunc("/login", u.loginHandler).Methods("POST")
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

func(u* userRouter) profileHandler(w http.ResponseWriter, r *http.Request) {
  claims := r.Context().Value(contextKeyAuthtoken).(Claims)
  username := claims.Username

  err, user := u.service.getUserByUsername(username)
  if err != nil {
    response.Error(w, http.StatusInternalServerError, err.Error())
    return
  }

  response.Json(w, http.StatusOK, user)
}

func(u* userRouter) getUserHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  username := vars["username"]
  
  err, user := u.service.getUserByUsername(username)
  if err != nil {
    response.Error(w, http.StatusInternalServerError, err.Error())
    return
  }

  response.Json(w, http.StatusOK, user)
}


func(u *userRouter) loginHandler(w http.ResponseWriter, r *http.Request) {
  err, credentials := decodeCredentials(r)
  if err != nil {
    response.Error(w, http.StatusBadRequest, "Invalid request payload")
    return
  }

  var user User
  err, user = u.service.login(credentials)
  if err == nil {
    cookie := newAuthCookie(user)
    response.JsonWithCookie(w, http.StatusOK, user, cookie)
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