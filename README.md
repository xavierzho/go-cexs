## Multi-Crypto Exchange client

etc...




## bybit

### subscribe
```go
type Response struct {
	Topic string `json:"topic"`
	Data any `json:"data"`
	Timestamp int64 `json:"ts"`
	Type string `json:"type"`
}
//{"symbol":"PEPEUSDT","lastPrice":"0.000020097","highPrice24h":"0.000023461","lowPrice24h":"0.000018806","prevPrice24h":"0.000022684","volume24h":"17445127596312","turnover24h":"372320928.283501236","price24hPcnt":"-0.1140","usdIndexPrice":"0.0000200946"}
// Ticker reference: https://bybit-exchange.github.io/docs/v5/websocket/public/ticker
type Ticker struct {
    Symbol    string    `json:"symbol"`
    UsdIndexPrice    string    `json:"usdIndexPrice"`
    LowPrice24h    string    `json:"lowPrice24h"`
    PrevPrice24h    string    `json:"prevPrice24h"`
    Volume24h    string    `json:"volume24h"`
    Price24hPcnt    string    `json:"price24hPcnt"`
    HighPrice24h    string    `json:"highPrice24h"`
    Turnover24h    string    `json:"turnover24h"`
    LastPrice    string    `json:"lastPrice"`
}

//{"s":"PEPEUSDT","b":[["0.000020097","86076321"]],"a":[],"u":22834042,"seq":74283820551}
// OrderBook reference: https://bybit-exchange.github.io/docs/v5/websocket/public/orderbook
type OrderBook struct {
    Asks    [][]string    `json:"a"`
    Bids   [][]string `json:"b"`
    Symbol    string    `json:"s"`
    UpdateId    int    `json:"u"`
    Seq    int64    `json:"seq"`
}

//[{
//"T": 1672304486865,
//"s": "BTCUSDT",
//"S": "Buy",
//"v": "0.001",
//"p": "16578.50",
//"L": "PlusTick",
//"i": "20f43950-d8dd-5b31-9112-a178eb6023af",
//"BT": false
//}
//]

type Trade struct {
    Price    string    `json:"p"`
    IsBlockTrade    bool    `json:"BT"`
    Symbol    string    `json:"s"`
    Side    string    `json:"S"`
    Timestamp    int64    `json:"T"`
    Volume    string    `json:"v"`
    Id    string    `json:"i"`
    L    string    `json:"L"`
}
```


## info
Binance stream 有spot和features，



```go
{"stream":"l2jPv9zZjYNdt8gF5QJRwhx2pImg7bp6W9LZD1ngrVH5BPtIuXt3IWjxwX03","data":{"e":"executionReport","E":1734347791772,"s":"VTHOUSDT","c":"web_4a04001897314cd384cab9bdc49a242a","S":"SELL","o":"MARKET","f":"GTC","q":"2000.00000000","p":"0.00000000","P":"0.00000000","F":"0.00000000","g":-1,"C":"","x":"NEW","X":"NEW","r":"NONE","i":486036557,"l":"0.00000000","z":"0.00000000","L":"0.00000000","n":"0","N":null,"T":1734347791772,"t":-1,"I":1025300685,"w":true,"m":false,"M":false,"O":1734347791772,"Z":"0.00000000","Y":"0.00000000","Q":"0.00000000","W":1734347791772,"V":"EXPIRE_MAKER"}}

{"e":"kline","E":1734452294014,"s":"BTCUSDT","k":{"t":1734452280000,"T":1734452339999,"s":"BTCUSDT","i":"1m","f":4288916204,"L":4288918697,"o":"106927.27000000","c":"106861.14000000","h":"106932.00000000","l":"106857.97000000","v":"12.90944000","n":2494,"x":false,"q":"1379809.39927630","V":"1.36987000","Q":"146411.52169050","B":"0"}}

{
"e": "depthUpdate", 
"E": 1672515782136, 
"s": "BNBBTC",      
"U": 157,           
"u": 160,           
"b": [              
[
"0.0024",       
"10",           
           
]
],
"a": [             
[
"0.0026",     
"100",          
]
]
}

{"symbol":"VTHOUSDT","orderId":487067425,"orderListId":-1,"clientOrderId":"8b7b29ee-e84d-4c76-a3ed-84f379276c80","transactTime":1734515352605,"price":"0.00320000","origQty":"2463.00000000","executedQty":"0.00000000","origQuoteOrderQty":"0.00000000","cummulativeQuoteQty":"0.00000000","status":"NEW","timeInForce":"GTC","type":"LIMIT","side":"SELL","workingTime":1734515352605,"fills":[],"selfTradePreventionMode":"EXPIRE_MAKER"}

```
```json

{
"symbol": "LTCBTC",
"origClientOrderId": "myOrder1",
"orderId": 4,
"orderListId": -1, 
"clientOrderId": "cancelMyOrder1",
"transactTime": 1684804350068,
"price": "2.00000000",
"origQty": "1.00000000",
"executedQty": "0.00000000",
"cummulativeQuoteQty": "0.00000000",
"status": "CANCELED",
"timeInForce": "GTC",
"type": "LIMIT",
"side": "BUY",
"selfTradePreventionMode": "NONE"
}
```
