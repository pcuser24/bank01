package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/user2410/simplebank/db/sqlc"
	"github.com/user2410/simplebank/storage"
	"github.com/user2410/simplebank/token"
	"github.com/user2410/simplebank/util"
)

// Server serve HTTP requests for banking services
type Server struct {
	config      util.Config
	store       db.Store
	fileStorage storage.Storage
	tokenMaker  token.Maker
	router      *gin.Engine
}

// NewServer creates a new HTTP server instance and setup routing
func NewServer(config util.Config, store db.Store, fileStorage storage.Storage) (*Server, error) {
	// Setup token maker
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		config:      config,
		tokenMaker:  tokenMaker,
		store:       store,
		fileStorage: fileStorage,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter(config.Environment)
	return server, nil
}

func (server *Server) setupRouter(env string) {
	router := gin.Default()

	// Setup logger
	if strings.ToLower(env) == "production" {
		gin.DisableConsoleColor()
		// make directory /tmp/simplebank if not exists
		err := os.MkdirAll("/tmp/simplebank", os.ModePerm)
		if err != nil {
			log.Fatal("cannot create log directory:", err)
		}
		// create new log file
		fname := fmt.Sprintf("/tmp/simplebank/server-%s.log", time.Now().Format("2006-01-02"))
		f, err := os.OpenFile(fname, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			log.Fatal("cannot create log file:", err)
		}
		// set logger
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	}

	router.GET("/healthcheck", server.healthcheck)

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccount)

	authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

// Start runs the HTTP server on given address
func (server *Server) Start(address string) {
	srv := &http.Server{
		Addr:    address,
		Handler: server.router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		log.Println("Listening on", address)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server forced to shutdown: ", err)
		return
	}

	log.Println("Server exiting")
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
