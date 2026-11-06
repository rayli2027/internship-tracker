package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                uint64    `json:"id,string"`
	Email             string    `json:"email"`
	Username          string    `json:"username"`
	PreferredLanguage string    `json:"preferred_language"`
	Role              string    `json:"role"`
	IsEnabled         bool      `json:"is_enabled"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type Application struct {
	ID               uint64    `json:"id,string"`
	CompanyName      string    `json:"company_name"`
	PositionName     string    `json:"position_name"`
	UserID           uint64    `json:"user_id,string"`
	UserName         string    `json:"user_name"`
	Status           string    `json:"status"`
	ApplyLink        string    `json:"apply_link"`
	Location         string    `json:"location"`
	DeadlineAt       time.Time `json:"deadline_at"`
	JDText           string    `json:"jd_text"`
	Notes            string    `json:"notes"`
	ResumeVersion    string    `json:"resume_version"`
	Priority         string    `json:"priority"`
	SalaryMin        float64   `json:"salary_min"`
	SalaryMax        float64   `json:"salary_max"`
	SalaryUnit       string    `json:"salary_unit"`
	ExperienceYears  int       `json:"experience_years"`
	Education        string    `json:"education"`
	WorkMode         string    `json:"work_mode"`
	DaysPerWeek      int       `json:"days_per_week"`
	InternMonths     int       `json:"intern_months"`
	CanConvert       bool      `json:"can_convert"`
	TechStack        string    `json:"tech_stack"`
	IsLiked          bool      `json:"is_liked"`
	SourceType       string    `json:"source_type"`
	SourceName       string    `json:"source_name"`
	SourceURL        string    `json:"source_url"`
	EmploymentType   string    `json:"employment_type"`
	Department       string    `json:"department"`
	RequirementsText string    `json:"requirements"`
	Responsibilities string    `json:"responsibilities"`
	ContactName      string    `json:"contact_name"`
	ContactEmail     string    `json:"contact_email"`
	AppliedAt        time.Time `json:"applied_at"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type registerRequest struct {
	Email             string `json:"email" binding:"required,email"`
	Username          string `json:"username"`
	Password          string `json:"password" binding:"required,min=6"`
	PreferredLanguage string `json:"preferred_language" binding:"required,oneof=zh en"`
}

type loginRequest struct {
	Email             string `json:"email" binding:"required,email"`
	Password          string `json:"password" binding:"required"`
	PreferredLanguage string `json:"preferred_language" binding:"omitempty,oneof=zh en"`
}

type applicationRequest struct {
	CompanyName      string  `json:"company_name" binding:"required"`
	PositionName     string  `json:"position_name" binding:"required"`
	Status           string  `json:"status"`
	ApplyLink        string  `json:"apply_link"`
	Location         string  `json:"location"`
	DeadlineAt       string  `json:"deadline_at"`
	JDText           string  `json:"jd_text"`
	Notes            string  `json:"notes"`
	ResumeVersion    string  `json:"resume_version"`
	Priority         string  `json:"priority"`
	SalaryMin        float64 `json:"salary_min"`
	SalaryMax        float64 `json:"salary_max"`
	SalaryUnit       string  `json:"salary_unit"`
	ExperienceYears  int     `json:"experience_years"`
	Education        string  `json:"education"`
	WorkMode         string  `json:"work_mode"`
	DaysPerWeek      int     `json:"days_per_week"`
	InternMonths     int     `json:"intern_months"`
	CanConvert       bool    `json:"can_convert"`
	TechStack        string  `json:"tech_stack"`
	IsLiked          bool    `json:"is_liked"`
	SourceType       string  `json:"source_type"`
	SourceName       string  `json:"source_name"`
	SourceURL        string  `json:"source_url"`
	EmploymentType   string  `json:"employment_type"`
	Department       string  `json:"department"`
	RequirementsText string  `json:"requirements"`
	Responsibilities string  `json:"responsibilities"`
	ContactName      string  `json:"contact_name"`
	ContactEmail     string  `json:"contact_email" binding:"omitempty,email"`
	AppliedAt        string  `json:"applied_at"`
}

