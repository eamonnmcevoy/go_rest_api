package mongo

import (
  "gopkg.in/mgo.v2/bson"
  "gopkg.in/mgo.v2"
  "golang.org/x/crypto/bcrypt"
  "github.com/google/uuid"
)

type userModel struct {
  Id           bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
  Username     string        `json:"username"`
  PasswordHash string        `json:"-"`
  Salt         string        `json:"-"`
}

func userModelIndex() mgo.Index {
  return mgo.Index{
    Key:        []string{"username"},
    Unique:     true,
    DropDups:   true,
    Background: true,
    Sparse:     true,
  }
}

func(u *userModel) comparePassword(password string) error { 
  incoming := []byte(password+u.Salt)
  existing := []byte(u.PasswordHash)
  err := bcrypt.CompareHashAndPassword(existing, incoming)
  return err
}

func(u *userModel) addSaltedPassword(password string) error { 
  salt := uuid.New().String()
  passwordBytes := []byte(password + salt)
  hash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
  if err != nil {
    return err
  }

  u.PasswordHash = string(hash[:])
  u.Salt = salt

  return nil
}