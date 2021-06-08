package app

import (
	"github.com/abhilashdk2016/bookstore_users_api/controllers/users"
)

func mapUrls() {
	router.GET("/users/:user_id", users.GetUser)
	router.GET("/internal/users/search", users.Search)
	router.POST("/users", users.CreateUser)
	router.PUT("/users/:user_id", users.UpdateUser)
	router.PATCH("/users/:user_id", users.UpdateUser)
	router.DELETE("/users/:user_id", users.DeleteUser)
	router.POST("/users/login", users.Login)
}
