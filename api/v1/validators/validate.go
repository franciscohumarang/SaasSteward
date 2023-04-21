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

func ValidateAccount(account *models.Account) error {
	return validate.Struct(account)
}

func ValidateBilling(billing *models.Billing) error {
	return validate.Struct(billing)
}

func ValidateCompany(company *models.Company) error {
	return validate.Struct(company)
}
