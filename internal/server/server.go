package server

import (
	"github.com/gladiusio/gladius-application-server/internal/controller"
	asrouting "github.com/gladiusio/gladius-application-server/internal/routing"
	"github.com/gladiusio/gladius-common/pkg/manager"
	"github.com/gladiusio/gladius-common/pkg/routing"
	"github.com/gladiusio/gladius-common/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"strings"
)

var Database *gorm.DB
var Router *mux.Router

func initializeConfiguration() {
	// TODO: Use the gladius-common definition when it's in
	base, err := utils.GetGladiusBase()
	if err != nil {
		log.Warn().Err(err).Msg("Error retrieving base directory")
	}

	// Add config file name and searching
	viper.SetConfigName("gladius-application-server")
	viper.AddConfigPath(base)

	// Setup env variable handling
	viper.SetEnvPrefix("APPLICATION-SERVER")
	r := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(r)
	viper.AutomaticEnv()

	// Load config
	err = viper.ReadInConfig()
	if err != nil {
		log.Warn().Err(err).Msg("Error reading config file, it may not exist or is corrupted. Using defaults.")
	}

	// Build our config options
	buildOptions(base)
}

func buildOptions(base string) {
	// Log options
	ConfigOption("Log.Level", "info")
	ConfigOption("Log.Pretty", true)

	// Blockchain options
	ConfigOption("Blockchain.Provider", "https://mainnet.infura.io/tjqLYxxGIUp0NylVCiWw")
	ConfigOption("Blockchain.MarketAddress", "0x27a9390283236f836a0b3c8dfdbed2ed854322fc")
	ConfigOption("Blockchain.PoolUrl", "http://174.138.111.1/api/")
	ConfigOption("Blockchain.PoolManagerAddress", "0x9717EaDbfE344457135a4f1fA8AE3B11B4CAB0b7")

	// API options
	ConfigOption("API.Port", 3333)
	ConfigOption("API.DebugRequests", true)
	ConfigOption("API.RemoteConnectionsAllowed", true)

	// Database
	ConfigOption("Database.Type", "sqlite3")
	ConfigOption("Database.Connection", "local.db")

	// Service
	ConfigOption("Service.Name", "GladiusApplicationServer")
	ConfigOption("Service.DisplayName", "Application Server")
	ConfigOption("Service.Description", "Gladius Application Server")
	ConfigOption("Service.Debug", true)

	// Misc.
	ConfigOption("GladiusBase", base) // Convenient option to have, not needed though
}

func ConfigOption(key string, defaultValue interface{}) string {
	viper.SetDefault(key, defaultValue)

	return key
}

func initializeService() {
	// Router Setup
	cRouter := routing.ControlRouter{
		Router: Router,
		Port:   viper.GetString("API.Port"),
		Debug:  viper.GetBool("API.DebugRequests"),
	}

	manager.RunService(
		viper.GetString("Service.Name"),
		viper.GetString("Service.DisplayName"),
		viper.GetString("Service.Description"),
		cRouter.Start)
}

func initializeDatabase(dbType, dbConnection string) {
	db, err := gorm.Open(dbType, dbConnection)
	if err != nil {
		log.Fatal()
	}

	Database, err = controller.Initialize(db)
	if err != nil {
		log.Fatal()
	}
}

func InitializeApplicationServer() {
	// Grab default configuration
	initializeConfiguration()

	initializeDatabase(viper.GetString("Database.Type"), viper.GetString("Database.Connection"))
	Router = asrouting.ApplicationServerRouter(Database)

	initializeService()
}
