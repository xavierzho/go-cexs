package platforms

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/xavierzho/go-cexs/types"
)

type StreamBase struct {
	dialer     *websocket.Dialer
	conn       *websocket.Conn
	url        string
	ctx        context.Context
	cancelFunc context.CancelFunc
	payload    map[string]any
}

func NewStream() *StreamBase {
	ctx, cancel := context.WithCancel(context.Background())

	return &StreamBase{
		dialer:     websocket.DefaultDialer,
		ctx:        ctx,
		cancelFunc: cancel,
	}
}
func (stream *StreamBase) Connect(url string) error {
	conn, _, err := stream.dialer.DialContext(stream.ctx, url, nil)
	if err != nil {
		return err
	}
	stream.conn = conn
	return nil
}
func (stream *StreamBase) SendMessage(payload map[string]any) error {
	stream.payload = payload
	err := stream.conn.WriteJSON(payload)
	if err != nil {
		return err
	}
	_, err = stream.ReadMessage()
	return err
}
func (stream *StreamBase) Close() error {
	stream.cancelFunc()
	return stream.conn.Close()
}
func (stream *StreamBase) Reconnect() error {
	err := stream.Connect(stream.url)
	if err != nil {
		return err
	}
	return stream.SendMessage(stream.payload)
}
func (stream *StreamBase) ReadMessage() ([]byte, error) {
	_, msg, err := stream.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return msg, nil
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
