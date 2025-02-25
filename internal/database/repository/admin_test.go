package repository

import (
	"context"
	"database/sql"
	"library-management/backend/internal/api/model"
	"library-management/backend/internal/database/transaction"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AdminRepositoryTestSuite struct {
	suite.Suite
	mock  sqlmock.Sqlmock
	db    *gorm.DB
	admin *AdminRepository
	sqlDB *sql.DB
	ctx   context.Context
}

func (s *AdminRepositoryTestSuite) SetupTest() {
	var err error
	s.sqlDB, s.mock, err = sqlmock.New()
	assert.NoError(s.T(), err)

	dialector := postgres.New(postgres.Config{
		Conn:       s.sqlDB,
		DriverName: "postgres",
	})

	s.db, err = gorm.Open(dialector, &gorm.Config{})
	assert.NoError(s.T(), err)

	txManager := transaction.NewTxManager(s.db)
	s.admin = NewAdminRepository(s.db, txManager)
	s.ctx = context.Background()
}

func (s *AdminRepositoryTestSuite) TearDownTest() {
	s.sqlDB.Close()
}

func TestAdminRepositorySuite(t *testing.T) {
	suite.Run(t, new(AdminRepositoryTestSuite))
}

func (s *AdminRepositoryTestSuite) TestAddBook() {
	book := &model.BookInventory{
		ISBN:            "1234567890",
		Title:           "Test Book",
		Authors:         "Test Author",
		Publisher:       "Test Publisher",
		Version:         "1.0",
		TotalCopies:     1,
		AvailableCopies: 1,
	}

	// Test case 1: Successful book addition
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
		WithArgs("admin@test.com").
		WillReturnRows(sqlmock.NewRows([]string{"email", "role"}).
			AddRow("admin@test.com", "admin"))

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "book_inventories"`)).
		WithArgs("1234567890").
		WillReturnRows(sqlmock.NewRows([]string{}))

	s.mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "book_inventories"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectCommit()

	err := s.admin.AddBook(s.ctx, book, "admin@test.com")
	assert.NoError(s.T(), err)
}

func (s *AdminRepositoryTestSuite) TestApproveIssueRequest() {
	requestID := "req123"
	approverID := "admin123"

	s.mock.ExpectBegin()

	// Mock existing request query
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "request_events"`)).
		WithArgs(requestID).
		WillReturnRows(sqlmock.NewRows([]string{"req_id", "book_id", "reader_id"}).
			AddRow(requestID, "1234567890", "reader123"))

	// Mock book inventory query
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "book_inventories"`)).
		WithArgs("1234567890").
		WillReturnRows(sqlmock.NewRows([]string{"isbn", "available_copies"}).
			AddRow("1234567890", 1))

	// Mock update available copies
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "book_inventories"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock update request events
	s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "request_events"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock create issue registry
	s.mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "issue_registries"`)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectCommit()

	err := s.admin.ApproveIssueRequest(s.ctx, requestID, approverID)
	assert.NoError(s.T(), err)
}

func (s *AdminRepositoryTestSuite) TestRejectIssueRequest() {
	requestID := "req123"

	s.mock.ExpectBegin()

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "request_events"`)).
		WithArgs(requestID).
		WillReturnRows(sqlmock.NewRows([]string{"req_id"}).
			AddRow(requestID))

	s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "request_events"`)).
		WithArgs(requestID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectCommit()

	err := s.admin.RejectIssueRequest(s.ctx, requestID)
	assert.NoError(s.T(), err)
}

func (s *AdminRepositoryTestSuite) TestListIssueRequests() {
	var requests []IssueRequestDetails

	s.mock.ExpectBegin()

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT r.*, b.title as book_title, b.available_copies FROM request_events r JOIN book_inventories b`)).
		WillReturnRows(sqlmock.NewRows([]string{"req_id", "book_title", "available_copies"}).
			AddRow("req123", "Test Book", 1))

	s.mock.ExpectCommit()

	err := s.admin.ListIssueRequests(s.ctx, &requests, "admin123")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, len(requests))
}
