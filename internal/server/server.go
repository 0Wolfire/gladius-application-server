package server

import (
	"os"
	"github.com/gladiusio/gladius-application-server/internal/controller"
	asrouting "github.com/gladiusio/gladius-application-server/internal/routing"
	"github.com/gladiusio/gladius-common/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"strings"
)

var Database *gorm.DB
var Router *mux.Router

func initializeConfiguration() (string, error) {
	// TODO: Use the gladius-common definition when it's in
	base, err := utils.GetGladiusBase()
	if err != nil {
		return "Error retrieving base directory", err
	}

	// Add config file name and searching
	viper.SetConfigName("gladius-application-server")
	viper.AddConfigPath(base)

	// Setup env variable handling
	viper.SetEnvPrefix("APPLICATIONSERVER")
	r := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(r)
	viper.AutomaticEnv()

	// Load config
	err = viper.ReadInConfig()
	var message = "Using provided config file and overriding default"
	
	if err != nil {
		return "Error reading config file, it may not exist or is corrupted. Using defaults.", err
	}

	// Build our config options
	buildOptions(base)

	return message, err
}

func buildOptions(base string) {
	// Log options
	ConfigOption("Log.Level", "info")
	ConfigOption("Log.Pretty", true)

	// Blockchain options
	ConfigOption("Blockchain.Provider", "https://mainnet.infura.io/tjqLYxxGIUp0NylVCiWw")
	ConfigOption("Blockchain.MarketAddress", "0x27a9390283236f836a0b3c8dfdbed2ed854322fc")

	// API options
	ConfigOption("API.Port", 3333)
	ConfigOption("API.DebugRequests", true)
	ConfigOption("API.RemoteConnectionsAllowed", true)

	// Database
	ConfigOption("Database.Type", "sqlite3")
	ConfigOption("Database.Connection", "local.db")
	ConfigOption("Database.InitializePoolInfo", true)

	// Service
	ConfigOption("Service.Name", "GladiusApplicationServer")
	ConfigOption("Service.DisplayName", "Application Server")
	ConfigOption("Service.Description", "Gladius Application Server")
	ConfigOption("Service.Debug", true)

	// Applications
	ConfigOption("Apllications.AutoAccept", false)

	// Misc.
	ConfigOption("GladiusBase", base) // Convenient option to have, not needed though
}

func ConfigOption(key string, defaultValue interface{}) string {
	viper.SetDefault(key, defaultValue)

	return key
}

func initializeService() {
	// Router Setup
	as := asrouting.New(viper.GetString("API.Port"), Database)
	as.Start()

	select {}
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

	if viper.GetBool("Database.InitializePoolInfo") {
		controller.InitializePoolManager(db)
	}
}

func initializeLogging() {
	// Setup logging level
	switch loglevel := viper.GetString("Log.Level"); strings.ToLower(loglevel) {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warning":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "disabled":
		zerolog.SetGlobalLevel(zerolog.Disabled)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	if !viper.IsSet("Log.Pretty") || (viper.IsSet("Log.Pretty") && viper.GetBool("Log.Pretty")) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

func InitializeApplicationServer() {
	// Grab default configuration
	message, err := initializeConfiguration()

	initializeLogging()

	if err != nil {
		log.Warn().Err(err).Msg(message)
	}

	initializeDatabase(viper.GetString("Database.Type"), viper.GetString("Database.Connection"))

	initializeService()
}
