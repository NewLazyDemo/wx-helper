package main

import (
	"flag"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/spf13/viper"
	"wx-helper/entity"
	"wx-helper/service"
)

func main() {
	// 命令行参数读取
	// -OPENAI_API_KEY 机器人API_KEY
	// -BASE_URL 请求地址
	// -PROXY_URL 代理地址
	var openAIKey, baseURL, proxyURL string
	flag.StringVar(&openAIKey, "OPENAI_API_KEY", "", "机器人API_KEY，可以使用英文逗号分隔多个key")
	flag.StringVar(&baseURL, "BASE_URL", "", "请求地址: https://api.openai.com/v1（默认）")
	flag.StringVar(&proxyURL, "PROXY_URL", "", "代理地址")
	// 解析命令行参数
	flag.Parse()

	if openAIKey == "" {
		panic("OPENAI_API_KEY 不能为空")
	}

	// 基础配置读取
	viper.AddConfigPath("./")
	viper.SetConfigName("default.config")
	viper.SetConfigType("toml")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	if err := viper.Unmarshal(&entity.CommonConfig); err != nil {
		panic(fmt.Errorf("Fatal error format config to json: %s \n", err))
	}

	// 数据回写
	entity.CommonConfig.OpenApi.OpenApiKey = openAIKey
	entity.CommonConfig.OpenApi.BaseUrl = baseURL
	entity.CommonConfig.OpenApi.ProxyUrl = proxyURL

	// 微信登录
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式
	// 注册消息处理函数
	bot.MessageHandler = func(msg *openwechat.Message) {
		if err := service.NewMessageService().Handle(msg); err != nil {
			fmt.Println(fmt.Sprintf("消息处理错误:%v", err))
		}
	}

	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 创建热存储容器对象
	reloadStorage := openwechat.NewFileHotReloadStorage(entity.CommonConfig.Wechat.HotReloadStorageDir)
	defer reloadStorage.Close()

	// 这里使用热登录
	if err := bot.HotLogin(reloadStorage, openwechat.NewRetryLoginOption()); err != nil {
		fmt.Println(err)
		return
	}

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
