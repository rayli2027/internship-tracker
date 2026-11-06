package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Resume struct {
	UserID     uint64    `json:"user_id,string"`
	Title      string    `json:"title"`
	SourceName string    `json:"source_name"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ResumeAnalysisResult struct {
	Score           int      `json:"score"`
	Summary         string   `json:"summary"`
	Strengths       []string `json:"strengths"`
	Gaps            []string `json:"gaps"`
	MissingKeywords []string `json:"missing_keywords"`
	Improvements    []string `json:"improvements"`
}

type ResumeAnalysis struct {
	ID            uint64               `json:"id,string"`
	ApplicationID uint64               `json:"application_id,string"`
	CompanyName   string               `json:"company_name"`
	PositionName  string               `json:"position_name"`
	Model         string               `json:"model"`
	Result        ResumeAnalysisResult `json:"result"`
	CreatedAt     time.Time            `json:"created_at"`
}

func (app *App) getResumeHandler(c *gin.Context) {
	if !app.requireDB(c) {
		return
	}
	user, _ := currentUser(c)
	var item Resume
	err := app.DB.QueryRow(`SELECT user_id,title,source_name,content,created_at,updated_at FROM resumes WHERE user_id=?`, user.ID).
		Scan(&item.UserID, &item.Title, &item.SourceName, &item.Content, &item.CreatedAt, &item.UpdatedAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, gin.H{"item": nil})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load resume"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"item": item})
}

func (app *App) saveResumeHandler(c *gin.Context) {
	if !app.requireDB(c) {
		return
	}
	user, _ := currentUser(c)
	var req struct {
		Title      string `json:"title"`
		SourceName string `json:"source_name"`
		Content    string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid resume"})
		return
	}
	req.Title, req.SourceName, req.Content = strings.TrimSpace(req.Title), strings.TrimSpace(req.SourceName), strings.TrimSpace(req.Content)
	if req.Title == "" {
		req.Title = "My Resume"
	}
	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "resume content is required"})
		return
	}
	if len(req.Content) > 100000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "resume content is too long"})
		return
	}
	var count int
	err := app.DB.QueryRow(`SELECT 1 FROM resumes WHERE user_id = ?`, user.ID).Scan(&count)
	if err == sql.ErrNoRows {
		_, err = app.DB.Exec(`INSERT INTO resumes (user_id, title, source_name, content, created_at, updated_at) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`, user.ID, req.Title, req.SourceName, req.Content)
	} else if err == nil {
		_, err = app.DB.Exec(`UPDATE resumes SET title = ?, source_name = ?, content = ?, updated_at = CURRENT_TIMESTAMP WHERE user_id = ?`, req.Title, req.SourceName, req.Content, user.ID)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save resume"})
		return
	}
	app.getResumeHandler(c)
}

func (app *App) analyzeResumeHandler(c *gin.Context) {
	if !app.requireDB(c) {
		return
	}
	if app.Config.DeepSeekAPIKey == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "AI analysis is not configured"})
		return
	}
	user, _ := currentUser(c)
	var req struct {
		ApplicationID uint64 `json:"application_id,string" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ApplicationID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "select an application"})
		return
	}
	var resume Resume
	if err := app.DB.QueryRow(`SELECT user_id,title,source_name,content,created_at,updated_at FROM resumes WHERE user_id=?`, user.ID).Scan(&resume.UserID, &resume.Title, &resume.SourceName, &resume.Content, &resume.CreatedAt, &resume.UpdatedAt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "save a resume before running analysis"})
		return
	}
	application, err := app.findApplicationForUser(req.ApplicationID, user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}
	jobText := strings.TrimSpace(strings.Join([]string{application.JDText, application.RequirementsText, application.Responsibilities, application.TechStack}, "\n"))
	if jobText == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "add a job description or requirements to this application first"})
		return
	}
	result, err := app.requestDeepSeekAnalysis(resumePlainText(resume.Content), application, jobText)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	if result.Score < 0 {
		result.Score = 0
	}
	if result.Score > 100 {
		result.Score = 100
	}
	encoded, _ := json.Marshal(result)
	id := newID()
	if _, err = app.DB.Exec(`INSERT INTO resume_analyses(id,user_id,application_id,score,result_json,model) VALUES(?,?,?,?,?,?)`, id, user.ID, application.ID, result.Score, string(encoded), app.Config.DeepSeekModel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save analysis"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"item": ResumeAnalysis{ID: id, ApplicationID: application.ID, CompanyName: application.CompanyName, PositionName: application.PositionName, Model: app.Config.DeepSeekModel, Result: result, CreatedAt: time.Now()}})
}

