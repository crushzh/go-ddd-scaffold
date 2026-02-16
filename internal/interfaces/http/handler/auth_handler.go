package handler

import (
	"strings"
	"time"

	"go-ddd-scaffold/pkg/config"
	"go-ddd-scaffold/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	jwtCfg *config.JWTConfig
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(jwtCfg *config.JWTConfig) *AuthHandler {
	return &AuthHandler{jwtCfg: jwtCfg}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type tokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

// Pre-computed hash for default password "admin123" (demo only)
var defaultAdminHash, _ = bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

// Login handles user login
// @Summary  User login
// @Tags     Auth
// @Accept   json
// @Produce  json
// @Param    body body loginRequest true "Login credentials"
// @Success  200  {object} response.Response{data=tokenResponse}
// @Router   /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, "invalid parameters")
		return
	}

	// TODO: Query user from database; using default admin for demo
	if req.Username != "admin" || bcrypt.CompareHashAndPassword(defaultAdminHash, []byte(req.Password)) != nil {
		response.Unauthorized(c, "invalid username or password")
		return
	}

	token, expiresAt, err := h.generateToken(req.Username, "admin")
	if err != nil {
		response.ServerError(c, "failed to generate token")
		return
	}

	response.Success(c, tokenResponse{Token: token, ExpiresAt: expiresAt})
}

// RefreshToken refreshes the JWT token
// @Summary  Refresh token
// @Tags     Auth
// @Security Bearer
// @Success  200 {object} response.Response{data=tokenResponse}
// @Router   /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	token, expiresAt, err := h.generateToken(username.(string), role.(string))
	if err != nil {
		response.ServerError(c, "failed to generate token")
		return
	}

	response.Success(c, tokenResponse{Token: token, ExpiresAt: expiresAt})
}

func (h *AuthHandler) generateToken(username, role string) (string, int64, error) {
	expiresAt := time.Now().Add(time.Duration(h.jwtCfg.Expire) * time.Hour)
	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      expiresAt.Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(h.jwtCfg.Secret))
	return tokenStr, expiresAt.Unix(), err
}

// AuthMiddleware validates JWT tokens
func AuthMiddleware(jwtCfg *config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			response.Unauthorized(c, "valid authentication token required")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtCfg.Secret), nil
		})

		if err != nil || !token.Valid {
			response.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Unauthorized(c, "failed to parse token")
			c.Abort()
			return
		}

		c.Set("username", claims["username"])
		c.Set("role", claims["role"])
		c.Next()
	}
}
