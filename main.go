package main

// @title TeraRush API
// @version 1.0
// @description Modular blogging platform API
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

import (
	"flag"
	"go-modular/internal/app"
	"go-modular/internal/pkg/config"
	"go-modular/internal/pkg/logger"
	"go-modular/internal/pkg/middleware"
	"go-modular/modules/activity"
	"go-modular/modules/agent"
	"go-modular/modules/analytics"
	"go-modular/modules/assets"
	"go-modular/modules/auth"
	"go-modular/modules/blogs"
	"go-modular/modules/bookmarks"
	"go-modular/modules/comments"
	"go-modular/modules/moderation"
	"go-modular/modules/notifications"
	"go-modular/modules/subscriptions"
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
	app.RegisterModule(agent.NewModule())
	app.RegisterModule(activity.NewModule())
	app.RegisterModule(analytics.NewModule())
	app.RegisterModule(bookmarks.NewModule())
	app.RegisterModule(moderation.NewModule())
	app.RegisterModule(notifications.NewModule())
	app.RegisterModule(subscriptions.NewModule())
	app.RegisterModule(user.NewModule())
	app.RegisterModule(auth.NewModule())
	app.RegisterModule(blogs.NewModule())
	app.RegisterModule(assets.NewModule())
	app.RegisterModule(comments.NewModule())

	// initialize the application
	if err := app.Initialize(); err != nil {
		log.Fatalf("Error initializing application : %v", err)
		os.Exit(1)
	}

	// Start the application
	app.Start()
}
