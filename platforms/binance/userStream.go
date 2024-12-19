package binance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/xavierzho/go-cexs/utils"
	"io"
	"log"
	"net/http"
)

var ListenKeyEndpoint = RestAPI + "/api/v3/userDataStream"

type UserDataStream struct {
	APIKey    string
	APISecret string
	dialer    *websocket.Dialer
	conn      *websocket.Conn

	account chan StreamResponse[AccountUpdate]
	balance chan StreamResponse[BalanceUpdate]
	order   chan StreamResponse[OrderUpdate]
	//expired chan map[string]any
	listenKey string
}

type StreamEvent struct {
	Time int64  `json:"E"` // Event time
	Type string `json:"e"` // Event type
}
type Event interface {
	GetTime() int64
	GetType() string
}

func (e StreamEvent) GetTime() int64 {
	return e.Time
}

func (e StreamEvent) GetType() string {
	return e.Type
}

type OrderUpdate struct {
	StreamEvent
	OriginalOrderId         string `json:"C,omitempty"` // Original client order ID; This is the ID of the order being canceled
	Iceberg                 string `json:"F,omitempty"` // Iceberg quantity
	Ignore                  int    `json:"I,omitempty"` // Ignore
	LastExecutedPrice       string `json:"L,omitempty"` // Last executed price
	M                       bool   `json:"M,omitempty"` // Ignore
	CommissionAsset         string `json:"N,omitempty"` // Commission asset
	OrderCreatedTime        int64  `json:"O,omitempty"` // Order creation time
	StopPrice               string `json:"P,omitempty"` // stop price
	QuoteOrderQty           string `json:"Q,omitempty"` // Quote Order Quantity
	Side                    string `json:"S,omitempty"` // order side
	TransactionTime         int64  `json:"T,omitempty"` // Transaction Time
	SelfTradePreventionMode string `json:"V,omitempty"` // SelfTradePreventionMode
	WorkingTime             int64  `json:"W,omitempty"` // Working Time; This is only visible if the order has been placed on the book.
	OrderStatus             string `json:"X,omitempty"` // Current order status
	LastQuoteQty            string `json:"Y,omitempty"` // Last quote asset transacted quantity (i.e. lastPrice * lastQty)
	CumulativeQty           string `json:"Z,omitempty"` // Cumulative quote asset transacted quantity
	ClientOrderId           string `json:"c,omitempty"` // clientOrderId
	TimeInForce             string `json:"f,omitempty"` // Time in force
	OrderListId             int    `json:"g,omitempty"` // OrderListId
	OrderId                 int    `json:"i,omitempty"` // OrderId
	LastQuantity            string `json:"l,omitempty"` // Last executed quantity
	TradeMakerSide          bool   `json:"m,omitempty"` // Is this trade the maker side?
	CommissionAmount        string `json:"n,omitempty"` // Commission amount
	OrderType               string `json:"o,omitempty"` // order type
	Price                   string `json:"p,omitempty"` // order price
	Quantity                string `json:"q,omitempty"` // Order quantity
	RejectReason            string `json:"r,omitempty"` // Order reject reason; will be an error code.
	Symbol                  string `json:"s,omitempty"` // symbol
	TradeId                 int    `json:"t,omitempty"` // Trade ID
	OnOrderBook             bool   `json:"w,omitempty"` // Is the order on the book?
	ExecutionType           string `json:"x,omitempty"` // Current execution type
	CumulativeQuantity      string `json:"z,omitempty"` // Cumulative filled quantity
	STP                     int64  `json:"v,omitempty"` // Prevented Match Id; This is only visible if the order expired due to STP
}

type BalanceUpdate struct {
	StreamEvent
	Asset     string `json:"a"`
	Delta     string `json:"d"`
	ClearTime int64  `json:"T"`
}

type AccountUpdate struct {
	StreamEvent
	Balances        []StreamAsset `json:"B"`
	UpdateTimestamp int64         `json:"u"`
}

type StreamAsset struct {
	Asset string `json:"a"`
	Free  string `json:"f"`
	Lock  string `json:"l"`
}
type ListenKeyResponse struct {
	ListenKey string `json:"listenKey"`
}

type PostListenKey struct {
	ListenKey string `json:"listenKey"`
}

type StreamResponse[T Event] struct {
	Stream string `json:"stream"`
	Data   T      `json:"data"`
}

func NewUserStream(apikey, apiSecret string) *UserDataStream {
	return &UserDataStream{
		APIKey:    apikey,
		APISecret: apiSecret,
		dialer:    websocket.DefaultDialer,
	}
}

