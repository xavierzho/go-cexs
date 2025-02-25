package platforms

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/xavierzho/go-cexs/types"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type StreamBase struct {
	dialer        *websocket.Dialer
	conn          *websocket.Conn
	url           string
	ctx           context.Context
	cancelFunc    context.CancelFunc
	payload       atomic.Value
	mux           sync.RWMutex
	reconnectChan chan struct{}
	closeOne      sync.Once
}

const (
	defaultWriteWait  = 30 * time.Second
	defaultPingPeriod = 15 * time.Second
)

func NewStream() *StreamBase {
	ctx, cancel := context.WithCancel(context.Background())

	return &StreamBase{
		dialer:        websocket.DefaultDialer,
		ctx:           ctx,
		cancelFunc:    cancel,
		reconnectChan: make(chan struct{}, 1),
	}
}
func (stream *StreamBase) getConn() *websocket.Conn {
	stream.mux.RLock()
	defer stream.mux.RUnlock()
	return stream.conn
}
func (stream *StreamBase) Connect(url string) error {
	stream.mux.Lock()
	defer stream.mux.Unlock()
	// close old conn
	if stream.conn != nil {
		_ = stream.conn.Close()
	}
	conn, _, err := stream.dialer.DialContext(stream.ctx, url, nil)
	if err != nil {
		return err
	}
	stream.conn = conn
	stream.url = url
	// 设置超时
	deadline := time.Now().Add(defaultWriteWait)

	err = stream.conn.SetReadDeadline(deadline)
	err = stream.conn.SetWriteDeadline(deadline)
	// 设置 Pong 处理器，响应服务器的 Ping
	stream.conn.SetPongHandler(func(string) error {
		return stream.conn.SetReadDeadline(time.Now().Add(defaultPingPeriod)) // 更新读取超时
	})

	go stream.KeepAlive(defaultPingPeriod)
	return nil
}
func (stream *StreamBase) SendMessage(payload map[string]any) error {
	stream.payload.Store(payload)
	conn := stream.getConn()
	if conn == nil {
		return fmt.Errorf("not connected")
	}
	stream.mux.Lock()
	defer stream.mux.Unlock()
	err := stream.conn.WriteJSON(payload)
	if err != nil {
		return err
	}
	_, _, err = stream.conn.ReadMessage()
	return err
}
func (stream *StreamBase) Close() error {

	stream.closeOne.Do(func() {
		stream.cancelFunc()
		stream.mux.Lock()
		defer stream.mux.Unlock()
		if stream.conn != nil {
			_ = stream.conn.Close()
			stream.conn = nil
		}
	})
	return nil
}
func (stream *StreamBase) Reconnect() error {
	stream.mux.Lock()
	defer stream.mux.Unlock()
	// 指数退避重试
	backoff := time.Second
	maxBackoff := 30 * time.Second

	for retry := 0; ; retry++ {
		if err := stream.Connect(stream.url); err != nil {
			if payload := stream.payload.Load(); payload != nil {
				return stream.SendMessage(payload.(map[string]any))
			}
			return nil
		}
		if backoff > maxBackoff {
			return fmt.Errorf("max reconnect attempts reached")
		}
		select {
		case <-stream.ctx.Done():
			return errors.New("context canceled")
		case <-time.After(backoff):
			backoff *= 2
		}
		time.Sleep(backoff)
	}
}
func (stream *StreamBase) ReadMessage() ([]byte, error) {
	conn := stream.getConn()
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// KeepAlive keep alive must be with go keyword.
func (stream *StreamBase) KeepAlive(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	//stream.conn
	for {
		select {
		case <-stream.ctx.Done():
			return
		case <-ticker.C:
			conn := stream.getConn()
			if conn == nil {
				continue
			}
			err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(5*time.Second))
			if err != nil {
				log.Printf("Ping failed: %v, attempting reconnect", err)
				stream.scheduleReconnect()
			}
		case <-stream.reconnectChan:
			go stream.Reconnect()
		}
	}
}
func (stream *StreamBase) scheduleReconnect() {
	select {
	case stream.reconnectChan <- struct{}{}:
		fmt.Println("Reconnect scheduled")
	default:

	}
}

type StreamClient interface {
	// Connect dial to server
	Connect(url string) error
	// Close connection
	Close() error
	// SendMessage send json payload
	SendMessage(payload map[string]any) error
	// ReadMessage read ws channel message
	ReadMessage() ([]byte, error)
	Reconnect
}
type Reconnect interface {
	// Reconnect reset ws connection
	Reconnect() error
}
type MarketStreamer interface {
	DepthStream(ctx context.Context, symbol string, channel chan<- types.DepthEntry) error
	CandleStream(ctx context.Context, symbol, interval string, channel chan<- types.CandleEntry) error
}

type UserDataStreamer interface {
	Login() error
	Reconnect
	OrderStream(ctx context.Context, channel chan<- types.OrderUpdateEntry) error
	BalanceStream(ctx context.Context, channel chan<- types.BalanceUpdateEntry) error
	AccountStream(ctx context.Context, channel chan<- types.AccountUpdateEntry) error
}
