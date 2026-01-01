package main

import (
	"fmt"
	"os"
	"playtz-api/database"
	"playtz-api/handlers"
	"playtz-api/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: .env file not found or couldn't be loaded: %v\n", err)
	}

	// Initialize database connection with retries
	maxRetries := 5
	retryDelay := 5 * time.Second
	var dbErr error

	for i := 0; i < maxRetries; i++ {
		if err := database.InitDB(); err != nil {
			dbErr = err
			if i < maxRetries-1 {
				fmt.Printf("Database connection failed (attempt %d/%d): %v. Retrying in %v...\n", i+1, maxRetries, err, retryDelay)
				time.Sleep(retryDelay)
				continue
			}
		} else {
			dbErr = nil
			break
		}
	}

	if dbErr != nil {
		fmt.Printf("FATAL: Failed to connect to database after %d attempts: %v\n", maxRetries, dbErr)
		os.Exit(1)
	}
	defer database.CloseDB()

	// Run database migrations
	if err := database.Migrate(); err != nil {
		fmt.Printf("WARNING: Database migration failed: %v. Continuing anyway...\n", err)
		// Don't exit on migration failure - tables might already exist
	}

	// Initialize Gin router
	r := gin.Default()

	// Configure CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Serve static files (HTML pages)
	r.Static("/static", "./static")
	r.LoadHTMLGlob("static/*.html")

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	
	// Debug: Check Cloudinary env vars (remove in production)
	r.GET("/debug/cloudinary", func(c *gin.Context) {
		cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
		apiKey := os.Getenv("CLOUDINARY_API_KEY")
		apiSecret := os.Getenv("CLOUDINARY_API_SECRET")
		secretPreview := ""
		if len(apiSecret) > 8 {
			secretPreview = apiSecret[:8] + "..."
		} else {
			secretPreview = apiSecret
		}
		c.JSON(200, gin.H{
			"cloud_name_set": cloudName != "",
			"cloud_name": cloudName,
			"api_key_set": apiKey != "",
			"api_key": apiKey,
			"api_secret_set": apiSecret != "",
			"api_secret_length": len(apiSecret),
			"api_secret_preview": secretPreview,
		})
	})

	// Admin pages (public)
	r.GET("/admin/login", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})

	// Admin page routes (protected)
	admin := r.Group("/admin")
	admin.Use(middleware.RequireAuth())
	{
		admin.GET("/dashboard", func(c *gin.Context) {
			c.HTML(200, "dashboard.html", nil)
		})
	}

	// API v1 routes
	api := r.Group("/api/v1")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
			auth.POST("/logout", handlers.Logout)
			auth.GET("/me", middleware.RequireAuth(), handlers.GetCurrentUser)
			auth.POST("/change-password", middleware.RequireAuth(), handlers.ChangePassword)
		}

		// Admin API routes (protected)
		apiAdmin := api.Group("/admin")
		apiAdmin.Use(middleware.RequireAuth())
		{
			apiAdmin.GET("/dashboard", handlers.GetAdminDashboard)
		}

		// News routes
		api.GET("/news", handlers.GetNews)
		api.POST("/news", handlers.CreateNews)
		api.GET("/news/:id", handlers.GetNewsByID)
		api.PUT("/news/:id", handlers.UpdateNews)
		api.DELETE("/news/:id", handlers.DeleteNews)

		// Events routes
		api.GET("/events", handlers.GetEvents)
		api.POST("/events", handlers.CreateEvent)
		api.GET("/events/:id", handlers.GetEventByID)
		api.PUT("/events/:id", handlers.UpdateEvent)
		api.DELETE("/events/:id", handlers.DeleteEvent)

		// Merchandise routes
		api.GET("/merch", handlers.GetMerch)
		api.POST("/merch", handlers.CreateMerch)
		api.GET("/merch/:id", handlers.GetMerchByID)
		api.PUT("/merch/:id", handlers.UpdateMerch)
		api.DELETE("/merch/:id", handlers.DeleteMerch)

		// Careers routes
		api.GET("/careers", handlers.GetCareers)
		api.POST("/careers", handlers.CreateCareer)
		api.GET("/careers/:id", handlers.GetCareerByID)
		api.PUT("/careers/:id", handlers.UpdateCareer)
		api.DELETE("/careers/:id", handlers.DeleteCareer)

		// Shopping Cart routes
		api.GET("/cart", handlers.GetCart)
		api.POST("/cart/add", handlers.AddToCart)
		api.PUT("/cart/update", handlers.UpdateCartItem)
		api.DELETE("/cart/remove", handlers.RemoveFromCart)
		api.DELETE("/cart/clear", handlers.ClearCart)

		// Checkout and Orders routes
		api.POST("/checkout", handlers.CreateOrder)
		api.GET("/orders", handlers.GetOrders)
		api.GET("/orders/:id", handlers.GetOrder)
		api.PUT("/orders/:id/status", handlers.UpdateOrderStatus)

		// Rooms routes
		api.GET("/rooms", handlers.GetRooms)
		api.POST("/rooms", handlers.CreateRoom)
		api.GET("/rooms/:id", handlers.GetRoomByID)
		api.PUT("/rooms/:id", handlers.UpdateRoom)
		api.DELETE("/rooms/:id", handlers.DeleteRoom)

		// Mixes routes
		api.GET("/mixes", handlers.GetMixes)
		api.POST("/mixes", handlers.CreateMix)
		api.GET("/mixes/:id", handlers.GetMixByID)
		api.PUT("/mixes/:id", handlers.UpdateMix)
		api.DELETE("/mixes/:id", handlers.DeleteMix)
		api.POST("/mixes/:id/tracks", handlers.AddTrackToMix)
		api.POST("/mixes/:id/tracks/bulk", handlers.AddTracksToMix)
		api.DELETE("/mixes/:id/tracks", handlers.RemoveTrackFromMix)

		// Upload routes
		api.POST("/upload", handlers.UploadImage)
		api.POST("/upload/multiple", handlers.UploadMultipleImages)

		// Users routes
		api.GET("/users", handlers.GetUsers)
		api.POST("/users", handlers.CreateUser)
		api.GET("/users/:id", handlers.GetUserByID)
		api.PUT("/users/:id", handlers.UpdateUser)
		api.DELETE("/users/:id", handlers.DeleteUser)
		api.PUT("/users/:id/role", handlers.UpdateUserRole)

		// Roles routes
		api.GET("/roles", handlers.GetRoles)
		api.POST("/roles", handlers.CreateRole)
		api.GET("/roles/:id", handlers.GetRoleByID)
		api.PUT("/roles/:id", handlers.UpdateRole)
		api.DELETE("/roles/:id", handlers.DeleteRole)
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Server starting on port %s\n", port)
	fmt.Printf("✅ Health check available at http://0.0.0.0:%s/health\n", port)

	if err := r.Run(":" + port); err != nil {
		fmt.Printf("FATAL: Server failed to start: %v\n", err)
		os.Exit(1)
	}
}
