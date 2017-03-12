package mongo

import (
  "gopkg.in/mgo.v2"
)

type Session struct {
  session *mgo.Session
}

func(s *Session) Open() error {
  var err error
  s.session, err = mgo.Dial("127.0.0.1:27017")
  if err != nil {
    return err
  }
  s.session.SetMode(mgo.Monotonic, true)
  return nil
}

func(s *Session) Copy() *mgo.Session {
  return s.session.Copy()
}

func(s *Session) Close() {
  if(s.session != nil) {
    s.session.Close()
  }
}
