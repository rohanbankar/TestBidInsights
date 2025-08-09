package auth

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	db                 *sql.DB
	jwtSecret          string
	jwtExpiry          time.Duration
	refreshTokenExpiry time.Duration
}

func NewHandler(db *sql.DB, jwtSecret string, jwtExpiry, refreshTokenExpiry time.Duration) *Handler {
	return &Handler{
		db:                 db,
		jwtSecret:          jwtSecret,
		jwtExpiry:          jwtExpiry,
		refreshTokenExpiry: refreshTokenExpiry,
	}
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Get user from database
	user, err := h.getUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	}

	// Generate tokens
	accessToken, err := h.generateAccessToken(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate access token",
		})
		return
	}

	refreshToken, err := h.generateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate refresh token",
		})
		return
	}

	// Set cookies
	c.SetCookie("access_token", accessToken, int(h.jwtExpiry.Seconds()), "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, int(h.refreshTokenExpiry.Seconds()), "/", "", false, true)

	// Remove password from response
	user.Password = ""

	c.JSON(http.StatusOK, LoginResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *Handler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Refresh token required",
		})
		return
	}

	token, err := jwt.ParseWithClaims(refreshToken, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid or expired refresh token",
		})
		return
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid refresh token claims",
		})
		return
	}

	// Get user from database
	user, err := h.getUserByID(claims.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}

	// Generate new access token
	accessToken, err := h.generateAccessToken(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate access token",
		})
		return
	}

	// Set new cookie
	c.SetCookie("access_token", accessToken, int(h.jwtExpiry.Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	// Clear cookies
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

func (h *Handler) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User information not found",
		})
		return
	}

	user, err := h.getUserByID(userID.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Remove password from response
	user.Password = ""

	c.JSON(http.StatusOK, user)
}

func (h *Handler) getUserByUsername(username string) (*User, error) {
	var user User
	query := "SELECT id, username, password_hash, role FROM users WHERE username = ?"
	err := h.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (h *Handler) getUserByID(id int) (*User, error) {
	var user User
	query := "SELECT id, username, password_hash, role FROM users WHERE id = ?"
	err := h.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (h *Handler) generateAccessToken(user User) (string, error) {
	claims := NewClaims(user, h.jwtExpiry)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}

func (h *Handler) generateRefreshToken(userID int) (string, error) {
	claims := NewRefreshTokenClaims(userID, h.refreshTokenExpiry)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}