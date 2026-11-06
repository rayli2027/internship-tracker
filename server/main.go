package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
)

type App struct {
	DB     *sql.DB
	Config Config
}

func main() {
	cfg := loadConfig()

	db, dbStatus := openDatabase(cfg)
	if db != nil {
		defer db.Close()

		if dbStatus.Connected {
			if err := migrate(db, cfg.DBDriver); err != nil {
				log.Fatalf("database migration failed: %v", err)
			}
			if err := seedAdmin(db, cfg); err != nil {
				log.Fatalf("admin seed failed: %v", err)
			}
			if err := seedDemo(db); err != nil {
				log.Fatalf("demo seed failed: %v", err)
			}
		}
	}

	app := &App{DB: db, Config: cfg}

	router := gin.Default()
	router.Use(corsMiddleware())

	// API health check
	router.GET("/api/health", app.healthHandler(dbStatus))

	// Auth routes
	auth := router.Group("/api/auth")
	auth.POST("/register", app.registerHandler)
	auth.POST("/login", app.loginHandler)
	auth.GET("/me", app.authMiddleware(), app.meHandler)
	auth.PATCH("/profile", app.authMiddleware(), app.updateProfileHandler)

	// Application routes
	applications := router.Group("/api/applications")
	applications.Use(app.authMiddleware())
	applications.GET("", app.listApplicationsHandler)
	applications.GET("/timeline", app.listAllApplicationEventsHandler)
	applications.POST("", app.createApplicationHandler)
	applications.GET("/:id", app.getApplicationHandler)
	applications.PUT("/:id", app.updateApplicationHandler)
	applications.PATCH("/:id/status", app.changeApplicationStatusHandler)
	applications.PATCH("/:id/liked", app.setApplicationLikedHandler)
	applications.GET("/:id/events", app.listApplicationEventsHandler)
	applications.POST("/:id/events", app.createApplicationEventHandler)
	applications.PUT("/:id/events/:event_id", app.updateApplicationEventHandler)
	applications.DELETE("/:id/events/:event_id", app.deleteApplicationEventHandler)
	applications.DELETE("/:id", app.deleteApplicationHandler)

	// Resume routes
	resume := router.Group("/api/resume")
	resume.Use(app.authMiddleware())
	resume.GET("", app.getResumeHandler)
	resume.PUT("", app.saveResumeHandler)
	resume.POST("/extract", app.extractResumeHandler)
	resume.GET("/analyses", app.listResumeAnalysesHandler)
	resume.POST("/analyze", app.analyzeResumeHandler)

	// Admin routes
	admin := router.Group("/api/admin")
	admin.Use(app.authMiddleware(), app.adminOnlyMiddleware())
	admin.GET("/overview", app.adminOverviewHandler)
	admin.GET("/users", app.adminListUsersHandler)
	admin.PATCH("/users/:id/enabled", app.adminSetUserEnabledHandler)
	admin.GET("/applications", app.adminListApplicationsHandler)

	// Static frontend serving (for single-container or unified cloud deployment)
	staticDir := resolveStaticDir(cfg.StaticDir)
	if staticDir != "" {
		log.Printf("serving static frontend from: %s", staticDir)
		router.Static("/assets", filepath.Join(staticDir, "assets"))

		faviconPath := filepath.Join(staticDir, "favicon.ico")
		if _, err := os.Stat(faviconPath); err == nil {
			router.StaticFile("/favicon.ico", faviconPath)
		}

		router.NoRoute(func(c *gin.Context) {
			if strings.HasPrefix(c.Request.URL.Path, "/api") {
				c.JSON(http.StatusNotFound, gin.H{"error": "api endpoint not found"})
				return
			}
			c.File(filepath.Join(staticDir, "index.html"))
		})
	} else {
		landing := func(c *gin.Context) {
			c.Redirect(http.StatusFound, "http://127.0.0.1:5173/")
		}
		router.GET("/", landing)
		router.HEAD("/", landing)
	}

	listenAddr := fmt.Sprintf("%s:%s", cfg.BindAddr, cfg.Port)
	log.Printf("server listening on http://%s (db_driver=%s)", listenAddr, cfg.DBDriver)
	if err := router.Run(listenAddr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func resolveStaticDir(custom string) string {
	if custom != "" {
		if info, err := os.Stat(custom); err == nil && info.IsDir() {
			return custom
		}
	}

	candidates := []string{
		"./dist",
		"../web/dist",
		"./web/dist",
	}

	for _, c := range candidates {
		indexFile := filepath.Join(c, "index.html")
		if _, err := os.Stat(indexFile); err == nil {
			return c
		}
	}

	return ""
}
