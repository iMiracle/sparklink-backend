package response

const (
	ErrBalanceInsufficient = "ERR_001"
	ErrConnectionTimeout   = "ERR_002"
	ErrProtocolUnsupported = "ERR_003"
	ErrNetworkUnreachable  = "ERR_004"
	ErrPermissionDenied    = "ERR_005"
	ErrNodeUnavailable     = "ERR_006"
	ErrServerError         = "ERR_007"
	ErrSubscriptionExpired = "ERR_008"
	ErrDeviceLimit         = "ERR_009"
	ErrAuthFailed          = "ERR_010"
	ErrRiskBlocked         = "ERR_011"
	ErrRegionRestricted    = "ERR_012"
	ErrCooldown            = "ERR_013"
	ErrInvalidParams       = "ERR_014"
	ErrNotFound            = "ERR_015"
	ErrRateLimited         = "ERR_016"
)

var ErrorMessages = map[string]string{
	ErrBalanceInsufficient: "余额不足",
	ErrConnectionTimeout:   "连接超时",
	ErrProtocolUnsupported: "协议不支持",
	ErrNetworkUnreachable:  "网络不可达",
	ErrPermissionDenied:    "权限被拒",
	ErrNodeUnavailable:     "节点不可用",
	ErrServerError:         "服务器错误",
	ErrSubscriptionExpired: "订阅过期",
	ErrDeviceLimit:         "设备超限",
	ErrAuthFailed:          "认证失败",
	ErrRiskBlocked:         "风控拦截",
	ErrRegionRestricted:    "地区限制",
	ErrCooldown:            "冷却中",
	ErrInvalidParams:       "参数错误",
	ErrNotFound:            "资源不存在",
	ErrRateLimited:         "请求过于频繁",
}
