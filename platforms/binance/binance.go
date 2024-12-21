package binance

import (
	"github.com/xavierzho/go-cexs/platforms"
	"net/http"

	"github.com/xavierzho/go-cexs/constants"
	"github.com/xavierzho/go-cexs/types"
)

type Connector struct {
	*platforms.Credentials
	Client *http.Client
}

func NewConnector(base *platforms.Credentials, client *http.Client) *Connector {
	if client == nil {
		client = http.DefaultClient
	}
	return &Connector{Credentials: base, Client: client}
}

type AccountResp struct {
	MakerCommission  int    `json:"makerCommission"`
	BuyerCommission  int    `json:"buyerCommission"`
	CanWithdraw      bool   `json:"canWithdraw"`
	AccountType      string `json:"accountType"`
	SellerCommission int    `json:"sellerCommission"`
	UpdateTime       int    `json:"updateTime"`
	CanTrade         bool   `json:"canTrade"`
	Brokered         bool   `json:"brokered"`
	PreventSor       bool   `json:"preventSor"`
	Uid              int    `json:"uid"`
	Balances         []struct {
		Asset  string `json:"asset"`
		Free   string `json:"free"`
		Locked string `json:"locked"`
	} `json:"balances"`
	Permissions     []string `json:"permissions"`
	CommissionRates struct {
		Seller string `json:"seller"`
		Maker  string `json:"maker"`
		Taker  string `json:"taker"`
		Buyer  string `json:"buyer"`
	} `json:"commissionRates"`
	CanDeposit                 bool `json:"canDeposit"`
	TakerCommission            int  `json:"takerCommission"`
	RequireSelfTradePrevention bool `json:"requireSelfTradePrevention"`
}

func (c *Connector) AccountInfo() (*AccountResp, error) {
	var account = new(AccountResp)

	err := c.Call(http.MethodGet, AccountEndpoint, map[string]any{}, constants.Signed, account)
	return account, err
}

func (c *Connector) Balance(symbols []string) (map[string]types.BalanceEntry, error) {
	var result = make(map[string]types.BalanceEntry)
	accountInfo, err := c.AccountInfo()
	if err != nil {
		return nil, err
	}

	if len(symbols) == 0 {
		for _, balance := range accountInfo.Balances {
			result[balance.Asset] = types.BalanceEntry{
				Free:     balance.Free,
				Locked:   balance.Locked,
				Currency: balance.Asset,
			}
		}
	} else {
		// 如果指定了 symbols，返回指定的余额
		symbolSet := make(map[string]struct{}, len(symbols))
		for _, symbol := range symbols {
			symbolSet[symbol] = struct{}{}
		}

		for _, balance := range accountInfo.Balances {
			if _, exists := symbolSet[balance.Asset]; exists {
				result[balance.Asset] = types.BalanceEntry{
					Free:     balance.Free,
					Locked:   balance.Locked,
					Currency: balance.Asset,
				}
			}
		}
	}

	return result, nil
}

func (c *Connector) Name() constants.Platform {
	return constants.Binance
}

func (c *Connector) SymbolPattern(symbol string) string {
	symbol, _ = constants.StandardizeSymbol(symbol)
	return symbol
}
