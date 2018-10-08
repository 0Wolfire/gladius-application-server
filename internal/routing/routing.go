package routing

import (
	"github.com/gladiusio/gladius-application-server/internal/handlers"
	"github.com/gladiusio/gladius-common/pkg/blockchain"
	chandlers "github.com/gladiusio/gladius-common/pkg/handlers"
	"github.com/gladiusio/gladius-common/pkg/routing"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

var apiRouter *mux.Router
var db *gorm.DB

func AppendServerEndpoints(router *mux.Router) error {
	// Initialize Base API sub-route
	routing.InitializeAPISubRoutes(router)
	// Applications
	applicationRouter := apiRouter.PathPrefix("/server").Subrouter()
	applicationRouter.HandleFunc("/info", handlers.PublicPoolInformationHandler(db)).
		Methods(http.MethodGet)

	return nil
}

func AppendStatusEndpoints(router *mux.Router) error {
	// Initialize Base API sub-route
	routing.InitializeAPISubRoutes(router)

	// TxHash Status Sub-Routes
	statusRouter := apiRouter.PathPrefix("/status").Subrouter()
	statusRouter.HandleFunc("/", chandlers.StatusHandler).
		Methods(http.MethodGet, http.MethodPut).
		Name("status")
	statusRouter.HandleFunc("/tx/{tx:0[xX][0-9a-fA-F]{64}}", chandlers.StatusTxHandler).
		Methods(http.MethodGet).
		Name("status-tx")

	return nil
}

func AppendApplicationEndpoints(router *mux.Router, db *gorm.DB) error {
	// Initialize Base API sub-route
	routing.InitializeAPISubRoutes(router)

	// Applications
	applicationRouter := apiRouter.PathPrefix("/applications").Subrouter()
	applicationRouter.HandleFunc("/new", handlers.PoolNewApplicationHandler(db)).
		Methods(http.MethodPost)
	applicationRouter.HandleFunc("/edit", handlers.PoolEditApplicationHandler(db)).
		Methods(http.MethodPost)
	applicationRouter.HandleFunc("/view", handlers.PoolViewApplicationHandler(db)).
		Methods(http.MethodPost)
	applicationRouter.HandleFunc("/status", handlers.PoolStatusViewHandler(db)).
		Methods(http.MethodPost)
	applicationRouter.HandleFunc("/pool/contains/{walletAddress:0[xX][0-9a-fA-F]{40}}", handlers.PoolContainsNode(db))
	applicationRouter.HandleFunc("/nodes", handlers.PoolNodes(db))

	return nil
}

func setupRouter() (*mux.Router, *blockchain.GladiusAccountManager) {
	router, err := routing.InitializeRouter()
	if err != nil {
		println("Failed to initialized router")
	}

	// Create a new GladiusAccountManager for the routes, this is so we can have
	// a shared account system between all endpoints
	ga := blockchain.NewGladiusAccountManager()

	return router, ga
}

func ApplicationServerRouter(db *gorm.DB) *mux.Router {
	router, _ := setupRouter()

	//err := AppendAccountManagementEndpoints(router)
	//if err != nil {
	//	log.Fatalln("Failed to append Account Management Endpoints")
	//}

	err := AppendStatusEndpoints(router)
	if err != nil {
		log.Fatalln("Failed to append Status Endpoints")
	}

	err = AppendServerEndpoints(router)
	if err != nil {
		log.Fatalln("Failed to append Server Endpoints")
	}

	err = AppendApplicationEndpoints(router, db)
	if err != nil {
		log.Fatalln("Failed to append Application Endpoints")
	}

	return router
}
