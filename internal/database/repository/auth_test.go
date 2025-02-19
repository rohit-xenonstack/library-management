package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"library-management/backend/internal/api/model"
	"library-management/backend/internal/database/transaction"
	"library-management/backend/internal/util"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, error) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	dialector := postgres.New(postgres.Config{
		Conn:       sqlDB,
		DriverName: "postgres",
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	return db, mock, nil
}

func TestNewAuthRepository(t *testing.T) {
	db, _, err := setupTestDB(t)
	assert.NoError(t, err)

	txManager := transaction.NewTxManager(db)
	repo := NewAuthRepository(db, txManager)

	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.DB)
	assert.Equal(t, txManager, repo.txManager)
}

func TestAuthRepository_Login_Success_WithAllFields(t *testing.T) {
	db, mock, err := setupTestDB(t)
	assert.NoError(t, err)

	txManager := transaction.NewTxManager(db)
	repo := NewAuthRepository(db, txManager)
	ctx := context.Background()

	email := "test@example.com"
	name := "Test User"
	contactNumber := "1234567890"
	role := "admin"
	libID := util.RandomUUID()

	expectedUser := &model.Users{
		ID:            util.RandomUUID(),
		Name:          name,
		Email:         email,
		ContactNumber: contactNumber,
		Role:          role,
		LibID:         &libID,
	}

	// Expect transaction operations
	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"id", "name", "email", "contact_number", "role", "lib_id", "created_at", "updated_at", "deleted_at"}).
		AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.ContactNumber, expectedUser.Role, expectedUser.LibID, time.Now(), time.Now(), nil)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(email, 1).
		WillReturnRows(rows)
	mock.ExpectCommit()

	user, err := repo.Login(ctx, email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Name, user.Name)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.ContactNumber, user.ContactNumber)
	assert.Equal(t, expectedUser.Role, user.Role)
	assert.Equal(t, expectedUser.LibID, user.LibID)
}

func TestAuthRepository_Login_Success_WithoutLibrary(t *testing.T) {
	db, mock, err := setupTestDB(t)
	assert.NoError(t, err)

	txManager := transaction.NewTxManager(db)
	repo := NewAuthRepository(db, txManager)
	ctx := context.Background()

	email := "test@example.com"
	name := "Test User"
	contactNumber := "1234567890"
	role := "reader"

	expectedUser := &model.Users{
		ID:            util.RandomUUID(),
		Name:          name,
		Email:         email,
		ContactNumber: contactNumber,
		Role:          role,
		LibID:         nil, // explicitly set to nil for reader role
	}

	// Expect transaction operations
	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"id", "name", "email", "contact_number", "role", "lib_id", "created_at", "updated_at", "deleted_at"}).
		AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.ContactNumber, expectedUser.Role, nil, time.Now(), time.Now(), nil)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(email, 1).
		WillReturnRows(rows)
	mock.ExpectCommit()

	user, err := repo.Login(ctx, email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Name, user.Name)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.ContactNumber, user.ContactNumber)
	assert.Equal(t, expectedUser.Role, user.Role)
	assert.Equal(t, expectedUser.LibID, user.LibID)
}

func TestAuthRepository_Login_DBError(t *testing.T) {
	db, mock, err := setupTestDB(t)
	assert.NoError(t, err)

	txManager := transaction.NewTxManager(db)
	repo := NewAuthRepository(db, txManager)
	ctx := context.Background()
	email := "test@example.com"

	// Expect transaction operations
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(email, 1).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	user, err := repo.Login(ctx, email)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, sql.ErrConnDone, err)
}

func TestAuthRepository_Login_TransactionError(t *testing.T) {
	db, mock, err := setupTestDB(t)
	assert.NoError(t, err)

	txManager := transaction.NewTxManager(db)
	repo := NewAuthRepository(db, txManager)
	ctx := context.Background()
	email := "test@example.com"

	// Expect failed transaction begin
	mock.ExpectBegin().WillReturnError(sql.ErrConnDone)

	user, err := repo.Login(ctx, email)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, sql.ErrConnDone, err)
}
