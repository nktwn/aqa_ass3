package app

import (
	"context"
	"diploma/docs"
	"diploma/internal/config"
	"diploma/modules/auth"
	"diploma/modules/cart"
	"diploma/modules/contract"
	"diploma/modules/notification"
	"diploma/modules/order"
	"diploma/modules/product"
	"diploma/modules/user"
	"diploma/pkg/logger"
	"diploma/pkg/metrics"

	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	serviceProvider *serviceProvider
	httpServer      *gin.Engine
}

func NewApp(ctx context.Context) (*App, error) {
	// Initialize logger first
	env := config.GetEnv("APP_ENV", "development")
	if err := logger.Initialize(env); err != nil {
		return nil, err
	}
	logger.Info("Logger initialized", logger.Field("environment", env))

	app := &App{}
	err := app.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) Run() error {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := a.httpServerRun(); err != nil {
			logger.Fatal("HTTP server error", logger.Field("error", err))
		}
	}()

	wg.Wait()
	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	// Try to load .env file but don't fail if it doesn't exist
	_ = config.Load(".env")
	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	router := gin.Default()

	// Add Prometheus middleware
	router.Use(metrics.MetricsMiddleware())

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	router.Use(func(c *gin.Context) {
		corsMiddleware.HandlerFunc(c.Writer, c.Request)
		c.Next()
	})

	// Add Prometheus metrics endpoint
	metrics.RegisterMetricsEndpoint(router)

	docs.SwaggerInfo.BasePath = ""
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiGroup := router.Group("/api")

	authHandler := a.serviceProvider.AuthHandler(ctx)
	auth.RegisterRoutes(apiGroup, authHandler)
	authMiddleware := a.serviceProvider.AuthMiddleware(ctx)

	productHandler := a.serviceProvider.ProductHandler(ctx)
	product.RegisterRoutes(apiGroup, productHandler)

	cartCallbackHander := a.serviceProvider.CartHandler(ctx)
	cart.RegisterRoutesCallback(apiGroup, cartCallbackHander)
	cart.RegisterPublicRoutes(apiGroup, cartCallbackHander)

	secureGroup := router.Group("/api")
	secureGroup.Use(authMiddleware.AuthMiddleware())

	product.RegisterProtectedRoutes(secureGroup, productHandler)
	product.RegisterAdminRoutes(secureGroup, productHandler)

	cartGroup := a.serviceProvider.CartHandler(ctx)
	cart.RegisterRoutes(secureGroup, cartGroup)

	notificationHandler := a.serviceProvider.NotificationHandler(ctx)
	notification.RegisterRoutes(secureGroup, notificationHandler)

	orderHandler := a.serviceProvider.OrderHandler(ctx)
	order.RegisterRoutes(secureGroup, orderHandler)

	contractHandler := a.serviceProvider.ContractHandler(ctx)
	contract.RegisterRoutes(secureGroup, contractHandler)

	userHandler := a.serviceProvider.UserHandler(ctx)
	user.RegisterRoutes(secureGroup, userHandler)

	a.httpServer = router
	return nil
}

func (a *App) httpServerRun() error {
	address := a.serviceProvider.HTTPConfig().Address()
	log.Printf("HTTP server is running on %s", address)
	return a.httpServer.Run(address)
}