const applicationSelectCols = `id, company_name, position_name, user_id, user_name, status, apply_link, location, deadline_at, jd_text, notes, resume_version, priority, salary_min, salary_max, salary_unit, experience_years, education, work_mode, days_per_week, intern_months, can_convert, tech_stack, is_liked, source_type, source_name, source_url, employment_type, department, requirements_text, responsibilities, contact_name, contact_email, applied_at, created_at, updated_at`

func (app *App) healthHandler(initial DatabaseStatus) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := initial
		if app.DB != nil {
			status = pingDatabase(app.DB)
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "hello from Internship Tracker API",
			"database": status,
		})
	}
}

func (app *App) registerHandler(c *gin.Context) {
	if !app.requireDB(c) {
		return
	}

	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	username := strings.TrimSpace(req.Username)
	var err error
	if username == "" {
		username = email
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	userID := newID()
	_, err = app.DB.Exec(
		`INSERT INTO users (id, email, username, password_hash, preferred_language, role, is_enabled) VALUES (?, ?, ?, ?, ?, 'user', 1)`,
		userID,
		email,
		username,
		string(hash),
		req.PreferredLanguage,
	)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}

	user, err := app.findUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load user"})
		return
	}

	token, err := app.makeToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token, "user": user})
}

func (app *App) loginHandler(c *gin.Context) {
	if !app.requireDB(c) {
		return
	}

	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, passwordHash, err := app.findUserWithPasswordByEmail(strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}
	if !user.IsEnabled {
		c.JSON(http.StatusForbidden, gin.H{"error": "user is disabled"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}
	if req.PreferredLanguage != "" && req.PreferredLanguage != user.PreferredLanguage {
		if _, err := app.DB.Exec(`UPDATE users SET preferred_language = ? WHERE id = ?`, req.PreferredLanguage, user.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update preferred language"})
			return
		}
		user.PreferredLanguage = req.PreferredLanguage
	}

	token, err := app.makeToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

func (app *App) meHandler(c *gin.Context) {
	user, _ := currentUser(c)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (app *App) updateProfileHandler(c *gin.Context) {
	user, _ := currentUser(c)
	var err error
	var req struct {
		Username          string `json:"username"`
		Password          string `json:"password"`
		PreferredLanguage string `json:"preferred_language" binding:"required,oneof=zh en"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	username := strings.TrimSpace(req.Username)
	if username == "" {
		username = user.Email
	}
	if len(req.Password) > 0 && len(req.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 6 characters"})
		return
	}
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}
		_, err = app.DB.Exec(`UPDATE users SET username = ?, password_hash = ?, preferred_language = ? WHERE id = ?`, username, string(hash), req.PreferredLanguage, user.ID)
	} else {
		_, err = app.DB.Exec(`UPDATE users SET username = ?, preferred_language = ? WHERE id = ?`, username, req.PreferredLanguage, user.ID)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}
	updated, err := app.findUserByID(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load profile"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": updated})
}

func (app *App) listApplicationsHandler(c *gin.Context) {
	if !app.requireDB(c) {
		return
	}
	user, _ := currentUser(c)
	keyword := "%" + strings.TrimSpace(c.Query("keyword")) + "%"
	likedOnly := c.Query("liked") == "1"

	query := `SELECT ` + applicationSelectCols + `
FROM applications
WHERE user_id = ? AND is_del = 0 AND (? = '%%' OR company_name LIKE ? OR position_name LIKE ?)`
	args := []any{user.ID, keyword, keyword, keyword}
	if likedOnly {
		query += ` AND is_liked = 1`
	}
	query += ` ORDER BY created_at DESC`

	rows, err := app.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	applications, err := scanApplications(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": applications})
}

func (app *App) createApplicationHandler(c *gin.Context) {
	if !app.requireDB(c) {
		return
	}
	user, _ := currentUser(c)

	req, ok := bindApplicationRequest(c)
	if !ok {
		return
	}

	id := newID()
	args := append([]any{id}, applicationWriteArgs(req, user)...)
	_, err := app.DB.Exec(
		`INSERT INTO applications (id, company_name, position_name, user_id, user_name, status, apply_link, location, deadline_at, jd_text, notes, resume_version, priority, salary_min, salary_max, salary_unit, experience_years, education, work_mode, days_per_week, intern_months, can_convert, tech_stack, is_liked, source_type, source_name, source_url, employment_type, department, requirements_text, responsibilities, contact_name, contact_email, applied_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		args...,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_ = insertApplicationEvent(app.DB, newID(), id, user.ID, "created", "", normalizeStatus(req.Status), "创建岗位", "", nil, time.Now())

	application, err := app.findApplicationForUser(id, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load application"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"item": application})
}

func (app *App) getApplicationHandler(c *gin.Context) {
	if !app.requireDB(c) {
		return
	}
	user, _ := currentUser(c)
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	application, err := app.findApplicationForUser(id, user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": application})
}

func (app *App) updateApplicationHandler(c *gin.Context) {
	if !app.requireDB(c) {
		return
	}
	user, _ := currentUser(c)
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	req, ok := bindApplicationRequest(c)
	if !ok {
		return
	}

	updateArgs := append(applicationUpdateArgs(req), id, user.ID)
	result, err := app.DB.Exec(
		`UPDATE applications
SET company_name = ?, position_name = ?, status = ?, apply_link = ?, location = ?, deadline_at = ?, jd_text = ?, notes = ?, resume_version = ?, priority = ?, salary_min = ?, salary_max = ?, salary_unit = ?, experience_years = ?, education = ?, work_mode = ?, days_per_week = ?, intern_months = ?, can_convert = ?, tech_stack = ?, is_liked = ?, source_type = ?, source_name = ?, source_url = ?, employment_type = ?, department = ?, requirements_text = ?, responsibilities = ?, contact_name = ?, contact_email = ?, applied_at = ?
WHERE id = ? AND user_id = ? AND is_del = 0`,
		updateArgs...,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}

	application, err := app.findApplicationForUser(id, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": application})
}

func (app *App) deleteApplicationHandler(c *gin.Context) {
	if !app.requireDB(c) {
		return
	}
	user, _ := currentUser(c)
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	result, err := app.DB.Exec(`UPDATE applications SET is_del = 1 WHERE id = ? AND user_id = ? AND is_del = 0`, id, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deleted": true})
}

func (app *App) setApplicationLikedHandler(c *gin.Context) {
	if !app.requireDB(c) {
		return
	}
	user, _ := currentUser(c)
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req struct {
		IsLiked bool `json:"is_liked"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := app.DB.Exec(`UPDATE applications SET is_liked = ? WHERE id = ? AND user_id = ? AND is_del = 0`, boolToTinyint(req.IsLiked), id, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}

	application, err := app.findApplicationForUser(id, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": application})
}

func (app *App) adminOverviewHandler(c *gin.Context) {
	if !app.requireDB(c) {
		return
	}

	var overview struct {
		UserCount             int64 `json:"user_count"`
		ApplicationCount      int64 `json:"application_count"`
		TodayUserCount        int64 `json:"today_user_count"`
		TodayApplicationCount int64 `json:"today_application_count"`
	}

	today := time.Now().Format("2006-01-02")
	err := app.DB.QueryRow(`
SELECT
    (SELECT COUNT(*) FROM users WHERE is_del = 0),
    (SELECT COUNT(*) FROM applications WHERE is_del = 0),
    (SELECT COUNT(*) FROM users WHERE is_del = 0 AND DATE(created_at) = ?),
    (SELECT COUNT(*) FROM applications WHERE is_del = 0 AND DATE(created_at) = ?)`,
		today,
		today,
	).Scan(
		&overview.UserCount,
		&overview.ApplicationCount,
		&overview.TodayUserCount,
		&overview.TodayApplicationCount,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, overview)
}

func (app *App) adminListUsersHandler(c *gin.Context) {
	if !app.requireDB(c) {
		return
	}
	keyword := "%" + strings.TrimSpace(c.Query("keyword")) + "%"
	page, pageSize, offset := parsePagination(c)

	var total int64
	if err := app.DB.QueryRow(
		`SELECT COUNT(*) FROM users WHERE is_del = 0 AND (? = '%%' OR email LIKE ? OR username LIKE ?)`,
		keyword,
		keyword,
		keyword,
	).Scan(&total); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rows, err := app.DB.Query(
		`SELECT id, email, username, preferred_language, role, is_enabled, created_at, updated_at
FROM users
WHERE is_del = 0 AND (? = '%%' OR email LIKE ? OR username LIKE ?)
ORDER BY created_at DESC
LIMIT ? OFFSET ?`,
		keyword,
		keyword,
		keyword,
		pageSize,
		offset,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	users, err := scanUsers(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": users, "total": total, "page": page, "page_size": pageSize})
}

func (app *App) adminListApplicationsHandler(c *gin.Context) {
	if !app.requireDB(c) {
		return
	}

	keyword := "%" + strings.TrimSpace(c.Query("keyword")) + "%"
	userID := strings.TrimSpace(c.Query("user_id"))
	page, pageSize, offset := parsePagination(c)

	where := `is_del = 0 AND (? = '%%' OR company_name LIKE ? OR position_name LIKE ? OR user_name LIKE ?)`
	args := []any{keyword, keyword, keyword, keyword}
	if userID != "" {
		parsed, err := strconv.ParseUint(userID, 10, 64)
		if err != nil || parsed == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
			return
		}
		where += ` AND user_id = ?`
		args = append(args, parsed)
	}

	var total int64
	countArgs := append([]any{}, args...)
	if err := app.DB.QueryRow(`SELECT COUNT(*) FROM applications WHERE `+where, countArgs...).Scan(&total); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	queryArgs := append(args, pageSize, offset)
	rows, err := app.DB.Query(
		`SELECT `+applicationSelectCols+`
FROM applications
WHERE `+where+`
ORDER BY updated_at DESC
LIMIT ? OFFSET ?`,
		queryArgs...,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	applications, err := scanApplications(rows)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": applications, "total": total, "page": page, "page_size": pageSize})
}

func (app *App) adminSetUserEnabledHandler(c *gin.Context) {
	if !app.requireDB(c) {
		return
	}
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req struct {
		IsEnabled bool `json:"is_enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := app.DB.Exec(`UPDATE users SET is_enabled = ? WHERE id = ? AND is_del = 0`, boolToTinyint(req.IsEnabled), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	user, err := app.findUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (app *App) requireDB(c *gin.Context) bool {
	if app.DB == nil {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": "database is not configured"})
		return false
	}

	if status := pingDatabase(app.DB); !status.Connected {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": status.Error})
		return false
	}

	return true
}

func (app *App) findUserByID(id uint64) (User, error) {
	row := app.DB.QueryRow(
		`SELECT id, email, username, preferred_language, role, is_enabled, created_at, updated_at FROM users WHERE id = ? AND is_del = 0`,
		id,
	)
	return scanUser(row)
}

func (app *App) findUserWithPasswordByEmail(email string) (User, string, error) {
	var user User
	var passwordHash string
	var enabled uint8
	err := app.DB.QueryRow(
		`SELECT id, email, username, password_hash, preferred_language, role, is_enabled, created_at, updated_at FROM users WHERE email = ? AND is_del = 0`,
		email,
	).Scan(&user.ID, &user.Email, &user.Username, &passwordHash, &user.PreferredLanguage, &user.Role, &enabled, &user.CreatedAt, &user.UpdatedAt)
	user.IsEnabled = enabled == 1
	return user, passwordHash, err
}

func (app *App) findApplicationForUser(id uint64, userID uint64) (Application, error) {
	row := app.DB.QueryRow(
		`SELECT `+applicationSelectCols+`
FROM applications
WHERE id = ? AND user_id = ? AND is_del = 0`,
		id,
		userID,
	)
	return scanApplication(row)
}

func applicationWriteArgs(req applicationRequest, user User) []any {
	return append([]any{
		strings.TrimSpace(req.CompanyName),
		strings.TrimSpace(req.PositionName),
		user.ID,
		user.Username,
	}, applicationUpdateArgs(req)[2:]...)
}

func applicationUpdateArgs(req applicationRequest) []any {
	return []any{
		strings.TrimSpace(req.CompanyName),
		strings.TrimSpace(req.PositionName),
		normalizeStatus(req.Status),
		strings.TrimSpace(req.ApplyLink),
		strings.TrimSpace(req.Location),
		parseDateTime(req.DeadlineAt),
		req.JDText,
		req.Notes,
		strings.TrimSpace(req.ResumeVersion),
		strings.TrimSpace(req.Priority),
		req.SalaryMin,
		req.SalaryMax,
		strings.TrimSpace(req.SalaryUnit),
		req.ExperienceYears,
		strings.TrimSpace(req.Education),
		strings.TrimSpace(req.WorkMode),
		req.DaysPerWeek,
		req.InternMonths,
		boolToTinyint(req.CanConvert),
		strings.TrimSpace(req.TechStack),
		boolToTinyint(req.IsLiked),
		strings.TrimSpace(req.SourceType),
		strings.TrimSpace(req.SourceName),
		strings.TrimSpace(req.SourceURL),
		strings.TrimSpace(req.EmploymentType),
		strings.TrimSpace(req.Department),
		req.RequirementsText,
		req.Responsibilities,
		strings.TrimSpace(req.ContactName),
		strings.ToLower(strings.TrimSpace(req.ContactEmail)),
		parseDateTime(req.AppliedAt),
	}
}

func bindApplicationRequest(c *gin.Context) (applicationRequest, bool) {
	var req applicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return req, false
	}
	return req, true
}

func scanUser(scanner interface {
	Scan(dest ...any) error
}) (User, error) {
	var user User
	var enabled uint8
	err := scanner.Scan(&user.ID, &user.Email, &user.Username, &user.PreferredLanguage, &user.Role, &enabled, &user.CreatedAt, &user.UpdatedAt)
	user.IsEnabled = enabled == 1
	return user, err
}

func scanUsers(rows *sql.Rows) ([]User, error) {
	users := make([]User, 0)
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func scanApplication(scanner interface {
	Scan(dest ...any) error
}) (Application, error) {
	var application Application
	var canConvert uint8
	var isLiked uint8
	err := scanner.Scan(
		&application.ID,
		&application.CompanyName,
		&application.PositionName,
		&application.UserID,
		&application.UserName,
		&application.Status,
		&application.ApplyLink,
		&application.Location,
		&application.DeadlineAt,
		&application.JDText,
		&application.Notes,
		&application.ResumeVersion,
		&application.Priority,
		&application.SalaryMin,
		&application.SalaryMax,
		&application.SalaryUnit,
		&application.ExperienceYears,
		&application.Education,
		&application.WorkMode,
		&application.DaysPerWeek,
		&application.InternMonths,
		&canConvert,
		&application.TechStack,
		&isLiked,
		&application.SourceType,
		&application.SourceName,
		&application.SourceURL,
		&application.EmploymentType,
		&application.Department,
		&application.RequirementsText,
		&application.Responsibilities,
		&application.ContactName,
		&application.ContactEmail,
		&application.AppliedAt,
		&application.CreatedAt,
		&application.UpdatedAt,
	)
	application.CanConvert = canConvert == 1
	application.IsLiked = isLiked == 1
	return application, err
}

func scanApplications(rows *sql.Rows) ([]Application, error) {
	applications := make([]Application, 0)
	for rows.Next() {
		application, err := scanApplication(rows)
		if err != nil {
			return nil, err
		}
		applications = append(applications, application)
	}
	return applications, rows.Err()
}

func parseDateTime(value string) time.Time {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Date(1970, 1, 1, 8, 0, 0, 0, time.Local)
	}

	for _, layout := range []string{time.RFC3339, "2006-01-02T15:04", "2006-01-02 15:04:05", "2006-01-02 15:04", "2006-01-02"} {
		parsed, err := time.ParseInLocation(layout, value, time.Local)
		if err == nil {
			return parsed
		}
	}

	return time.Date(1970, 1, 1, 8, 0, 0, 0, time.Local)
}

func normalizeStatus(status string) string {
	switch strings.TrimSpace(status) {
	case "Applied", "OA", "Interview", "Offer", "Rejected", "Withdrawn":
		return strings.TrimSpace(status)
	default:
		return "Wishlist"
	}
}

func boolToTinyint(value bool) uint8 {
	if value {
		return 1
	}
	return 0
}

func parsePagination(c *gin.Context) (int, int, int) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if err != nil || pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return page, pageSize, (page - 1) * pageSize
}
