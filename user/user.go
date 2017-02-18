package user

import (
  "gopkg.in/mgo.v2/bson"
)

type NewUserTransport struct {
  Username string `json:"username"`
  Password string `json:"password"`
}

type User struct {
  Id       bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
  Username string        `json:"username"`
  PasswordHash string    `json:"passwordHash"`
  Salt         string    `json:"salt"`
}

func NewUser(username string, hash string, salt string) User {
  user := User{}
  user.Username = username
  user.PasswordHash = hash
  user.Salt = salt
  return user
}