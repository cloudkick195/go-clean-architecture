package server

import (
	"context"
	"go_clean_architecture/config"
	"go_clean_architecture/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	db     *gorm.DB
	router *gin.Engine
}

func NewServer() *Server {
	router := gin.Default()

	// Cấu hình CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"https://appnap247.com", "http://58.84.2.206"}
	config.AddAllowHeaders("Authorization")
	router.Use(cors.New(config))

	return &Server{
		db:     InitDatabase(),
		router: router,
	}
}

func (s *Server) Run() error {
	s.InitRouter()
	utils.NewValidator()
	server := &http.Server{
		Addr:    ":" + config.Env.PORT,
		Handler: s.router,
	}

	go func() {
		log.Printf("Server is running on %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %s\n", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %s\n", err.Error())
	}

	log.Println("Server gracefully stopped")

	return nil
}
