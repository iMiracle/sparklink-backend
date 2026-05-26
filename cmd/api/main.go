package main

import (
	"log"
	"os"
<<<<<<< HEAD
	"time"
=======
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f

	"sparklink-backend/config"
	"sparklink-backend/handler"
	"sparklink-backend/middleware"
<<<<<<< HEAD
	"sparklink-backend/mock"
	"sparklink-backend/model"
=======
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	"sparklink-backend/repository"
	"sparklink-backend/service"

	"github.com/gin-gonic/gin"
<<<<<<< HEAD
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()
	useMock := os.Getenv("USE_MOCK") == "true"

	var repos *service.Repos

	if useMock {
		mockData := mock.NewMockData()
		mockRepos := repository.NewMockRepository(mockData)
		repos = &service.Repos{
			User:         mockRepos.User,
			Node:         mockRepos.Node,
			Reward:       mockRepos.Reward,
			Subscription: mockRepos.Subscription,
			Connect:      mockRepos.Connect,
			Verification: mockRepos.Verification,
		}
		log.Println("Using mock data store")
	} else {
		db, err := config.InitDB(cfg)
		if err != nil {
			log.Fatalf("数据库连接失败: %v", err)
		}

		if err := db.AutoMigrate(
			&model.User{},
			&model.VerificationCode{},
			&model.QRSession{},
			&model.Device{},
			&model.Node{},
			&model.ConnectSession{},
			&model.Subscription{},
			&model.AdLog{},
			&model.DailyCheckin{},
			&model.Invite{},
			&model.Plan{},
		); err != nil {
			log.Fatalf("数据库迁移失败: %v", err)
		}
		log.Println("数据库迁移完成")

		seedData(db)

		repos = &service.Repos{
			User:         repository.NewUserRepository(db),
			Node:         repository.NewNodeRepository(db),
			Reward:       repository.NewRewardRepository(db),
			Subscription: repository.NewSubscriptionRepository(db),
			Connect:      repository.NewConnectRepository(db),
			Verification: repository.NewVerificationRepository(db),
		}
	}

	svc := service.NewServices(repos, cfg)

	authHandler := handler.NewAuthHandler(svc.Auth)
	nodeHandler := handler.NewNodeHandler(svc.Node)
	rewardHandler := handler.NewRewardHandler(svc.Reward)
	subHandler := handler.NewSubscriptionHandler(svc.Subscription)
	connectHandler := handler.NewConnectHandler(svc.Connect)
	userHandler := handler.NewUserHandler(svc.User)

	r := gin.Default()
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimiter())

	authMw := middleware.Auth(cfg.JWTSecret)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "timestamp": time.Now().Unix()})
	})

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/sendcode", authHandler.SendCode)
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authMw, authHandler.RefreshToken)
			auth.POST("/logout", authMw, authHandler.Logout)
			auth.POST("/qrcode", authHandler.QrCode)
			auth.GET("/qrcode/:sessionId", authHandler.QrCodeStatus)
		}

		nodes := api.Group("/nodes")
		{
			nodes.GET("/list", nodeHandler.List)
			nodes.GET("/:nodeId", nodeHandler.Get)
			nodes.POST("/speedtest", nodeHandler.Speedtest)
			nodes.POST("/ping", authMw, nodeHandler.UpdatePing)
		}

		rewards := api.Group("/rewards")
		{
			rewards.POST("/claim", authMw, rewardHandler.Claim)
			rewards.GET("/cooldown", authMw, rewardHandler.GetCooldown)
			rewards.GET("/invite", authMw, rewardHandler.GetInviteInfo)
		}

		sub := api.Group("/subscription")
		{
			sub.GET("/plans", subHandler.ListPlans)
			sub.POST("/create", authMw, subHandler.Create)
			sub.POST("/verify", authMw, subHandler.Verify)
		}

		conn := api.Group("/connect")
		{
			conn.POST("/start", authMw, connectHandler.Start)
			conn.POST("/stop", authMw, connectHandler.Stop)
			conn.POST("/report", authMw, connectHandler.Report)
		}

		user := api.Group("/user")
		{
			user.GET("/profile", authMw, userHandler.Profile)
			user.GET("/devices", authMw, userHandler.Devices)
			user.DELETE("/devices/:deviceId", authMw, userHandler.RemoveDevice)
		}
	}

	log.Printf("Server starting on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func seedData(db *gorm.DB) {
	var count int64
	db.Model(&model.Plan{}).Count(&count)
	if count > 0 {
		return
	}

	plans := []model.Plan{
		{PlanID: "weekly", Name: "周卡", Price: 2.99, OriginalPrice: 4.99, Currency: "USD", DurationDays: 7, Popular: true, Features: "VIP nodes"},
		{PlanID: "monthly", Name: "月卡", Price: 9.99, OriginalPrice: 14.99, Currency: "USD", DurationDays: 30, Popular: true, Features: "VIP nodes + No ads"},
		{PlanID: "quarterly", Name: "季卡", Price: 24.99, OriginalPrice: 29.99, Currency: "USD", DurationDays: 90, Features: "VIP nodes + No ads"},
		{PlanID: "yearly", Name: "年卡", Price: 79.99, OriginalPrice: 99.99, Currency: "USD", DurationDays: 365, Features: "All features + Priority support"},
		{PlanID: "traffic_3d", Name: "3天流量包", Price: 0.99, OriginalPrice: 1.99, Currency: "USD", DurationDays: 3, Features: "Unlimited traffic 3 days"},
		{PlanID: "traffic_30d", Name: "30天流量包", Price: 12.99, OriginalPrice: 14.99, Currency: "USD", DurationDays: 30, Features: "Unlimited traffic 30 days"},
	}
	for _, p := range plans {
		db.Create(&p)
	}

	nodes := []model.Node{
		{NodeId: "node_tokyo_001", Name: "东京-01", Protocol: "wireguard", Latency: 45, Load: 35, RegionCode: "JP", RegionName: "日本·东京", Tags: "fast", VisibilityLevel: "free", Priority: 1, Host: "jp1.sparklink.app", Port: 443, Status: "online"},
		{NodeId: "node_tokyo_002", Name: "东京-02", Protocol: "wireguard", Latency: 52, Load: 28, RegionCode: "JP", RegionName: "日本·东京", Tags: "game,vip", VisibilityLevel: "vip", Priority: 1, Host: "jp2.sparklink.app", Port: 51820, Status: "online"},
		{NodeId: "node_singapore_001", Name: "新加坡-01", Protocol: "wireguard", Latency: 78, Load: 42, RegionCode: "SG", RegionName: "新加坡", Tags: "stable", VisibilityLevel: "free", Priority: 2, Host: "sg1.sparklink.app", Port: 443, Status: "online"},
		{NodeId: "node_singapore_002", Name: "新加坡-02", Protocol: "wireguard", Latency: 65, Load: 55, RegionCode: "SG", RegionName: "新加坡", Tags: "game,vip", VisibilityLevel: "vip", Priority: 2, Host: "sg2.sparklink.app", Port: 51820, Status: "online"},
		{NodeId: "node_la_001", Name: "洛杉矶-01", Protocol: "wireguard", Latency: 120, Load: 67, RegionCode: "US", RegionName: "美国·洛杉矶", Tags: "video", VisibilityLevel: "free", Priority: 3, Host: "us1.sparklink.app", Port: 443, Status: "online"},
		{NodeId: "node_la_002", Name: "洛杉矶-02", Protocol: "wireguard", Latency: 110, Load: 45, RegionCode: "US", RegionName: "美国·洛杉矶", Tags: "game,vip,video", VisibilityLevel: "vip", Priority: 3, Host: "us2.sparklink.app", Port: 51820, Status: "online"},
		{NodeId: "node_frankfurt_001", Name: "法兰克福-01", Protocol: "wireguard", Latency: 150, Load: 38, RegionCode: "DE", RegionName: "德国·法兰克福", Tags: "stable", VisibilityLevel: "free", Priority: 4, Host: "de1.sparklink.app", Port: 443, Status: "online"},
		{NodeId: "node_sydney_001", Name: "悉尼-01", Protocol: "wireguard", Latency: 180, Load: 22, RegionCode: "AU", RegionName: "澳大利亚·悉尼", Tags: "vip", VisibilityLevel: "vip", Priority: 5, Host: "au1.sparklink.app", Port: 443, Status: "online"},
		{NodeId: "node_london_001", Name: "伦敦-01", Protocol: "v2ray", Latency: 140, Load: 50, RegionCode: "GB", RegionName: "英国·伦敦", Tags: "vip,video", VisibilityLevel: "vip", Priority: 4, Host: "uk1.sparklink.app", Port: 443, Status: "online"},
		{NodeId: "node_seoul_001", Name: "首尔-01", Protocol: "wireguard", Latency: 38, Load: 31, RegionCode: "KR", RegionName: "韩国·首尔", Tags: "fast,game", VisibilityLevel: "free", Priority: 1, Host: "kr1.sparklink.app", Port: 443, Status: "online"},
	}
	for _, n := range nodes {
		db.Create(&n)
	}

	log.Println("初始数据已写入")
}
=======
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
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
