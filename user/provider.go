package user

import (
  "gopkg.in/mgo.v2/bson"
  "gopkg.in/mgo.v2"
  "go_rest_api/mongo"
)

type provider struct {
  Collection *mgo.Collection
}

func NewUserProvider() *provider {
  p := provider{}
  session := mongo.Session()
  p.Collection = session.DB("test").C("user")
  return &p
}

func(p *provider) InsertUser(user User) error {
  return p.Collection.Insert(&user)
}

func (p *provider) GetUserByUsername(username string) (error, User) {
  result := User{}
  err := p.Collection.Find(bson.M{"username": username}).One(&result)
  return err, result
}