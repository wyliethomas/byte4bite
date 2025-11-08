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
	pantryRepo := repositories.NewPantryRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	itemRepo := repositories.NewItemRepository(db)
	cartRepo := repositories.NewCartRepository(db)
	orderRepo := repositories.NewOrderRepository(db)
	donationRepo := repositories.NewDonationRepository(db)

	// Initialize services
	jwtService := auth.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpiryHours)
	authService := services.NewAuthService(userRepo, jwtService)
	pantryService := services.NewPantryService(pantryRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	itemService := services.NewItemService(itemRepo)
	cartService := services.NewCartService(cartRepo, itemRepo)
	orderService := services.NewOrderService(orderRepo, itemRepo)
	donationService := services.NewDonationService(donationRepo, pantryRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userRepo)
	pantryHandler := handlers.NewPantryHandler(pantryService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	itemHandler := handlers.NewItemHandler(itemService)
	cartHandler := handlers.NewCartHandler(cartService, orderRepo)
	orderHandler := handlers.NewOrderHandler(orderService)
	donationHandler := handlers.NewDonationHandler(donationService)

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
		// Public pantry routes (no authentication required)
		pantries := v1.Group("/pantries")
		{
			pantries.GET("", pantryHandler.GetPantries)
			pantries.GET("/search", pantryHandler.SearchPantries)
			pantries.GET("/by-city", pantryHandler.GetPantriesByCity)
			pantries.GET("/by-zip", pantryHandler.GetPantriesByZipCode)
			pantries.GET("/:id", pantryHandler.GetPantry)
		}

		// Public donation route (no authentication required)
		v1.POST("/donations", donationHandler.CreateDonation)

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

			// Orders routes
			orders := protected.Group("/orders")
			{
				orders.GET("", orderHandler.GetOrders)
				orders.GET("/:id", orderHandler.GetOrder)
				orders.DELETE("/:id", orderHandler.CancelOrder)
			}
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

			// Admin order management routes
			adminOrders := admin.Group("/orders")
			{
				adminOrders.GET("", orderHandler.GetOrders)
				adminOrders.GET("/:id", orderHandler.GetOrder)
				adminOrders.PUT("/:id/status", orderHandler.UpdateOrderStatus)
				adminOrders.PUT("/:id/assign", orderHandler.AssignStaff)
				adminOrders.DELETE("/:id", orderHandler.CancelOrder)
			}

			// Admin pantry management routes
			adminPantries := admin.Group("/pantries")
			{
				adminPantries.POST("", pantryHandler.CreatePantry)
				adminPantries.PUT("/:id", pantryHandler.UpdatePantry)
				adminPantries.DELETE("/:id", pantryHandler.DeletePantry)
				adminPantries.PATCH("/:id/toggle", pantryHandler.TogglePantryStatus)
			}

			// Admin donation management routes
			adminDonations := admin.Group("/donations")
			{
				adminDonations.GET("", donationHandler.GetDonations)
				adminDonations.GET("/search", donationHandler.SearchDonations)
				adminDonations.GET("/stats", donationHandler.GetDonationStats)
				adminDonations.GET("/by-donor", donationHandler.GetDonationsByDonor)
				adminDonations.GET("/:id", donationHandler.GetDonation)
				adminDonations.PUT("/:id", donationHandler.UpdateDonation)
				adminDonations.DELETE("/:id", donationHandler.DeleteDonation)
				adminDonations.PATCH("/:id/receipt", donationHandler.MarkReceiptSent)
			}
		}
	}
}
