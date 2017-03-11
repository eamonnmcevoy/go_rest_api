package mongo

import (
  "io"
  "fmt"
  "gopkg.in/mgo.v2/bson"
  "gopkg.in/mgo.v2"
  "go_rest_api/pkg"
  "crypto/rand"
  "golang.org/x/crypto/bcrypt"
)

type UserService struct {
 collection *mgo.Collection
}

type userModel struct {
  Id           bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
  Username     string        `json:"username"`
  PasswordHash string        `json:"-"`
  Salt         string        `json:"-"`
}

func(u *userModel) comparePassword(password string) error { 
  incoming := []byte(password+u.Salt)
  existing := []byte(u.PasswordHash)
  err := bcrypt.CompareHashAndPassword(existing, incoming)
  return err
}

func NewUserService() *UserService {
  return &UserService {collection: Session().DB("test").C("user")}
}

func(p *UserService) CreateUser(u *root.User) error {

  hash, salt, err := salt(u.Password)
  if err != nil {
    return err
  }

  return p.collection.Insert(&userModel{
    Username: u.Username,
    PasswordHash: hash,
    Salt: salt })
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


func salt(str string) (string, string, error) { 
  uuid, uuidErr := uuid()
  if uuidErr != nil {
    return "", "", uuidErr
  }

  hash, err := bcrypt.GenerateFromPassword([]byte(str+uuid), bcrypt.DefaultCost)
  
  hashString := string(hash[:])
  return hashString, uuid, err
}

func uuid() (string, error) {
  uuid := make([]byte, 16)
  n, err := io.ReadFull(rand.Reader, uuid)
  if n != len(uuid) || err != nil {
    return "", err
  }
  // variant bits; see section 4.1.1
  uuid[8] = uuid[8]&^0xc0 | 0x80
  // version 4 (pseudo-random); see section 4.1.3
  uuid[6] = uuid[6]&^0xf0 | 0x40
  return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}