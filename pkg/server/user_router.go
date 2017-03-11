package server

import (
  "log"
  "errors"
  "net/http"
  "encoding/json"
  "github.com/gorilla/mux"
  "go_rest_api/pkg"
)

type userRouter struct {
  userService root.UserService
}

func NewUserRouter(u root.UserService, router *mux.Router) *mux.Router {
  userRouter := userRouter{u}

  router.HandleFunc("/", userRouter.createUserHandler).Methods("PUT")
  router.HandleFunc("/profile", validate(userRouter.profileHandler)).Methods("GET")
  router.HandleFunc("/{username}", userRouter.getUserHandler).Methods("GET")
  router.HandleFunc("/login", userRouter.loginHandler).Methods("POST")
  return router
}

func(s* userRouter) createUserHandler(w http.ResponseWriter, r *http.Request) {
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

func(s* userRouter) profileHandler(w http.ResponseWriter, r *http.Request) {
  claim, ok := r.Context().Value(contextKeyAuthtoken).(claims)
  if !ok {
    Error(w, http.StatusBadRequest, "no context")
    return
  }
  username := claim.Username

  err, user := s.userService.GetUserByUsername(username)
  if err != nil {
    Error(w, http.StatusNotFound, err.Error())
    return
  }

  Json(w, http.StatusOK, user)
}

func(s *userRouter) getUserHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  log.Println(vars)
  username := vars["username"]
  
  err, user := s.userService.GetUserByUsername(username)
  if err != nil {
    Error(w, http.StatusNotFound, err.Error())
    return
  }

  Json(w, http.StatusOK, user)
}

func(s* userRouter) loginHandler(w http.ResponseWriter, r *http.Request) {
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
  var u root.User
  if r.Body == nil {
    return errors.New("no request body"), u
  }
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&u)
  return err, u
}

func decodeCredentials(r *http.Request) (error,root.Credentials) {
  var c root.Credentials
  if r.Body == nil {
    return errors.New("no request body"), c
  }
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&c)
  return err, c
}