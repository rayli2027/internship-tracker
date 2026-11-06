package main

import (
	"database/sql"
	"os"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func TestFullApplicationFlowWithSQLite(t *testing.T) {
	dbFile := "test_flow.db"
	defer os.Remove(dbFile)

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	defer db.Close()

	// 1. Migrate
	if err := migrate(db, "sqlite"); err != nil {
		t.Fatalf("sqlite migrate failed: %v", err)
	}

	cfg := Config{
		AdminEmail:    "admin@example.com",
		AdminPassword: "Admin@123456",
		DBDriver:      "sqlite",
	}

	// 2. Seed Admin & Demo
	if err := seedAdmin(db, cfg); err != nil {
		t.Fatalf("seedAdmin failed: %v", err)
	}
	if err := seedDemo(db); err != nil {
		t.Fatalf("seedDemo failed: %v", err)
	}

	// 3. Verify admin exists
	var adminRole string
	err = db.QueryRow(`SELECT role FROM users WHERE email = ?`, "admin@example.com").Scan(&adminRole)
	if err != nil {
		t.Fatalf("admin user not found: %v", err)
	}
	if adminRole != "admin" {
		t.Fatalf("expected admin role, got %s", adminRole)
	}

	// 4. Verify demo user and preloaded applications exist
	var demoID uint64
	err = db.QueryRow(`SELECT id FROM users WHERE email = ?`, "demo@example.com").Scan(&demoID)
	if err != nil {
		t.Fatalf("demo user not found: %v", err)
	}

	var appCount int
	err = db.QueryRow(`SELECT COUNT(*) FROM applications WHERE user_id = ? AND is_del = 0`, demoID).Scan(&appCount)
	if err != nil {
		t.Fatalf("failed to count demo applications: %v", err)
	}
	if appCount == 0 {
		t.Fatalf("expected demo applications to be seeded, got 0")
	}

	// 5. Test overview query with DATE(created_at) = ?
	today := time.Now().Format("2006-01-02")
	var userCount, totalApps, todayUsers, todayApps int64
	err = db.QueryRow(`
SELECT
    (SELECT COUNT(*) FROM users WHERE is_del = 0),
    (SELECT COUNT(*) FROM applications WHERE is_del = 0),
    (SELECT COUNT(*) FROM users WHERE is_del = 0 AND DATE(created_at) = ?),
    (SELECT COUNT(*) FROM applications WHERE is_del = 0 AND DATE(created_at) = ?)`,
		today, today,
	).Scan(&userCount, &totalApps, &todayUsers, &todayApps)
	if err != nil {
		t.Fatalf("overview query failed: %v", err)
	}
	if userCount < 2 || totalApps < 2 {
		t.Fatalf("unexpected counts: users=%d, apps=%d", userCount, totalApps)
	}

	// 6. Test resume save & update
	_, err = db.Exec(`INSERT INTO resumes (user_id, title, source_name, content, created_at, updated_at) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`, demoID, "My Resume", "resume.pdf", "sample resume content")
	if err != nil {
		t.Fatalf("resume insert failed: %v", err)
	}

	var content string
	err = db.QueryRow(`SELECT content FROM resumes WHERE user_id = ?`, demoID).Scan(&content)
	if err != nil || content != "sample resume content" {
		t.Fatalf("failed to query resume: %v", err)
	}
}

