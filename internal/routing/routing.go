package routing

import (
	"github.com/gladiusio/gladius-application-server/internal/handlers"
	"github.com/gladiusio/gladius-common/pkg/blockchain"
	chandlers "github.com/gladiusio/gladius-common/pkg/handlers"
	"github.com/gladiusio/gladius-common/pkg/routing"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	"net/http"
)

type ApplicationServer struct {
	ga     	*blockchain.GladiusAccountManager
	router 	*mux.Router
	port   	string
	db 		*gorm.DB
}

func New(port string, db *gorm.DB) *ApplicationServer {
	return &ApplicationServer{
		ga:     blockchain.NewGladiusAccountManager(),
		router: mux.NewRouter(),
		port:   port,
		db: 	db,
	}
}

func (as *ApplicationServer) Start() {
	as.addMiddleware()
	as.addRoutes()

	as.router.StrictSlash(true)
	// Listen locally and setup CORS
	go func() {
		err := http.ListenAndServe(":"+as.port, as.router)
		if err != nil {
			log.Fatal().Err(err).Msg("Error starting API")
		}
	}()

	log.Info().Msg("Started API at http://localhost:" + as.port)
}

func (as *ApplicationServer) addMiddleware() {
	// as.addLogging(as.router)
	// as.router.Use(responseMiddleware)
}

func (as *ApplicationServer) addRoutes() {
	routing.AppendVersionEndpoints(as.router, "0.7.2")
	
	baseRouter := as.router.PathPrefix("/api").Subrouter().StrictSlash(true)
	baseRouter.NotFoundHandler = http.HandlerFunc(chandlers.NotFoundHandler)

	// TxHash Status Sub-Routes
	statusRouter := baseRouter.PathPrefix("/status").Subrouter().StrictSlash(true)
	statusRouter.HandleFunc("/", chandlers.StatusHandler).
		Methods(http.MethodGet, http.MethodPut).
		Name("status")
	statusRouter.HandleFunc("/tx/{tx:0[xX][0-9a-fA-F]{64}}", chandlers.StatusTxHandler).
		Methods(http.MethodGet).
		Name("status-tx")

	// Server
	serverRouter := baseRouter.PathPrefix("/server").Subrouter().StrictSlash(true)
	serverRouter.HandleFunc("/info", handlers.PublicPoolInformationHandler(as.db)).
		Methods(http.MethodGet)

	// Applications
	applicationRouter := baseRouter.PathPrefix("/applications").Subrouter().StrictSlash(true)
	applicationRouter.HandleFunc("/new", handlers.PoolNewApplicationHandler(as.db)).
		Methods(http.MethodPost)
	applicationRouter.HandleFunc("/edit", handlers.PoolEditApplicationHandler(as.db)).
		Methods(http.MethodPost)
	applicationRouter.HandleFunc("/view", handlers.PoolViewApplicationHandler(as.db)).
		Methods(http.MethodPost)
	applicationRouter.HandleFunc("/status", handlers.PoolStatusViewHandler(as.db)).
		Methods(http.MethodPost)
	applicationRouter.HandleFunc("/pool/contains/{walletAddress:0[xX][0-9a-fA-F]{40}}", handlers.PoolContainsNode(as.db))
	applicationRouter.HandleFunc("/nodes", handlers.PoolNodes(as.db))
}
