// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	"github.com/cadmiumcat/books-api/config"
	"github.com/cadmiumcat/books-api/interfaces"
	"github.com/cadmiumcat/books-api/models"
	"sync"
)

// Ensure, that DataStoreMock does implement interfaces.DataStore.
// If this is not the case, regenerate this file with moq.
var _ interfaces.DataStore = &DataStoreMock{}

// DataStoreMock is a mock implementation of interfaces.DataStore.
//
//     func TestSomethingThatUsesDataStore(t *testing.T) {
//
//         // make and configure a mocked interfaces.DataStore
//         mockedDataStore := &DataStoreMock{
//             AddBookFunc: func(ctx context.Context, book *models.Book) error {
// 	               panic("mock out the AddBook method")
//             },
//             CloseFunc: func(ctx context.Context) error {
// 	               panic("mock out the Close method")
//             },
//             GetBookFunc: func(ctx context.Context, id string) (*models.Book, error) {
// 	               panic("mock out the GetBook method")
//             },
//             GetBooksFunc: func(ctx context.Context) (models.Books, error) {
// 	               panic("mock out the GetBooks method")
//             },
//             GetReviewFunc: func(ctx context.Context, reviewID string) (*models.Review, error) {
// 	               panic("mock out the GetReview method")
//             },
//             GetReviewsFunc: func(ctx context.Context, bookID string) (models.Reviews, error) {
// 	               panic("mock out the GetReviews method")
//             },
//             InitFunc: func(in1 config.MongoConfig) error {
// 	               panic("mock out the Init method")
//             },
//         }
//
//         // use mockedDataStore in code that requires interfaces.DataStore
//         // and then make assertions.
//
//     }
type DataStoreMock struct {
	// AddBookFunc mocks the AddBook method.
	AddBookFunc func(ctx context.Context, book *models.Book) error

	// CloseFunc mocks the Close method.
	CloseFunc func(ctx context.Context) error

	// GetBookFunc mocks the GetBook method.
	GetBookFunc func(ctx context.Context, id string) (*models.Book, error)

	// GetBooksFunc mocks the GetBooks method.
	GetBooksFunc func(ctx context.Context) (models.Books, error)

	// GetReviewFunc mocks the GetReview method.
	GetReviewFunc func(ctx context.Context, reviewID string) (*models.Review, error)

	// GetReviewsFunc mocks the GetReviews method.
	GetReviewsFunc func(ctx context.Context, bookID string) (models.Reviews, error)

	// InitFunc mocks the Init method.
	InitFunc func(in1 config.MongoConfig) error

	// calls tracks calls to the methods.
	calls struct {
		// AddBook holds details about calls to the AddBook method.
		AddBook []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Book is the book argument value.
			Book *models.Book
		}
		// Close holds details about calls to the Close method.
		Close []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
		// GetBook holds details about calls to the GetBook method.
		GetBook []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the id argument value.
			ID string
		}
		// GetBooks holds details about calls to the GetBooks method.
		GetBooks []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
		// GetReview holds details about calls to the GetReview method.
		GetReview []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ReviewID is the reviewID argument value.
			ReviewID string
		}
		// GetReviews holds details about calls to the GetReviews method.
		GetReviews []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// BookID is the bookID argument value.
			BookID string
		}
		// Init holds details about calls to the Init method.
		Init []struct {
			// In1 is the in1 argument value.
			In1 config.MongoConfig
		}
	}
	lockAddBook    sync.RWMutex
	lockClose      sync.RWMutex
	lockGetBook    sync.RWMutex
	lockGetBooks   sync.RWMutex
	lockGetReview  sync.RWMutex
	lockGetReviews sync.RWMutex
	lockInit       sync.RWMutex
}

