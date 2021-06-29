package mongo

import (
	"context"
	"fmt"
	dpMongodb "github.com/ONSdigital/dp-mongodb"
	dpMongoLock "github.com/ONSdigital/dp-mongodb/dplock"
	"github.com/ONSdigital/log.go/log"
	"github.com/cadmiumcat/books-api/config"
	"github.com/cadmiumcat/books-api/models"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"time"
)

// Mongo contains the information needed to create and interact with a mongo session
type Mongo struct {
	BooksCollection   string
	ReviewsCollection string
	Database          string
	Session           *mgo.Session
	URI               string
	lockClient        *dpMongoLock.Lock
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

	m.BooksCollection = mongoConfig.BooksCollection
	m.ReviewsCollection = mongoConfig.ReviewsCollection
	m.Database = mongoConfig.Database

	return nil
}

// Close closes the mongo session and returns any error
func (m *Mongo) Close(ctx context.Context) (err error) {
	m.lockClient.Close(ctx)
	return dpMongodb.Close(ctx, m.Session)
}

// AddBook adds a Book
func (m *Mongo) AddBook(ctx context.Context, book *models.Book) error {
	session := m.Session.Copy()
	defer session.Close()

	logData := log.Data{
		"book": book,
	}

	collection := session.DB(m.Database).C(m.BooksCollection)
	err := collection.Insert(book)
	if err != nil {
		log.Event(ctx, "unexpected error when adding a book", log.ERROR, log.Error(err), logData)
		return errors.Wrap(err, "unexpected error when adding a book")
	}

	return nil
}

// GetBook returns a models.Book for a given ID.
// It returns an error if the Book is not found
func (m *Mongo) GetBook(ctx context.Context, ID string) (*models.Book, error) {
	session := m.Session.Copy()
	defer session.Close()

	logData := log.Data{
		"book_id":    ID,
		"database":   m.Database,
		"collection": m.BooksCollection}

	var book models.Book
	err := session.DB(m.Database).C(m.BooksCollection).Find(bson.M{"_id": ID}).One(&book)

	if err != nil {
		if err == mgo.ErrNotFound {
			log.Event(ctx, ErrBookNotFound.Error(), log.ERROR, log.Error(err), logData)
			return nil, ErrBookNotFound
		}
		log.Event(ctx, "unexpected error when getting a book", log.ERROR, log.Error(err), logData)
		return nil, errors.Wrap(err, "unexpected error when getting a book")
	}

	return &book, nil
}

// GetBooks returns all the existing []models.Book.
// It returns an error if the []models.Book cannot be listed.
func (m *Mongo) GetBooks(ctx context.Context, offset, limit int) ([]models.Book, int, error) {

	session := m.Session.Copy()
	defer session.Close()

	logData := log.Data{
		"database":   m.Database,
		"collection": m.BooksCollection}

	list := session.DB(m.Database).C(m.BooksCollection).Find(nil)
	var books []models.Book

	totalCount, err := list.Count()
	if err != nil {
		if err == mgo.ErrNotFound {
			log.Event(ctx, "no book resources found", log.WARN, log.Error(err))
			return []models.Book{}, totalCount, nil
		}
		log.Event(ctx, "failure to retrieve list of books", log.ERROR, log.Error(err))
		return nil, totalCount, err
	}

	if limit > 0 {
		iter := list.Skip(offset).Limit(limit).Iter()

		defer func() {
			err := iter.Close()
			if err != nil {
				log.Event(ctx, "error closing iterator", log.ERROR, log.Error(err), logData)
			}
		}()

		if err := iter.All(&books); err != nil {
			if err == mgo.ErrNotFound {
				return []models.Book{}, totalCount, nil
			}
			log.Event(ctx, "unable to retrieve books", log.ERROR, log.Error(err), logData)
			return []models.Book{}, totalCount, errors.Wrap(err, "unexpected error when getting books")
		}
	}

	return books, totalCount, nil
}

