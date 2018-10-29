package routing

import (
	"github.com/gladiusio/gladius-application-server/internal/handlers"
	chandlers "github.com/gladiusio/gladius-common/pkg/handlers"
	"github.com/gladiusio/gladius-common/pkg/routing"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

func AppendServerEndpoints(router *mux.Router, db *gorm.DB) error {
	// Applications
	applicationRouter := router.PathPrefix("/server").Subrouter().StrictSlash(true)
	applicationRouter.HandleFunc("/info", handlers.PublicPoolInformationHandler(db)).
		Methods(http.MethodGet)

	return nil
}

func AppendStatusEndpoints(router *mux.Router) error {
	// TxHash Status Sub-Routes
	statusRouter := router.PathPrefix("/status").Subrouter().StrictSlash(true)
	statusRouter.HandleFunc("/", chandlers.StatusHandler).
		Methods(http.MethodGet, http.MethodPut).
		Name("status")
	statusRouter.HandleFunc("/tx/{tx:0[xX][0-9a-fA-F]{64}}", chandlers.StatusTxHandler).
		Methods(http.MethodGet).
		Name("status-tx")

	return nil
}

func AppendApplicationEndpoints(router *mux.Router, db *gorm.DB) error {
	// Applications
	applicationRouter := router.PathPrefix("/applications").Subrouter().StrictSlash(true)
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

func ApplicationServerRouter(db *gorm.DB) *mux.Router {
	router, err := routing.InitializeRouter()
	if err != nil {
		println("Failed to initialize router")
	}

	apiRouter := routing.InitializeAPISubRoutes(router)

	err = AppendStatusEndpoints(apiRouter)
	if err != nil {
		log.Fatalln("Failed to append Status Endpoints")
	}

	err = AppendServerEndpoints(apiRouter, db)
	if err != nil {
		log.Fatalln("Failed to append Server Endpoints")
	}

	err = AppendApplicationEndpoints(apiRouter, db)
	if err != nil {
		log.Fatalln("Failed to append Application Endpoints")
	}

	return router
}
