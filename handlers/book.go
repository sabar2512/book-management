package handlers

import (
	"book-management/config"
	"book-management/models"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetAllBooks retrieves all books
func GetAllBooks(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT id, title, description, image_url, release_year, price, 
		       total_page, thickness, category_id, created_at, created_by, 
		       modified_at, modified_by 
		FROM books 
		ORDER BY id DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch books",
		})
		return
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Description,
			&book.ImageURL,
			&book.ReleaseYear,
			&book.Price,
			&book.TotalPage,
			&book.Thickness,
			&book.CategoryID,
			&book.CreatedAt,
			&book.CreatedBy,
			&book.ModifiedAt,
			&book.ModifiedBy,
		)
		if err != nil {
			continue
		}
		books = append(books, book)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": books,
	})
}

// CreateBook creates a new book
func CreateBook(c *gin.Context) {
	var input models.BookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check if category exists
	var categoryExists bool
	err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)", input.CategoryID).Scan(&categoryExists)
	if err != nil || !categoryExists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid category ID - category does not exist",
		})
		return
	}

	username, _ := c.Get("username")
	usernameStr := username.(string)

	// Calculate thickness
	thickness := input.CalculateThickness()

	var bookID int
	err = config.DB.QueryRow(`
		INSERT INTO books (
			title, description, image_url, release_year, price, 
			total_page, thickness, category_id, 
			created_at, created_by, modified_at, modified_by
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`,
		input.Title,
		input.Description,
		input.ImageURL,
		input.ReleaseYear,
		input.Price,
		input.TotalPage,
		thickness,
		input.CategoryID,
		time.Now(),
		usernameStr,
		time.Now(),
		usernameStr,
	).Scan(&bookID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create book",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Book created successfully",
		"id":        bookID,
		"thickness": thickness,
	})
}

// GetBookByID retrieves a book by ID
func GetBookByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid book ID",
		})
		return
	}

	var book models.Book
	err = config.DB.QueryRow(`
		SELECT id, title, description, image_url, release_year, price, 
		       total_page, thickness, category_id, created_at, created_by, 
		       modified_at, modified_by 
		FROM books 
		WHERE id = $1
	`, id).Scan(
		&book.ID,
		&book.Title,
		&book.Description,
		&book.ImageURL,
		&book.ReleaseYear,
		&book.Price,
		&book.TotalPage,
		&book.Thickness,
		&book.CategoryID,
		&book.CreatedAt,
		&book.CreatedBy,
		&book.ModifiedAt,
		&book.ModifiedBy,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Book not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch book",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": book,
	})
}

// UpdateBook updates a book by ID
func UpdateBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid book ID",
		})
		return
	}

	var input models.BookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check if category exists
	var categoryExists bool
	err = config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)", input.CategoryID).Scan(&categoryExists)
	if err != nil || !categoryExists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid category ID - category does not exist",
		})
		return
	}

	username, _ := c.Get("username")
	usernameStr := username.(string)

	// Calculate thickness
	thickness := input.CalculateThickness()

	result, err := config.DB.Exec(`
		UPDATE books 
		SET title = $1, description = $2, image_url = $3, release_year = $4, 
		    price = $5, total_page = $6, thickness = $7, category_id = $8,
		    modified_at = $9, modified_by = $10
		WHERE id = $11
	`,
		input.Title,
		input.Description,
		input.ImageURL,
		input.ReleaseYear,
		input.Price,
		input.TotalPage,
		thickness,
		input.CategoryID,
		time.Now(),
		usernameStr,
		id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update book",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Book not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Book updated successfully",
		"thickness": thickness,
	})
}

// DeleteBook deletes a book by ID
func DeleteBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid book ID",
		})
		return
	}

	result, err := config.DB.Exec("DELETE FROM books WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete book",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Book not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Book deleted successfully",
	})
}
