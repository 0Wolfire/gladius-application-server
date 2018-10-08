package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gladiusio/gladius-application-server/internal/controller"
	asrouting "github.com/gladiusio/gladius-application-server/internal/routing"
	"github.com/gladiusio/gladius-common/pkg/manager"
	"github.com/gladiusio/gladius-common/pkg/routing"
	"github.com/gladiusio/gladius-common/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"log"
	"path/filepath"
	"strings"
)

var Database *gorm.DB

func initializeConfiguration() {
	// Set Defaults if no configuration file is found
	defaultConfiguration()
}

type ConfigurationOptions struct {
	Name        string
	DisplayName string
	Description string
	Debug       bool
	Port        string
}

type DatabaseConfig struct {
	Type     string
	Host     string
	Port     string
	User     string
	Name     string
	Password string
	SSL      bool
}

func (databaseConfig *DatabaseConfig) GormConnectionString() string {
	sslMode := "disable"

	if databaseConfig.SSL {
		sslMode = "require"
	}

	connection := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		databaseConfig.Host,
		databaseConfig.Port,
		databaseConfig.User,
		databaseConfig.Name,
		databaseConfig.Password,
		sslMode,
	)

	return connection
}

type BlockchainConfig struct {
	Provider           string
	MarketAddress      string
	PoolUrl            string
	PoolManagerAddress string
}

type Configuration struct {
	Version    string
	Build      int
	Blockchain BlockchainConfig
	Directory  struct {
		Base   string
		Wallet string
	}
	ApplicationServer struct {
		Database DatabaseConfig
		Config   ConfigurationOptions
	}
}

func (configuration Configuration) defaults() Configuration {
	baseDir, err := utils.GetGladiusBase()
	if err != nil {
		baseDir = ""
	}

	return Configuration{
		Version: "0.7.0",
		Build:   20181008,
		Blockchain: BlockchainConfig{
			Provider:           "https://mainnet.infura.io/tjqLYxxGIUp0NylVCiWw",
			MarketAddress:      "0x27a9390283236f836a0b3c8dfdbed2ed854322fc",
			PoolUrl:            "http://174.138.111.1/api/",
			PoolManagerAddress: "0x9717EaDbfE344457135a4f1fA8AE3B11B4CAB0b7",
		},
		Directory: struct {
			Base   string
			Wallet string
		}{
			Base:   baseDir,
			Wallet: filepath.Join(baseDir, "wallet"),
		},
	}
}

func defaultConfiguration() {
	var configuration Configuration

	// Path of used config value
	configFile := viper.ConfigFileUsed()

	if configFile == "" {
		log.Printf("\n\nUnable to find gladius-controld.toml in project root, or default directories below.\n")
		log.Printf("\n\nUsing Default Node Manager Configuration")

		configuration = configuration.defaults()

		viper.SetDefault("blockchain", configuration.Blockchain)

		jsonBytes, err := json.Marshal(configuration)
		if err != nil {
			log.Fatalf("Unable to marshal configuration struct to load defaults, %v", err)
		}

		viper.SetConfigType("json")
		viper.ReadConfig(bytes.NewBuffer(jsonBytes))
	}

	// Setup environment vars, they look like CONTROLD_OBJECT_KEY
	viper.SetEnvPrefix("CONTROLD")
	r := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(r)
	viper.AutomaticEnv()
}

func initializeService(router *mux.Router, configuration ConfigurationOptions) {
	// Router Setup
	cRouter := routing.ControlRouter{
		Router: router,
		Port:   configuration.Port,
		Debug:  configuration.Debug,
	}

	manager.RunService(configuration.Name, configuration.DisplayName, configuration.Description, cRouter.Start)
}

func initializeDatabase(databaseConfig DatabaseConfig) {
	db, err := gorm.Open(databaseConfig.Type, databaseConfig.GormConnectionString())
	if err != nil {
		log.Fatal("Could not open database")
	}

	Database, err = controller.Initialize(db)
	if err != nil {
		log.Fatal("Could not migrate database")
	}
}

// Returns the singleton viper config as a parsed struct
func ViperConfiguration() Configuration {
	var viperConfiguration Configuration
	viper.Unmarshal(&viperConfiguration)

	return viperConfiguration
}

func InitializeApplicationServer() {
	// Grab default configuration
	initializeConfiguration()

	configuration := ViperConfiguration()

	asConfig := configuration.ApplicationServer.Config

	initializeDatabase(configuration.ApplicationServer.Database)
	initializeService(asrouting.ApplicationServerRouter(Database), asConfig)
}
