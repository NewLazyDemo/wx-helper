package service

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"net/http"
	"net/url"
	"wx-helper/entity"
)

type OpenApiServiceInterface interface {
	Send(str string) (string, error)
}

// GptClient GPT 客户端
type GptClient struct {
	gpt     *openai.Client
	ctx     context.Context
	program string
}

// NewGptClient 创建一个新的 GptClient
func NewGptClient(ctx context.Context) *GptClient {
	openApi := entity.CommonConfig.OpenApi
	config := openai.DefaultConfig(openApi.OpenApiKey)
	if openApi.BaseUrl != "" {
		config.BaseURL = openApi.BaseUrl
	}
	if openApi.ProxyUrl != "" {
		httpClient := &http.Client{
			Transport: &http.Transport{
				Proxy: func(request *http.Request) (*url.URL, error) {
					return url.Parse(openApi.ProxyUrl)
				},
			},
		}
		config.HTTPClient = httpClient
	}
	client := openai.NewClientWithConfig(config)
	return &GptClient{
		ctx: ctx,
		gpt: client,
	}
}

// Send 发送消息并获取响应
func (g *GptClient) Send(str string) (string, error) {
	message, ok := g.ctx.Value("message").([]openai.ChatCompletionMessage)
	if !ok {
		message = []openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleSystem, Content: str}}
	} else {
		message = append(message, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleUser, Content: str})
	}

	resp, err := g.gpt.CreateChatCompletion(
		g.ctx,
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: message,
		},
	)
	if err != nil {
		return "", errors.Wrap(err, "ChatGpt Error")
	}

	message = append(message, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleAssistant, Content: resp.Choices[0].Message.Content})
	g.ctx = context.WithValue(g.ctx, "message", message)
	return resp.Choices[0].Message.Content, nil
}
