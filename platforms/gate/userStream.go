package gate

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha512"
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

func (u *UserDataStream) Login() error {
	err := u.Connect(StreamAPI)
	if err != nil {
		return err
	}
	return nil
}
func (u *UserDataStream) sign(channel, event string, timestamp int64) map[string]any {
	buf := bytes.NewBufferString(fmt.Sprintf("channel=%s&event=%s&timestamp=%d", channel, event, timestamp))
	mac := hmac.New(sha512.New, []byte(u.APISecret))
	mac.Write(buf.Bytes())
	return map[string]any{
		"method":   "api_key",
		KeyHeader:  u.APIKey,
		SignHeader: hex.EncodeToString(mac.Sum(nil)),
	}
}

type OrderUpdate []Order

func (u *UserDataStream) OrderStream(ctx context.Context, channels chan<- types.OrderUpdateEntry) error {
	const channel = "spot.usertrades"
	const event = "subscribe"
	var timestamp = time.Now().Unix()
	err := u.SendMessage(map[string]any{
		"id":      uuid.New().ID(),
		"time":    timestamp,
		"channel": channel,
		"event":   event,
		"auth":    u.sign(channel, event, timestamp),
	})
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				u.Close()
				return
			default:
				msg, err := u.ReadMessage()
				if err != nil {
					continue
				}
				var event Event[OrderUpdate]
				_ = utils.Json.Unmarshal(msg, &event)

				for _, order := range event.Result {
					channels <- types.OrderUpdateEntry{
						OrderId:       order.ID,
						Status:        OrderStatus(order.Status).Convert(),
						ClientOrderId: order.Text,
					}
				}
			}
		}
	}()
	return nil
}

type BalanceUpdate struct {
	Total        string `json:"total"`
	Freeze       string `json:"freeze"`
	Change       string `json:"change"`
	FreezeChange string `json:"freeze_change"`
	Available    string `json:"available"`
	TimestampMs  string `json:"timestamp_ms"`
	Currency     string `json:"currency"`
	ChangeType   string `json:"change_type"`
	User         string `json:"user"`
	Timestamp    string `json:"timestamp"`
}

func (u *UserDataStream) BalanceStream(ctx context.Context, channels chan<- types.BalanceUpdateEntry) error {
	const channel = "spot.balances"
	const event = "subscribe"
	var timestamp = time.Now().Unix()
	err := u.SendMessage(map[string]any{
		"id":      uuid.New().ID(),
		"time":    timestamp,
		"channel": channel,
		"event":   event,
		"auth":    u.sign(channel, event, timestamp),
	})
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				u.Close()
				return
			default:
				msg, err := u.ReadMessage()
				if err != nil {
					continue
				}
				var event Event[[]BalanceUpdate]
				_ = utils.Json.Unmarshal(msg, &event)

				for _, update := range event.Result {
					channels <- types.BalanceUpdateEntry{}
					fmt.Println(update)
				}

			}
		}
	}()
	return nil
}

func (u *UserDataStream) AccountStream(ctx context.Context, channels chan<- types.AccountUpdateEntry) error {
	return fmt.Errorf("not support stream")
}

func NewUserStream(cred *platforms.Credentials) platforms.UserDataStreamer {
	return &UserDataStream{
		StreamBase:  platforms.NewStream(),
		Credentials: cred,
	}
}
