package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	envServerPort            = "PORT"
	envDbAddress             = "DB_ADDRESS"
	envApiVersion            = "API_VERSION"
	envAccountsCollection    = "ACCOUNTS_COLLECTION"
	envTransactionCollection = "TRANSACTION_COLLECTION"
	envDBName                = "DB_NAME"
)

// Server struct for the api backend server
type Server struct {
	config                *Config
	server                *echo.Echo
	db                    *mongo.Client
	validator             *validator.Validate
	accountCollection     *mongo.Collection
	transactionCollection *mongo.Collection
}

// NewApiServer method to initiate a new instance of the ApiServer
func NewApiServer() *Server {
	server := &Server{}
	server.config = server.loadConfig()
	server.validator = validator.New()

	clientOptions := options.Client().ApplyURI(server.config.dbAddress)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	server.accountCollection = client.Database(server.config.dbName).Collection(server.config.accCollection)
	server.transactionCollection = client.Database(server.config.dbName).Collection(server.config.transCollection)

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")

	server.db = client

	return server
}

// Config struct for configuration of ApiServer
type Config struct {
	port            int
	dbAddress       string
	apiVersion      string
	dbName          string
	accCollection   string
	transCollection string
}

// loadConfig method to load config from the env variables
func (s *Server) loadConfig() *Config {
	err := godotenv.Load()
	port := os.Getenv(envServerPort)
	if port == "" {
		log.Fatal("unable to configure server port: ", envServerPort)
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal("unable to configure server port: ", envServerPort)
	}

	dbAddress := os.Getenv(envDbAddress)
	if dbAddress == "" {
		log.Fatal("unable to configure server port: ", dbAddress)
	}
	accCollection := os.Getenv(envAccountsCollection)
	if accCollection == "" {
		log.Fatal("unable to configure server port: ", envAccountsCollection)
	}
	transCollection := os.Getenv(envTransactionCollection)
	if transCollection == "" {
		log.Fatal("unable to configure server port: ", envTransactionCollection)
	}
	dbName := os.Getenv(envDBName)
	if dbName == "" {
		log.Fatal("unable to configure server port: ", envDBName)
	}

	apiVersion := os.Getenv(envApiVersion)
	if dbAddress == "" {
		// by default api version is v1
		apiVersion = "v1"
	}
	return &Config{
		port:            p,
		dbAddress:       dbAddress,
		apiVersion:      apiVersion,
		dbName:          dbName,
		accCollection:   accCollection,
		transCollection: transCollection,
	}
}

// Start method to start ApiServer
func (s *Server) Start() {
	s.server = echo.New()
	apiGroup := s.server.Group("v1")
	s.loadAccountsRoutes(apiGroup)
	s.loadTransactionRoutes(apiGroup)
	s.loadHealthRoutes(apiGroup)
	s.server.Use(middleware.Recover())
	s.server.Pre(middleware.RemoveTrailingSlash())

	routes := s.server.Routes()
	for _, route := range routes {
		log.Printf("%s %s\n", route.Method, route.Path)
	}

	go func(port int) {
		addr := fmt.Sprintf(":%d", port)
		s.server.Logger.Fatal(s.server.Start(addr))
	}(s.config.port)

}

// Stop shutdown both server and db
func (s *Server) Stop() {
	s.server.Shutdown(context.Background())
	s.db.Disconnect(context.Background())
}

// loadHealthRoutes load health api endpoint
func (s *Server) loadHealthRoutes(g *echo.Group) {
	healthGroup := g.Group("/health")
	healthGroup.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "service is up and running"})
	})
}
