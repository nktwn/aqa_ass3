package testutils

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// MockDB holds the mock database and mock instance
type MockDB struct {
	DB   *sql.DB
	Mock sqlmock.Sqlmock
}

// NewMockDB creates a new mock database
func NewMockDB(t *testing.T) *MockDB {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}

	return &MockDB{
		DB:   db,
		Mock: mock,
	}
}

// Close closes the mock database
func (m *MockDB) Close() error {
	return m.DB.Close()
}

// ExpectationsWereMet checks if all expectations were met
func (m *MockDB) ExpectationsWereMet() error {
	return m.Mock.ExpectationsWereMet()
}

// SetupTestDB sets up a test database environment
func SetupTestDB(t *testing.T) (*MockDB, func()) {
	mock := NewMockDB(t)

	cleanup := func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
		mock.Close()
	}

	return mock, cleanup
}
