package core

import (
	"douyin-backend/app/service/websocket/on_open_success"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type Client struct {
	Hub                *Hub            // 负责处理客户端注册、注销、在线管理
	Conn               *websocket.Conn // 一个ws连接
	Send               chan []byte     // 一个ws连接存储自己的消息管道
	PingPeriod         time.Duration
	ReadDeadline       time.Duration
	WriteDeadline      time.Duration
	HeartbeatFailTimes int
	ClientLastPongTime time.Time // 客户端最近一次响应服务端 ping 消息的时间
	State              uint8     // ws状态，1=ok；0=出错、掉线等
	sync.RWMutex
	on_open_success.ClientMoreParams // 这里追加一个结构体，方便开发者在成功上线后，可以自定义追加更多字段信息
}
