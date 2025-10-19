package main

import (
	_ "api/docs"
	"api/handlers"
	"api/internal/config"
	"api/internal/data/model"
	"api/repository"
	"api/services"
	"log"

	//"os/exec"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	//swaggerFiles "github.com/swaggo/files"
	//ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadConfig()
	cfg.Log()

	// Swagger docs only in dev
	/*if cfg.IsDev() {
		log.Println("Regenerating Swagger docs (development mode)...")
		cmd := exec.Command("swag", "init", "--generalInfo", "handlers/analytics_handler.go", "--output", "docs")
		if err := cmd.Run(); err != nil {
			log.Fatalf("failed to generate Swagger docs: %v", err)
		}
	}*/

	// Database connection
	db, err := gorm.Open(postgres.Open(cfg.PostgresDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&model.PortfolioVisit{}); err != nil {
		log.Fatalf("failed to auto-migrate tables: %v", err)
	}

	// Services & Handlers
	analyticsRepo := repository.NewAnalyticsRepository(db)
	analyticsService := services.NewAnalyticsService(analyticsRepo)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)

	// Gin setup
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{cfg.FrontendURL},
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
		MaxAge:       12 * time.Hour,
	}))

	// Routes
	r.POST("/visit", analyticsHandler.AddVisit)
	r.POST("/visit/github", analyticsHandler.AddGithubVisit)
	r.POST("/visit/linkedin", analyticsHandler.AddLinkedinVisit)
	r.POST("/visit/resume", analyticsHandler.AddResumeDownload)
	r.GET("/health", func(c *gin.Context) { c.String(200, "ok") })
	/*if cfg.IsDev() {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}*/

	// Async analytics worker
	go handlers.StartAnalyticsQueueWorker()

	// Start server
	log.Printf("Server starting on port %s in %s mode", cfg.Port, cfg.AppEnv)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}