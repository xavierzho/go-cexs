package cexconns

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

type WsChannel struct {
	dialer     *websocket.Dialer
	conn       *websocket.Conn
	ctx        context.Context
	cancelFunc context.CancelFunc
	Channels   map[string]chan any // id: data channel
	pingFunc   func() bool
	mu         sync.Mutex  // 用于保护 conn 和 Channels 的并发安全
	url        string      // 保存 WebSocket URL 以支持重连
	header     http.Header // 保存头信息以支持重连
}

// NewWsChannel 创建一个新的 WebSocket 工具实例
func NewWsChannel() *WsChannel {
	ctx, cancel := context.WithCancel(context.Background())
	dialer := websocket.DefaultDialer

	return &WsChannel{
		dialer:     dialer,
		ctx:        ctx,
		cancelFunc: cancel,
		Channels:   make(map[string]chan any),
	}
}

// Connect 连接到指定的 WebSocket URL
func (ws *WsChannel) Connect(url string, header http.Header) error {
	conn, _, err := ws.dialer.DialContext(ws.ctx, url, header)
	if err != nil {
		return err
	}
	ws.mu.Lock()
	ws.conn = conn
	ws.url = url
	ws.header = header
	ws.mu.Unlock()

	// 启动接收消息和 Ping 检测
	go ws.listen()
	go ws.keepAlive()
	return nil
}

// Reconnect 重连到 WebSocket
func (ws *WsChannel) Reconnect() error {
	ws.Close() // 关闭当前连接
	return ws.Connect(ws.url, ws.header)
}

// Close 关闭 WebSocket 连接并清理资源
func (ws *WsChannel) Close() error {
	ws.cancelFunc() // 取消上下文
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if ws.conn != nil {
		err := ws.conn.Close()
		ws.conn = nil
		return err
	}
	return nil
}

// Subscribe 订阅指定 topic，并返回消息通道
func (ws *WsChannel) Subscribe(topic string, payload any) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if ws.conn == nil {
		return fmt.Errorf("connection is not established")
	}

	// 发送订阅消息
	err := ws.conn.WriteJSON(payload)
	if err != nil {
		return err
	}

	// 创建通道以接收该 topic 的消息
	ws.Channels[topic] = make(chan any, 100)
	return nil
}

// Unsubscribe 取消订阅指定 topic
func (ws *WsChannel) Unsubscribe(topic string) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if _, exists := ws.Channels[topic]; exists {
		close(ws.Channels[topic])
		delete(ws.Channels, topic)
	}

	// 可根据协议实现取消订阅的消息发送逻辑
	return nil
}

// keepAlive 定期发送 Ping 消息
func (ws *WsChannel) keepAlive() {
	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ws.mu.Lock()
			if ws.conn == nil {
				ws.mu.Unlock()
				return
			}
			err := ws.conn.WriteMessage(websocket.PingMessage, nil)
			ws.mu.Unlock()
			if err != nil {
				log.Printf("Ping failed: %v. Reconnecting...", err)
				ws.Reconnect() // 断开后重连
			}
		case <-ws.ctx.Done():
			return
		}
	}
}

// listen 监听消息并分发到相应的 channel
func (ws *WsChannel) listen() {
	for {
		select {
		case <-ws.ctx.Done():
			return
		default:
			ws.mu.Lock()
			if ws.conn == nil {
				ws.mu.Unlock()
				return
			}
			_, message, err := ws.conn.ReadMessage()
			ws.mu.Unlock()
			if err != nil {
				log.Printf("ReadMessage error: %v. Reconnecting...", err)
				ws.Reconnect()
				continue
			}

			// 假设消息是 JSON 格式，解析并分发
			ws.dispatch(message)
		}
	}
}

// dispatch 分发消息到相应的 channel
func (ws *WsChannel) dispatch(message []byte) {
	// 根据具体协议解析消息，找到对应的 topic
	// 假设协议中消息格式是：{"topic":"<topicName>", "data":<payload>}
	var msg struct {
		Topic string `json:"topic"`
		Data  any    `json:"data"`
	}

	// 解码消息
	err := json.Unmarshal(message, &msg)
	if err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	// 将数据发送到对应的 channel
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if ch, exists := ws.Channels[msg.Topic]; exists {
		select {
		case ch <- msg.Data:
		default: // 如果通道阻塞，丢弃消息以防死锁
			log.Printf("Channel for topic %s is full. Dropping message.", msg.Topic)
		}
	}
}
