package services

import (
	"auth-api-go/models"
)

func CreateUser(username string, password string) (*models.User, error) {
	hash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	userEntry := &models.User{
		Username: username,
		Hash:     hash,
	}

	err = models.DB.Create(userEntry).Error
	if err != nil {
		return nil, err
	}

	return userEntry, nil
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := models.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func DeleteUserByUsername(username string) error {
	var user models.User
	result := models.DB.Where("username = ?", username).Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func AuthenticateUser(username string, password string) (bool, error) {
	user, err := GetUserByUsername(username)
	if err != nil {
		return false, err
	}
	return CheckPasswordHash(password, user.Hash), nil
}
