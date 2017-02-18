package user

import (
  "gopkg.in/mgo.v2/bson"
)

type User struct {
  Id       bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
  Username string `json:"username"`
}