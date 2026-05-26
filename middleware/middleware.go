package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"sparklink-backend/pkg/auth"
	"sparklink-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}

func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "authorization required")
			c.Abort()
			return
		}

		tokenStr := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = authHeader[7:]
		}

		token, err := jwt.ParseWithClaims(tokenStr, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil {
			response.Unauthorized(c, "invalid token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*auth.Claims)
		if !ok || !token.Valid {
			response.Unauthorized(c, "invalid token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

func RateLimiter() gin.HandlerFunc {
	type client struct {
		count    int
		resetAt time.Time
	}
	clients := make(map[string]*client)
	var mu sync.Mutex

	go func() {
		for {
			time.Sleep(5 * time.Minute)
			mu.Lock()
			now := time.Now()
			for ip, cl := range clients {
				if now.After(cl.resetAt) {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		mu.Lock()
		cl, exists := clients[ip]
		now := time.Now()
		if !exists || now.After(cl.resetAt) {
			clients[ip] = &client{count: 0, resetAt: now.Add(time.Minute)}
			mu.Unlock()
			c.Next()
			return
		}
		cl.count++
		if cl.count >= 60 {
			mu.Unlock()
			response.TooManyRequests(c, "请求过于频繁")
			c.Abort()
			return
		}
		mu.Unlock()
		c.Next()
	}
}
