package http

import (
  "log"
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

func(s* userRouter) profileHandler(w http.ResponseWriter, r *http.Request) {
  _claims := r.Context().Value(contextKeyAuthtoken).(claims)
  username := _claims.Username

  err, user := s.userService.GetUserByUsername(username)
  if err != nil {
    Error(w, http.StatusInternalServerError, err.Error())
    return
  }

  Json(w, http.StatusOK, user)
}

func(s *userRouter) getUserHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  username := vars["username"]
  
  err, user := s.userService.GetUserByUsername(username)
  if err != nil {
    Error(w, http.StatusInternalServerError, err.Error())
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