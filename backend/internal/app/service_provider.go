package app

import (
	"context"
	"log"

	"diploma/internal/config"

	"diploma/pkg/client/db"
	"diploma/pkg/client/db/pg"
	"diploma/pkg/client/db/transaction"
	"diploma/pkg/closer"

	authApi "diploma/modules/auth/handler"
	authJWT "diploma/modules/auth/jwt"
	authMiddlware "diploma/modules/auth/middleware"
	authRepository "diploma/modules/auth/repository/user"
	authService "diploma/modules/auth/service/auth"

	productApi "diploma/modules/product/handler"
	productRepository "diploma/modules/product/repository/product"
	productService "diploma/modules/product/service"

	cartOrderClient "diploma/modules/cart/client/order"
	cartPaymentClient "diploma/modules/cart/client/payment"
	cartSupplierClient "diploma/modules/cart/client/supplier"
	cartApi "diploma/modules/cart/handler"
	cartRedisClient "diploma/modules/cart/redis"
	cartRepository "diploma/modules/cart/repository"
	cartService "diploma/modules/cart/service"

	supplierRepository "diploma/modules/supplier/repo"
	supplierService "diploma/modules/supplier/service"

	orderProductClient "diploma/modules/order/client/product"
	orderSupplierCleint "diploma/modules/order/client/supplier"
	orderHander "diploma/modules/order/handler"
	orderRepository "diploma/modules/order/repo"
	// orderProductClient "diploma/modules/order/client/product"
	orderService "diploma/modules/order/service"
	notificationHandler "diploma/modules/notification/handler"
	notificationRepository "diploma/modules/notification/repository"
	notificationService "diploma/modules/notification/service"

	userApi "diploma/modules/user/handler"
	userRepository "diploma/modules/user/repository"
	userService "diploma/modules/user/service"

	contractHandler "diploma/modules/contract/handler"
	contractRepo "diploma/modules/contract/repo"
	contractService "diploma/modules/contract/service"

	"github.com/go-redis/redis/v8"
)

type serviceProvider struct {
	pgConfig      config.PGConfig
	jwtConfig     config.JWTConfig
	httpConfig    config.HTTPConfig
	swaggerConfig config.SwaggerConfig
	paymentConfig config.PaymentConfig
	redisConfig   config.RedisConfig

	dbClient    db.Client
	txManager   db.TxManager
	redisClient []*redis.Client

	// auth
	authRepository authService.IAuthRepository
	authJWT        *authJWT.JSONWebToken
	authService    authApi.IAuthService
	authHanlder    *authApi.AuthHandler
	authMiddlware  *authMiddlware.AuthMiddleware

	// product

	productRepository productService.IProductRepository
	productService    *productService.ProductService
	productHanlder    *productApi.CatalogHandler

	// cart
	cartSupplierClient cartService.ISupplierClient
	cartOrderClient    cartService.IOrderClient
	cartRepository     cartService.ICartRepository
	cartService        cartApi.ICartService
	cartPaymentClient  cartService.IPaymentClient
	cartRedisClient    cartService.IRedis
	cartHanlder        *cartApi.CartHandler

	// supplier
	supplierRepository supplierService.ISupplierRepository
	supplierService    *supplierService.SupplierService

	// order
	// orderHandler  *orderHandler.OrderHandler
	orderRepository     orderService.IOrderRepository
	orderSupplierCleint orderService.ISupplierClient
	orderProductClient  orderService.IProductClient
	orderService        *orderService.OrderService
	orderHandler        *orderHander.OrderHandler

	// user
	userRepository userService.IUserRepository
	userService    userApi.IUserService
	userHandler    *userApi.UserHandler

	// contract
	contractRepo    contractService.Repository
	contractService *contractService.Service
	contractHandler *contractHandler.Handler

	// notification
	notificationRepo    notificationService.Repository
	notificationService *notificationService.Service
	notificationHandler *notificationHandler.Handler
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) JWTConfig() config.JWTConfig {
	if s.jwtConfig == nil {
		cfg, err := config.NewJWTConfig()
		if err != nil {
			log.Fatalf("failed to get jwt config: %s", err.Error())
		}

		s.jwtConfig = cfg
	}

	return s.jwtConfig
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := config.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swagger config: %s", err.Error())
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

func (s *serviceProvider) PaymentConfig() config.PaymentConfig {
	if s.paymentConfig == nil {
		cfg, err := config.NewPaymentConfig()
		if err != nil {
			log.Fatalf("failed to get payment config: %s", err.Error())
		}

		s.paymentConfig = cfg
	}

	return s.paymentConfig
}

func (s *serviceProvider) RedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := config.NewRedisConfig()
		if err != nil {
			log.Fatalf("failed to get redis config: %s", err.Error())
		}

		s.redisConfig = cfg
	}

	return s.redisConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) RedisClient(ctx context.Context, dbNumber int) *redis.Client {
	if s.redisClient == nil {
		s.redisClient = make([]*redis.Client, 16)
	}
	if s.redisClient[dbNumber] == nil {
		rdb := redis.NewClient(&redis.Options{
			Addr:     s.RedisConfig().Addr(),
			Password: s.RedisConfig().Password(),
			DB:       dbNumber,
		})
		s.redisClient[dbNumber] = rdb
	}
	return s.redisClient[dbNumber]
}

