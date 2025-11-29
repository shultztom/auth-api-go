package services

import (
	"auth-api-go/models"
)

func GetRolesByUsername(username string) ([]models.Roles, error) {
	var roles []models.Roles
	result := models.DB.Find(&roles, "username = ?", username)
	if result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}

func RoleCheck(roleToCheck string, username string) (bool, error) {
	roles, err := GetRolesByUsername(username)
	if err != nil {
		return false, err
	}

	for _, role := range roles {
		if role.Role == roleToCheck {
			return true, nil
		}
	}

	return false, nil
}

func AddRole(username string, role string) error {
	roleEntry := &models.Roles{
		Username: username,
		Role:     role,
	}

	_, err := models.DB.Create(roleEntry).Rows()
	if err != nil {
		return err
	}

	return nil
}
