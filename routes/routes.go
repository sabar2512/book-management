package routes

import (
	"book-management/handlers"
	"book-management/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Root endpoint (optional, biar nggak 404 di "/")
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":    "Book Management API is running 🚀",
			"health":     "/health",
			"login":      "/api/login",
			"categories": "/api/categories",
			"books":      "/api/books",
		})
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Public routes
	api := router.Group("/api")
	{
		// Authentication
		api.POST("/login", handlers.Login)
	}

	// Protected routes (require JWT token)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		// Category routes
		categories := protected.Group("/categories")
		{
			categories.GET("", handlers.GetAllCategories)
			categories.POST("", handlers.CreateCategory)
			categories.GET("/:id", handlers.GetCategoryByID)
			categories.PUT("/:id", handlers.UpdateCategory)
			categories.DELETE("/:id", handlers.DeleteCategory)
			categories.GET("/:id/books", handlers.GetBooksByCategory)
		}

		// Book routes
		books := protected.Group("/books")
		{
			books.GET("", handlers.GetAllBooks)
			books.POST("", handlers.CreateBook)
			books.GET("/:id", handlers.GetBookByID)
			books.PUT("/:id", handlers.UpdateBook)
			books.DELETE("/:id", handlers.DeleteBook)
		}
	}
}
