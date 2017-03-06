package mongo

import (
  "gopkg.in/mgo.v2"
)

var (
  session *mgo.Session
)

func Session() *mgo.Session {
  if(session == nil) {
    createSession()
  }
  return session.Copy()
}

func Close() {
  session.Close()
}

func createSession() {
  var err error
  session, err = mgo.Dial("127.0.0.1:27017")
  if err != nil {
    panic(err)
  }
  session.SetMode(mgo.Monotonic, true)
}