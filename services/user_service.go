package services

import (
	"github.com/abhilashdk2016/bookstore_users_api/domain/users"
	"github.com/abhilashdk2016/bookstore_users_api/utils"
	"github.com/abhilashdk2016/bookstore_utils_go/rest_errors"
)

type usersServiceInterface interface {
	CreateUser(users.User) (*users.User, rest_errors.RestErr)
	GetUser(int64) (*users.User, rest_errors.RestErr)
	UpdateUser(bool,users.User) (*users.User, rest_errors.RestErr)
	DeleteUser(int64) rest_errors.RestErr
	Search(string) (users.Users, rest_errors.RestErr)
	LoginUser(request users.LoginRequest) (*users.User, rest_errors.RestErr)
}

type usersService struct {}

var (
	UsersService usersServiceInterface = &usersService{}
)

func (s *usersService) CreateUser(user users.User) (*users.User, rest_errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	user.Status = users.StatusActive
	user.Password = utils.GetMd5(user.Password)
	if err := user.Save(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *usersService) GetUser(userId int64) (*users.User, rest_errors.RestErr) {
	result := &users.User{Id: userId}
	if err := result.Get(); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *usersService) UpdateUser(isPartial bool, user users.User) (*users.User, rest_errors.RestErr) {
	current, err := s.GetUser(user.Id)
	if err != nil {
		return nil, err
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	if isPartial {
		if user.FirstName != "" {
			current.FirstName = user.FirstName
		}
		if user.LastName != "" {
			current.LastName = user.LastName
		}
		if user.Email != "" {
			current.Email = user.Email
		}
	} else {
		current.FirstName = user.FirstName
		current.LastName = user.LastName
		current.Email = user.Email
	}
	if err := current.Update(); err != nil {
		return nil, err
	}
	return current, nil
}

func (s *usersService) DeleteUser(userId int64) rest_errors.RestErr {
	user := &users.User{Id: userId}
	return user.Delete()
}

func (s *usersService) Search(status string) (users.Users, rest_errors.RestErr) {
	user := &users.User{}
	return user.FindByStatus(status)
}

func (s *usersService) LoginUser(request users.LoginRequest) (*users.User, rest_errors.RestErr) {
	dao := &users.User{
		Email: request.Email,
		Password: request.Password,
		Status: request.Status,
	}
	if err := dao.FindByEmailAndPassword(); err != nil {
		return nil, err
	}

	return dao, nil
}
