package mock

import (
	"sparklink-backend/model"
	"time"
)

type MockData struct {
	Users           []model.User
	Nodes           []model.Node
	Subscriptions   []model.Subscription
	Plans           []model.Plan
	AdLogs          []model.AdLog
	VerificationCodes []model.VerificationCode
	QRSessions      []model.QRSession
	ConnectSessions []model.ConnectSession
}

func NewMockData() *MockData {
	now := time.Now()
	weeklyPrice := 0.43
	monthlyPrice := 0.33
	quarterlyPrice := 0.28
	yearlyPrice := 0.22

	return &MockData{
		Users: []model.User{
			{
				ID:             1,
				Phone:          "+1234567890",
				Nickname:       "Demo User",
				DeviceID:       "device-001",
				VipStatus:      "inactive",
				BalanceMinutes: 120,
				InviteCode:     "SPARK123456",
				InvitedCount:   0,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		},
		Nodes: []model.Node{
			{ID: 1, NodeId: "node_tokyo_001", Name: "东京-01", RegionCode: "JP", RegionName: "日本·东京", Protocol: "wireguard", Protocols: "wireguard,v2ray", Latency: 45, Load: 35, VisibilityLevel: "free", Priority: 1, Tags: "fast", Host: "jp1.sparklink.app", Port: 443, PublicKey: "pubkey1", Distance: 2000, BandwidthLimit: 100, Status: "online", CreatedAt: now, UpdatedAt: now},
			{ID: 2, NodeId: "node_tokyo_002", Name: "东京-02", RegionCode: "JP", RegionName: "日本·东京", Protocol: "wireguard", Protocols: "wireguard", Latency: 52, Load: 28, VisibilityLevel: "vip", Priority: 1, Tags: "game,vip", Host: "jp2.sparklink.app", Port: 51820, PublicKey: "pubkey2", Distance: 2000, BandwidthLimit: 100, Status: "online", CreatedAt: now, UpdatedAt: now},
			{ID: 3, NodeId: "node_singapore_001", Name: "新加坡-01", RegionCode: "SG", RegionName: "新加坡", Protocol: "wireguard", Protocols: "wireguard,v2ray,openvpn", Latency: 78, Load: 42, VisibilityLevel: "free", Priority: 2, Tags: "stable", Host: "sg1.sparklink.app", Port: 443, PublicKey: "pubkey3", Distance: 3500, BandwidthLimit: 100, Status: "online", CreatedAt: now, UpdatedAt: now},
			{ID: 4, NodeId: "node_singapore_002", Name: "新加坡-02", RegionCode: "SG", RegionName: "新加坡", Protocol: "wireguard", Protocols: "wireguard", Latency: 65, Load: 55, VisibilityLevel: "vip", Priority: 2, Tags: "game,vip", Host: "sg2.sparklink.app", Port: 51820, PublicKey: "pubkey4", Distance: 3500, BandwidthLimit: 100, Status: "online", CreatedAt: now, UpdatedAt: now},
			{ID: 5, NodeId: "node_la_001", Name: "洛杉矶-01", RegionCode: "US", RegionName: "美国·洛杉矶", Protocol: "wireguard", Protocols: "wireguard,v2ray", Latency: 120, Load: 67, VisibilityLevel: "free", Priority: 3, Tags: "video", Host: "us1.sparklink.app", Port: 443, PublicKey: "pubkey5", Distance: 8000, BandwidthLimit: 100, Status: "online", CreatedAt: now, UpdatedAt: now},
			{ID: 6, NodeId: "node_la_002", Name: "洛杉矶-02", RegionCode: "US", RegionName: "美国·洛杉矶", Protocol: "wireguard", Protocols: "wireguard,v2ray", Latency: 110, Load: 45, VisibilityLevel: "vip", Priority: 3, Tags: "game,vip,video", Host: "us2.sparklink.app", Port: 51820, PublicKey: "pubkey6", Distance: 8000, BandwidthLimit: 100, Status: "online", CreatedAt: now, UpdatedAt: now},
			{ID: 7, NodeId: "node_frankfurt_001", Name: "法兰克福-01", RegionCode: "DE", RegionName: "德国·法兰克福", Protocol: "wireguard", Protocols: "wireguard,openvpn", Latency: 150, Load: 38, VisibilityLevel: "free", Priority: 4, Tags: "stable", Host: "de1.sparklink.app", Port: 443, PublicKey: "pubkey7", Distance: 12000, BandwidthLimit: 100, Status: "online", CreatedAt: now, UpdatedAt: now},
			{ID: 8, NodeId: "node_sydney_001", Name: "悉尼-01", RegionCode: "AU", RegionName: "澳大利亚·悉尼", Protocol: "wireguard", Protocols: "wireguard", Latency: 180, Load: 22, VisibilityLevel: "vip", Priority: 5, Tags: "vip", Host: "au1.sparklink.app", Port: 443, PublicKey: "pubkey8", Distance: 15000, BandwidthLimit: 100, Status: "online", CreatedAt: now, UpdatedAt: now},
			{ID: 9, NodeId: "node_london_001", Name: "伦敦-01", RegionCode: "GB", RegionName: "英国·伦敦", Protocol: "v2ray", Protocols: "v2ray,wireguard", Latency: 140, Load: 50, VisibilityLevel: "vip", Priority: 4, Tags: "vip,video", Host: "uk1.sparklink.app", Port: 443, PublicKey: "pubkey9", Distance: 10000, BandwidthLimit: 100, Status: "online", CreatedAt: now, UpdatedAt: now},
			{ID: 10, NodeId: "node_seoul_001", Name: "首尔-01", RegionCode: "KR", RegionName: "韩国·首尔", Protocol: "wireguard", Protocols: "wireguard,v2ray", Latency: 38, Load: 31, VisibilityLevel: "free", Priority: 1, Tags: "fast,game", Host: "kr1.sparklink.app", Port: 443, PublicKey: "pubkey10", Distance: 1000, BandwidthLimit: 100, Status: "online", CreatedAt: now, UpdatedAt: now},
		},
		Plans: []model.Plan{
			{PlanID: "weekly", Name: "周卡", Price: 2.99, OriginalPrice: 4.99, Currency: "USD", DurationDays: 7, DailyPrice: &weeklyPrice, Popular: true, Tag: "hot", Features: "VIP nodes"},
			{PlanID: "monthly", Name: "月卡", Price: 9.99, OriginalPrice: 14.99, Currency: "USD", DurationDays: 30, DailyPrice: &monthlyPrice, Popular: true, Tag: "recommended", Features: "VIP nodes + No ads"},
			{PlanID: "quarterly", Name: "季卡", Price: 24.99, OriginalPrice: 29.99, Currency: "USD", DurationDays: 90, DailyPrice: &quarterlyPrice, Tag: "", Features: "VIP nodes + No ads"},
			{PlanID: "yearly", Name: "年卡", Price: 79.99, OriginalPrice: 99.99, Currency: "USD", DurationDays: 365, DailyPrice: &yearlyPrice, Tag: "best_value", Features: "All features + Priority support"},
			{PlanID: "traffic_3d", Name: "3天流量包", Price: 0.99, OriginalPrice: 1.99, Currency: "USD", DurationDays: 3, Tag: "", Features: "Unlimited traffic 3 days"},
			{PlanID: "traffic_30d", Name: "30天流量包", Price: 12.99, OriginalPrice: 14.99, Currency: "USD", DurationDays: 30, Tag: "", Features: "Unlimited traffic 30 days"},
		},
	}
}

func Now() time.Time {
	return time.Now()
}
