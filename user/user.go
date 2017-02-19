package user

import (
  "gopkg.in/mgo.v2/bson"
  "golang.org/x/crypto/bcrypt"
)

type User struct {
  Id       bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
  Username string        `json:"username"`
  PasswordHash string    `json:"passwordHash"`
  Salt         string    `json:"salt"`
}

func NewUser(c Credentials) (User, error) {
  user := User{}
  user.Username = c.Username

  hash, salt, err := c.salt()
  if err != nil {
    return user, err
  }

  user.PasswordHash = hash
  user.Salt = salt

  return user, err
}

func(u* User) comparePassword(password string) bool { 
  incoming := []byte(password+u.Salt)
  existing := []byte(u.PasswordHash)
  err := bcrypt.CompareHashAndPassword(existing, incoming)
  return err == nil
}