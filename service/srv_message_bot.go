package service

import (
	"context"
	"github.com/eatmoreapple/openwechat"
	"sync"
)

// MessageSession 消息bot 服务 (偷懒直接丢内存)
// todo 后续优化为其他存储，保证可靠性
// key 为微信用户的From+To
var MessageSession = make(map[string]ChatMessage)
var lock = sync.RWMutex{}

// todo 新增定时器，定时清理过期的bot
const defaultLifeTime = 180 // 默认生命时长 3分钟

// ChatMessage 对话bot
// 需要配置生命时长（例如3分钟），超过生命时长直接销毁
type ChatMessage struct {
	createTime int64 // 创建时间 按秒计算
	diedTime   int64 // 死亡时间 按秒计算
	ai         OpenApiServiceInterface
}

// Handle 处理消息
func (c *ChatMessage) Handle(msg *openwechat.Message) error {
	lock.Lock()
	defer lock.Unlock()
	// 判断是否已经存在对话bot
	// 这里一旦出现死亡时间超过现在时间，就直接回复
	if client, ok := MessageSession[msg.FromUserName+msg.ToUserName]; ok {
		if client.diedTime < msg.CreateTime {
			// 销毁
			delete(MessageSession, msg.FromUserName+msg.ToUserName)
			msg.ReplyText(c.Error(""))
			return nil
		}
		// 刷新死亡时间
		client.diedTime = msg.CreateTime + defaultLifeTime
		//内容处理
		if reply, err := client.ai.Send(msg.Content); err != nil {
			return err
		} else {
			msg.ReplyText(reply)
		}
		return nil
	}
	ctx := context.Background()
	// 先默认用gpt
	ai := NewGptClient(ctx)
	MessageSession[msg.FromUserName+msg.ToUserName] = ChatMessage{
		createTime: msg.CreateTime,
		diedTime:   msg.CreateTime + defaultLifeTime, // 3分钟
		ai:         ai,
	}
	//内容处理
	if reply, err := ai.Send(msg.Content); err != nil {
		return err
	} else {
		msg.ReplyText(reply)
	}
	return nil
}

func (c *ChatMessage) Close(msg *openwechat.Message) string {
	lock.Lock()
	defer lock.Unlock()
	// 判断是否已经存在对话bot
	if _, ok := MessageSession[msg.FromUserName+msg.ToUserName]; ok {
		// 销毁
		delete(MessageSession, msg.FromUserName+msg.ToUserName)
	}
	return c.End()
}

func (c *ChatMessage) Error(str string) string {
	if str == "" {
		return "我不明白你在说什么"
	}
	return str
}

func (c *ChatMessage) End() string {
	return "再见，和你聊天很愉快，有时间再聊"
}
