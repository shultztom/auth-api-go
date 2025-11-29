package services

import (
	"testing"

	"github.com/golang-jwt/jwt"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "hash simple password",
			password: "password123",
		},
		{
			name:     "hash complex password",
			password: "C0mpl3x!P@ssw0rd#2024",
		},
		{
			name:     "hash empty password",
			password: "",
		},
		{
			name:     "hash long password",
			password: "thisIsAVeryLongPasswordThatShouldStillWorkCorrectly12345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if err != nil {
				t.Errorf("HashPassword() error = %v", err)
				return
			}
			if hash == "" {
				t.Error("HashPassword() returned empty hash")
			}
			if hash == tt.password {
				t.Error("HashPassword() returned unhashed password")
			}
		})
	}
}

func TestHashPassword_UniqueSalts(t *testing.T) {
	password := "samePassword"
	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	if hash1 == hash2 {
		t.Error("HashPassword() should generate unique hashes for same password (different salts)")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "correct password",
			password: "password123",
			want:     true,
		},
		{
			name:     "complex password",
			password: "C0mpl3x!P@ssw0rd#2024",
			want:     true,
		},
		{
			name:     "empty password",
			password: "",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if err != nil {
				t.Fatalf("HashPassword() error = %v", err)
			}

			got := CheckPasswordHash(tt.password, hash)
			if got != tt.want {
				t.Errorf("CheckPasswordHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckPasswordHash_WrongPassword(t *testing.T) {
	password := "correctPassword"
	wrongPassword := "wrongPassword"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	got := CheckPasswordHash(wrongPassword, hash)
	if got != false {
		t.Error("CheckPasswordHash() should return false for wrong password")
	}
}

func TestCheckPasswordHash_InvalidHash(t *testing.T) {
	password := "password123"
	invalidHash := "notavalidhash"

	got := CheckPasswordHash(password, invalidHash)
	if got != false {
		t.Error("CheckPasswordHash() should return false for invalid hash")
	}
}

func TestCheckPasswordHash_EmptyHash(t *testing.T) {
	password := "password123"
	emptyHash := ""

	got := CheckPasswordHash(password, emptyHash)
	if got != false {
		t.Error("CheckPasswordHash() should return false for empty hash")
	}
}

func TestParseToken_EmptyToken(t *testing.T) {
	jwtKey := []byte("testsecret")

	_, err := ParseToken("", jwtKey)
	if err == nil {
		t.Error("ParseToken() should return error for empty token")
	}
	if err.Error() != "missing token" {
		t.Errorf("ParseToken() error = %v, want 'missing token'", err)
	}
}

func TestParseToken_InvalidToken(t *testing.T) {
	jwtKey := []byte("testsecret")

	_, err := ParseToken("invalid.token.here", jwtKey)
	if err == nil {
		t.Error("ParseToken() should return error for invalid token")
	}
}

func TestParseToken_WrongSigningKey(t *testing.T) {
	// Create a token with one key
	claims := &Claims{
		Username: "testuser",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 9999999999,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("correctkey"))
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	// Try to parse with wrong key
	_, err = ParseToken(tokenString, []byte("wrongkey"))
	// This should fail because the signature won't match, or Redis lookup will fail
	// Either way, we expect an error
	if err == nil {
		t.Error("ParseToken() should return error when token is signed with different key or session doesn't exist")
	}
}

func TestVerifyToken_EmptyToken(t *testing.T) {
	jwtKey := []byte("testsecret")

	valid, err := VerifyToken("", jwtKey)
	if err == nil {
		t.Error("VerifyToken() should return error for empty token")
	}
	if valid != false {
		t.Error("VerifyToken() should return false for empty token")
	}
}

func TestVerifyToken_InvalidToken(t *testing.T) {
	jwtKey := []byte("testsecret")

	valid, err := VerifyToken("invalid.token.here", jwtKey)
	if err == nil {
		t.Error("VerifyToken() should return error for invalid token")
	}
	if valid != false {
		t.Error("VerifyToken() should return false for invalid token")
	}
}

func TestGetUsernameFromToken_EmptyToken(t *testing.T) {
	jwtKey := []byte("testsecret")

	username, err := GetUsernameFromToken("", jwtKey)
	if err == nil {
		t.Error("GetUsernameFromToken() should return error for empty token")
	}
	if username != "" {
		t.Error("GetUsernameFromToken() should return empty string for empty token")
	}
}

func TestGetUsernameFromToken_InvalidToken(t *testing.T) {
	jwtKey := []byte("testsecret")

	username, err := GetUsernameFromToken("invalid.token.here", jwtKey)
	if err == nil {
		t.Error("GetUsernameFromToken() should return error for invalid token")
	}
	if username != "" {
		t.Error("GetUsernameFromToken() should return empty string for invalid token")
	}
}
