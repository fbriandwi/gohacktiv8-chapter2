package controller

import (
	"book/database"
	"book/entity"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllBooks(ctx *gin.Context) {
	db := database.GetDB()
	var books []entity.Book

	db.Find(&books)

	ctx.JSON(http.StatusOK, books)
}

func GetBookById(ctx *gin.Context) {
	db := database.GetDB()
	var book entity.Book

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"status":  "Bad Request",
			"message": "The id must be a number",
		})
		return
	}

	errFind := db.First(&book, "id = ?", id).Error
	if errFind != nil {
		if errors.Is(errFind, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"code":    http.StatusNotFound,
				"status":  "Not Found",
				"message": fmt.Sprintf("The book with id = %d is not found in the database", id),
			})
			return
		}
	}

	ctx.JSON(http.StatusFound, book)
}

func CreateNewBook(ctx *gin.Context) {
	db := database.GetDB()
	var newBook entity.Book

	if err := ctx.ShouldBindJSON(&newBook); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"status":  "Bad Request",
			"message": err,
		})
		return
	}

	err := db.Create(&newBook).Error
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"status":  "Internal Server Error",
			"message": "Error creating new book record, there is something error with the server",
		})
		return
	}

	ctx.JSON(http.StatusCreated, newBook)
}

func UpdateBookById(ctx *gin.Context) {
	db := database.GetDB()
	var oldBook entity.Book
	var updatedBook entity.Book

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"status":  "Bad Request",
			"message": "The id must be a number",
		})
		return
	}

	if err := ctx.ShouldBindJSON(&updatedBook); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"status":  "Bad Request",
			"message": err,
		})
		return
	}

	errUpdate := db.Model(&oldBook).Where("id = ?", id).Updates(&updatedBook).Error
	if errUpdate != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"status":  "Internal Server Error",
			"message": "Error updating the book record, there is something error with the server",
		})
		return
	}

	ctx.JSON(http.StatusOK, updatedBook)
}

func DeleteBookById(ctx *gin.Context) {
	db := database.GetDB()
	var book entity.Book

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"status":  "Bad Request",
			"message": "The id must be a number",
		})
		return
	}

	errDelete := db.Where("id = ?", id).Delete(&book).Error
	if errDelete != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"status":  "Not Found",
			"message": fmt.Sprintf("The book with id = %d is not found in the database", id),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Book deleted successfully",
	})
}
