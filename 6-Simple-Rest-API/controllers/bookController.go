package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
	Desc   string `json:"desc" binding:"required"`
}

var BookDatas = []Book{
	{ID: 1, Title: "Golang", Author: "Gopher", Desc: "A book for go"},
}

func GetBooks(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"books": BookDatas,
	})
}

func GetBook(ctx *gin.Context) {
	bookID, err := strconv.Atoi(ctx.Param("id"))

	var bookData Book

	for i, book := range BookDatas {
		if bookID == book.ID {
			bookData = BookDatas[i]
			break
		}
	}

	// check if book is exist and parameter is valid int
	if bookData == (Book{}) || err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error_status":  "book not found",
			"error_message": fmt.Sprintf("book with id %v not found", bookID),
		})
		return
	}

	ctx.JSON(http.StatusOK, bookData)
}

func CreateBook(ctx *gin.Context) {
	var newBook Book

	// input validation
	if err := ctx.ShouldBindJSON(&newBook); err != nil || strings.TrimSpace(newBook.Title) == "" || strings.TrimSpace(newBook.Author) == "" || strings.TrimSpace(newBook.Desc) == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid inputs, please recheck your inputs",
		})
		return
	}

	//generate book id
	if len(BookDatas) > 0 {
		newBook.ID = BookDatas[len(BookDatas)-1].ID + 1
	} else {
		newBook.ID = 1
	}
	BookDatas = append(BookDatas, newBook) // add new book to book data

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "new book succesfully created",
		"book":    newBook,
	})
}

func UpdateBook(ctx *gin.Context) {
	bookID, err := strconv.Atoi(ctx.Param("id"))
	isExist := false
	var updateBook Book

	// input validation and parameter validation
	if checkBody := ctx.ShouldBindJSON(&updateBook); checkBody != nil || err != nil || strings.TrimSpace(updateBook.Title) == "" || strings.TrimSpace(updateBook.Author) == "" || strings.TrimSpace(updateBook.Desc) == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid inputs, please recheck your inputs",
		})
		return
	}

	updateBook.ID = bookID
	for i, book := range BookDatas {
		if bookID == book.ID {
			isExist = true
			BookDatas[i] = updateBook
			break
		}
	}

	// check if book is exist
	if !isExist {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error_status":  "data not found",
			"error_message": fmt.Sprintf("the book with id %v that you're trying to update is not exist", bookID),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("book with id %v successfully updated", bookID),
	})
}

func DeleteBook(ctx *gin.Context) {
	bookID, err := strconv.Atoi(ctx.Param("id"))
	isExist := false
	var bookIndex int

	for i, book := range BookDatas {
		if bookID == book.ID {
			isExist = true
			bookIndex = i
			break
		}
	}

	// check if book is exist and parameter is valid int
	if !isExist || err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error_status":  "book not found",
			"error_message": fmt.Sprintf("the book you're trying to delete with id %v is not exist", bookID),
		})
		return
	}

	// delete book by cutting the position index
	BookDatas = append(BookDatas[:bookIndex], BookDatas[bookIndex+1:]...)

	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("book with id %v has been succesfully deleted", bookID),
	})
}