// AddReview adds a Review to a Book
func (m *Mongo) AddReview(ctx context.Context, review *models.Review) error {
	session := m.Session.Copy()
	defer session.Close()

	logData := log.Data{
		"review": review,
	}

	collection := session.DB(m.Database).C(m.ReviewsCollection)
	err := collection.Insert(review)

	if err != nil {
		log.Event(ctx, "unexpected error when adding a book", log.ERROR, log.Error(err), logData)
		return errors.Wrap(err, "unexpected error when adding a book")
	}

	return nil
}

// UpdateReview updates an existing Review.
// Only the message and user can be updated.
// It returns an error if the review is not found
func (m *Mongo) UpdateReview(ctx context.Context, reviewID string, review *models.Review) error {
	s := m.Session.Copy()
	defer s.Close()

	logData := log.Data{
		"review_id":  reviewID,
		"database":   m.Database,
		"collection": m.ReviewsCollection}

	updates := make(bson.M)
	if review.Message != "" {
		updates["message"] = review.Message
	}

	if review.User != (models.User{}) {
		if review.User.Forenames != "" {
			updates["user.forenames"] = review.User.Forenames
		}
		if review.User.Surname != "" {
			updates["user.surname"] = review.User.Surname
		}
	}

	if len(updates) == 0 {
		return nil
	}

	updates["last_updated"] = time.Now().UTC()

	update := bson.M{"$set": updates}
	if err := s.DB(m.Database).C(m.ReviewsCollection).UpdateId(reviewID, update); err != nil {
		if err == mgo.ErrNotFound {
			log.Event(ctx, ErrReviewNotFound.Error(), log.ERROR, log.Error(err), logData)
			return ErrReviewNotFound
		}
		return err
	}

	return nil
}

// GetReview returns a models.Review for a given reviewID.
// It returns an error if the review is not found.
func (m *Mongo) GetReview(ctx context.Context, reviewID string) (*models.Review, error) {
	session := m.Session.Copy()
	defer session.Close()

	logData := log.Data{
		"review_id":  reviewID,
		"database":   m.Database,
		"collection": m.ReviewsCollection}

	var review models.Review
	err := session.DB(m.Database).C(m.ReviewsCollection).Find(bson.M{"_id": reviewID}).One(&review)

	if err != nil {
		if err == mgo.ErrNotFound {
			log.Event(ctx, ErrReviewNotFound.Error(), log.ERROR, log.Error(err), logData)
			return nil, ErrReviewNotFound
		}
		return nil, errors.Wrap(err, "unexpected error when getting a review")
	}

	return &review, nil
}

// GetReviews returns all the existing models.Reviews.
// It returns an error if the models.Reviews cannot be listed.
func (m *Mongo) GetReviews(ctx context.Context, bookID string, offset, limit int) ([]models.Review, int, error) {

	session := m.Session.Copy()
	defer session.Close()

	logData := log.Data{
		"database":   m.Database,
		"collection": m.ReviewsCollection}

	list := session.DB(m.Database).C(m.ReviewsCollection).Find(bson.M{"links.book": fmt.Sprintf("/books/%s", bookID)})
	var reviews []models.Review

	totalCount, err := list.Count()
	if err != nil {
		if err == mgo.ErrNotFound {
			log.Event(ctx, "no reviews resources found for the book", log.WARN, log.Error(err))
			return []models.Review{}, totalCount, nil
		}
		log.Event(ctx, "failure to retrieve list of reviews", log.ERROR, log.Error(err))
		return nil, totalCount, err
	}

	if limit > 0 {
		iter := list.Skip(offset).Limit(limit).Iter()

		defer func() {
			err := iter.Close()
			if err != nil {
				log.Event(ctx, "error closing iterator", log.ERROR, log.Error(err), logData)
			}
		}()

		if err := iter.All(&reviews); err != nil {
			if err == mgo.ErrNotFound {
				return []models.Review{}, totalCount, nil
			}
			log.Event(ctx, "unable to retrieve reviews", log.ERROR, log.Error(err), logData)
			return []models.Review{}, totalCount, errors.Wrap(err, "unexpected error when getting reviews")
		}
	}

	return reviews, totalCount, nil
}
