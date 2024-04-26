package middlewares

import (
	"errors"
	"go_clean_architecture/commons"
	"go_clean_architecture/commons/models"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func (m *middlewareManager) PerClientRateLimiterForAll() gin.HandlerFunc {
	return m.perClientRateLimiter(m.clients, rate.Limit(2), 4, 3*time.Minute, nil)
}

func (m *middlewareManager) PerClientRateLimiterForSynchronized() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqMember, _ := c.Get("member")
		if reqMember == nil {
			panic(commons.ErrInternal(errors.New("member not found")))
		}

		member, ok := reqMember.(*models.Member)
		if !ok {
			panic(commons.ErrInternal(errors.New("member not found")))
		}
		memberID := strconv.Itoa(int(member.ID))

		m.perClientRateLimiter(m.synchronized, rate.Limit(0.2), 1, 5*time.Second, &memberID)(c)
	}
}

func (m *middlewareManager) PerClientSingleRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqMember, _ := c.Get("member")
		if reqMember == nil {
			panic(commons.ErrInternal(errors.New("member not found")))
		}

		member, ok := reqMember.(*models.Member)
		if !ok {
			panic(commons.ErrInternal(errors.New("member not found")))
		}
		memberID := strconv.Itoa(int(member.ID))
		m.mu.Lock()

		if _, found := m.singleRequest[memberID]; found {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"status":  "Request Failed",
				"message": "The API is at capacity, try again later.",
				"key":     "ErrSingleRequest",
			})

			c.Abort()
			m.mu.Unlock()
			return
		}
		m.singleRequest.Add(memberID)
		m.mu.Unlock()
		c.Next()
	}
}

func (m *middlewareManager) perClientRateLimiter(clients map[string]*Client, r rate.Limit, b int, ttl time.Duration, field *string) gin.HandlerFunc {
	go func() {
		for {
			time.Sleep(time.Minute)
			// Lock the mutex to protect this section from race conditions.
			m.mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > ttl {
					delete(clients, ip)
				}
			}
			m.mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		var key string
		if field != nil {
			key = *field
		} else {
			ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "Error",
					"message": "Internal Server Error",
				})
				c.Abort()
				return
			}
			key = ip
		}

		m.mu.Lock()
		if _, found := clients[key]; !found {
			clients[key] = &Client{limiter: rate.NewLimiter(r, b)}
		}
		clients[key].lastSeen = time.Now()
		limiter := clients[key].limiter

		if !limiter.Allow() {
			m.mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"status":  "Request Failed",
				"message": "The API is at capacity, try again later.",
				"key":     "ErrManyRequestsEvery_5Second",
			})
			c.Abort()
			return
		}
		m.mu.Unlock()
		c.Next()
	}
}
