package mongo

import (
	"errors"
	"github.com/cadmiumcat/books-api/config"
	"github.com/globalsign/mgo"
)

type Mongo struct {
	Collection string
	Database   string
	Session    *mgo.Session
	URI        string
}

func (m *Mongo) Init(mongoConfig config.MongoConfig) (err error) {
	if m.Session != nil {
		return errors.New("session already exists")
	}

	if m.Session, err = mgo.Dial(mongoConfig.BindAddr); err != nil {
		return err
	}

	return nil
}
