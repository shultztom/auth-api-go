package controllers

import (
	"auth-api-go/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

// Structs

type userRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// JWT secret

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

// Utils

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateToken(username string) string {
	expirationTime := time.Now().Add(8 * time.Hour)

	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// TODO improve error handling
		return "error"
	}
	return tokenString
}

// Register POST /register
func Register(c *gin.Context) {
	var newUser userRequest

	if err := c.BindJSON(&newUser); err != nil {
		return
	}
	hash, _ := HashPassword(newUser.Password)

	userEntry := &models.User{
		Username: newUser.Username,
		Hash:     hash,
	}

	err := models.DB.Create(userEntry).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	token := CreateToken(userEntry.Username)

	c.JSON(http.StatusCreated, gin.H{"token": token})

}

// Login POST /login
func Login(c *gin.Context) {
	var userReq userRequest

	if err := c.BindJSON(&userReq); err != nil {
		return
	}

	var user models.User

	if err := models.DB.Where("username = ?", userReq.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "User not found!"})
		return
	}

	isMatch := CheckPasswordHash(userReq.Password, user.Hash)
	if isMatch {
		token := CreateToken(userReq.Username)
		c.JSON(http.StatusOK, gin.H{"token": token})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Password is incorrect!"})
	}
}

// Verify GET /verify
func Verify(c *gin.Context) {
	tokenHeader := c.GetHeader("x-auth-token")

	if tokenHeader == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Missing Token!"})
		return
	}

	//claims := jwt.MapClaims{}
	token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unable to Parse Token!"})
		return
	}

	isValid := token.Valid

	if isValid {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Token!"})
	}

}
