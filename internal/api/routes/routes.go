package routes

import (
	"net/http"

	"github.com/byte4bite/byte4bite/internal/api/handlers"
	"github.com/byte4bite/byte4bite/internal/api/middleware"
	"github.com/byte4bite/byte4bite/internal/auth"
	"github.com/byte4bite/byte4bite/internal/config"
	"github.com/byte4bite/byte4bite/internal/repositories"
	"github.com/byte4bite/byte4bite/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Setup configures all application routes
func Setup(router *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// Apply CORS middleware
	router.Use(middleware.CORSMiddleware())

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	itemRepo := repositories.NewItemRepository(db)
	cartRepo := repositories.NewCartRepository(db)
	orderRepo := repositories.NewOrderRepository(db)

	// Initialize services
	jwtService := auth.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpiryHours)
	authService := services.NewAuthService(userRepo, jwtService)
	categoryService := services.NewCategoryService(categoryRepo)
	itemService := services.NewItemService(itemRepo)
	cartService := services.NewCartService(cartRepo, itemRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	itemHandler := handlers.NewItemHandler(itemService)
	cartHandler := handlers.NewCartHandler(cartService, orderRepo)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Byte4Bite API is running",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		v1.GET("/pantries", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "List pantries - Coming in Phase 6"})
		})

		// Auth routes
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
			authRoutes.POST("/refresh", authHandler.RefreshToken)
			authRoutes.POST("/logout", authHandler.Logout)

			// Protected auth routes
			authRoutes.GET("/me", middleware.AuthMiddleware(jwtService), authHandler.Me)
		}

		// Protected routes (authentication required)
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware(jwtService))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/profile", userHandler.GetProfile)
				users.PUT("/profile", userHandler.UpdateProfile)
				users.PUT("/password", userHandler.UpdatePassword)
			}

			// Items routes - public browsing for authenticated users
			protected.GET("/items", itemHandler.ListItemsPublic)
			protected.GET("/items/:id", itemHandler.GetItem)

			// Cart routes
			carts := protected.Group("/carts")
			{
				carts.GET("/current", cartHandler.GetCurrentCart)
				carts.POST("/items", cartHandler.AddItem)
				carts.PUT("/items/:id", cartHandler.UpdateItemQuantity)
				carts.DELETE("/items/:id", cartHandler.RemoveItem)
				carts.DELETE("/current", cartHandler.ClearCart)
				carts.POST("/checkout", cartHandler.Checkout)
			}

			// Orders routes (coming in Phase 5)
			protected.GET("/orders", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "List orders - Coming in Phase 5"})
			})
		}

		// Admin routes (admin authentication required)
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(jwtService), middleware.AdminMiddleware())
		{
			admin.GET("/dashboard", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Admin dashboard",
					"status":  "active",
				})
			})

			// Category routes
			categories := admin.Group("/categories")
			{
				categories.GET("", categoryHandler.ListCategories)
				categories.POST("", categoryHandler.CreateCategory)
				categories.GET("/:id", categoryHandler.GetCategory)
				categories.PUT("/:id", categoryHandler.UpdateCategory)
				categories.DELETE("/:id", categoryHandler.DeleteCategory)
			}

			// Item routes
			items := admin.Group("/items")
			{
				items.GET("", itemHandler.ListItems)
				items.POST("", itemHandler.CreateItem)
				items.GET("/low-stock", itemHandler.GetLowStockItems)
				items.GET("/:id", itemHandler.GetItem)
				items.PUT("/:id", itemHandler.UpdateItem)
				items.DELETE("/:id", itemHandler.DeleteItem)
				items.PATCH("/:id/quantity", itemHandler.UpdateItemQuantity)
			}

			// Orders routes (coming in Phase 5)
			admin.GET("/orders", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Admin list orders - Coming in Phase 5"})
			})
		}
	}
}
