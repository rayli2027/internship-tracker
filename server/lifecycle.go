package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ApplicationEvent struct {
	ID            uint64         `json:"id,string"`
	ApplicationID uint64         `json:"application_id,string"`
	EventType     string         `json:"event_type"`
	FromStatus    string         `json:"from_status"`
	ToStatus      string         `json:"to_status"`
	Title         string         `json:"title"`
	Content       string         `json:"content"`
	Details       map[string]any `json:"details"`
	EventAt       time.Time      `json:"event_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	CompanyName   string         `json:"company_name,omitempty"`
	PositionName  string         `json:"position_name,omitempty"`
}

func (app *App) listAllApplicationEventsHandler(c *gin.Context) {
	user, _ := currentUser(c)
	rows, err := app.DB.Query(`SELECT e.id, e.application_id, e.event_type, e.from_status, e.to_status, e.title, e.content, e.details_json, e.event_at, e.created_at, e.updated_at, a.company_name, a.position_name FROM application_events e JOIN applications a ON a.id=e.application_id AND a.user_id=e.user_id WHERE e.user_id=? AND e.is_del=0 AND a.is_del=0 ORDER BY e.event_at DESC, e.created_at DESC`, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	items := make([]ApplicationEvent, 0)
	for rows.Next() {
		var event ApplicationEvent
		var encoded string
		if err := rows.Scan(&event.ID, &event.ApplicationID, &event.EventType, &event.FromStatus, &event.ToStatus, &event.Title, &event.Content, &encoded, &event.EventAt, &event.CreatedAt, &event.UpdatedAt, &event.CompanyName, &event.PositionName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		event.Details = map[string]any{}
		_ = json.Unmarshal([]byte(encoded), &event.Details)
		items = append(items, event)
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

type eventRequest struct {
	EventType string         `json:"event_type"`
	Title     string         `json:"title" binding:"required"`
	Content   string         `json:"content"`
	EventAt   string         `json:"event_at"`
	Details   map[string]any `json:"details"`
}

type statusChangeRequest struct {
	Status  string         `json:"status" binding:"required"`
	Title   string         `json:"title"`
	Content string         `json:"content"`
	EventAt string         `json:"event_at"`
	Details map[string]any `json:"details"`
}

type eventExecer interface {
	Exec(query string, args ...any) (sql.Result, error)
}

func insertApplicationEvent(exec eventExecer, id, applicationID, userID uint64, eventType, fromStatus, toStatus, title, content string, details map[string]any, eventAt time.Time) error {
	encoded, _ := json.Marshal(details)
	_, err := exec.Exec(`INSERT INTO application_events (id, application_id, user_id, event_type, from_status, to_status, title, content, details_json, event_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, id, applicationID, userID, eventType, fromStatus, toStatus, strings.TrimSpace(title), content, string(encoded), eventAt)
	return err
}

