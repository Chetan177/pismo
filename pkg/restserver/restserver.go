package restserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Chetan177/pismo/pkg/db"
	"github.com/Chetan177/pismo/pkg/db/mongodb"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	envServerPort            = "PORT"
	envDbAddress             = "DB_ADDRESS"
	envApiVersion            = "API_VERSION"
	envAccountsCollection    = "ACCOUNTS_COLLECTION"
	envTransactionCollection = "TRANSACTION_COLLECTION"
	envDBName                = "DB_NAME"
)

var server *RestServer

// RestServer struct for the api backend server
type RestServer struct {
	config *Config
	server *echo.Echo
	db     db.DB
}

// NewRestServer method to initiate a new instance of the RestServer
func NewRestServer() *RestServer {
	server = &RestServer{}
	server.config = server.loadConfig()
	server.db = mongodb.NewMongoClient(server.config.dbAddress, server.config.dbName, server.config.accCollection, server.config.transCollection)

	return server
}

// Config struct for configuration of RestServer
type Config struct {
	port            int
	dbAddress       string
	apiVersion      string
	dbName          string
	accCollection   string
	transCollection string
}

// loadConfig method to load config from the env variables
func (r *RestServer) loadConfig() *Config {
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

// Start method to start RestServer
func (r *RestServer) Start() {
	r.server = echo.New()
	apiGroup := r.server.Group("v1")
	r.loadAccountsRoutes(apiGroup)
	r.loadTransactionRoutes(apiGroup)
	r.loadHealthRoutes(apiGroup)
	r.server.Use(middleware.Recover())
	r.server.Pre(middleware.RemoveTrailingSlash())

	routes := r.server.Routes()
	for _, route := range routes {
		log.Printf("%s %s\n", route.Method, route.Path)
	}

	go func(port int) {
		addr := fmt.Sprintf(":%d", port)
		r.server.Logger.Fatal(r.server.Start(addr))
	}(r.config.port)

}

// Stop shutdown both server and db
func (r *RestServer) Stop() {
	r.server.Shutdown(context.Background())
	r.db.ShutDown()
}

// loadHealthRoutes load health api endpoint
func (r *RestServer) loadHealthRoutes(g *echo.Group) {
	healthGroup := g.Group("/health")
	healthGroup.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "service is up and running"})
	})
}
