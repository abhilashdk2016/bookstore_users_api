package users

import (
	"errors"
	"fmt"
	"github.com/abhilashdk2016/bookstore_users_api/datasources/mysql/users_db"
	"github.com/abhilashdk2016/bookstore_users_api/logger"
	"github.com/abhilashdk2016/bookstore_users_api/utils/date_utils"
	"github.com/abhilashdk2016/bookstore_users_api/utils/mysql_utils"
	"github.com/abhilashdk2016/bookstore_utils_go/rest_errors"
	"strings"
)

const (
	queryInsertUser = "INSERT INTO users(first_name, last_name, email, date_created, password, status) VALUES(?, ?, ?, ?, ?, ?)"
	queryGetUser = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id=?"
	queryUpdateUser = "UPDATE users SET first_name = ?, last_name = ?, email = ? WHERE id = ?"
	queryDeleteUser = "DELETE FROM users where id = ?"
	queryFindUserByStatus = "SELECT id, first_name, last_name, email, date_created, status from users WHERE status = ?"
	queryFindByEmailAndPassword = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE email=? AND password=? AND status = ?"
)

func (user *User) Save() rest_errors.RestErr {
	stmt, err := users_db.DbClient.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("Error when trying to prepare get user statement", err)
		return rest_errors.NewInternalServerError("Error while trying to get user", errors.New("Something wrong with the database. Please try again"))
	}
	defer stmt.Close()
	user.DateCreated = date_utils.GetNowDbFormat()
	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Password, user.Status)

	if saveErr != nil {
		//return mysql_utils.ParseError(saveErr)
		logger.Error("Error when trying to save user", err)
		return rest_errors.NewInternalServerError("Error while trying to save user", errors.New("Something wrong with the database. Please try again"))
	}
	userId, err := insertResult.LastInsertId()
	if err != nil {
		//return errors.NewInternalServerError(fmt.Sprintf("Error while fetching LastInsertId: %s", err.Error()))
		logger.Error("Error while fetching LastInsertId", err)
		return rest_errors.NewInternalServerError("Error while trying to save user", errors.New("Something wrong with the database. Please try again"))
	}
	user.Id = userId
	return nil
}

func (user *User) Get()  rest_errors.RestErr {
	stmt, err := users_db.DbClient.Prepare(queryGetUser)
	if err != nil {
		logger.Error("Error when trying to prepare get user statement", err)
		return rest_errors.NewInternalServerError("Error while trying to get user", errors.New("Something wrong with the database. Please try again"))
	}
	defer stmt.Close()
	result := stmt.QueryRow(user.Id)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		//return mysql_utils.ParseError(getErr)
		logger.Error("Error when trying fetch user", getErr)
		return rest_errors.NewInternalServerError("Error while trying to get user", errors.New("Something wrong with the database. Please try again"))
	}
	return nil
}

func (user *User) Update() rest_errors.RestErr {
	stmt, err := users_db.DbClient.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("Error when trying to prepare get user statement", err)
		return rest_errors.NewInternalServerError("Error while trying to get user", errors.New("Something wrong with the database. Please try again"))
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)

	if err != nil {
		logger.Error("Error when trying to update user", err)
		return rest_errors.NewInternalServerError("Error while trying to update user", errors.New("Something wrong with the database. Please try again"))
		//return mysql_utils.ParseError(err)
	}
	return nil
}

func (user *User) Delete() rest_errors.RestErr {
	stmt, err := users_db.DbClient.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("Error when trying to prepare get user statement", err)
		return rest_errors.NewInternalServerError("Error while trying to get user", errors.New("Something wrong with the database. Please try again"))
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.Id)

	if err != nil {
		//return mysql_utils.ParseError(err)
		logger.Error("Error when trying delete user", err)
		return rest_errors.NewInternalServerError("Error while trying to delete user", errors.New("Something wrong with the database. Please try again"))
	}
	return nil
}

func (user *User) FindByStatus(status string) ([]User, rest_errors.RestErr) {
	stmt, err := users_db.DbClient.Prepare(queryFindUserByStatus)
	if err != nil {
		logger.Error("Error when trying to prepare get user statement", err)
		return nil, rest_errors.NewInternalServerError("Error while trying to get user by status", errors.New("Something wrong with the database. Please try again"))
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("Error when trying fetch user by status", err)
		return nil, rest_errors.NewInternalServerError("Error while trying to get user by status", errors.New("Something wrong with the database. Please try again"))
	}

	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			//return nil, mysql_utils.ParseError(err)
			logger.Error("Error when trying to populate local user slice by status", err)
			return nil, rest_errors.NewInternalServerError("Error while trying to populate user slice by status", errors.New("Something wrong with the database. Please try again"))
		}
		results = append(results, user)
	}

	if len(results) == 0 {
		return nil, rest_errors.NewNotFoundError(fmt.Sprintf("No users matching status %s", status))
	}
	return results, nil
}

func (user *User) FindByEmailAndPassword()  rest_errors.RestErr {
	stmt, err := users_db.DbClient.Prepare(queryFindByEmailAndPassword)
	if err != nil {
		logger.Error("Error when trying to prepare get user by email and password statement", err)
		return rest_errors.NewInternalServerError("Error while trying to prepare get user statememt", errors.New("Something wrong with the database. Please try again"))
	}
	defer stmt.Close()
	result := stmt.QueryRow(user.Email, user.Password, user.Status)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		if strings.Contains(getErr.Error(), mysql_utils.ErrorNoRows) {
			return rest_errors.NewNotFoundError("Invalid Credentials!!!")
		}
		logger.Error("Error when trying fetch user by email and password", getErr)
		return rest_errors.NewInternalServerError("Error while trying to get user by email id and status", errors.New("Something wrong with the database. Please try again"))
	}
	return nil
}
