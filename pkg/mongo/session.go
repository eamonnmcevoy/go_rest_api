package mongo

import (
  "go_rest_api/pkg"
  "gopkg.in/mgo.v2"
)

type Session struct {
  session *mgo.Session
}

func NewSession(config *root.MongoConfig) (*Session,error) {
  //var err error
  session, err := mgo.Dial(config.Ip)
  if err != nil {
    return nil,err
  }
  session.SetMode(mgo.Monotonic, true)
  return &Session{session}, err
}

func(s *Session) Copy() *mgo.Session {
  return s.session.Copy()
}

func(s *Session) Close() {
  if(s.session != nil) {
    s.session.Close()
  }
}

func(s *Session) DropDatabase(db string) error {
  if(s.session != nil) {
    return s.session.DB(db).DropDatabase()
  }
  return nil
}
