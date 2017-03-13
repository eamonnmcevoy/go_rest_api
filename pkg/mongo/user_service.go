package mongo

import (
  "gopkg.in/mgo.v2/bson"
  "gopkg.in/mgo.v2"
  "go_rest_api/pkg"
)

type UserService struct {
 collection *mgo.Collection
}

func NewUserService(session *mgo.Session) *UserService {
  collection := session.DB("test").C("user")
  collection.EnsureIndex(userModelIndex())
  return &UserService {collection}
}

func(p *UserService) CreateUser(u *root.User) error {
  user, err := newUserModel(u)
  if err != nil {
    return err
  }

  return p.collection.Insert(&user)
}

func (p *UserService) GetUserByUsername(username string) (error, root.User) {
  model := userModel{}
  err := p.collection.Find(bson.M{"username": username}).One(&model)
  return err, root.User{
    Id: model.Id.Hex(),
    Username: model.Username,
    Password: "-" }
}

func (p *UserService) Login(c root.Credentials) (error, root.User) {
  model := userModel{}
  err := p.collection.Find(bson.M{"username": c.Username}).One(&model)

  err = model.comparePassword(c.Password)
  if(err != nil) {
    return err, root.User{}
  }

  return err, root.User{
    Id: model.Id.Hex(),
    Username: model.Username,
    Password: "-" }
}

