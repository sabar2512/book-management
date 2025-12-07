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

// GetAllCategories retrieves all categories
func GetAllCategories(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT id, name, created_at, created_by, modified_at, modified_by 
		FROM categories 
		ORDER BY id DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch categories",
		})
		return
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.CreatedAt,
			&category.CreatedBy,
			&category.ModifiedAt,
			&category.ModifiedBy,
		)
		if err != nil {
			continue
		}
		categories = append(categories, category)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": categories,
	})
}

// CreateCategory creates a new category
func CreateCategory(c *gin.Context) {
	var input models.CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	username, _ := c.Get("username")
	usernameStr := username.(string)

	var categoryID int
	err := config.DB.QueryRow(`
		INSERT INTO categories (name, created_at, created_by, modified_at, modified_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, input.Name, time.Now(), usernameStr, time.Now(), usernameStr).Scan(&categoryID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create category",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Category created successfully",
		"id":      categoryID,
	})
}

// GetCategoryByID retrieves a category by ID
func GetCategoryByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid category ID",
		})
		return
	}

	var category models.Category
	err = config.DB.QueryRow(`
		SELECT id, name, created_at, created_by, modified_at, modified_by 
		FROM categories 
		WHERE id = $1
	`, id).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.CreatedBy,
		&category.ModifiedAt,
		&category.ModifiedBy,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Category not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": category,
	})
}

// UpdateCategory updates a category by ID
func UpdateCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid category ID",
		})
		return
	}

	var input models.CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	username, _ := c.Get("username")
	usernameStr := username.(string)

	result, err := config.DB.Exec(`
		UPDATE categories 
		SET name = $1, modified_at = $2, modified_by = $3
		WHERE id = $4
	`, input.Name, time.Now(), usernameStr, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update category",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Category not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category updated successfully",
	})
}

// DeleteCategory deletes a category by ID
func DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid category ID",
		})
		return
	}

	result, err := config.DB.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete category",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Category not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category deleted successfully",
	})
}

// GetBooksByCategory retrieves all books in a specific category
func GetBooksByCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid category ID",
		})
		return
	}

	// Check if category exists
	var exists bool
	err = config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)", id).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Category not found",
		})
		return
	}

	rows, err := config.DB.Query(`
		SELECT id, title, description, image_url, release_year, price, 
		       total_page, thickness, category_id, created_at, created_by, 
		       modified_at, modified_by 
		FROM books 
		WHERE category_id = $1
		ORDER BY id DESC
	`, id)
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
