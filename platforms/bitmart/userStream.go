package bitmart

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/xavierzho/go-cexs/platforms"
	"github.com/xavierzho/go-cexs/types"
	"github.com/xavierzho/go-cexs/utils"
)

type UserDataStream struct {
	credentials *platforms.Credentials

	*platforms.StreamBase
}

func NewUserStream(credentials *platforms.Credentials) *UserDataStream {
	stream := platforms.NewStream()
	return &UserDataStream{
		credentials: credentials,
		StreamBase:  stream,
	}
}
func (stream *UserDataStream) Login() error {
	var timestamp = strconv.FormatInt(time.Now().UnixMilli(), 10)
	err := stream.Connect(PrivateChannel)
	if err != nil {
		return err
	}

	signature := stream.Sign(timestamp)
	return stream.SendMessage(map[string]any{
		"op": "login",
		"args": []string{
			stream.credentials.APIKey,
			timestamp,
			signature,
		},
	})
}

func (stream *UserDataStream) Sign(timestamp string) string {
	var message = new(bytes.Buffer)
	message.WriteString(timestamp)
	message.WriteRune('#')
	message.WriteString(*stream.credentials.Option)
	message.WriteString("#bitmart.WebSocket")
	mac := hmac.New(sha256.New, []byte(stream.credentials.APISecret))
	mac.Write(message.Bytes())
	return hex.EncodeToString(mac.Sum(nil))
}

type OrderUpdate struct {
	Symbol         string `json:"symbol"`
	Notional       string `json:"notional"`
	LastFillTime   string `json:"last_fill_time"`
	Type           string `json:"type"`
	FilledNotional string `json:"filled_notional"`
	LastFillPrice  string `json:"last_fill_price"`
	UpdateTime     string `json:"update_time"`
	Price          string `json:"price"`
	LastFillCount  string `json:"last_fill_count"`
	State          string `json:"state"`
	OrderType      string `json:"order_type"`
	Side           string `json:"side"`
	ClientOrderID  string `json:"client_order_id"`
	CreateTime     string `json:"create_time"`
	OrderMode      string `json:"order_mode"`
	MsT            string `json:"ms_t"`
	ExecType       string `json:"exec_type"`
	DealFee        string `json:"dealFee"`
	DetailID       string `json:"detail_id"`
	Size           string `json:"size"`
	FilledSize     string `json:"filled_size"`
	MarginTrading  string `json:"margin_trading"`
	OrderID        string `json:"order_id"`
	OrderState     string `json:"order_state"`
	EntrustType    string `json:"entrust_type"`
}

func (o OrderUpdate) GetSymbol() string {
	return o.Symbol
}

type BalanceUpdate struct {
	EventType      string          `json:"event_type"`
	BalanceDetails []BalanceDetail `json:"balance_details"`
	EventTime      string          `json:"event_time"`
}

func (b BalanceUpdate) GetSymbol() string {
	return b.EventType
}

type BalanceDetail struct {
	Lock  string `json:"fz_bal"`
	Free  string `json:"av_bal"`
	Asset string `json:"ccy"`
}

func (stream *UserDataStream) BalanceStream(ctx context.Context, channel chan<- types.BalanceUpdateEntry) error {
	err := stream.SendMessage(map[string]any{
		"op": "subscribe",
		"args": []string{
			"spot/user/balance:BALANCE_UPDATE",
		},
	})
	if err != nil {
		return err
	}

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := stream.ReadMessage()
				if err != nil {
					continue
				}
				var event StreamResp[BalanceUpdate]
				_ = utils.Json.Unmarshal(msg, &event)
				channel <- types.BalanceUpdateEntry{}
			}
		}
	}(ctx)
	return nil
}

func (stream *UserDataStream) OrderStream(ctx context.Context, channel chan<- types.OrderUpdateEntry) error {
	err := stream.SendMessage(map[string]any{
		"op": "subscribe",
		"args": []string{
			"spot/user/order:ALL_SYMBOLS",
		},
	})
	if err != nil {
		return err
	}

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := stream.ReadMessage()
				if err != nil {
					continue
				}
				var event StreamResp[OrderUpdate]
				_ = utils.Json.Unmarshal(msg, &event)

				for _, e := range event.Data {
					channel <- types.OrderUpdateEntry{
						OrderId:       e.OrderID,
						ClientOrderId: e.ClientOrderID,
						Status:        OrderStatus(e.OrderState).Convert(),
					}
				}
			}
		}
	}(ctx)
	return nil
}

func (stream *UserDataStream) AccountStream(ctx context.Context, channel chan<- types.AccountUpdateEntry) error {
	return nil
}
