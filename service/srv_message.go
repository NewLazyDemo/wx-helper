package service

import (
	"github.com/eatmoreapple/openwechat"
	"regexp"
	"wx-helper/entity"
)

// MessageEnum 消息类型枚举
// key为正则表达式，value为具体消息bot的实现
var messageEnum = map[string]messageInterface{
	"(聊天|对话|询问|咨询)": &ChatMessage{},
}

// MessageInterface 消息实现接口
// todo 这里后续需要增加消息实现，例如列出目前所有bot、唤醒对话bot关键词、在指定时间后提醒我（提醒bot）、总结xxxx群从时间段-时间段的什么内容....
type messageInterface interface {
	Handle(msg *openwechat.Message) error
	Close(msg *openwechat.Message) string
}

// MessageService 消息处理服务
type MessageService struct {
}

func NewMessageService() *MessageService {
	return new(MessageService)
}

func checkInWhitelist(username string) bool {
	whitelist := entity.CommonConfig.Default.Whitelist
	if whitelist != nil && len(whitelist) > 0 {
		for _, v := range whitelist {
			if v == username {
				return true
			}
		}
	}
	return false
}

// Handle 实际消息处理
func (m *MessageService) Handle(msg *openwechat.Message) error {
	// todo 消息记录
	// 消息分析（是否符合白名单人发言，符合才可以触发机器人功能）
	// 获取发送者用户名
	from, _ := msg.Sender()
	if !checkInWhitelist(from.NickName) {
		return nil
	}
	// 消息处理
	for k, v := range messageEnum {
		if ok, _ := regexp.MatchString(k, msg.Content); ok {
			return v.Handle(msg)
		}
	}
	return nil
}
