package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

type Config struct {
	Port           string
	BindAddr       string
	DBDriver       string
	SQLitePath     string
	MySQLDSN       string
	AuthSecret     string
	AdminEmail     string
	AdminPassword  string
	DeepSeekAPIKey string
	DeepSeekModel  string
	StaticDir      string
}

type DatabaseStatus struct {
	Configured bool   `json:"configured"`
	Connected  bool   `json:"connected"`
	Driver     string `json:"driver,omitempty"`
	Error      string `json:"error,omitempty"`
}

func loadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("load .env skipped: %v", err)
	}

	port := env("PORT", "")
	if port == "" {
		port = env("SERVER_PORT", "8080")
	}

	driver := strings.ToLower(env("DB_DRIVER", ""))
	if driver == "" {
		// Default to sqlite for effortless local and cloud deployment.
		// Use mysql only if explicitly requested or if valid non-internal MySQL host is provided.
		mysqlHost := env("MYSQL_HOST", "")
		if mysqlHost != "" && mysqlHost != "10.99.2.201" {
			driver = "mysql"
		} else {
			driver = "sqlite"
		}
	}

	mysqlDSN := env("MYSQL_DSN", "")
	if mysqlDSN == "" {
		mysqlDSN = env("DATABASE_URL", "")
	}
	if mysqlDSN == "" && driver == "mysql" {
		mysqlDSN = mysqlDSNFromParts()
	}

	return Config{
		Port:           port,
		BindAddr:       env("BIND_ADDR", "0.0.0.0"),
		DBDriver:       driver,
		SQLitePath:     env("SQLITE_PATH", "./data/internship_tracker.db"),
		MySQLDSN:       mysqlDSN,
		AuthSecret:     env("AUTH_SECRET", "dev-auth-secret-change-me"),
		AdminEmail:     env("ADMIN_EMAIL", "admin@example.com"),
		AdminPassword:  env("ADMIN_PASSWORD", "Admin@123456"),
		DeepSeekAPIKey: env("DEEPSEEK_API_KEY", ""),
		DeepSeekModel:  env("DEEPSEEK_MODEL", "deepseek-chat"),
		StaticDir:      env("STATIC_DIR", ""),
	}
}

func openDatabase(cfg Config) (*sql.DB, DatabaseStatus) {
	if cfg.DBDriver == "sqlite" {
		dir := filepath.Dir(cfg.SQLitePath)
		if dir != "" && dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, DatabaseStatus{Configured: true, Connected: false, Driver: "sqlite", Error: err.Error()}
			}
		}

		db, err := sql.Open("sqlite", cfg.SQLitePath)
		if err != nil {
			return nil, DatabaseStatus{Configured: true, Connected: false, Driver: "sqlite", Error: err.Error()}
		}

		db.SetMaxOpenConns(1) // SQLite works best with serialized writes
		status := pingDatabase(db)
		status.Driver = "sqlite"
		if !status.Connected {
			log.Printf("sqlite ping failed: %s", status.Error)
		} else {
			// Enable WAL mode for better concurrency
			_, _ = db.Exec("PRAGMA journal_mode=WAL;")
			_, _ = db.Exec("PRAGMA foreign_keys=ON;")
		}
		return db, status
	}

	if cfg.MySQLDSN == "" {
		return nil, DatabaseStatus{Configured: false, Connected: false, Driver: "mysql", Error: "MYSQL_DSN is empty"}
	}

	db, err := sql.Open("mysql", cfg.MySQLDSN)
	if err != nil {
		return nil, DatabaseStatus{Configured: true, Connected: false, Driver: "mysql", Error: err.Error()}
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	status := pingDatabase(db)
	status.Driver = "mysql"
	if !status.Connected {
		log.Printf("mysql ping failed, server will keep running: %s", status.Error)
	}

	return db, status
}

func pingDatabase(db *sql.DB) DatabaseStatus {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return DatabaseStatus{Configured: true, Connected: false, Error: err.Error()}
	}

	return DatabaseStatus{Configured: true, Connected: true}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func mysqlDSNFromParts() string {
	cfg := mysql.Config{
		User:                 env("MYSQL_USER", "root"),
		Passwd:               env("MYSQL_PASSWORD", ""),
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", env("MYSQL_HOST", "127.0.0.1"), env("MYSQL_PORT", "3306")),
		DBName:               env("MYSQL_DATABASE", "internship_tracker"),
		Params:               map[string]string{"charset": "utf8mb4"},
		ParseTime:            true,
		Loc:                  time.Local,
		AllowNativePasswords: true,
	}

	return cfg.FormatDSN()
}