// AddBook calls AddBookFunc.
func (mock *DataStoreMock) AddBook(ctx context.Context, book *models.Book) error {
	if mock.AddBookFunc == nil {
		panic("DataStoreMock.AddBookFunc: method is nil but DataStore.AddBook was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Book *models.Book
	}{
		Ctx:  ctx,
		Book: book,
	}
	mock.lockAddBook.Lock()
	mock.calls.AddBook = append(mock.calls.AddBook, callInfo)
	mock.lockAddBook.Unlock()
	return mock.AddBookFunc(ctx, book)
}

// AddBookCalls gets all the calls that were made to AddBook.
// Check the length with:
//     len(mockedDataStore.AddBookCalls())
func (mock *DataStoreMock) AddBookCalls() []struct {
	Ctx  context.Context
	Book *models.Book
} {
	var calls []struct {
		Ctx  context.Context
		Book *models.Book
	}
	mock.lockAddBook.RLock()
	calls = mock.calls.AddBook
	mock.lockAddBook.RUnlock()
	return calls
}

// Close calls CloseFunc.
func (mock *DataStoreMock) Close(ctx context.Context) error {
	if mock.CloseFunc == nil {
		panic("DataStoreMock.CloseFunc: method is nil but DataStore.Close was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockClose.Lock()
	mock.calls.Close = append(mock.calls.Close, callInfo)
	mock.lockClose.Unlock()
	return mock.CloseFunc(ctx)
}

// CloseCalls gets all the calls that were made to Close.
// Check the length with:
//     len(mockedDataStore.CloseCalls())
func (mock *DataStoreMock) CloseCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockClose.RLock()
	calls = mock.calls.Close
	mock.lockClose.RUnlock()
	return calls
}

// GetBook calls GetBookFunc.
func (mock *DataStoreMock) GetBook(ctx context.Context, id string) (*models.Book, error) {
	if mock.GetBookFunc == nil {
		panic("DataStoreMock.GetBookFunc: method is nil but DataStore.GetBook was just called")
	}
	callInfo := struct {
		Ctx context.Context
		ID  string
	}{
		Ctx: ctx,
		ID:  id,
	}
	mock.lockGetBook.Lock()
	mock.calls.GetBook = append(mock.calls.GetBook, callInfo)
	mock.lockGetBook.Unlock()
	return mock.GetBookFunc(ctx, id)
}

// GetBookCalls gets all the calls that were made to GetBook.
// Check the length with:
//     len(mockedDataStore.GetBookCalls())
func (mock *DataStoreMock) GetBookCalls() []struct {
	Ctx context.Context
	ID  string
} {
	var calls []struct {
		Ctx context.Context
		ID  string
	}
	mock.lockGetBook.RLock()
	calls = mock.calls.GetBook
	mock.lockGetBook.RUnlock()
	return calls
}

// GetBooks calls GetBooksFunc.
func (mock *DataStoreMock) GetBooks(ctx context.Context) (models.Books, error) {
	if mock.GetBooksFunc == nil {
		panic("DataStoreMock.GetBooksFunc: method is nil but DataStore.GetBooks was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockGetBooks.Lock()
	mock.calls.GetBooks = append(mock.calls.GetBooks, callInfo)
	mock.lockGetBooks.Unlock()
	return mock.GetBooksFunc(ctx)
}

// GetBooksCalls gets all the calls that were made to GetBooks.
// Check the length with:
//     len(mockedDataStore.GetBooksCalls())
func (mock *DataStoreMock) GetBooksCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockGetBooks.RLock()
	calls = mock.calls.GetBooks
	mock.lockGetBooks.RUnlock()
	return calls
}

// GetReview calls GetReviewFunc.
func (mock *DataStoreMock) GetReview(ctx context.Context, reviewID string) (*models.Review, error) {
	if mock.GetReviewFunc == nil {
		panic("DataStoreMock.GetReviewFunc: method is nil but DataStore.GetReview was just called")
	}
	callInfo := struct {
		Ctx      context.Context
		ReviewID string
	}{
		Ctx:      ctx,
		ReviewID: reviewID,
	}
	mock.lockGetReview.Lock()
	mock.calls.GetReview = append(mock.calls.GetReview, callInfo)
	mock.lockGetReview.Unlock()
	return mock.GetReviewFunc(ctx, reviewID)
}

// GetReviewCalls gets all the calls that were made to GetReview.
// Check the length with:
//     len(mockedDataStore.GetReviewCalls())
func (mock *DataStoreMock) GetReviewCalls() []struct {
	Ctx      context.Context
	ReviewID string
} {
	var calls []struct {
		Ctx      context.Context
		ReviewID string
	}
	mock.lockGetReview.RLock()
	calls = mock.calls.GetReview
	mock.lockGetReview.RUnlock()
	return calls
}

// GetReviews calls GetReviewsFunc.
func (mock *DataStoreMock) GetReviews(ctx context.Context, bookID string) (models.Reviews, error) {
	if mock.GetReviewsFunc == nil {
		panic("DataStoreMock.GetReviewsFunc: method is nil but DataStore.GetReviews was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		BookID string
	}{
		Ctx:    ctx,
		BookID: bookID,
	}
	mock.lockGetReviews.Lock()
	mock.calls.GetReviews = append(mock.calls.GetReviews, callInfo)
	mock.lockGetReviews.Unlock()
	return mock.GetReviewsFunc(ctx, bookID)
}

// GetReviewsCalls gets all the calls that were made to GetReviews.
// Check the length with:
//     len(mockedDataStore.GetReviewsCalls())
func (mock *DataStoreMock) GetReviewsCalls() []struct {
	Ctx    context.Context
	BookID string
} {
	var calls []struct {
		Ctx    context.Context
		BookID string
	}
	mock.lockGetReviews.RLock()
	calls = mock.calls.GetReviews
	mock.lockGetReviews.RUnlock()
	return calls
}

// Init calls InitFunc.
func (mock *DataStoreMock) Init(in1 config.MongoConfig) error {
	if mock.InitFunc == nil {
		panic("DataStoreMock.InitFunc: method is nil but DataStore.Init was just called")
	}
	callInfo := struct {
		In1 config.MongoConfig
	}{
		In1: in1,
	}
	mock.lockInit.Lock()
	mock.calls.Init = append(mock.calls.Init, callInfo)
	mock.lockInit.Unlock()
	return mock.InitFunc(in1)
}

// InitCalls gets all the calls that were made to Init.
// Check the length with:
//     len(mockedDataStore.InitCalls())
func (mock *DataStoreMock) InitCalls() []struct {
	In1 config.MongoConfig
} {
	var calls []struct {
		In1 config.MongoConfig
	}
	mock.lockInit.RLock()
	calls = mock.calls.Init
	mock.lockInit.RUnlock()
	return calls
}
