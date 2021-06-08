package users

import (
	"github.com/abhilashdk2016/bookstore_utils_go/rest_errors"
	"strings"
)

const (
	StatusActive = "active"
)

type User struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
	Password    string `json:"password"`
}

type Users []User

func (user *User) Validate() rest_errors.RestErr {
	user.LastName = strings.TrimSpace(strings.ToLower(user.LastName))
	user.FirstName = strings.TrimSpace(strings.ToLower(user.FirstName))
	user.Password = strings.TrimSpace(user.Password)
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == "" {
		return rest_errors.NewBadRequestError("Invalid Email Address")
	}
	if user.Password == "" {
		return rest_errors.NewBadRequestError("Invalid Password")
	}
	return nil
}
