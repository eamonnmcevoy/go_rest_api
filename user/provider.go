package user

import (
  "gopkg.in/mgo.v2/bson"
  "gopkg.in/mgo.v2"
)

type provider struct {
  Collection *mgo.Collection
}

func NewUserProvider(mongoSession *mgo.Session) *provider {
  p := provider{}
  p.Collection = mongoSession.DB("test").C("user")
  return &p
}

func(p *provider) InsertUser(user User) error {
  return p.Collection.Insert(&user)
}

func (p *provider) GetUser(username string) (error, User) {
  result := User{}
  err := p.Collection.Find(bson.M{"username": username}).One(&result)
  return err, result
}