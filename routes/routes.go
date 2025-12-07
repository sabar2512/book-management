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
			"message": "Book Management API is running 🚀",
			"endpoints": gin.H{
				"Books": gin.H{
					"GET /api/books":        "Menampilkan seluruh buku",
					"POST /api/books":       "Menambahkan buku baru",
					"GET /api/books/:id":    "Menampilkan detail buku berdasarkan ID",
					"DELETE /api/books/:id": "Menghapus buku berdasarkan ID",
				},
				"Categories": gin.H{
					"GET /api/categories":        "Menampilkan semua kategori",
					"POST /api/categories":       "Menambahkan kategori baru",
					"GET /api/categories/:id":    "Menampilkan detail kategori by ID",
					"PUT /api/categories/:id":    "Update kategori berdasarkan ID",
					"DELETE /api/categories/:id": "Hapus kategori berdasarkan ID",
				},
				"Auth": gin.H{
					"POST /api/login": "Login dan mendapatkan JWT token",
				},
				"Health Check": gin.H{
					"GET /health": "Menampilkan status API",
				},
			},
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
