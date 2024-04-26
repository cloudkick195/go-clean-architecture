package middlewares

import (
	"errors"
	"go_clean_architecture/commons"
	"go_clean_architecture/config"
	"go_clean_architecture/modules/log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
)

type IMiddlewareManager interface {
	ErrorMiddleware() gin.HandlerFunc
	AuthMiddleware() gin.HandlerFunc
	PermissionMiddleware(roleStr string) gin.HandlerFunc
	PerClientRateLimiterForAll() gin.HandlerFunc
	PerClientRateLimiterForSynchronized() gin.HandlerFunc
	PerClientSingleRequest() gin.HandlerFunc
	SingleRequestAddUser(Id uint)
	SingleRequestReleaseUser(Id uint)
	IpHasPermissions() gin.HandlerFunc
}

type Client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type singleRequest map[string]*Client

func (sr singleRequest) Add(key string) {
	sr[key] = &Client{}
}
func (sr singleRequest) Del(key string) {
	delete(sr, key)
}

type UserSingleRequest struct {
	Count int
	Chan  chan bool
}

type middlewareManager struct {
	logUsecase        log.IUsecase
	db                *gorm.DB
	mu                sync.Mutex
	clients           map[string]*Client
	synchronized      map[string]*Client
	singleRequest     singleRequest
	userSingleRequest map[string]*UserSingleRequest
}

func NewMiddlewareManager(logUsecase log.IUsecase, db *gorm.DB) IMiddlewareManager {
	return &middlewareManager{logUsecase: logUsecase, db: db,
		clients:           make(map[string]*Client),
		synchronized:      make(map[string]*Client),
		singleRequest:     make(map[string]*Client),
		userSingleRequest: make(map[string]*UserSingleRequest),
	}
}
func (u *middlewareManager) SingleRequestAddUser(Uid uint) {
	u.mu.Lock()
	Id := strconv.FormatUint(uint64(Uid), 10)
	if _, found := u.userSingleRequest[Id]; !found {
		user := &UserSingleRequest{
			Count: 0,
			Chan:  make(chan bool, 1),
		}
		u.userSingleRequest[Id] = user
	}
	u.userSingleRequest[Id].Count++
	u.mu.Unlock()
	u.userSingleRequest[Id].Chan <- true

}
func (u *middlewareManager) SingleRequestReleaseUser(Uid uint) {
	u.mu.Lock()
	Id := strconv.FormatUint(uint64(Uid), 10)
	if _, found := u.userSingleRequest[Id]; found {
		<-u.userSingleRequest[Id].Chan
		u.userSingleRequest[Id].Count--
		if u.userSingleRequest[Id].Count < 1 {
			close(u.userSingleRequest[Id].Chan)
			delete(u.userSingleRequest, Id)
		}
	}
	u.mu.Unlock()
}

// func (m *middlewareManager) SingleRequestAddUser(Id uint) {
// 	m.singleRequest.Add(strconv.FormatUint(uint64(Id), 10))
// }

// func (m *middlewareManager) SingleRequestReleaseUser(Id uint) {
// 	m.singleRequest.Del(strconv.FormatUint(uint64(Id), 10))
// }

func (m *middlewareManager) IpHasPermissions() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP() // Lấy IP của client từ request
		allowedIPs := strings.Split(config.Env.AUTHORIZED_IPS, ";")
		// Kiểm tra xem IP của client có trong danh sách IP được phép hay không
		allowed := false
		for _, ip := range allowedIPs {
			if clientIP == ip {
				allowed = true
				break
			}
		}

		if !allowed {
			err := errors.New("UnAuthorized")
			panic(commons.NewErrIpPermission(err, "UnAuthorized", "ip intentionally infiltrated", "UnAuthorized"))
		}

		c.Next()
	}
}