// ========= authentication =========
func (s *serviceProvider) AuthRepository(ctx context.Context) authService.IAuthRepository {
	if s.authRepository == nil {
		s.authRepository = authRepository.NewRepository(s.DBClient(ctx))
	}

	return s.authRepository
}

func (s *serviceProvider) JWT(ctx context.Context) *authJWT.JSONWebToken {
	if s.authJWT == nil {
		s.authJWT = authJWT.NewJSONWebToken(s.JWTConfig().GetSecretKey())
	}

	return s.authJWT
}

func (s *serviceProvider) AuthService(ctx context.Context) authApi.IAuthService {
	if s.authService == nil {
		s.authService = authService.NewService(s.AuthRepository(ctx), s.JWT(ctx), s.TxManager(ctx))
	}

	return s.authService
}

func (s *serviceProvider) AuthHandler(ctx context.Context) *authApi.AuthHandler {
	if s.authHanlder == nil {
		s.authHanlder = authApi.NewHandler(s.AuthService(ctx))
	}

	return s.authHanlder
}

func (s *serviceProvider) AuthMiddleware(ctx context.Context) *authMiddlware.AuthMiddleware {
	if s.authMiddlware == nil {
		s.authMiddlware = authMiddlware.NewAuthMiddleware(s.JWT(ctx))
	}

	return s.authMiddlware
}

// ========= product =========

func (s *serviceProvider) ProductRepository(ctx context.Context) productService.IProductRepository {
	if s.productRepository == nil {
		s.productRepository = productRepository.NewRepository(s.DBClient(ctx))
	}

	return s.productRepository
}

func (s *serviceProvider) ProductService(ctx context.Context) *productService.ProductService {
	if s.productService == nil {
		s.productService = productService.NewService(s.ProductRepository(ctx), s.TxManager(ctx))
	}

	return s.productService
}

func (s *serviceProvider) ProductHandler(ctx context.Context) *productApi.CatalogHandler {
	if s.productHanlder == nil {
		s.productHanlder = productApi.NewHandler(s.ProductService(ctx))
	}

	return s.productHanlder
}

// ========= suppliers =========

func (s *serviceProvider) SupplierRepo(ctx context.Context) supplierService.ISupplierRepository {
	if s.supplierRepository == nil {
		s.supplierRepository = supplierRepository.NewRepository(s.DBClient(ctx))
	}

	return s.supplierRepository
}

func (s *serviceProvider) SupplierService(ctx context.Context) *supplierService.SupplierService {
	if s.supplierService == nil {
		s.supplierService = supplierService.NewService(s.SupplierRepo(ctx), s.TxManager(ctx))
	}

	return s.supplierService
}

// ========= cart =========

func (s *serviceProvider) CartRepo(ctx context.Context) cartService.ICartRepository {
	if s.cartRepository == nil {
		s.cartRepository = cartRepository.NewRepository(s.DBClient(ctx))
	}

	return s.cartRepository

}

func (s *serviceProvider) CartSupplierClient(ctx context.Context) cartService.ISupplierClient {
	if s.cartSupplierClient == nil {
		s.cartSupplierClient = cartSupplierClient.NewClient(s.SupplierService(ctx))
	}

	return s.cartSupplierClient
}

func (s *serviceProvider) CartOrderClient(ctx context.Context) cartService.IOrderClient {
	if s.cartOrderClient == nil {
		s.cartOrderClient = cartOrderClient.NewClient(s.OrderService(ctx))
	}

	return s.cartOrderClient
}

func (s *serviceProvider) CartPaymentClient(ctx context.Context) cartService.IPaymentClient {
	if s.cartPaymentClient == nil {
		s.cartPaymentClient = cartPaymentClient.NewPaymentClient(s.PaymentConfig().CheckoutURL(), s.PaymentConfig().MerchantID(), s.PaymentConfig().MerchantPassword(), s.PaymentConfig().CallbackURL())
	}

	return s.cartPaymentClient
}

