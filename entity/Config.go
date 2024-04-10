package entity

var CommonConfig Config

type Config struct {
	Default DefaultConfig // 基础配置
	Wechat  WechatConfig  // 微信配置
	OpenApi OpenApiConfig // OpenApi配置 (这里是从命令行读取的)
}

type OpenApiConfig struct {
	OpenApiKey string
	BaseUrl    string
	ProxyUrl   string
}

type DefaultConfig struct {
	Whitelist []string // 白名单
}

type WechatConfig struct {
	HotReloadStorageDir string
}
