package main

import (
	"log"
	"os"

	"sparklink-backend/config"
	"sparklink-backend/handler"
	"sparklink-backend/middleware"
	"sparklink-backend/repository"
	"sparklink-backend/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// 初始化Redis
	rdb := config.InitRedis(cfg)

	// 初始化仓库
	userRepo := repository.NewUserRepository(db)
	nodeRepo := repository.NewNodeRepository(db)
	rewardRepo := repository.NewRewardRepository(db)
	subscriptionRepo := repository.NewSubscriptionRepository(db)

	// 初始化服务
	authService := service.NewAuthService(userRepo, cfg)
	nodeService := service.NewNodeService(nodeRepo)
	rewardService := service.NewRewardService(rewardRepo, rdb)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)

	// 初始化处理器
	authHandler := handler.NewAuthHandler(authService)
	nodeHandler := handler.NewNodeHandler(nodeService)
	rewardHandler := handler.NewRewardHandler(rewardService)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)

	// 创建路由
	r := gin.Default()

	// 中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// API路由
	api := r.Group("/api/v1")
	{
		// 认证路由
		auth := api.Group("/auth")
		{
			auth.POST("/sendcode", authHandler.SendCode)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", middleware.Auth(cfg), authHandler.RefreshToken)
			auth.POST("/logout", middleware.Auth(cfg), authHandler.Logout)
		}

		// 节点路由
		nodes := api.Group("/nodes")
		{
			nodes.GET("/list", middleware.Auth(cfg), nodeHandler.List)
			nodes.GET("/:id", middleware.Auth(cfg), nodeHandler.Get)
			nodes.POST("/ping", middleware.Auth(cfg), nodeHandler.UpdatePing)
			nodes.GET("/favorites", middleware.Auth(cfg), nodeHandler.Favorites)
			nodes.POST("/favorites", middleware.Auth(cfg), nodeHandler.AddFavorite)
			nodes.DELETE("/favorites/:id", middleware.Auth(cfg), nodeHandler.RemoveFavorite)
		}

		// 激励路由
		rewards := api.Group("/rewards")
		{
			rewards.POST("/video", middleware.Auth(cfg), rewardHandler.VideoReward)
			rewards.POST("/daily", middleware.Auth(cfg), rewardHandler.DailyCheckin)
			rewards.GET("/invite", middleware.Auth(cfg), rewardHandler.GetInviteInfo)
			rewards.POST("/invite", middleware.Auth(cfg), rewardHandler.BindInvite)
			rewards.GET("/credits", middleware.Auth(cfg), rewardHandler.GetCredits)
		}

		// 订阅路由
		subscription := api.Group("/subscription")
		{
			subscription.GET("/plans", subscriptionHandler.ListPlans)
			subscription.POST("/create", middleware.Auth(cfg), subscriptionHandler.Create)
			subscription.POST("/verify", middleware.Auth(cfg), subscriptionHandler.Verify)
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}