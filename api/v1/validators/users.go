package validators

import (
	"saasmanagement/models"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateUser(user *models.User) error {
	return validate.Struct(user)
}

func ValidateLogin(loginData *models.Login) error {
	return validate.Struct(loginData)
}
