package app

import (
	"fmt"
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/config"
	"go-modular/internal/pkg/database"
	"go-modular/internal/pkg/logger"
	"go-modular/internal/pkg/server"
	_validator "go-modular/internal/pkg/validator"
	"time"

	_ "go-modular/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
)

// App represents the application
type App struct {
	db        *gorm.DB
	server    *server.ServerContext
	modules   []Module
	r         *echo.Echo
	logger    *logger.Logger
	startTime time.Time
}

// NewApp creates a new application
func NewApp(cfg *logger.Config) (*App, error) {
	appLogger, err := logger.NewLogger(*cfg, config.GetString("APP_NAME"))
	if err != nil {
		return nil, err
	}
	defer appLogger.Sync()
	return &App{
		modules: make([]Module, 0),
		logger:  appLogger,
	}, nil
}

func (a *App) SetRouter() *echo.Echo {
	return echo.New()
}

// RegisterModule registers a module with the application
func (a *App) RegisterModule(module Module) {
	a.modules = append(a.modules, module)
	a.logger.Info("Registered module: %s", module.Name())
}

// Initialize initializes the application
func (a *App) Initialize() error {
	a.logger.Info("Initializing application...")

	// Initialize database
	var err *error
	a.db, err = a.SetDatabase().OpenDB()
	if err != nil {
		a.logger.Error("Failed to initialize database: %v", err)
		return *err
	}

	// Set database instance for all modules
	database.DB = a.db

	// event bus initialization
	event := bus.NewEventBus()

	// initialize router
	a.r = a.SetRouter()
	a.r.Use(middleware.Logger())
	a.r.Use(middleware.Recover())

	// Setup CORS
	corsOrigins := config.GetString("CORS_ORIGINS")
	if corsOrigins == "" {
		corsOrigins = "*"
	}
	a.r.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{corsOrigins},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		MaxAge:       3600,
	}))

	// validate request
	a.r.Validator = _validator.NewCustomValidator()

	// Initialize modules
	for _, module := range a.modules {
		a.logger.Info("Initializing module: %s", module.Name())

		// Create module-specific logger
		moduleLogger := a.logger.WithPrefix(module.Name())
		if err := module.Initialize(a.db, moduleLogger, event); err != nil {
			a.logger.Error("Failed to initialize module %s: %v", module.Name(), err)
			return err
		}

		a.logger.Info("Module initialized: %s", module.Name())
	}

	// Run migrations for all modules
	for _, module := range a.modules {
		err := module.Migrations()
		if err != nil {
			a.logger.Error("Failed to run migrations for module %s: %v", module.Name(), err)
		}
		a.logger.Info("Migrations completed for module: %s", module.Name())
	}

	// Initialize HTTP server
	a.server = a.SetServer()

	// api version
	version := fmt.Sprintf("/api/v%s", config.GetString("API_VERSION"))

	// Serve static files from public directory
	a.r.Static("/public", "./public")
	a.logger.Info("Static file serving enabled for /public")

	a.r.GET("/api/docs/*", echoSwagger.WrapHandler)

	// Capture application start time
	a.startTime = time.Now()

	a.r.GET("/", func(c echo.Context) error {
		uptime := time.Since(a.startTime).String()
		return c.JSON(200, map[string]interface{}{
			"message": "Welcome to the API",
			"version": config.GetString("API_VERSION"),
			"status":  "running",
			"time":    time.Now().Format(time.RFC3339),
			"uptime":  uptime,
		})
	})

	// Add metadata endpoint with API configuration
	// @Summary Get API Metadata
	// @Description Returns API configuration, version, base URL, endpoints, and available features
	// @Tags System
	// @Produce json
	// @Success 200 {object} map[string]interface{} "API metadata"
	// @Router /api/meta [get]
	a.r.GET("/api/meta", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"appName":       config.GetString("APP_NAME"),
			"baseUrl":       fmt.Sprintf("http://%s", c.Request().Host),
			"apiVersion":    config.GetString("API_VERSION"),
			"apiPath":       version,
			"environment":   config.GetString("SERVER_MODE"),
			"status":        "running",
			"timestamp":     time.Now().Format(time.RFC3339),
			"uptime":        time.Since(a.startTime).String(),
			"accessToken":   "Bearer <JWT_TOKEN>",
			"tokenLocation": "Authorization header",
			"features": map[string]bool{
				"auth":      true,
				"users":     true,
				"blogs":     true,
				"comments":  true,
				"favorites": true,
				"assets":    true,
			},
			"endpoints": map[string]interface{}{
				"docs":     "/api/docs",
				"meta":     "/api/meta",
				"static":   "/public",
				"auth":     fmt.Sprintf("%s/auth", version),
				"users":    fmt.Sprintf("%s/users", version),
				"blogs":    fmt.Sprintf("%s/blogs", version),
				"comments": fmt.Sprintf("%s/comments", version),
				"assets":   fmt.Sprintf("%s/assets", version),
			},
			"defines": map[string]interface{}{
				"apiDefinition": map[string]string{
					"name":        "TeraRush API",
					"description": "A modular blogging platform",
					"version":     config.GetString("API_VERSION"),
					"environment": config.GetString("SERVER_MODE"),
				},
				"authentication": map[string]string{
					"type":   "Bearer Token",
					"scheme": "JWT",
					"header": "Authorization",
					"format": "Bearer <token>",
				},
				"rateLimit": map[string]interface{}{
					"enabled":  true,
					"requests": 1000,
					"window":   "1 hour",
				},
				"pagination": map[string]interface{}{
					"enabled": true,
					"maxSize": 100,
				},
			},
		})
	})

	// Register routes for all modules
	for _, module := range a.modules {
		a.logger.Info("Registering routes for module: %s", module.Name())
		module.RegisterRoutes(a.r, version)
		a.logger.Info("Routes registered for module: %s", module.Name())
	}

	// append handler to server
	a.server.Handler = a.r

	a.logger.Info("Application initialization completed")

	for _, v := range a.r.Routes() {
		fmt.Printf("PATH: %v | METHOD: %v\n", v.Path, v.Method)
	}

	return nil
}

// Start starts the application
func (a *App) Start() {
	a.logger.Info("Starting server on %s", a.server.Host)
	a.server.Run()
}

// setup database model
func (a *App) SetDatabase() *database.DBModel {
	return &database.DBModel{
		ServerMode:   config.GetString("SERVER_MODE"),
		Driver:       config.GetString("DB_DRIVER"),
		Host:         config.GetString("DB_HOST"),
		Port:         config.GetString("DB_PORT"),
		Name:         config.GetString("DB_NAME"),
		Username:     config.GetString("DB_USERNAME"),
		Password:     config.GetString("DB_PASSWORD"),
		MaxIdleConn:  config.GetInt("CONN_IDLE"),
		MaxOpenConn:  config.GetInt("CONN_MAX"),
		ConnLifeTime: config.GetInt("CONN_LIFETIME"),
	}
}

// Setup Web Server
func (a *App) SetServer() *server.ServerContext {
	return &server.ServerContext{
		Host:         ":" + config.GetString("SERVER_PORT"),
		ReadTimeout:  time.Duration(config.GetInt("HTTP_TIMEOUT")),
		WriteTimeout: time.Duration(config.GetInt("HTTP_TIMEOUT")),
	}
}