// https://developers.binance.com/docs/binance-spot-api-docs/user-data-stream#create-a-listenkey-user_stream
func (stream *UserDataStream) getListenKey() error {
	//payload := map[string]string{}
	req, err := http.NewRequest(http.MethodPost, ListenKeyEndpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(HeaderAPIKEY, stream.APIKey)
	//http.Post(fmt.Sprintf("%s/api/v3/userDataStream", HttpBaseEndpoint), "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	var buff bytes.Buffer
	var response ListenKeyResponse
	_, err = io.Copy(&buff, resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buff.Bytes(), &response)
	if err != nil {
		return err
	}
	stream.listenKey = response.ListenKey
	return nil

}

// https://developers.binance.com/docs/binance-spot-api-docs/user-data-stream#close-a-listenkey-user_stream
func (stream *UserDataStream) closeListenKey(key string) error {
	var body = PostListenKey{
		ListenKey: key,
	}
	bytesBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodDelete, ListenKeyEndpoint, bytes.NewReader(bytesBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(HeaderAPIKEY, stream.APIKey)
	_, err = http.DefaultClient.Do(req)
	return err
}

// https://developers.binance.com/docs/binance-spot-api-docs/user-data-stream#pingkeep-alive-a-listenkey-user_stream
func (stream *UserDataStream) keepAlive(listenKey string) error {
	var body = PostListenKey{
		ListenKey: listenKey,
	}
	bytesBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, ListenKeyEndpoint, bytes.NewReader(bytesBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(HeaderAPIKEY, stream.APIKey)
	_, err = http.DefaultClient.Do(req)
	return err
}
func (stream *UserDataStream) Reconnect() error {
	var attempt int
	for attempt = 0; attempt < 3; attempt++ {
		_ = stream.closeListenKey(stream.listenKey)
		log.Printf("Reconnection attempt %d of %d...\n", attempt+1, 3)

		// 获取新的 listenKey 并尝试连接
		if err := stream.getListenKey(); err != nil {
			log.Printf("Failed to get listenKey: %v\n", err)
			//time.Sleep(stream.reconnectInterval)
			continue
		}

		// 尝试建立 WebSocket 连接
		conn, _, err := stream.dialer.Dial(StreamAPI+"?streams="+stream.listenKey, nil)
		if err != nil {
			log.Printf("Failed to reconnect WebSocket: %v\n", err)
			//time.Sleep(stream.reconnectInterval)
			continue
		}

		// 更新连接
		stream.conn = conn
		log.Println("Reconnected successfully")

		// 重新开始接收消息
		go stream.listenMessages()

		return nil
	}

	return fmt.Errorf("failed to reconnect after %d attempts", 3)
}
func (stream *UserDataStream) Login() error {
	err := stream.getListenKey()
	if err != nil {
		return err
	}
	//
	conn, _, err := stream.dialer.Dial(StreamAPI+"?streams="+stream.listenKey, nil)
	if err != nil {
		return err
	}
	stream.conn = conn
	go stream.listenMessages()
	return nil
}

func (stream *UserDataStream) listenMessages() {
	for {
		_, msg, err := stream.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading WebSocket message: %v\n", err)
			// 在连接失败时尝试重连
			if err := stream.Reconnect(); err != nil {
				log.Printf("Reconnect failed: %v\n", err)
				return
			}
		}

		//if err := json.Unmarshal(msg, &event); err != nil {
		//	log.Printf("Error unmarshalling WebSocket message: %v\n", err)
		//	continue
		//}

		// 根据事件类型分发处理
		switch utils.Json.Get(msg, "data", "e").ToString() {
		case "executionReport":
			var event StreamResponse[OrderUpdate]

			_ = utils.Json.Unmarshal(msg, &event)
			stream.order <- event
		case "outboundAccountPosition":
			var event StreamResponse[AccountUpdate]

			_ = utils.Json.Unmarshal(msg, &event)
			stream.account <- event
		case "balanceUpdate":
			var event StreamResponse[BalanceUpdate]

			_ = utils.Json.Unmarshal(msg, &event)
			stream.balance <- event
		case "listenKeyExpired":
			log.Println("ListenKey expired, reconnecting...")
			// 当 listenKey 过期时，重新连接
			if err := stream.Reconnect(); err != nil {
				log.Printf("Failed to reconnect after listenKey expired: %v\n", err)
				return
			}
		}
	}
}

func (stream *UserDataStream) GetOrderUpdate() <-chan StreamResponse[OrderUpdate] {
	return stream.order
}

func (stream *UserDataStream) GetAccountUpdate() <-chan StreamResponse[AccountUpdate] {
	return stream.account
}

func (stream *UserDataStream) GetBalanceUpdate() <-chan StreamResponse[BalanceUpdate] {
	return stream.balance
}
