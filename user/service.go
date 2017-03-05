package user

import (
  "gopkg.in/mgo.v2"
)

type IuserService interface {
  createUser(c Credentials) error
  getUserByUsername(username string) (error, User)
  login(c Credentials) (error, User)
}

type userService struct {
  provider *provider
}

func NewUserService(session *mgo.Session) *userService{
  s := userService{NewUserProvider(session)}
  return &s
}

func(u userService) createUser(c Credentials) error {
  user, err := NewUser(c)
  if err == nil {
    err = u.provider.InsertUser(user)
  }
  return err
}

func(u userService) getUserByUsername(username string) (error, User) {
  err, user := u.provider.GetUserByUsername(username)
  return err, user
}

func(u userService) login(c Credentials) (error, User) {
  err, user := u.provider.GetUserByUsername(c.Username)
  if err == nil {
    err = user.comparePassword(c.Password) 
  }
  return err, user  
}