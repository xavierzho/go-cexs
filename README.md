## Multi-Crypto Exchange client


## Exchanges Supports

- [x] [Binance](https://developers.binance.com/docs/zh-CN/binance-spot-api-docs/rest-api/public-rest-api-for-binance)
- [x] [Bitmart](https://developer-pro.bitmart.com/en/quick)
- [x] [Bybit](https://bybit-exchange.github.io/docs/v5/intro)
- [x] [Okx](https://www.okx.com/docs-v5/)
- [x] [Mexc](https://mexcdevelop.github.io/apidocs/spot_v3_en/#introduction)
- [x] [Gate](https://www.gate.io/docs/developers/apiv4/ws/en/)
- [ ] etc....


# FAQ
## symbol, trading_pair format
All symbol formats are uppercase `{base}{quote}`, The converter is in [symbol.go](constants/symbol.go)


## The difference between balance stream and account stream
| Feature             | Balance Stream                         | Account Stream                      |
|---------------------|----------------------------------------|-------------------------------------|
| **Primary Focus**   | Non-transactional fund changes         | Transaction-related asset changes   |
| **Event Types**     | Deposits, withdrawals, transfers       | Order executions, fee deductions    |
| **Account Changes** | Direct increase/decrease in balance    | Asset conversion due to trading     |
| **Typical Scenario**| Transferring funds to futures account  | Buying BTC with USDT                |
| **Data Source**     | Balance update notifications           | Order update notifications          |
