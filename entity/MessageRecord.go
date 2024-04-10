package entity

import "time"

// MessageRecord 消息记录
type MessageRecord struct {
	Id         uint64    `json:"id"`
	Content    string    `json:"content"`
	From       string    `json:"from"`
	To         string    `json:"to"`
	CreateTime time.Time `json:"createTime"`
}
