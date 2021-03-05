package mongo

import (
	"context"
	"errors"
	dpHealthcheck "github.com/ONSdigital/dp-healthcheck/healthcheck"
	dpMongodb "github.com/ONSdigital/dp-mongodb"
	dpMongoLock "github.com/ONSdigital/dp-mongodb/dplock"
	"github.com/ONSdigital/log.go/log"
	"github.com/cadmiumcat/books-api/config"
	"github.com/cadmiumcat/books-api/models"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Mongo contains the information needed to create and interact with a mongo session
type Mongo struct {
	Collection string
	Database   string
	Session    *mgo.Session
	URI        string
	lockClient *dpMongoLock.Lock
}

// Init initialises a mongo session with the given configuration.
// It returns an error if the session already exists or if it cannot connect.
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

// AddBook adds a Book
func (m *Mongo) AddBook(book *models.Book) {
	session := m.Session.Copy()
	defer session.Close()

	collection := session.DB(m.Database).C(m.Collection)
	collection.Insert(book)

	return
}

// GetBook returns a models.Book for a given ID.
// It returns an error if the Book is not found
func (m *Mongo) GetBook(ID string) (*models.Book, error) {
	session := m.Session.Copy()
	defer session.Close()

	var book models.Book
	err := session.DB(m.Database).C(m.Collection).Find(bson.M{"id": ID}).One(&book)

	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, errors.New("book not found")
		}
		return nil, err
	}

	return &book, err
}

// GetBooks returns all the existing models.Books.
// It returns an error if the models.Books cannot be listed.
func (m *Mongo) GetBooks() (models.Books, error) {
	session := m.Session.Copy()
	defer session.Close()

	list := session.DB(m.Database).C(m.Collection).Find(nil)

	books := &models.Books{}
	if err := list.All(&books.Items); err != nil {
		log.Event(nil, "can't get it", log.FATAL, log.Error(err))
	}

	return *books, nil
}

// Checker calls an api health endpoint and updates the provided CheckState
func (m *Mongo) Checker(ctx context.Context, state *dpHealthcheck.CheckState) error {
	if err := m.Healthcheck(ctx); err != nil {
		state.Update(dpHealthcheck.StatusCritical, err.Error(), 0)
		return nil
	}
	state.Update(dpHealthcheck.StatusOK, "Mongodb is ok", 0)
	return nil
}

// Healthcheck calls the service to check its health status
func (m *Mongo) Healthcheck(ctx context.Context) error {
	s := m.Session.Copy()
	defer s.Close()
	err := s.Ping()
	if err != nil {
		log.Event(ctx, "Ping mongo", log.ERROR, log.Error(err))
		return err
	}
	return nil
}
