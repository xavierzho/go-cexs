package bybit

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"github.com/xavierzho/go-cexs/platforms"
	"github.com/xavierzho/go-cexs/types"
	"github.com/xavierzho/go-cexs/utils"
	"time"
)

type UserDataStream struct {
	*platforms.Credentials
	*platforms.StreamBase
}
type DataStream[T fmt.Stringer] struct {
	ID           string `json:"id"`
	Topic        string `json:"topic"`
	CreationTime int64  `json:"creationTime"`
	Data         T      `json:"data"`
}

func (stream *UserDataStream) Sign(expires int64) string {
	mac := hmac.New(sha256.New, []byte(stream.APISecret))
	mac.Write([]byte(fmt.Sprintf("GET/realtime%d", expires)))
	return hex.EncodeToString(mac.Sum(nil))
}
func (stream *UserDataStream) Login() error {
	err := stream.Connect(fmt.Sprintf("%s%s", StreamAPI, PrivateChannel))
	if err != nil {
		return err
	}
	expires := time.Now().Add(time.Second * 10)
	exp := expires.UnixNano() / 1e6
	err = stream.SendMessage(map[string]any{
		"req_id": uuid.New().String(), // optional
		"op":     "auth",
		"args": []any{
			stream.APIKey,
			exp, // expires; is greater than your current timestamp
			stream.Sign(exp),
		},
	})
	if err != nil {
		return err
	}
	_, err = stream.ReadMessage()
	return err
}

type OrderEvent []OrderInfo

func (OrderEvent) String() string {
	return ""
}
func (stream *UserDataStream) OrderStream(ctx context.Context, channel chan<- types.OrderUpdateEntry) error {
	err := stream.SendMessage(map[string]any{
		"op":   "subscribe",
		"args": []string{"order"},
	})
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				_ = stream.Close()
				return
			default:
				msg, err := stream.ReadMessage()
				if err != nil {
					continue
				}
				var event DataStream[OrderEvent]
				_ = utils.Json.Unmarshal(msg, &event)
				for _, e := range event.Data {
					channel <- types.OrderUpdateEntry{
						OrderId:       e.OrderId,
						ClientOrderId: e.OrderLinkId,
						Status:        OrderStatus(e.OrderStatus).Convert(),
					}
				}

			}
		}
	}()
	return nil
}

type BalanceEvent []WalletAccountInfo

func (BalanceEvent) String() string {
	return ""
}

func (stream *UserDataStream) BalanceStream(ctx context.Context, channel chan<- types.BalanceUpdateEntry) error {
	err := stream.SendMessage(map[string]any{})
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				_ = stream.Close()
				return
			default:
				msg, err := stream.ReadMessage()
				if err != nil {
					continue
				}
				var event DataStream[BalanceEvent]
				_ = utils.Json.Unmarshal(msg, &event)
				channel <- types.BalanceUpdateEntry{}
			}
		}
	}()

	return nil
}

func (stream *UserDataStream) AccountStream(ctx context.Context, channel chan<- types.AccountUpdateEntry) error {
	return fmt.Errorf("not support this method")
}

func NewUserStream(cred *platforms.Credentials) platforms.UserDataStreamer {
	return &UserDataStream{
		StreamBase:  platforms.NewStream(),
		Credentials: cred,
	}
}
