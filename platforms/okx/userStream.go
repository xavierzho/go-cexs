package okx

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/xavierzho/go-cexs/platforms"
	"github.com/xavierzho/go-cexs/types"
	"github.com/xavierzho/go-cexs/utils"
	"strconv"
	"time"
)

type UserDataStream struct {
	*platforms.Credentials
	*platforms.StreamBase
}

func (stream *UserDataStream) Sign(timestamp int64) string {
	mac := hmac.New(sha256.New, []byte(stream.APISecret))
	mac.Write([]byte(fmt.Sprintf("%dGET/users/self/verify", timestamp)))
	return hex.EncodeToString(mac.Sum(nil))
}

type LoginReturn struct {
	Event  string `json:"event"`
	Code   string `json:"code"`
	Msg    string `json:"msg"`
	ConnId string `json:"connId"`
}

func (stream *UserDataStream) Login() error {
	err := stream.Connect(StreamAPI + PrivateChannel)
	if err != nil {
		return err
	}
	timestamp := time.Now().Unix()
	err = stream.SendMessage(map[string]any{
		"op": "login",
		"args": []map[string]any{
			{
				"apiKey":     stream.APIKey,
				"passphrase": *stream.Option,
				"timestamp":  strconv.Itoa(int(timestamp)),
				"sign":       stream.Sign(timestamp),
			},
		},
	})
	if err != nil {
		return err
	}

	data, err := stream.ReadMessage()
	if err != nil {
		return err
	}
	var res LoginReturn
	_ = utils.Json.Unmarshal(data, &res)
	if res.Event == "error" {
		return fmt.Errorf("[Okx stream] code: %s, msg: %s", res.Code, res.Msg)
	}
	return nil
}

func (stream *UserDataStream) OrderStream(ctx context.Context, channel chan<- types.OrderUpdateEntry) error {
	//TODO implement me
	panic("implement me")
}

func (stream *UserDataStream) BalanceStream(ctx context.Context, channel chan<- types.BalanceUpdateEntry) error {
	//TODO implement me
	panic("implement me")
}

func (stream *UserDataStream) AccountStream(ctx context.Context, channel chan<- types.AccountUpdateEntry) error {
	//TODO implement me
	panic("implement me")
}

func NewUserStream(cred *platforms.Credentials) platforms.UserDataStreamer {
	return &UserDataStream{
		StreamBase:  platforms.NewStream(),
		Credentials: cred,
	}
}