func (s *serviceProvider) CartRedisClient(ctx context.Context) cartService.IRedis {
	if s.cartRedisClient == nil {
		s.cartRedisClient = cartRedisClient.NewCartRedis(s.RedisClient(ctx, 0))
	}

	return s.cartRedisClient
}

func (s *serviceProvider) CartService(ctx context.Context) cartApi.ICartService {
	if s.cartService == nil {
		s.cartService = cartService.NewService(s.CartRepo(ctx), s.ProductService(ctx), s.CartSupplierClient(ctx), s.CartOrderClient(ctx), s.UserRepo(ctx), s.CartPaymentClient(ctx), s.CartRedisClient(ctx), s.TxManager(ctx))
	}

	return s.cartService
}

func (s *serviceProvider) CartHandler(ctx context.Context) *cartApi.CartHandler {
	if s.cartHanlder == nil {
		s.cartHanlder = cartApi.NewHandler(s.CartService(ctx))
	}

	return s.cartHanlder
}

// order

func (s *serviceProvider) OrderSupplierClient(ctx context.Context) orderService.ISupplierClient {
	if s.orderSupplierCleint == nil {
		s.orderSupplierCleint = orderSupplierCleint.NewClient(s.SupplierService(ctx))
	}

	return s.orderSupplierCleint
}

func (s *serviceProvider) OrderProductClient(ctx context.Context) orderService.IProductClient {
	if s.orderProductClient == nil {
		s.orderProductClient = orderProductClient.NewClient(s.ProductService(ctx))
	}

	return s.orderProductClient
}

func (s *serviceProvider) OrderRepo(ctx context.Context) orderService.IOrderRepository {
	if s.orderRepository == nil {
		s.orderRepository = orderRepository.NewRepository(s.DBClient(ctx))
	}

	return s.orderRepository
}

func (s *serviceProvider) OrderService(ctx context.Context) *orderService.OrderService {
	if s.orderService == nil {
		s.orderService = orderService.NewService(
			s.OrderRepo(ctx),
			s.OrderSupplierClient(ctx),
			s.OrderProductClient(ctx),
			s.OrderContractClient(ctx),
			s.NotificationService(ctx),
			s.TxManager(ctx),
		)

	}

	return s.orderService
}

func (s *serviceProvider) OrderHandler(ctx context.Context) *orderHander.OrderHandler {
	if s.orderHandler == nil {
		s.orderHandler = orderHander.NewHandler(s.OrderService(ctx))
	}

	return s.orderHandler
}

func (s *serviceProvider) NotificationRepo(ctx context.Context) notificationService.Repository {
	if s.notificationRepo == nil {
		s.notificationRepo = notificationRepository.NewRepository(s.DBClient(ctx))
	}
	return s.notificationRepo
}

func (s *serviceProvider) NotificationService(ctx context.Context) *notificationService.Service {
	if s.notificationService == nil {
		s.notificationService = notificationService.NewService(s.NotificationRepo(ctx))
	}
	return s.notificationService
}

func (s *serviceProvider) NotificationHandler(ctx context.Context) *notificationHandler.Handler {
	if s.notificationHandler == nil {
		s.notificationHandler = notificationHandler.NewHandler(s.NotificationService(ctx))
	}
	return s.notificationHandler
}

func (s *serviceProvider) UserRepo(ctx context.Context) userService.IUserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) UserService(ctx context.Context) userApi.IUserService {
	if s.userService == nil {
		s.userService = userService.NewService(s.UserRepo(ctx), s.TxManager(ctx))
	}

	return s.userService
}

func (s *serviceProvider) UserHandler(ctx context.Context) *userApi.UserHandler {
	if s.userHandler == nil {
		s.userHandler = userApi.NewHandler(s.UserService(ctx))
	}

	return s.userHandler
}

func (s *serviceProvider) ContractRepo(ctx context.Context) contractService.Repository {
	if s.contractRepo == nil {
		s.contractRepo = contractRepo.NewRepository(s.DBClient(ctx))
	}
	return s.contractRepo
}

func (s *serviceProvider) ContractService(ctx context.Context) *contractService.Service {
	if s.contractService == nil {
		s.contractService = contractService.NewService(s.ContractRepo(ctx))
	}
	return s.contractService
}

func (s *serviceProvider) ContractHandler(ctx context.Context) *contractHandler.Handler {
	if s.contractHandler == nil {
		s.contractHandler = contractHandler.NewHandler(s.ContractService(ctx))
	}
	return s.contractHandler
}

func (s *serviceProvider) OrderContractClient(ctx context.Context) orderService.IContractService {
	if s.contractService == nil {
		s.contractService = contractService.NewService(s.ContractRepo(ctx)) // или что-то подобное
	}
	return s.contractService
}
