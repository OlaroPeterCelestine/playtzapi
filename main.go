package main

import (
	"fmt"
	"os"
	"playtz-api/database"
	"playtz-api/handlers"
	"playtz-api/middleware"
	"strings"
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

	// Seed admin user (creates if doesn't exist)
	if err := database.SeedAdmin(); err != nil {
		fmt.Printf("WARNING: Failed to seed admin user: %v. Continuing anyway...\n", err)
		// Don't exit on seed failure - admin might already exist
	}

	// Initialize Gin router
	r := gin.Default()

	// Configure CORS middleware
	config := cors.DefaultConfig()
	
	// Get allowed origins from environment or use defaults
	allowedOrigins := os.Getenv("CORS_ORIGINS")
	if allowedOrigins == "" {
		// Default: allow localhost for development and Vercel production
		allowedOrigins = "http://localhost:3000,http://localhost:3001,http://localhost:5173,http://localhost:8080,https://playtzadmin.vercel.app"
	}
	
	// Parse allowed origins
	origins := []string{}
	for _, origin := range strings.Split(allowedOrigins, ",") {
		origins = append(origins, strings.TrimSpace(origin))
	}
	config.AllowOrigins = origins
	
	// Enable credentials (cookies, authorization headers)
	config.AllowCredentials = true
	
	// Allow methods
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	
	// Allow headers
	config.AllowHeaders = []string{
		"Content-Type",
		"Authorization",
		"Accept",
		"Origin",
		"X-Requested-With",
		"Access-Control-Request-Method",
		"Access-Control-Request-Headers",
	}
	
	// Expose headers
	config.ExposeHeaders = []string{
		"Content-Length",
		"Content-Type",
		"Authorization",
	}
	
	// Max age for preflight requests (24 hours)
	config.MaxAge = 86400
	
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
		// Auth routes (public - no auth required)
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
			auth.POST("/logout", handlers.Logout)
			auth.GET("/me", handlers.GetCurrentUserOptional) // Optional auth - returns 200 with null if not authenticated
			auth.POST("/change-password", middleware.RequireAuth(), handlers.ChangePassword)
		}
	}

	// Protected API routes - All require authentication
	protected := r.Group("/api/v1")
	protected.Use(middleware.RequireAuth())
	{
		// Admin API routes (protected)
		protected.GET("/admin/dashboard", handlers.GetAdminDashboard)

		// News routes - All protected
		protected.GET("/news", handlers.GetNews)
		protected.GET("/news/:id", handlers.GetNewsByID)
		protected.POST("/news", handlers.CreateNews)
		protected.PUT("/news/:id", handlers.UpdateNews)
		protected.DELETE("/news/:id", handlers.DeleteNews)

		// Events routes - All protected
		protected.GET("/events", handlers.GetEvents)
		protected.GET("/events/:id", handlers.GetEventByID)
		protected.POST("/events", handlers.CreateEvent)
		protected.PUT("/events/:id", handlers.UpdateEvent)
		protected.DELETE("/events/:id", handlers.DeleteEvent)

		// Merchandise routes - All protected
		protected.GET("/merch", handlers.GetMerch)
		protected.GET("/merch/:id", handlers.GetMerchByID)
		protected.POST("/merch", handlers.CreateMerch)
		protected.PUT("/merch/:id", handlers.UpdateMerch)
		protected.DELETE("/merch/:id", handlers.DeleteMerch)

		// Careers routes - All protected
		protected.GET("/careers", handlers.GetCareers)
		protected.GET("/careers/:id", handlers.GetCareerByID)
		protected.POST("/careers", handlers.CreateCareer)
		protected.PUT("/careers/:id", handlers.UpdateCareer)
		protected.DELETE("/careers/:id", handlers.DeleteCareer)

		// News routes - All protected
		api.GET("/news", handlers.GetNews)
		api.GET("/news/:id", handlers.GetNewsByID)
		api.POST("/news", handlers.CreateNews)
		api.PUT("/news/:id", handlers.UpdateNews)
		api.DELETE("/news/:id", handlers.DeleteNews)

		// Events routes - All protected
		api.GET("/events", handlers.GetEvents)
		api.GET("/events/:id", handlers.GetEventByID)
		api.POST("/events", handlers.CreateEvent)
		api.PUT("/events/:id", handlers.UpdateEvent)
		api.DELETE("/events/:id", handlers.DeleteEvent)

		// Merchandise routes - All protected
		api.GET("/merch", handlers.GetMerch)
		api.GET("/merch/:id", handlers.GetMerchByID)
		api.POST("/merch", handlers.CreateMerch)
		api.PUT("/merch/:id", handlers.UpdateMerch)
		api.DELETE("/merch/:id", handlers.DeleteMerch)

		// Careers routes - All protected
		api.GET("/careers", handlers.GetCareers)
		api.GET("/careers/:id", handlers.GetCareerByID)
		api.POST("/careers", handlers.CreateCareer)
		api.PUT("/careers/:id", handlers.UpdateCareer)
		api.DELETE("/careers/:id", handlers.DeleteCareer)

		// Shopping Cart routes - All protected
		protected.GET("/cart", handlers.GetCart)
		protected.POST("/cart/add", handlers.AddToCart)
		protected.PUT("/cart/update", handlers.UpdateCartItem)
		protected.DELETE("/cart/remove", handlers.RemoveFromCart)
		protected.DELETE("/cart/clear", handlers.ClearCart)

		// Checkout and Orders routes - All protected
		protected.POST("/checkout", handlers.CreateOrder)
		protected.GET("/orders", handlers.GetOrders)
		protected.GET("/orders/:id", handlers.GetOrder)
		protected.PUT("/orders/:id/status", handlers.UpdateOrderStatus)

		// Rooms routes - All protected
		protected.GET("/rooms", handlers.GetRooms)
		protected.GET("/rooms/:id", handlers.GetRoomByID)
		protected.POST("/rooms", handlers.CreateRoom)
		protected.PUT("/rooms/:id", handlers.UpdateRoom)
		protected.DELETE("/rooms/:id", handlers.DeleteRoom)

		// Mixes routes - All protected
		protected.GET("/mixes", handlers.GetMixes)
		protected.GET("/mixes/:id", handlers.GetMixByID)
		protected.POST("/mixes", handlers.CreateMix)
		protected.PUT("/mixes/:id", handlers.UpdateMix)
		protected.DELETE("/mixes/:id", handlers.DeleteMix)
		protected.POST("/mixes/:id/tracks", handlers.AddTrackToMix)
		protected.POST("/mixes/:id/tracks/bulk", handlers.AddTracksToMix)
		protected.DELETE("/mixes/:id/tracks", handlers.RemoveTrackFromMix)

		// Upload routes - All protected
		protected.POST("/upload", handlers.UploadImage)
		protected.POST("/upload/multiple", handlers.UploadMultipleImages)

		// Users routes - All protected
		protected.GET("/users", handlers.GetUsers)
		protected.POST("/users", handlers.CreateUser)
		protected.GET("/users/:id", handlers.GetUserByID)
		protected.PUT("/users/:id", handlers.UpdateUser)
		protected.DELETE("/users/:id", handlers.DeleteUser)
		protected.PUT("/users/:id/role", handlers.UpdateUserRole)

		// Roles routes - All protected
		protected.GET("/roles", handlers.GetRoles)
		protected.POST("/roles", handlers.CreateRole)
		protected.GET("/roles/:id", handlers.GetRoleByID)
		protected.PUT("/roles/:id", handlers.UpdateRole)
		protected.DELETE("/roles/:id", handlers.DeleteRole)
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
