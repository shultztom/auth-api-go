package services

import (
	"auth-api-go/models"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupRoleMockDB(t *testing.T) (sqlmock.Sqlmock, func()) {
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

func TestGetRolesByUsername_Success(t *testing.T) {
	mock, cleanup := setupRoleMockDB(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "role"}).
		AddRow(1, nil, nil, nil, "testuser", "admin").
		AddRow(2, nil, nil, nil, "testuser", "user")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles" WHERE username = $1`)).
		WithArgs("testuser").
		WillReturnRows(rows)

	roles, err := GetRolesByUsername("testuser")
	if err != nil {
		t.Errorf("GetRolesByUsername() error = %v", err)
		return
	}

	if len(roles) != 2 {
		t.Errorf("GetRolesByUsername() returned %d roles, want 2", len(roles))
		return
	}

	if roles[0].Role != "admin" {
		t.Errorf("GetRolesByUsername() first role = %v, want admin", roles[0].Role)
	}

	if roles[1].Role != "user" {
		t.Errorf("GetRolesByUsername() second role = %v, want user", roles[1].Role)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetRolesByUsername_NoRoles(t *testing.T) {
	mock, cleanup := setupRoleMockDB(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "role"})

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles" WHERE username = $1`)).
		WithArgs("testuser").
		WillReturnRows(rows)

	roles, err := GetRolesByUsername("testuser")
	if err != nil {
		t.Errorf("GetRolesByUsername() error = %v", err)
		return
	}

	if len(roles) != 0 {
		t.Errorf("GetRolesByUsername() returned %d roles, want 0", len(roles))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetRolesByUsername_DBError(t *testing.T) {
	mock, cleanup := setupRoleMockDB(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles" WHERE username = $1`)).
		WithArgs("testuser").
		WillReturnError(gorm.ErrInvalidDB)

	roles, err := GetRolesByUsername("testuser")
	if err == nil {
		t.Error("GetRolesByUsername() should return error on DB failure")
	}

	if roles != nil {
		t.Error("GetRolesByUsername() should return nil roles on error")
	}
}

func TestRoleCheck_HasRole(t *testing.T) {
	mock, cleanup := setupRoleMockDB(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "role"}).
		AddRow(1, nil, nil, nil, "testuser", "admin").
		AddRow(2, nil, nil, nil, "testuser", "user")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles" WHERE username = $1`)).
		WithArgs("testuser").
		WillReturnRows(rows)

	hasRole, err := RoleCheck("admin", "testuser")
	if err != nil {
		t.Errorf("RoleCheck() error = %v", err)
		return
	}

	if !hasRole {
		t.Error("RoleCheck() should return true when user has the role")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestRoleCheck_DoesNotHaveRole(t *testing.T) {
	mock, cleanup := setupRoleMockDB(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "role"}).
		AddRow(1, nil, nil, nil, "testuser", "user")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles" WHERE username = $1`)).
		WithArgs("testuser").
		WillReturnRows(rows)

	hasRole, err := RoleCheck("admin", "testuser")
	if err != nil {
		t.Errorf("RoleCheck() error = %v", err)
		return
	}

	if hasRole {
		t.Error("RoleCheck() should return false when user does not have the role")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestRoleCheck_NoRoles(t *testing.T) {
	mock, cleanup := setupRoleMockDB(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "role"})

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles" WHERE username = $1`)).
		WithArgs("testuser").
		WillReturnRows(rows)

	hasRole, err := RoleCheck("admin", "testuser")
	if err != nil {
		t.Errorf("RoleCheck() error = %v", err)
		return
	}

	if hasRole {
		t.Error("RoleCheck() should return false when user has no roles")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestRoleCheck_DBError(t *testing.T) {
	mock, cleanup := setupRoleMockDB(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles" WHERE username = $1`)).
		WithArgs("testuser").
		WillReturnError(gorm.ErrInvalidDB)

	hasRole, err := RoleCheck("admin", "testuser")
	if err == nil {
		t.Error("RoleCheck() should return error on DB failure")
	}

	if hasRole {
		t.Error("RoleCheck() should return false on error")
	}
}

//func TestAddRole_Success(t *testing.T) {
//	mock, cleanup := setupRoleMockDB(t)
//	defer cleanup()
//
//	mock.ExpectBegin()
//	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "roles"`)).
//		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "testuser", "admin").
//		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
//	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles" WHERE "roles"."deleted_at" IS NULL AND "roles"."id" = $1`)).
//		WithArgs(1).
//		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "role"}).
//			AddRow(1, nil, nil, nil, "testuser", "admin"))
//	mock.ExpectCommit()
//
//	err := AddRole("testuser", "admin")
//	if err != nil {
//		t.Errorf("AddRole() error = %v", err)
//	}
//
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("Unfulfilled expectations: %v", err)
//	}
//}

func TestAddRole_DBError(t *testing.T) {
	mock, cleanup := setupRoleMockDB(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "roles"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "testuser", "admin").
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()

	err := AddRole("testuser", "admin")
	if err == nil {
		t.Error("AddRole() should return error on DB failure")
	}
}

//func TestAddRole_EmptyRole(t *testing.T) {
//	mock, cleanup := setupRoleMockDB(t)
//	defer cleanup()
//
//	mock.ExpectBegin()
//	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "roles"`)).
//		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "testuser", "").
//		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
//	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles" WHERE "roles"."deleted_at" IS NULL AND "roles"."id" = $1`)).
//		WithArgs(1).
//		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "role"}).
//			AddRow(1, nil, nil, nil, "testuser", ""))
//	mock.ExpectCommit()
//
//	err := AddRole("testuser", "")
//	if err != nil {
//		t.Errorf("AddRole() error = %v", err)
//	}
//
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("Unfulfilled expectations: %v", err)
//	}
//}
