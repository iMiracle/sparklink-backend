package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"sparklink-backend/pkg/auth"
	"sparklink-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
)

var (
	blacklist    = make(map[string]struct{})
	blacklistMu  sync.RWMutex
	redisClient  *redis.Client
	useRedis     bool
)

func InitRedis(rdb *redis.Client) {
	redisClient = rdb
	useRedis = rdb != nil
}

func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return "req_" + hex.EncodeToString(b)
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader("X-Request-Id")
		if rid == "" {
			rid = generateID()
		}
		c.Set("request_id", rid)
		c.Writer.Header().Set("X-Request-Id", rid)
		c.Next()
	}
}

func BlacklistToken(tokenString string) {
	if useRedis {
		redisClient.SetEX(context.Background(), "bl:"+tokenString, "1", 24*time.Hour)
		return
	}
	blacklistMu.Lock()
	defer blacklistMu.Unlock()
	blacklist[tokenString] = struct{}{}
}

func isTokenBlacklisted(tokenString string) bool {
	if useRedis {
		val, err := redisClient.Exists(context.Background(), "bl:"+tokenString).Result()
		return err == nil && val > 0
	}
	blacklistMu.RLock()
	defer blacklistMu.RUnlock()
	_, ok := blacklist[tokenString]
	return ok
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		rid, _ := c.Get("request_id")

		c.Next()

		status := c.Writer.Status()
		duration := time.Since(start)
		log.Printf("[%s] %s %s %d %v", rid, method, path, status, duration)
	}
}

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

		if isTokenBlacklisted(tokenStr) {
			response.Unauthorized(c, "token已被注销")
			c.Abort()
			return
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

type tokenBucket struct {
	tokens    float64
	lastRefill time.Time
	rate      float64
	burst     int
	mu        sync.Mutex
}

func newTokenBucket(rate float64, burst int) *tokenBucket {
	return &tokenBucket{
		tokens:    float64(burst),
		lastRefill: time.Now(),
		rate:      rate,
		burst:     burst,
	}
}

func (tb *tokenBucket) allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.tokens = min(tb.tokens+elapsed*tb.rate, float64(tb.burst))
	tb.lastRefill = now
	if tb.tokens >= 1 {
		tb.tokens--
		return true
	}
	return false
}

type endpointGroup struct {
	prefix string
	rate   float64
	burst  int
}

var endpointGroups = []endpointGroup{
	{prefix: "/api/v1/auth/sendcode", rate: 2, burst: 5},
	{prefix: "/api/v1/auth/register", rate: 3, burst: 5},
	{prefix: "/api/v1/auth/login", rate: 5, burst: 10},
	{prefix: "/api/v1/auth", rate: 10, burst: 20},
	{prefix: "/api/v1/connect", rate: 15, burst: 30},
	{prefix: "/api/v1", rate: 60, burst: 100},
}

func matchEndpoint(path string) (float64, int) {
	for _, g := range endpointGroups {
		if strings.HasPrefix(path, g.prefix) {
			return g.rate, g.burst
		}
	}
	return 60, 100
}

func RateLimiter() gin.HandlerFunc {
	if useRedis {
		return redisRateLimiter()
	}
	return memoryRateLimiter()
}

func redisRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		path := c.Request.URL.Path
		rate, burst := matchEndpoint(path)
		key := "rl:" + ip
		now := time.Now().Unix()

		pipe := redisClient.Pipeline()
		pipe.LTrim(context.Background(), key, 0, int64(burst)-1)
		count := pipe.LLen(context.Background(), key)
		pipe.Expire(context.Background(), key, time.Minute)
		_, err := pipe.Exec(context.Background())
		if err != nil {
			c.Next()
			return
		}

		if count.Val() >= int64(burst) {
			response.TooManyRequests(c, "请求过于频繁")
			c.Abort()
			return
		}

		redisClient.RPush(context.Background(), key, now)
		refillInterval := int64(1000 / rate)
		oldest := now - refillInterval*int64(burst)
		redisClient.LTrim(context.Background(), key, 0, int64(burst)-1)
		redisClient.LRem(context.Background(), key, 1, oldest)
		c.Next()
	}
}

func memoryRateLimiter() gin.HandlerFunc {
	buckets := make(map[string]*tokenBucket)
	var mu sync.RWMutex

	go func() {
		for {
			time.Sleep(10 * time.Minute)
			mu.Lock()
			for ip, tb := range buckets {
				tb.mu.Lock()
				if time.Since(tb.lastRefill) > 10*time.Minute {
					delete(buckets, ip)
				}
				tb.mu.Unlock()
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		path := c.Request.URL.Path
		rate, burst := matchEndpoint(path)

		mu.RLock()
		tb, exists := buckets[ip]
		mu.RUnlock()
		if !exists {
			tb = newTokenBucket(rate, burst)
			mu.Lock()
			buckets[ip] = tb
			mu.Unlock()
		}

		if !tb.allow() {
			response.TooManyRequests(c, "请求过于频繁")
			c.Abort()
			return
		}
		c.Next()
	}
}
