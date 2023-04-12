package router

import (
	"book/controller"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.Default()

	router.GET("/books", controller.GetAllBooks)
	router.GET("/books/:id", controller.GetBookById)
	router.POST("/books", controller.CreateNewBook)
	router.PUT("/books/:id", controller.UpdateBookById)
	router.DELETE("/books/:id", controller.DeleteBookById)

	return router
}
