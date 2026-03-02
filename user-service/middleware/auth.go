package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const sessionCookieName = "aiforge_session_id"

func AuthRequired(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(sessionCookieName)
		if err != nil || cookie == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}

		redisKey := fmt.Sprintf("session:%s", cookie)
		val, err := rdb.Get(c.Request.Context(), redisKey).Result()
		if err != nil {
			if err == redis.Nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "session expired"})
				return
			}
			log.Printf("Redis error during session lookup: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		userID, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			log.Printf("Failed to parse user ID from session: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
