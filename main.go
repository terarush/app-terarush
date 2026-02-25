package main

import (
	"flag"
	"go-modular/internal/app"
	"go-modular/internal/pkg/config"
	"go-modular/internal/pkg/logger"
	"go-modular/internal/pkg/middleware"
	"go-modular/modules/auth"
	"go-modular/modules/nodes"
	"go-modular/modules/products"
	"go-modular/modules/transactions"
	user "go-modular/modules/users"
	"log"
	"os"
)

var configFile *string

func init() {
	configFile = flag.String("c", ".env", "configuration file")
	flag.Parse()
}

func main() {

	// Load configuration
	cfg := config.NewConfig(*configFile)
	if err := cfg.Initialize(); err != nil {
		log.Fatalf("Error reading config : %v", err)
		os.Exit(1)
	}

	// initialize logger
	logCfg := logger.DefaultConfig()

	// Start the application
	app, err := app.NewApp(&logCfg)
	if err != nil {
		log.Fatalf("Error creating application : %v", err)
		os.Exit(1)
	}

	// Initialize Auth middleware
	jwtSignatureKey := config.GetJWTService()
	middleware.InitializeAuth(jwtSignatureKey)

	// register modules
	app.RegisterModule(user.NewModule())
	app.RegisterModule(auth.NewModule())
	app.RegisterModule(nodes.NewModule())
	app.RegisterModule(products.NewModule())
	app.RegisterModule(transactions.NewModule())

	// initialize the application
	if err := app.Initialize(); err != nil {
		log.Fatalf("Error initializing application : %v", err)
		os.Exit(1)
	}

	// Start the application
	app.Start()
}
