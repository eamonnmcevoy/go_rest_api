package user

import (
  "fmt"
  "net/http"
  "encoding/json"
  "github.com/gorilla/mux"
  "gopkg.in/mgo.v2"
  "go_rest_api/util"
  "golang.org/x/crypto/bcrypt")

type userService struct {
  provider *provider
}

func NewUserRouter(session *mgo.Session, router *mux.Router) *mux.Router {
  fmt.Println("set up routes - user")

  s := userService{NewUserProvider(session)}

  router.HandleFunc("/create", s.createUserHandler)
  router.HandleFunc("/get", s.getUserHandler)
  router.HandleFunc("/authenticate", s.authenticateHandler)
  return router
}

func(u* userService) createUserHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Println("createUserHandler")
  
  var newUser NewUserTransport
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&newUser)
  if err != nil {
    util.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
    return
  }

  hash, salt, saltErr := salt(newUser.Password)
  if saltErr != nil {
    util.RespondWithError(w, http.StatusInternalServerError, saltErr.Error())
    return
  }

  user := NewUser(newUser.Username, hash, salt)
  err = u.provider.InsertUser(user)
  if err != nil {
    util.RespondWithError(w, http.StatusInternalServerError, err.Error())
    return
  }

  util.RespondWithJSON(w, http.StatusOK, err)
}

func(u* userService) getUserHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Println("getUserHandler")
  
  err, findUser := decodeUser(r)
  if err != nil {
    util.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
    return
  }

  err, user := u.provider.GetUser(findUser.Username)
  if err != nil {
    util.RespondWithError(w, http.StatusInternalServerError, err.Error())
    return
  }

  util.RespondWithJSON(w, http.StatusOK, user)
}

func(u *userService) authenticateHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Println("authenticateHandler")
  
  var newUser NewUserTransport
  decoder := json.NewDecoder(r.Body)
  decodeErr := decoder.Decode(&newUser)
  if decodeErr != nil {
    util.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
    return
  }

  findErr, dbUser := u.provider.GetUser(newUser.Username)
  if findErr != nil {
    util.RespondWithError(w, http.StatusInternalServerError, "No user with that name")
    return
  }

  if err := bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(newUser.Password+dbUser.Salt)); err != nil {
    util.RespondWithError(w, http.StatusInternalServerError, "Incorrect password")
    return
  }
  util.RespondWithJSON(w, http.StatusOK, map[string]string{"Success: ":newUser.Username})
}

func decodeUser(r *http.Request) (error,User) {
  var user User
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&user)
  return err, user
}

func salt(password string) (string, string, error) { 
  uuid, uuidErr := util.NewUUID()
  if uuidErr != nil {
    return "", "", uuidErr
  }

  hash, err := bcrypt.GenerateFromPassword([]byte(password+uuid), bcrypt.DefaultCost)
  
  hashString := string(hash[:])
  return hashString, uuid, err
}