func (app *App) listResumeAnalysesHandler(c *gin.Context) {
	if !app.requireDB(c) {
		return
	}
	user, _ := currentUser(c)
	rows, err := app.DB.Query(`SELECT r.id,r.application_id,a.company_name,a.position_name,r.model,r.result_json,r.created_at FROM resume_analyses r JOIN applications a ON a.id=r.application_id WHERE r.user_id=? AND a.is_del=0 ORDER BY r.created_at DESC LIMIT 30`, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load analyses"})
		return
	}
	defer rows.Close()
	items := make([]ResumeAnalysis, 0)
	for rows.Next() {
		var item ResumeAnalysis
		var raw string
		if rows.Scan(&item.ID, &item.ApplicationID, &item.CompanyName, &item.PositionName, &item.Model, &raw, &item.CreatedAt) == nil {
			_ = json.Unmarshal([]byte(raw), &item.Result)
			items = append(items, item)
		}
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (app *App) requestDeepSeekAnalysis(resume string, application Application, jobText string) (ResumeAnalysisResult, error) {
	system := `You are a careful internship resume reviewer. Compare only evidence present in the resume with the job information. Return valid JSON only with this exact shape: {"score":0,"summary":"","strengths":[""],"gaps":[""],"missing_keywords":[""],"improvements":[""]}. Score from 0 to 100. Do not invent candidate experience. Make improvements specific and honest.`
	user := fmt.Sprintf("ROLE: %s at %s\nLOCATION: %s\nJOB INFORMATION:\n%s\n\nRESUME:\n%s", application.PositionName, application.CompanyName, application.Location, jobText, resume)
	payload := map[string]any{"model": app.Config.DeepSeekModel, "messages": []map[string]string{{"role": "system", "content": system}, {"role": "user", "content": user}}, "temperature": 0.2, "response_format": map[string]string{"type": "json_object"}}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "https://api.deepseek.com/chat/completions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+app.Config.DeepSeekAPIKey)
	client := &http.Client{Timeout: 60 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		return ResumeAnalysisResult{}, fmt.Errorf("AI service could not be reached")
	}
	defer response.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(response.Body, 2<<20))
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return ResumeAnalysisResult{}, fmt.Errorf("AI service returned an error (%d)", response.StatusCode)
	}
	var envelope struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if json.Unmarshal(raw, &envelope) != nil || len(envelope.Choices) == 0 {
		return ResumeAnalysisResult{}, fmt.Errorf("AI service returned an invalid response")
	}
	content := strings.TrimSpace(envelope.Choices[0].Message.Content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)
	var result ResumeAnalysisResult
	if json.Unmarshal([]byte(content), &result) != nil {
		return ResumeAnalysisResult{}, fmt.Errorf("AI analysis was not valid JSON; please try again")
	}
	return result, nil
}

func resumePlainText(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	var doc map[string]any
	if err := json.Unmarshal([]byte(raw), &doc); err != nil {
		return raw
	}
	version, _ := doc["v"].(float64)
	if version != 2 {
		return raw
	}
	get := func(key string) string {
		value, _ := doc[key].(string)
		return strings.TrimSpace(value)
	}
	lines := make([]string, 0, 16)
	if name := get("full_name"); name != "" {
		lines = append(lines, name)
	}
	contact := make([]string, 0, 4)
	for _, key := range []string{"email", "phone", "location", "links"} {
		if value := get(key); value != "" {
			contact = append(contact, value)
		}
	}
	if len(contact) > 0 {
		lines = append(lines, strings.Join(contact, " · "))
	}
	sections := []struct {
		key   string
		label string
	}{
		{"summary", "Summary"},
		{"education", "Education"},
		{"experience", "Experience"},
		{"projects", "Projects"},
		{"skills", "Skills"},
		{"extras", "Additional"},
	}
	for _, section := range sections {
		value := get(section.key)
		if value == "" {
			continue
		}
		lines = append(lines, "", section.label, value)
	}
	text := strings.TrimSpace(strings.Join(lines, "\n"))
	if text == "" {
		return raw
	}
	return text
}
