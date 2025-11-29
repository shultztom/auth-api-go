package services

import (
	"auth-api-go/models"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	originalDB := models.DB
	models.DB = gormDB

	cleanup := func() {
		models.DB = originalDB
		db.Close()
	}

	return mock, cleanup
}

func TestCreateUser_Success(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "testuser", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	user, err := CreateUser("testuser", "password123")
	if err != nil {
		t.Errorf("CreateUser() error = %v", err)
		return
	}

	if user == nil {
		t.Error("CreateUser() returned nil user")
		return
	}

	if user.Username != "testuser" {
		t.Errorf("CreateUser() username = %v, want %v", user.Username, "testuser")
	}

	if user.Hash == "" {
		t.Error("CreateUser() hash is empty")
	}

	if user.Hash == "password123" {
		t.Error("CreateUser() should hash the password")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCreateUser_DBError(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "testuser", sqlmock.AnyArg()).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()

	user, err := CreateUser("testuser", "password123")
	if err == nil {
		t.Error("CreateUser() should return error on DB failure")
	}

	if user != nil {
		t.Error("CreateUser() should return nil user on error")
	}
}

func TestGetUserByUsername_Success(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "hash"}).
		AddRow(1, nil, nil, nil, "testuser", "$2a$14$hashedpassword")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("testuser", 1).
		WillReturnRows(rows)

	user, err := GetUserByUsername("testuser")
	if err != nil {
		t.Errorf("GetUserByUsername() error = %v", err)
		return
	}

	if user == nil {
		t.Error("GetUserByUsername() returned nil user")
		return
	}

	if user.Username != "testuser" {
		t.Errorf("GetUserByUsername() username = %v, want %v", user.Username, "testuser")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetUserByUsername_NotFound(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("nonexistent", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := GetUserByUsername("nonexistent")
	if err == nil {
		t.Error("GetUserByUsername() should return error for non-existent user")
	}

	if user != nil {
		t.Error("GetUserByUsername() should return nil user for non-existent user")
	}
}

func TestDeleteUserByUsername_Success(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "deleted_at"=$1 WHERE username = $2`)).
		WithArgs(sqlmock.AnyArg(), "testuser").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := DeleteUserByUsername("testuser")
	if err != nil {
		t.Errorf("DeleteUserByUsername() error = %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDeleteUserByUsername_DBError(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "deleted_at"=$1 WHERE username = $2`)).
		WithArgs(sqlmock.AnyArg(), "testuser").
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()

	err := DeleteUserByUsername("testuser")
	if err == nil {
		t.Error("DeleteUserByUsername() should return error on DB failure")
	}
}

func TestAuthenticateUser_Success(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	// First hash a password to use in the test
	hashedPassword, _ := HashPassword("password123")

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "hash"}).
		AddRow(1, nil, nil, nil, "testuser", hashedPassword)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("testuser", 1).
		WillReturnRows(rows)

	authenticated, err := AuthenticateUser("testuser", "password123")
	if err != nil {
		t.Errorf("AuthenticateUser() error = %v", err)
		return
	}

	if !authenticated {
		t.Error("AuthenticateUser() should return true for correct password")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestAuthenticateUser_WrongPassword(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	hashedPassword, _ := HashPassword("correctpassword")

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "hash"}).
		AddRow(1, nil, nil, nil, "testuser", hashedPassword)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("testuser", 1).
		WillReturnRows(rows)

	authenticated, err := AuthenticateUser("testuser", "wrongpassword")
	if err != nil {
		t.Errorf("AuthenticateUser() error = %v", err)
		return
	}

	if authenticated {
		t.Error("AuthenticateUser() should return false for incorrect password")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestAuthenticateUser_UserNotFound(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("nonexistent", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	authenticated, err := AuthenticateUser("nonexistent", "password123")
	if err == nil {
		t.Error("AuthenticateUser() should return error for non-existent user")
	}

	if authenticated {
		t.Error("AuthenticateUser() should return false for non-existent user")
	}
}
