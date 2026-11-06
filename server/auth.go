package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type tokenClaims struct {
	UserID uint64 `json:"user_id"`
	Exp    int64  `json:"exp"`
}

func (app *App) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !app.requireDB(c) {
			return
		}

		header := c.GetHeader("Authorization")
		token := strings.TrimPrefix(header, "Bearer ")
		if token == header || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		claims, err := app.parseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		user, err := app.findUserByID(claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}
		if !user.IsEnabled {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "user is disabled"})
			return
		}

		c.Set("current_user", user)
		c.Next()
	}
}

func (app *App) adminOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := currentUser(c)
		if !ok || user.Role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin only"})
			return
		}

		c.Next()
	}
}

func (app *App) makeToken(userID uint64) (string, error) {
	claims := tokenClaims{
		UserID: userID,
		Exp:    time.Now().Add(24 * time.Hour).Unix(),
	}

	body, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	payload := base64.RawURLEncoding.EncodeToString(body)
	return payload + "." + app.sign(payload), nil
}

func (app *App) parseToken(token string) (tokenClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return tokenClaims{}, errors.New("invalid token format")
	}

	if !hmac.Equal([]byte(app.sign(parts[0])), []byte(parts[1])) {
		return tokenClaims{}, errors.New("invalid token signature")
	}

	body, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return tokenClaims{}, err
	}

	var claims tokenClaims
	if err := json.Unmarshal(body, &claims); err != nil {
		return tokenClaims{}, err
	}

	if claims.UserID == 0 || claims.Exp < time.Now().Unix() {
		return tokenClaims{}, errors.New("token expired")
	}

	return claims, nil
}

func (app *App) sign(payload string) string {
	mac := hmac.New(sha256.New, []byte(app.Config.AuthSecret))
	mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func currentUser(c *gin.Context) (User, bool) {
	value, exists := c.Get("current_user")
	if !exists {
		return User{}, false
	}

	user, ok := value.(User)
	return user, ok
}

func parseIDParam(c *gin.Context, key string) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param(key), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return 0, false
	}

	return id, true
}