func (app *App) listApplicationEventsHandler(c *gin.Context) {
	user, _ := currentUser(c)
	applicationID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	if _, err := app.findApplicationForUser(applicationID, user.ID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}
	rows, err := app.DB.Query(`SELECT id, application_id, event_type, from_status, to_status, title, content, details_json, event_at, created_at, updated_at FROM application_events WHERE application_id = ? AND user_id = ? AND is_del = 0 ORDER BY event_at DESC, created_at DESC`, applicationID, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	items := make([]ApplicationEvent, 0)
	for rows.Next() {
		event, err := scanApplicationEvent(rows)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		items = append(items, event)
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (app *App) createApplicationEventHandler(c *gin.Context) {
	user, _ := currentUser(c)
	applicationID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	if _, err := app.findApplicationForUser(applicationID, user.ID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}
	var req eventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := newID()
	eventAt := parseEventTime(req.EventAt)
	if err := insertApplicationEvent(app.DB, id, applicationID, user.ID, normalizeEventType(req.EventType), "", "", req.Title, req.Content, req.Details, eventAt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	event, _ := app.findApplicationEvent(id, applicationID, user.ID)
	c.JSON(http.StatusCreated, gin.H{"item": event})
}

func (app *App) updateApplicationEventHandler(c *gin.Context) {
	user, _ := currentUser(c)
	applicationID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	eventID, ok := parseIDParam(c, "event_id")
	if !ok {
		return
	}
	var req eventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	encoded, _ := json.Marshal(req.Details)
	result, err := app.DB.Exec(`UPDATE application_events SET event_type=?, title=?, content=?, details_json=?, event_at=? WHERE id=? AND application_id=? AND user_id=? AND is_del=0`, normalizeEventType(req.EventType), strings.TrimSpace(req.Title), req.Content, string(encoded), parseEventTime(req.EventAt), eventID, applicationID, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}
	event, _ := app.findApplicationEvent(eventID, applicationID, user.ID)
	c.JSON(http.StatusOK, gin.H{"item": event})
}

func (app *App) deleteApplicationEventHandler(c *gin.Context) {
	user, _ := currentUser(c)
	applicationID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	eventID, ok := parseIDParam(c, "event_id")
	if !ok {
		return
	}
	result, err := app.DB.Exec(`UPDATE application_events SET is_del=1 WHERE id=? AND application_id=? AND user_id=? AND is_del=0`, eventID, applicationID, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}

func (app *App) changeApplicationStatusHandler(c *gin.Context) {
	user, _ := currentUser(c)
	applicationID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req statusChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	next := normalizeStatus(req.Status)
	tx, err := app.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer tx.Rollback()
	var previous string
	selectStatusQuery := `SELECT status FROM applications WHERE id=? AND user_id=? AND is_del=0`
	if app.Config.DBDriver == "mysql" {
		selectStatusQuery += " FOR UPDATE"
	}
	if err = tx.QueryRow(selectStatusQuery, applicationID, user.ID).Scan(&previous); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}
	if previous != next {
		occurredAt := parseEventTime(req.EventAt)
		if next == "Applied" {
			_, err = tx.Exec(`UPDATE applications SET status=?, applied_at=? WHERE id=? AND user_id=? AND is_del=0`, next, occurredAt, applicationID, user.ID)
		} else {
			_, err = tx.Exec(`UPDATE applications SET status=? WHERE id=? AND user_id=? AND is_del=0`, next, applicationID, user.ID)
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		title := strings.TrimSpace(req.Title)
		if title == "" {
			title = previous + " → " + next
		}
		if err = insertApplicationEvent(tx, newID(), applicationID, user.ID, "status_change", previous, next, title, req.Content, req.Details, occurredAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	item, err := app.findApplicationForUser(applicationID, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load application"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"item": item})
}

func scanApplicationEvent(scanner interface{ Scan(...any) error }) (ApplicationEvent, error) {
	var event ApplicationEvent
	var encoded string
	err := scanner.Scan(&event.ID, &event.ApplicationID, &event.EventType, &event.FromStatus, &event.ToStatus, &event.Title, &event.Content, &encoded, &event.EventAt, &event.CreatedAt, &event.UpdatedAt)
	event.Details = map[string]any{}
	_ = json.Unmarshal([]byte(encoded), &event.Details)
	return event, err
}

func (app *App) findApplicationEvent(id, applicationID, userID uint64) (ApplicationEvent, error) {
	return scanApplicationEvent(app.DB.QueryRow(`SELECT id, application_id, event_type, from_status, to_status, title, content, details_json, event_at, created_at, updated_at FROM application_events WHERE id=? AND application_id=? AND user_id=? AND is_del=0`, id, applicationID, userID))
}

func parseEventTime(value string) time.Time {
	parsed := parseDateTime(value)
	if value == "" || parsed.Year() == 1970 {
		return time.Now()
	}
	return parsed
}
func normalizeEventType(value string) string {
	switch strings.TrimSpace(value) {
	case "applied", "oa", "interview", "offer", "rejected", "withdrawn", "memo":
		return strings.TrimSpace(value)
	default:
		return "memo"
	}
}
