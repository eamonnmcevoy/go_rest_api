package user

import (
  "go_rest_api/util/uuid"
  "golang.org/x/crypto/bcrypt"
)

type Credentials struct {
  Username string `json:"username"`
  Password string `json:"password"`
}

func(c* Credentials) salt() (string, string, error) { 
  uuid, uuidErr := uuid.New()
  if uuidErr != nil {
    return "", "", uuidErr
  }

  hash, err := bcrypt.GenerateFromPassword([]byte(c.Password+uuid), bcrypt.DefaultCost)
  
  hashString := string(hash[:])
  return hashString, uuid, err
}