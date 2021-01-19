package mongo

import (
	"context"
	"errors"
	dpMongodb "github.com/ONSdigital/dp-mongodb"
	dpMongoLock "github.com/ONSdigital/dp-mongodb/dplock"
	"github.com/cadmiumcat/books-api/config"
	"github.com/cadmiumcat/books-api/models"
	"github.com/globalsign/mgo"
)

type Mongo struct {
	Collection string
	Database   string
	Session    *mgo.Session
	URI        string
	lockClient   *dpMongoLock.Lock
}

func (m *Mongo) Init(mongoConfig config.MongoConfig) (err error) {
	if m.Session != nil {
		return errors.New("session already exists")
	}

	if m.Session, err = mgo.Dial(mongoConfig.BindAddr); err != nil {
		return err
	}

	m.Collection = mongoConfig.Collection
	m.Database = mongoConfig.Database

	return nil
}


// Close closes the mongo session and returns any error
func (m *Mongo) Close(ctx context.Context) (err error) {
	m.lockClient.Close(ctx)
	return dpMongodb.Close(ctx, m.Session)
}

func (m *Mongo) AddBook(book *models.Book) {
	session := m.Session.Copy()
	defer session.Close()

	collection := session.DB(m.Database).C(m.Collection)
	collection.Insert(book)
}