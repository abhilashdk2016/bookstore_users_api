package mysql_utils

import (
	"errors"
	"fmt"
	"github.com/abhilashdk2016/bookstore_utils_go/rest_errors"
	"github.com/go-sql-driver/mysql"
	"strings"
)

const (
	ErrorNoRows = "no rows in result set"
)

func ParseError(err error) rest_errors.RestErr {
	sqlError, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), ErrorNoRows) {
			return rest_errors.NewNotFoundError("No Record Found!!!")
		}
		return rest_errors.NewInternalServerError("Error parsing database response", errors.New(""))
	}

	switch  sqlError.Number {
		case 1062: // Duplicated key
			return rest_errors.NewBadRequestError(fmt.Sprintf("Invalid Data"))
	}
	return rest_errors.NewInternalServerError("Error parsing database response", errors.New(""))
}
