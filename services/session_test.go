package services

import (
	"auth-api-go/redis"
	"errors"
	"testing"

	"github.com/go-redis/redismock/v9"
)

func TestDeleteSessionInRedis_Success(t *testing.T) {
	db, mock := redismock.NewClientMock()
	originalRedis := redis.REDIS
	redis.REDIS = db
	defer func() {
		redis.REDIS = originalRedis
	}()

	mock.ExpectDel("testuser-token").SetVal(1)

	success, err := DeleteSessionInRedis("testuser")
	if err != nil {
		t.Errorf("DeleteSessionInRedis() error = %v", err)
		return
	}

	if !success {
		t.Error("DeleteSessionInRedis() should return true on success")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDeleteSessionInRedis_KeyNotExists(t *testing.T) {
	db, mock := redismock.NewClientMock()
	originalRedis := redis.REDIS
	redis.REDIS = db
	defer func() {
		redis.REDIS = originalRedis
	}()

	// Del returns 0 when key doesn't exist but no error
	mock.ExpectDel("nonexistent-token").SetVal(0)

	success, err := DeleteSessionInRedis("nonexistent")
	if err != nil {
		t.Errorf("DeleteSessionInRedis() error = %v", err)
		return
	}

	if !success {
		t.Error("DeleteSessionInRedis() should return true even if key doesn't exist")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDeleteSessionInRedis_RedisError(t *testing.T) {
	db, mock := redismock.NewClientMock()
	originalRedis := redis.REDIS
	redis.REDIS = db
	defer func() {
		redis.REDIS = originalRedis
	}()

	mock.ExpectDel("testuser-token").SetErr(errors.New("redis connection error"))

	success, err := DeleteSessionInRedis("testuser")
	if err == nil {
		t.Error("DeleteSessionInRedis() should return error on Redis failure")
	}

	if success {
		t.Error("DeleteSessionInRedis() should return false on error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDeleteSessionInRedis_EmptyUsername(t *testing.T) {
	db, mock := redismock.NewClientMock()
	originalRedis := redis.REDIS
	redis.REDIS = db
	defer func() {
		redis.REDIS = originalRedis
	}()

	mock.ExpectDel("-token").SetVal(0)

	success, err := DeleteSessionInRedis("")
	if err != nil {
		t.Errorf("DeleteSessionInRedis() error = %v", err)
		return
	}

	if !success {
		t.Error("DeleteSessionInRedis() should return true")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestDeleteSessionInRedis_SpecialCharacters(t *testing.T) {
	db, mock := redismock.NewClientMock()
	originalRedis := redis.REDIS
	redis.REDIS = db
	defer func() {
		redis.REDIS = originalRedis
	}()

	mock.ExpectDel("user@email.com-token").SetVal(1)

	success, err := DeleteSessionInRedis("user@email.com")
	if err != nil {
		t.Errorf("DeleteSessionInRedis() error = %v", err)
		return
	}

	if !success {
		t.Error("DeleteSessionInRedis() should return true on success")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
