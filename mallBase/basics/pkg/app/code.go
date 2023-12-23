package app

import "time"

// 返回自定义状态码
const (
	Success             = 0   // 响应正确
	Fail                = 2   // 响应错误
	Initialization      = 3   // 初始化错误
	NotExist            = 4   // 数据不存在
	RequestError        = 100 // 请求错误
	StoreError          = 201 // 店铺错误
	Validator           = 400 // 表单验证错误
	AuthFail            = 401 // 登录验证失败
	FailedToAcquireLock = 402 // 请求限制尝试过快
	PermissionDenied    = 403 // 无权限
	NotFound            = 404 // 404页面不存在
)

// 自定义的一些错误消息的返回
const (
	SuccessMessage             = "操作成功"
	NotExistMessage            = "数据不存在"
	RequestMessage             = "请求出现出错"
	ValidatorMessage           = "验证错误"
	AuthFailMessage            = "认证失败"
	PermissionDeniedMessage    = "无操作权限"
	FailedToAcquireLockMessage = "请求限制尝试过快"
)

// Redis Key
const (
	// Platformr
	RedisViewsStatisticsPlatformPvKey = "pv-platform:%d:%d:%d:%d" //Platform Pv统计key - 企业ID,统计日期,店铺ID,平台标识
	RedisViewsStatisticsPlatformUvKey = "uv-platform:%d:%d:%d:%d" //Platform Uv统计key - 企业ID,统计日期,店铺ID,平台标识
	RedisViewsStatisticsPlatformPv    = "pv-platform"             //Platform Pv统计key
	RedisViewsStatisticsPlatformUv    = "uv-platform"             //Platform Uv统计key

	// Area
	RedisViewsStatisticsAreaPvKey = "pv-area:%d:%d:%d:%d" //Area Pv统计key - 企业ID,统计日期,店铺ID,页面ID
	RedisViewsStatisticsAreaUvKey = "uv-area:%d:%d:%d:%d" //Area Uv统计key - 企业ID,统计日期,店铺ID,页面ID
	RedisViewsStatisticsAreaPv    = "pv-area"             //Area Pv统计key
	RedisViewsStatisticsAreaUv    = "uv-area"             //Area Uv统计key

	// Goods
	RedisViewsStatisticsGoodsPvKey = "pv-goods:%d:%d:%d:%d" //Goods Pv统计key - 企业ID,统计日期,店铺ID,商品ID
	RedisViewsStatisticsGoodsUvKey = "uv-goods:%d:%d:%d:%d" //Goods Uv统计key - 企业ID,统计日期,店铺ID,商品ID
	RedisViewsStatisticsGoodsPv    = "pv-goods"             //Goods Pv统计key
	RedisViewsStatisticsGoodsUv    = "uv-goods"             //Goods Uv统计key

	// View
	RedisViewsStatisticsViewPvKey = "pv-view:%d:%d" //View Pv统计key - 企业ID,统计日期
	RedisViewsStatisticsViewUvKey = "uv-view:%d:%d" //View Uv统计key - 企业ID,统计日期
	RedisViewsStatisticsViewPv    = "pv-view"       //View Pv统计key
	RedisViewsStatisticsViewUv    = "uv-view"       //View Uv统计key
)

const SmsRedisStoreKey = "msg-auth:%s"           //短信code
const SmsCostRedisStoreKey = "sms_cost:cache:%d" // 短信扣费

// Message Type

const (
	MessageTypeSms       = 1 //1：短信
	MessageTypeEmail     = 2 //2：邮件
	MessageTypeWxAccount = 3 //3：微信公众号
	MessageTypeApp       = 4 //4：APP消息推送
)

// 极光推送消息类型

const (
	MessageJpushTypeAll    = 0 //全部
	MessageJpushTypeNotice = 1 //通知
	MessageJpushTypeMsg    = 2 //消息
)

//audienceType 推送对象类型

const (
	AudienceTypeAll   = 0 //所有
	AudienceTypeAlias = 1 //别名
)

const (
	SmsRuleTagsAli     = "smsConfigAli"   //短信规则标签 - 阿里云
	SmsRuleTagsGuoYang = "smsConfigGy"    //短信规则标签 - 国阳云
	AuthCodeExpir      = 10 * time.Minute //验证码过期时间
)
