package mexc

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/platforms"
	"github.com/xavierzho/go-cexs/types"
	"github.com/xavierzho/go-cexs/utils"
	"io"
	"net/http"
	"time"
)

const listenKeyEndpoint = RestAPI + ListenKeyEndpoint

type UserDataStream struct {
	*platforms.Credentials
	base      *platforms.StreamBase
	listenKey string
}

func (stream *UserDataStream) getListenKey() error {
	body, err := stream.request(http.MethodPost, map[string]any{})
	if err != nil {
		return err
	}
	stream.listenKey = utils.Json.Get(body, "listenKey").ToString()
	return nil
}
func (stream *UserDataStream) keepAlive(listenKey string) {
	resp, err := stream.request(http.MethodPut, map[string]any{
		"listenKey": listenKey,
	})
	if err != nil {
		return
	}
	stream.listenKey = utils.Json.Get(resp, "listenKey").ToString()
}
func (stream *UserDataStream) closeListenKey(listenKey string) {
	_, _ = stream.request(http.MethodDelete, map[string]any{
		"listenKey": listenKey,
	})
}

func (stream *UserDataStream) request(method string, params map[string]any) ([]byte, error) {
	params["timestamp"] = time.Now().UnixMilli()
	queryString := utils.EncodeParams(params)

	signature := stream.Sign([]byte(queryString))
	req, err := http.NewRequest(method, fmt.Sprintf("%s?%s&signature=%s", listenKeyEndpoint, params, signature), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(KeyHeader, stream.APIKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (stream *UserDataStream) listKeys() ([]string, error) {
	var result struct {
		ListenKey []string `json:"listenKey"`
	}
	resp, err := stream.request(http.MethodGet, map[string]any{})
	if err != nil {
		return nil, err
	}

	_ = utils.Json.Unmarshal(resp, &result)
	return result.ListenKey, nil
}
func (stream *UserDataStream) Login() error {
	err := stream.getListenKey()
	if err != nil {
		return err
	}
	return stream.base.Connect(StreamAPI + "?listenKey=" + stream.listenKey)
}
func (stream *UserDataStream) Sign(params []byte) string {
	mac := hmac.New(sha256.New, []byte(stream.APISecret))
	mac.Write(params)
	return hex.EncodeToString(mac.Sum(nil))
}

func (stream *UserDataStream) Reconnect() error {
	var attempt int
	for attempt = 0; attempt < 3; attempt++ {
		stream.closeListenKey(stream.listenKey)
		fmt.Printf("Reconnection attempt %d of %d...\n", attempt+1, 3)

		// 获取新的 listenKey 并尝试连接
		if err := stream.getListenKey(); err != nil {
			fmt.Printf("Failed to get listenKey: %v\n", err)
			//time.Sleep(stream.reconnectInterval)
			continue
		}

		// 尝试建立 WebSocket 连接
		err := stream.base.Connect(StreamAPI + "?streams=" + stream.listenKey)
		if err != nil {
			fmt.Printf("Failed to reconnect WebSocket: %v\n", err)
			//time.Sleep(stream.reconnectInterval)
			continue
		}

		fmt.Println("Reconnected successfully")

		// 重新开始接收消息
		return nil
	}

	return fmt.Errorf("failed to reconnect after %d attempts", 3)
}

type OrderUpdate struct {
	RemainAmount       string `json:"A"`
	CreateTime         int64  `json:"O"`
	Side               int    `json:"S"` // 1= buy 2= sell
	RemainQuantity     string `json:"V"`
	Amount             string `json:"a"`
	TradeNo            string `json:"c"`
	OrderId            string `json:"i"`
	IsMaker            int    `json:"m"`
	OrderType          int    `json:"o"` // LIMIT_ORDER(1),POST_ONLY(2),IMMEDIATE_OR_CANCEL(3), FILL_OR_KILL(4),MARKET_ORDER(5);STOP_LIMIT(100)
	Price              string `json:"p"`
	Status             int    `json:"s"`
	Quantity           string `json:"v"`
	AvgPrice           string `json:"ap"`
	CumulativeQuantity string `json:"cv"`
	CumulativeAmount   string `json:"ca"`
}

func (d OrderUpdate) GetSymbol() string {
	return ""
}

func (stream *UserDataStream) OrderStream(ctx context.Context, channel chan<- types.OrderUpdateEntry) error {
	err := stream.base.SendMessage(map[string]any{
		"method": SubscribeOp,
		"params": []string{"spot@private.orders.v3.api"},
	})
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
			default:
				msg, err := stream.base.ReadMessage()
				if err != nil {
					continue
				}
				var resp StreamResp[OrderUpdate]

				_ = utils.Json.Unmarshal(msg, &resp)
				var status constants.OrderStatus = constants.Error

				switch resp.Data.Status {
				case 1:
					status = constants.Open
				case 2:
					status = constants.Filled
				case 3:
					status = constants.PartiallyFilled
				case 4:
					status = constants.Canceled
				case 5:
					status = constants.PartiallyCanceled
				}
				channel <- types.OrderUpdateEntry{
					OrderId:       resp.Data.OrderId,
					ClientOrderId: resp.Data.TradeNo,
					Status:        status,
				}
			}
		}
	}()
	return nil
}

type BalanceUpdate struct {
	Asset      string `json:"a"`
	Timestamp  int64  `json:"c"`
	Free       string `json:"f"`
	LockChange string `json:"ld"`
	Lock       string `json:"l"`
	FreeChange string `json:"fd"`
	Type       string `json:"o"`
}

func (BalanceUpdate) GetSymbol() string {
	return ""
}

func (stream *UserDataStream) BalanceStream(ctx context.Context, channel chan<- types.BalanceUpdateEntry) error {
	err := stream.base.SendMessage(map[string]any{
		"method": SubscribeOp,
		"params": []string{"spot@private.account.v3.api"},
	})
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := stream.base.ReadMessage()
				if err != nil {
					continue
				}
				var resp StreamResp[BalanceUpdate]

				_ = utils.Json.Unmarshal(msg, &resp)

				channel <- types.BalanceUpdateEntry{}
			}
		}
	}()
	return nil
}

type AccountUpdate struct {
	Price           string `json:"p"`
	Amount          string `json:"a"`
	IsSelfTrade     int    `json:"st"`
	Side            int    `json:"S"`
	TradeNo         string `json:"c"`
	TradeTime       int64  `json:"T"`
	TradeId         string `json:"t"`
	Volume          string `json:"v"`
	OrderId         string `json:"i"`
	IsMaker         int    `json:"m"`
	CommissionFee   string `json:"n"`
	CommissionAsset string `json:"N"`
}

func (AccountUpdate) GetSymbol() string {
	return ""
}

func (stream *UserDataStream) AccountStream(ctx context.Context, channel chan<- types.AccountUpdateEntry) error {
	err := stream.base.SendMessage(map[string]any{
		"method": SubscribeOp,
		"params": []string{"spot@private.account.v3.api"},
	})
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := stream.base.ReadMessage()
				if err != nil {
					continue
				}

				var resp StreamResp[AccountUpdate]
				_ = utils.Json.Unmarshal(msg, &resp)
				channel <- types.AccountUpdateEntry{}
			}
		}
	}()
	return nil
}

func NewUserStream(cred *platforms.Credentials) platforms.UserDataStreamer {
	return &UserDataStream{
		base:        platforms.NewStream(),
		Credentials: cred,
	}
}
