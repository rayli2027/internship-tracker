package main

import (
	"database/sql"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func migrate(db *sql.DB, driver string) error {
	if driver == "sqlite" {
		return migrateSQLite(db)
	}
	return migrateMySQL(db)
}

func migrateSQLite(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			username TEXT NOT NULL DEFAULT '',
			password_hash TEXT NOT NULL DEFAULT '',
			preferred_language TEXT NOT NULL DEFAULT 'zh',
			role TEXT NOT NULL DEFAULT 'user',
			is_enabled INTEGER NOT NULL DEFAULT 1,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			is_del INTEGER NOT NULL DEFAULT 0
		);`,
		`CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);`,
		`CREATE INDEX IF NOT EXISTS idx_users_updated_at ON users(updated_at);`,

		`CREATE TABLE IF NOT EXISTS applications (
			id INTEGER PRIMARY KEY,
			company_name TEXT NOT NULL DEFAULT '',
			position_name TEXT NOT NULL DEFAULT '',
			user_id INTEGER NOT NULL DEFAULT 0,
			user_name TEXT NOT NULL DEFAULT '',
			status TEXT NOT NULL DEFAULT 'Wishlist',
			apply_link TEXT NOT NULL DEFAULT '',
			location TEXT NOT NULL DEFAULT '',
			deadline_at DATETIME NOT NULL DEFAULT '1970-01-01 08:00:00',
			jd_text TEXT NOT NULL DEFAULT '',
			notes TEXT NOT NULL DEFAULT '',
			resume_version TEXT NOT NULL DEFAULT '',
			priority TEXT NOT NULL DEFAULT '',
			salary_min REAL NOT NULL DEFAULT 0,
			salary_max REAL NOT NULL DEFAULT 0,
			salary_unit TEXT NOT NULL DEFAULT '',
			experience_years INTEGER NOT NULL DEFAULT 0,
			education TEXT NOT NULL DEFAULT '',
			work_mode TEXT NOT NULL DEFAULT '',
			days_per_week INTEGER NOT NULL DEFAULT 0,
			intern_months INTEGER NOT NULL DEFAULT 0,
			can_convert INTEGER NOT NULL DEFAULT 0,
			tech_stack TEXT NOT NULL DEFAULT '',
			is_liked INTEGER NOT NULL DEFAULT 0,
			source_type TEXT NOT NULL DEFAULT '',
			source_name TEXT NOT NULL DEFAULT '',
			source_url TEXT NOT NULL DEFAULT '',
			employment_type TEXT NOT NULL DEFAULT '',
			department TEXT NOT NULL DEFAULT '',
			requirements_text TEXT NOT NULL DEFAULT '',
			responsibilities TEXT NOT NULL DEFAULT '',
			contact_name TEXT NOT NULL DEFAULT '',
			contact_email TEXT NOT NULL DEFAULT '',
			applied_at DATETIME NOT NULL DEFAULT '1970-01-01 08:00:00',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			is_del INTEGER NOT NULL DEFAULT 0
		);`,
		`CREATE INDEX IF NOT EXISTS idx_app_user_id ON applications(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_app_status ON applications(status);`,
		`CREATE INDEX IF NOT EXISTS idx_app_company_name ON applications(company_name);`,
		`CREATE INDEX IF NOT EXISTS idx_app_is_liked ON applications(is_liked);`,
		`CREATE INDEX IF NOT EXISTS idx_app_updated_at ON applications(updated_at);`,

		`CREATE TABLE IF NOT EXISTS application_events (
			id INTEGER PRIMARY KEY,
			application_id INTEGER NOT NULL DEFAULT 0,
			user_id INTEGER NOT NULL DEFAULT 0,
			event_type TEXT NOT NULL DEFAULT 'memo',
			from_status TEXT NOT NULL DEFAULT '',
			to_status TEXT NOT NULL DEFAULT '',
			title TEXT NOT NULL DEFAULT '',
			content TEXT NOT NULL DEFAULT '',
			details_json TEXT NOT NULL DEFAULT '{}',
			event_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			is_del INTEGER NOT NULL DEFAULT 0
		);`,
		`CREATE INDEX IF NOT EXISTS idx_event_application ON application_events(application_id);`,
		`CREATE INDEX IF NOT EXISTS idx_event_user ON application_events(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_event_at ON application_events(event_at);`,

		`CREATE TABLE IF NOT EXISTS resumes (
			user_id INTEGER PRIMARY KEY,
			title TEXT NOT NULL DEFAULT 'My Resume',
			source_name TEXT NOT NULL DEFAULT '',
			content TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS resume_analyses (
			id INTEGER PRIMARY KEY,
			user_id INTEGER NOT NULL DEFAULT 0,
			application_id INTEGER NOT NULL DEFAULT 0,
			score INTEGER NOT NULL DEFAULT 0,
			result_json TEXT NOT NULL DEFAULT '',
			model TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE INDEX IF NOT EXISTS idx_analysis_user ON resume_analyses(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_analysis_application ON resume_analyses(application_id);`,
		`CREATE INDEX IF NOT EXISTS idx_analysis_created ON resume_analyses(created_at);`,
	}

	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return err
		}
	}

	return nil
}

func migrateMySQL(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS ` + "`users`" + ` (
    id BIGINT(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '主键id',
    email VARCHAR(300) NOT NULL DEFAULT '' COMMENT '邮箱',
    username VARCHAR(300) NOT NULL DEFAULT '' COMMENT '用户名',
    password_hash VARCHAR(300) NOT NULL DEFAULT '' COMMENT '密码哈希',
    role VARCHAR(64) NOT NULL DEFAULT 'user' COMMENT '角色',
    is_enabled TINYINT(3) UNSIGNED NOT NULL DEFAULT '1' COMMENT '是否启用',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    is_del TINYINT(3) UNSIGNED NOT NULL DEFAULT '0' COMMENT '是否删除',
    PRIMARY KEY ( id ),
    UNIQUE KEY ` + "`uk_email`" + ` ( ` + "`email`" + ` ) COMMENT '邮箱唯一索引',
    KEY ` + "`idx_role`" + ` ( ` + "`role`" + ` ) COMMENT '角色索引',
    KEY ` + "`idx_updated_at`" + ` ( ` + "`updated_at`" + ` ) COMMENT '更新时间索引'
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '用户表';`,
		`CREATE TABLE IF NOT EXISTS ` + "`applications`" + ` (
    id BIGINT(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '主键id',
    company_name VARCHAR(300) NOT NULL DEFAULT '' COMMENT '公司名称',
    position_name VARCHAR(300) NOT NULL DEFAULT '' COMMENT '岗位名称',
    user_id BIGINT(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '所属用户id',
    user_name VARCHAR(300) NOT NULL DEFAULT '' COMMENT '所属用户名',
    status VARCHAR(64) NOT NULL DEFAULT 'Wishlist' COMMENT '投递状态',
    apply_link VARCHAR(1000) NOT NULL DEFAULT '' COMMENT '申请链接',
    location VARCHAR(300) NOT NULL DEFAULT '' COMMENT '地点',
    deadline_at DATETIME NOT NULL DEFAULT '1970-01-01 08:00:00' COMMENT '截止时间',
    jd_text VARCHAR(3000) NOT NULL DEFAULT '' COMMENT 'JD文本',
    notes VARCHAR(3000) NOT NULL DEFAULT '' COMMENT '备注',
    resume_version VARCHAR(300) NOT NULL DEFAULT '' COMMENT '简历版本备注',
    priority VARCHAR(64) NOT NULL DEFAULT '' COMMENT '优先级',
    salary_min DECIMAL(20,2) NOT NULL DEFAULT '0' COMMENT '薪酬最低',
    salary_max DECIMAL(20,2) NOT NULL DEFAULT '0' COMMENT '薪酬最高',
    salary_unit VARCHAR(64) NOT NULL DEFAULT '' COMMENT '薪酬单位',
    experience_years INT(11) NOT NULL DEFAULT '0' COMMENT '要求年限',
    education VARCHAR(64) NOT NULL DEFAULT '' COMMENT '学历要求',
    work_mode VARCHAR(64) NOT NULL DEFAULT '' COMMENT '办公方式',
    days_per_week INT(11) NOT NULL DEFAULT '0' COMMENT '每周到岗天数',
    intern_months INT(11) NOT NULL DEFAULT '0' COMMENT '最短实习月数',
    can_convert TINYINT(3) UNSIGNED NOT NULL DEFAULT '0' COMMENT '可否转正',
    tech_stack VARCHAR(500) NOT NULL DEFAULT '' COMMENT '硬性要求技术栈',
    is_liked TINYINT(3) UNSIGNED NOT NULL DEFAULT '0' COMMENT '是否喜爱',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    is_del TINYINT(3) UNSIGNED NOT NULL DEFAULT '0' COMMENT '是否删除',
    PRIMARY KEY ( id ),
    KEY ` + "`idx_user_id`" + ` ( ` + "`user_id`" + ` ) COMMENT '所属用户id索引',
    KEY ` + "`idx_status`" + ` ( ` + "`status`" + ` ) COMMENT '投递状态索引',
    KEY ` + "`idx_company_name`" + ` ( ` + "`company_name`" + ` ) COMMENT '公司名称索引',
    KEY ` + "`idx_is_liked`" + ` ( ` + "`is_liked`" + ` ) COMMENT '喜爱索引',
    KEY ` + "`idx_updated_at`" + ` ( ` + "`updated_at`" + ` ) COMMENT '更新时间索引'
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '投递岗位表';`,
		`CREATE TABLE IF NOT EXISTS application_events (
    id BIGINT(20) UNSIGNED NOT NULL DEFAULT '0', application_id BIGINT(20) UNSIGNED NOT NULL DEFAULT '0', user_id BIGINT(20) UNSIGNED NOT NULL DEFAULT '0',
    event_type VARCHAR(64) NOT NULL DEFAULT 'memo', from_status VARCHAR(64) NOT NULL DEFAULT '', to_status VARCHAR(64) NOT NULL DEFAULT '',
    title VARCHAR(300) NOT NULL DEFAULT '', content TEXT NOT NULL, details_json TEXT NOT NULL, event_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    is_del TINYINT(3) UNSIGNED NOT NULL DEFAULT '0', PRIMARY KEY (id), KEY idx_event_application (application_id), KEY idx_event_user (user_id), KEY idx_event_at (event_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = 'application lifecycle events';`,
		`CREATE TABLE IF NOT EXISTS resumes (
    user_id BIGINT(20) UNSIGNED NOT NULL, title VARCHAR(300) NOT NULL DEFAULT 'My Resume', source_name VARCHAR(300) NOT NULL DEFAULT '',
    content LONGTEXT NOT NULL, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = 'user resume profiles';`,
		`CREATE TABLE IF NOT EXISTS resume_analyses (
    id BIGINT(20) UNSIGNED NOT NULL, user_id BIGINT(20) UNSIGNED NOT NULL, application_id BIGINT(20) UNSIGNED NOT NULL,
    score INT NOT NULL DEFAULT 0, result_json LONGTEXT NOT NULL, model VARCHAR(100) NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (id), KEY idx_analysis_user (user_id),
    KEY idx_analysis_application (application_id), KEY idx_analysis_created (created_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = 'AI resume compatibility analyses';`,
		`INSERT INTO application_events (id, application_id, user_id, event_type, from_status, to_status, title, content, details_json, event_at, created_at, updated_at)
SELECT a.id, a.id, a.user_id, 'created', '', a.status, '创建岗位', '', '{}', a.created_at, a.created_at, a.created_at
FROM applications a
WHERE a.is_del = 0 AND NOT EXISTS (SELECT 1 FROM application_events e WHERE e.application_id = a.id AND e.user_id = a.user_id AND e.event_type = 'created' AND e.is_del = 0)`,
	}

	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return err
		}
	}

	return ensureApplicationColumns(db)
}

func ensureApplicationColumns(db *sql.DB) error {
	alters := []string{
		`ALTER TABLE users ADD COLUMN preferred_language VARCHAR(8) NOT NULL DEFAULT 'zh' COMMENT 'preferred language'`,
		`ALTER TABLE applications MODIFY COLUMN jd_text TEXT NOT NULL COMMENT 'JD文本'`,
		`ALTER TABLE applications MODIFY COLUMN notes TEXT NOT NULL COMMENT '备注'`,
		`ALTER TABLE applications ADD COLUMN salary_min DECIMAL(20,2) NOT NULL DEFAULT '0' COMMENT '薪酬最低'`,
		`ALTER TABLE applications ADD COLUMN salary_max DECIMAL(20,2) NOT NULL DEFAULT '0' COMMENT '薪酬最高'`,
		`ALTER TABLE applications ADD COLUMN salary_unit VARCHAR(64) NOT NULL DEFAULT '' COMMENT '薪酬单位'`,
		`ALTER TABLE applications ADD COLUMN experience_years INT(11) NOT NULL DEFAULT '0' COMMENT '要求年限'`,
		`ALTER TABLE applications ADD COLUMN education VARCHAR(64) NOT NULL DEFAULT '' COMMENT '学历要求'`,
		`ALTER TABLE applications ADD COLUMN work_mode VARCHAR(64) NOT NULL DEFAULT '' COMMENT '办公方式'`,
		`ALTER TABLE applications ADD COLUMN days_per_week INT(11) NOT NULL DEFAULT '0' COMMENT '每周到岗天数'`,
		`ALTER TABLE applications ADD COLUMN intern_months INT(11) NOT NULL DEFAULT '0' COMMENT '最短实习月数'`,
		`ALTER TABLE applications ADD COLUMN can_convert TINYINT(3) UNSIGNED NOT NULL DEFAULT '0' COMMENT '可否转正'`,
		`ALTER TABLE applications ADD COLUMN tech_stack VARCHAR(500) NOT NULL DEFAULT '' COMMENT '硬性要求技术栈'`,
		`ALTER TABLE applications ADD COLUMN is_liked TINYINT(3) UNSIGNED NOT NULL DEFAULT '0' COMMENT '是否喜爱'`,
		`ALTER TABLE applications ADD KEY idx_is_liked (is_liked) COMMENT '喜爱索引'`,
		`ALTER TABLE applications ADD COLUMN source_type VARCHAR(64) NOT NULL DEFAULT '' COMMENT '岗位来源类型'`,
		`ALTER TABLE applications ADD COLUMN source_name VARCHAR(300) NOT NULL DEFAULT '' COMMENT '岗位来源名称'`,
		`ALTER TABLE applications ADD COLUMN source_url VARCHAR(1000) NOT NULL DEFAULT '' COMMENT '岗位来源链接'`,
		`ALTER TABLE applications ADD COLUMN employment_type VARCHAR(64) NOT NULL DEFAULT '' COMMENT '用工类型'`,
		`ALTER TABLE applications ADD COLUMN department VARCHAR(300) NOT NULL DEFAULT '' COMMENT '部门团队'`,
		`ALTER TABLE applications ADD COLUMN requirements_text TEXT NOT NULL COMMENT '岗位要求'`,
		`ALTER TABLE applications ADD COLUMN responsibilities TEXT NOT NULL COMMENT '岗位职责'`,
		`ALTER TABLE applications ADD COLUMN contact_name VARCHAR(300) NOT NULL DEFAULT '' COMMENT '招聘联系人'`,
		`ALTER TABLE applications ADD COLUMN contact_email VARCHAR(300) NOT NULL DEFAULT '' COMMENT '招聘联系邮箱'`,
		`ALTER TABLE applications ADD COLUMN applied_at DATETIME NOT NULL DEFAULT '1970-01-01 08:00:00' COMMENT '实际投递时间'`,
	}

	for _, statement := range alters {
		if _, err := db.Exec(statement); err != nil && !isIgnorableAlterError(err) {
			return err
		}
	}

	return nil
}

func isIgnorableAlterError(err error) bool {
	message := err.Error()
	return strings.Contains(message, "Duplicate column") || strings.Contains(message, "Duplicate key")
}

func seedAdmin(db *sql.DB, cfg Config) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var existingID uint64
	err = db.QueryRow(`SELECT id FROM users WHERE email = ? AND is_del = 0`, cfg.AdminEmail).Scan(&existingID)
	if err == sql.ErrNoRows {
		_, err = db.Exec(
			`INSERT INTO users (id, email, username, password_hash, preferred_language, role, is_enabled) VALUES (?, ?, 'Administrator', ?, 'en', 'admin', 1)`,
			newID(), cfg.AdminEmail, string(passwordHash),
		)
		return err
	} else if err != nil {
		return err
	}

	_, err = db.Exec(
		`UPDATE users SET password_hash = ?, role = 'admin', is_enabled = 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		string(passwordHash), existingID,
	)
	return err
}

func seedDemo(db *sql.DB) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("Demo@123456"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var userID uint64
	var username string
	err = db.QueryRow(`SELECT id, username FROM users WHERE email = ? AND is_del = 0`, "demo@example.com").Scan(&userID, &username)
	if err == sql.ErrNoRows {
		userID = newID()
		username = "Alex Chen"
		if _, err = db.Exec(
			`INSERT INTO users (id, email, username, password_hash, preferred_language, role, is_enabled) VALUES (?, ?, ?, ?, 'en', 'user', 1)`,
			userID,
			"demo@example.com",
			username,
			string(passwordHash),
		); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	if err := seedDemoWishlistApplications(db, userID, username); err != nil {
		return err
	}

	return nil
}

func seedDemoWishlistApplications(db *sql.DB, userID uint64, username string) error {
	seeds := []struct {
		company      string
		position     string
		status       string
		location     string
		workMode     string
		salaryMin    float64
		salaryMax    float64
		salaryUnit   string
		deadline     string
		notes        string
		techStack    string
		appliedAt    string
		eventTitle   string
	}{
		{
			company:    "Google",
			position:   "Software Engineering Intern",
			status:     "Interview",
			location:   "Sunnyvale, CA",
			workMode:   "Hybrid",
			salaryMin:  9000,
			salaryMax:  12000,
			salaryUnit: "USD/month",
			deadline:   "2026-11-30 23:59:59",
			notes:      "Preparing for technical phone screen focusing on algorithms and system design.",
			techStack:  "Go, C++, Distributed Systems",
			appliedAt:  "2026-09-01 10:00:00",
			eventTitle: "Technical Interview Scheduled",
		},
		{
			company:    "Stripe",
			position:   "Backend Engineering Intern",
			status:     "OA",
			location:   "San Francisco, CA",
			workMode:   "Remote",
			salaryMin:  8500,
			salaryMax:  11000,
			salaryUnit: "USD/month",
			deadline:   "2026-11-15 23:59:59",
			notes:      "Online assessment completed, waiting for reviewer feedback.",
			techStack:  "Ruby, Go, API Design, SQL",
			appliedAt:  "2026-09-10 14:30:00",
			eventTitle: "Online Assessment Submitted",
		},
		{
			company:    "ByteDance",
			position:   "Frontend Development Intern",
			status:     "Applied",
			location:   "San Jose, CA",
			workMode:   "Hybrid",
			salaryMin:  8000,
			salaryMax:  10500,
			salaryUnit: "USD/month",
			deadline:   "2026-10-31 23:59:59",
			notes:      "Applied through employee referral. Team works on creator tools.",
			techStack:  "React, TypeScript, Vite",
			appliedAt:  "2026-09-15 09:00:00",
			eventTitle: "Application Submitted",
		},
		{
			company:    "Jane Street",
			position:   "Quantitative Trading Intern",
			status:     "Wishlist",
			location:   "New York, NY",
			workMode:   "On-site",
			salaryMin:  12000,
			salaryMax:  16000,
			salaryUnit: "USD/month",
			deadline:   "2026-12-15 23:59:59",
			notes:      "Resume tailored for trading team; reviewing probability & mental math questions.",
			techStack:  "OCaml, Python, Statistics",
			appliedAt:  "1970-01-01 08:00:00",
			eventTitle: "Added to Wishlist",
		},
	}

	for _, seed := range seeds {
		var exists int
		if err := db.QueryRow(
			`SELECT 1 FROM applications WHERE user_id = ? AND company_name = ? AND position_name = ? AND is_del = 0 LIMIT 1`,
			userID, seed.company, seed.position,
		).Scan(&exists); err == nil {
			continue
		} else if err != sql.ErrNoRows {
			return err
		}

		appID := newID()
		appliedAt := seed.appliedAt
		if appliedAt == "" {
			appliedAt = "1970-01-01 08:00:00"
		}

		_, err := db.Exec(
			`INSERT INTO applications (
				id, company_name, position_name, user_id, user_name, status, apply_link, location, deadline_at,
				jd_text, notes, resume_version, priority, salary_min, salary_max, salary_unit, experience_years,
				education, work_mode, days_per_week, intern_months, can_convert, tech_stack, is_liked,
				source_type, source_name, source_url, employment_type, department, requirements_text,
				responsibilities, contact_name, contact_email, applied_at
			) VALUES (?, ?, ?, ?, ?, ?, '', ?, ?, '', ?, '', 'High', ?, ?, ?, 0, 'Bachelor', ?, 5, 3, 1, ?, 1, 'Job board', 'Official Site', '', 'Internship', 'Engineering', '', '', '', '', ?)`,
			appID,
			seed.company,
			seed.position,
			userID,
			username,
			seed.status,
			seed.location,
			seed.deadline,
			seed.notes,
			seed.salaryMin,
			seed.salaryMax,
			seed.salaryUnit,
			seed.workMode,
			seed.techStack,
			appliedAt,
		)
		if err != nil {
			return err
		}

		eventID := newID()
		_, err = db.Exec(
			`INSERT INTO application_events (id, application_id, user_id, event_type, from_status, to_status, title, content, details_json, event_at, created_at, updated_at)
			VALUES (?, ?, ?, 'created', '', ?, ?, '', '{}', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
			eventID,
			appID,
			userID,
			seed.status,
			seed.eventTitle,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func newID() uint64 {
	return uint64(time.Now().UnixNano())
}
