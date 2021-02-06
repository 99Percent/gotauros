package taurosapi

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"
)

// TauAPI - main interface struct
type TauAPI struct {
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
	URL       string `json:"url"`
	Email     string `json:"email"`
}

// TauWsObject - Tauros Websocket message "object"
type TauWsObject struct {
	Amount         string `json:"amount"`
	AmountPaid     string `json:"amount_paid"`
	AmountReceived string `json:"amount_received"`
	ClosedAt       string `json:"closed_at"`
	CreatedAt      string `json:"created_at"`
	FeeAmountPaid  string `json:"fee_amount_paid"`
	FeeDecimal     string `json:"fee_decimal"`
	FeePercent     string `json:"fee_percent"`
	Filled         string `json:"filled"`
	ID             int64  `json:"id"`
	InitialAmount  string `json:"initial_amount"`
	InitialValue   string `json:"initial_value"`
	IsOpen         bool   `json:"is_open"`
	LeftCoin       string `json:"left_coin"`
	Market         string `json:"market"`
	Price          string `json:"price"`
	RightCoin      string `json:"right_coin"`
	Side           string `json:"side"`
	Value          string `json:"value"`
}

// TauWebHookObject - Taures Webhook message "object"
type TauWebHookObject struct {
	Market              string `json:"market"`
	Side                string `json:"side"`
	InitialAmount       string `json:"initial_amount"`
	Filled              string `json:"filled"`
	Value               string `json:"value"`
	InitialValue        string `json:"initial_value"`
	Price               string `json:"price"`
	FeeDecimal          string `json:"fee_decimal"` //todo: issue to correct too much data and overlapping names
	FeePercent          string `json:"fee_percent"`
	FeeAmountPaid       string `json:"fee_amount_paid"`
	IsOpen              bool   `json:"is_open"`
	AmountPaid          string `json:"amount_paid"`
	AmountReceived      string `json:"amount_received"`
	CreatedAt           string `json:"created_at"`
	ClosedAt            string `json:"closed_at"`
	LeftCoin            string `json:"left_coin"`
	RightCoin           string `json:"right_coin"`
	LeftCoinIcon        string `json:"left_coin_icon"`
	RightCoinIcon       string `json:"right_coin_icon"`
	Sender              string `json:"sender"`
	Receiver            string `json:"receiver"`
	Coin                string `json:"coin"`
	CoinName            string `json:"coin_name"`
	CoinIcon            string `json:"coin_icon"`
	Amount              string `json:"amount"`
	TxID                string `json:"txId"` //todo: github issue correcting json format to "tx_id"
	Confirmed           bool   `json:"confirmed"`
	ConfirmedAt         string `json:"confirmed_at"`
	IsInnerTransfer     bool   `json:"is_innerTransfer"` //todo: issue to correct json name to is_inner_transfer
	Address             string `json:"address"`
	ExplorerLink        string `json:"explorer_link"`
	FeeAmount           string `json:"fee_amount"`
	TotalAmount         string `json:"total_amount"`
	Type                string `json:"type"`
	Description         string `json:"description"`
	TradeAmountPaid     string `json:"trade_amount_paid"`     //the actual amounts of the trade
	TradeAmountReceived string `json:"trade_amount_received"` // in order filled messages
	ID                  int64  `json:"id"`
	TransactionType     string `json:"transaction_type"`
	DateTime            string `json:"datetime"`
}

// NewOrder - new order data
type NewOrder struct {
	Market        string `json:"market"`
	Side          string `json:"side"`
	Amount        string `json:"amount"`
	Type          string `json:"type"`
	Price         string `json:"price"`
	IsAmountValue bool   `json:"is_amount_value,omitempty"`
}

// TauWsMessage - Tauros Websocket message header
type TauWsMessage struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Type        string      `json:"type"`
	Date        string      `json:"date"`
	Object      TauWsObject `json:"object"`
}

// TauWebHookMessage - Tauros POST message received via webhooks
type TauWebHookMessage struct { //todo: unify this with TauWsMessage
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Type        string           `json:"type"`
	Date        string           `json:"date"`
	Object      TauWebHookObject `json:"object"`
}

// Message - main message struct
type Message struct {
	ID            int64  `json:"id,omitempty"`
	Market        string `json:"market,omitempty"`
	Amount        string `json:"amount,omitempty"`
	Side          string `json:"side,omitempty"`
	Type          string `json:"type,omitempty"`
	Price         string `json:"price,omitempty"`
	IsAmountValue bool   `json:"is_amount_value,omitempty"`
	Email         string `json:"email,omitempty"`
	Password      string `json:"password,omitempty"`
}

// TransferMsg - json for direct Tauros Transfer
type TransferMsg struct {
	Nip       string  `json:"nip"`
	Coin      string  `json:"coin"`
	Recipient string  `json:"recipient"`
	Amount    float64 `json:"amount"`
}

// Order - order message struct
type Order struct {
	ID            int64       `json:"id"`       //for PlaceOrder
	OrderID       int64       `json:"order_id"` //for GetOpenOrders
	Market        string      `json:"market"`
	Side          string      `json:"side"`
	Amount        json.Number `json:"amount,Number"`
	InitialAmount json.Number `json:"initial_amount,Number"`
	Filled        json.Number `json:"filled,Number"`
	Value         json.Number `json:"value,Number"`
	InitialValue  json.Number `json:"initial_value,Number"`
	Price         json.Number `json:"price,Number"`
	FeeDecimal    json.Number `json:"fee_decimal"`
	CreatedAt     string      `json:"created_at"`
}

// MarketOrders - market orders (bids and asks) struct
type MarketOrders struct {
	Market string `json:"market"`
	Asks   []Order
	Bids   []Order
}

// Coin - available coins
type Coin struct {
	Coin                  string      `json:"coin"`
	MinWithdrawal         json.Number `json:"min_withdraw"`
	FeeWithdrawal         json.Number `json:"fee_withdraw"`
	Country               string      `json:"country"`
	ConfirmationsRequired int         `json:"confirmations_required"`
}

// Balance - available balances
type Balance struct {
	Coin     string `json:"coin"`
	CoinName string `json:"coin_name"`
	Address  string `json:"address"`
	Balances struct {
		Available json.Number `json:"available"`
		Pending   json.Number `json:"pending"`
		Frozen    json.Number `json:"frozen"`
		InOrders  json.Number `json:"in_orders"`
	} `json:"balances"`
}

// Webhook - data of a Webhook
type Webhook struct {
	ID                   int64  `json:"id"`
	Name                 string `json:"name"`
	Endpoint             string `json:"endpoint"`
	NotifyDeposit        bool   `json:"notify_deposit"`
	NotifyWithdrawal     bool   `json:"notify_withdrawal"`
	NotifyOrderPlaced    bool   `json:"notify_order_place"`
	NotifyOrderFilled    bool   `json:"notify order_filled"`
	NotifyTrade          bool   `json:"notify_trade"`
	AuthorizationHeader  string `json:"authorization_header"`
	AuthorizationContent string `json:"authorization_content"`
	IsActive             bool   `json:"is_active"`
	CreatedAt            string `json:"created_at"`
	UpdatedAt            string `json:"updated_at"`
	Detail               string `json:"detail"`
}

// Market - data of a market
type Market struct {
	Name      string      `json:"name"`
	MinAmount json.Number `json:"min_amount"`
	MaxAmount json.Number `json:"max_amount"`
	MinValue  json.Number `json:"min_value"`
	MaxValue  json.Number `json:"max_value"`
	MinPrice  json.Number `json:"min_price"`
	MaxPrice  json.Number `json:"max_price"`
	IsOpen    bool        `json:"is_open"`
}

//TauReq - request parameters
type TauReq struct {
	Version   int
	Method    string
	Path      string
	NeedsAuth bool
	PostMsg   []byte
}

// GetWebhooks - get all the registered webhooks
func (t *TauAPI) GetWebhooks() (webhooks []Webhook, error error) {
	var w = []Webhook{}
	var d struct {
		Count    int64     `json:"count"`
		Webhooks []Webhook `json:"results"`
		Detail   string    `json:"detail"`
	}
	jsonData, err := t.doTauRequest(&TauReq{
		Version:   2,
		Method:    "GET",
		Path:      "webhooks/webhooks",
		NeedsAuth: true,
	})
	if err != nil {
		return w, err
	}
	if err := json.Unmarshal(jsonData, &d); err != nil {
		return w, err
	}
	if d.Detail == "Invalid token." { //todo: really use http status code instead.
		return w, fmt.Errorf("Tauros API: %s", d.Detail)
	}
	return d.Webhooks, nil
}

// CreateWebhook - add a webhook
func (t *TauAPI) CreateWebhook(webhook Webhook) (ID int64, error error) {
	jsonPostMsg, _ := json.Marshal(webhook)
	jsonData, err := t.doTauRequest(&TauReq{
		Version:   2,
		Method:    "POST",
		Path:      "webhooks/webhooks",
		NeedsAuth: true,
		PostMsg:   jsonPostMsg,
	})
	if err != nil {
		return 0, nil
	}
	if string(jsonData) == "[\"Limit reached\"]" {
		return 0, fmt.Errorf("Limit of webhooks reached (5)")
	}
	var d struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(jsonData, &d); err != nil {
		return 0, fmt.Errorf("CreateWebhook -> unmarshal jsonData %v", err)
	}
	return d.ID, nil
}

// DeleteWebhook - delete one webhook according to the webhook ID
func (t *TauAPI) DeleteWebhook(ID int64) error {
	_, err := t.doTauRequest(&TauReq{
		Version:   2,
		Method:    "DELETE",
		Path:      "webhooks/webhooks/" + strconv.FormatInt(ID, 10),
		NeedsAuth: true,
	})
	return err
}

// DeleteWebhooks - delete all currently registered webhooks
func (t *TauAPI) DeleteWebhooks() (error error) {
	webhooks, err := t.GetWebhooks()
	if err != nil {
		return err
	}
	for _, w := range webhooks {
		err := t.DeleteWebhook(w.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetCoins - get all available coins handled by the exchange
func (t *TauAPI) GetCoins() (coins []Coin, error error) {
	var c = []Coin{}
	var d struct {
		Crypto []Coin `json:"cryto"` //typo from api
		Fiat   []Coin `json:"fiat"`
	}
	jsonData, err := t.doTauRequest(&TauReq{
		Version: 2,
		Method:  "GET",
		Path:    "coins",
	})
	if err != nil {
		return c, err
	}
	if err := json.Unmarshal(jsonData, &d); err != nil {
		return c, err
	}
	c = append(d.Crypto, d.Fiat...)
	return c, nil
}

// GetMarkets - get current available markets
func (t *TauAPI) GetMarkets() (markets []Market, error error) {
	var m []Market
	jsonData, err := t.doTauRequest(&TauReq{
		Version: 2,
		Method:  "GET",
		Path:    "trading/markets",
	})
	if err != nil {
		return nil, fmt.Errorf("TauGetMarkets ->%v", err)
	}
	if err := json.Unmarshal(jsonData, &m); err != nil {
		return nil, fmt.Errorf("TauGetMarkets json.Unmarshall->%v", err)
	}
	return m, nil
}

// GetMarketOrders - get current market orders for one market
func (t *TauAPI) GetMarketOrders(market string) (MarketOrders, error) {
	var mo MarketOrders
	jsonData, err := t.doTauRequest(&TauReq{
		Version: 1,
		Method:  "GET",
		Path:    "trading/orders?market=" + strings.ToLower(market),
	})
	if err != nil {
		return mo, fmt.Errorf("TauGetMarketOrders ->%s", err.Error())
	}
	if err := json.Unmarshal(jsonData, &mo); err != nil {
		return mo, err
	}
	return mo, nil
}

// GetBalances - get available balances of the user
func (t *TauAPI) GetBalances() (balances []Balance, error error) {
	var b []Balance
	var w struct {
		Wallets []Balance `json:"wallets"`
	}
	jsonData, err := t.doTauRequest(&TauReq{
		Version:   1,
		Method:    "GET",
		Path:      "data/listbalances",
		NeedsAuth: true,
	})
	if err != nil {
		return b, err
	}
	if err := json.Unmarshal(jsonData, &w); err != nil {
		return b, err
	}
	return w.Wallets, nil
}

// GetDepositAddress - get the deposit address of the user for the specified coin
func (t *TauAPI) GetDepositAddress(coin string) (address string, error error) {
	jsonData, err := t.doTauRequest(&TauReq{
		Version:   1,
		Method:    "GET",
		Path:      "data/getdepositaddress?coin=" + coin,
		NeedsAuth: true,
	})
	if err != nil {
		return "", fmt.Errorf("TauDepositAddress-> %v", err)
	}
	var d struct {
		Coin    string `json:"coin"`
		Address string `json:"address"`
	}
	if err := json.Unmarshal(jsonData, &d); err != nil {
		return "", fmt.Errorf("TauDepositAddress-> %v", err)
	}
	return d.Address, nil
}

// PlaceOrder - add a new order
func (t *TauAPI) PlaceOrder(newOrder NewOrder) (Order, error) {
	jsonPostMsg, _ := json.Marshal(newOrder)
	jsonData, err := t.doTauRequest(&TauReq{
		Version:   1,
		Method:    "POST",
		Path:      "trading/placeorder",
		NeedsAuth: true,
		PostMsg:   jsonPostMsg,
	})
	//log.Printf("raw json of placeorder: %s", string(jsonData))
	var o Order
	if err != nil {
		return o, fmt.Errorf("PlaceOrder-> %v", err)
	}
	if err := json.Unmarshal(jsonData, &o); err != nil {
		return o, fmt.Errorf("PlaceOrder-> unmarshal jsonData %v", err)
	}
	return o, nil
}

// GetOpenOrders - get all open orders by the user
func (t *TauAPI) GetOpenOrders() (orders []Order, error error) {
	jsonData, err := t.doTauRequest(&TauReq{
		Version:   1,
		Method:    "GET",
		Path:      "trading/myopenorders",
		NeedsAuth: true,
	})
	if err != nil {
		return nil, fmt.Errorf("GetOpenOrders->%v", err)
	}
	if err := json.Unmarshal(jsonData, &orders); err != nil {
		return nil, fmt.Errorf("GetOpenOrders->%v", err)
	}
	return orders, nil
}

// CloseAllOrders - close all currently open orders
func (t *TauAPI) CloseAllOrders() error {
	orders, err := t.GetOpenOrders()
	if err != nil {
		return fmt.Errorf("CloseAllOrders ->%v", err)
	}
	for _, o := range orders {
		if err := t.CloseOrder(o.OrderID); err != nil {
			return fmt.Errorf("CloseAllOrders Deleting Order %d ->%v", o.ID, err)
		}
	}
	return nil
}

// CloseOrder - close the order specified by the order ID
func (t *TauAPI) CloseOrder(orderID int64) error {
	jsonPostMsg, _ := json.Marshal(Message{ID: orderID})
	_, err := t.doTauRequest(&TauReq{
		Version:   1,
		Method:    "POST",
		Path:      "trading/closeorder",
		NeedsAuth: true,
		PostMsg:   jsonPostMsg,
	})
	if err != nil {
		return fmt.Errorf("CloseOrder->%v", err)
	}
	return nil
}

// Login - simulate a login to get the jwt token
func (t *TauAPI) Login(email string, password string) (jwtToken string, err error) {
	jsonPostMsg, _ := json.Marshal(&Message{Email: email, Password: password})
	jsonData, err := t.doTauRequest(&TauReq{
		Version:   2,
		Method:    "POST",
		Path:      "auth/signin",
		NeedsAuth: false,
		PostMsg:   jsonPostMsg,
	})
	if err != nil {
		return "", fmt.Errorf("Login->%v", err)
	}
	var d struct {
		Token     string `json:"token"`
		TwoFactor bool   `json:"two_factor"`
	}
	if err := json.Unmarshal(jsonData, &d); err != nil {
		return "", fmt.Errorf("Login->%v", err)
	}
	return d.Token, nil
}

// Transfer - direct transfer of funds to another Tauros account
func (t *TauAPI) Transfer(transfer TransferMsg) error {
	jsonPostMsg, _ := json.Marshal(&transfer)
	_, err := t.doTauRequest(&TauReq{
		Version:   2,
		Method:    "POST",
		Path:      "wallets/inner-transfer",
		NeedsAuth: true,
		PostMsg:   jsonPostMsg,
	})
	if err != nil {
		return fmt.Errorf("Transfer->%v", err)
	}
	return nil //no need to see return post
}

func (t *TauAPI) doTauRequest(tauReq *TauReq) (msgdata json.RawMessage, e error) {
	var httpReq *http.Request
	var signatureDebugInfo string
	var err error
	apiVersion := fmt.Sprintf("v%1d", tauReq.Version)
	if tauReq.NeedsAuth {
		tauReq.Path += "/"
	}
	httpReq, err = http.NewRequest(tauReq.Method, t.URL+"/api/"+apiVersion+"/"+tauReq.Path, bytes.NewBuffer(tauReq.PostMsg))
	httpReq.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("doTauRequest-> Error on http.NewRequest: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	if tauReq.NeedsAuth {

		var postMsg string
		var nonce string
		var path string
		var message string
		var messageHash [32]byte
		var decodedAPISecret []byte

		postMsg = string(tauReq.PostMsg)
		if postMsg == "" {
			postMsg = "{}"
		}
		nonce = strconv.FormatInt(time.Now().UnixNano()/1e6, 10) //todo: check divide by time.millisecond
		path = "/api/" + apiVersion + "/" + tauReq.Path          //trailing backslash must be added at each post request in path
		message = nonce + tauReq.Method + path + postMsg
		messageHash = sha256.Sum256([]byte(message))
		if d, err := base64.StdEncoding.DecodeString(t.APISecret); err != nil {
			return nil, fmt.Errorf("doTauRequest -> Error decoding APiSecret base64: %v", err)
		} else {
			decodedAPISecret = d
		}
		h := hmac.New(sha512.New, decodedAPISecret)
		h.Write(messageHash[:])
		signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
		httpReq.Header.Set("Authorization", "Bearer "+t.APIKey)
		httpReq.Header.Set("Taur-Nonce", nonce)
		httpReq.Header.Set("Taur-Signature", signature)

		dumpRequest, _ := httputil.DumpRequestOut(httpReq, true)
		signatureDebugInfo = fmt.Sprintf("\nPost message: %q", string(tauReq.PostMsg)) +
			fmt.Sprintf("\nPost Message body: %s", postMsg) +
			fmt.Sprintf("\nNonce=%s", nonce) +
			fmt.Sprintf("\nmessage: %s", message) +
			fmt.Sprintf("\nmessageHash: %x", messageHash) +
			fmt.Sprintf("\nApi_key: %s", t.APIKey) +
			fmt.Sprintf("\nApi_secret: %s", t.APISecret) +
			fmt.Sprintf("\nDecoded b64 api_secret: %x", decodedAPISecret) +
			fmt.Sprintf("\nsignature=%s", signature) +
			fmt.Sprintf("\nbytes.NewBuffer(tauReq.PostMsg: %s", bytes.NewBuffer(tauReq.PostMsg).String()) +
			fmt.Sprintf("\nDump Request [%s]\n", dumpRequest)
	}

	client := http.Client{Timeout: time.Second * 3}
	start := time.Now()
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("Elapsed: %s | doTauRequest-> Error reading response: %v", time.Since(start), err)
	}
	defer resp.Body.Close()
	//todo: check StatusCode
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("doTauRequest-> Error ioutil body: %v", err)
	}
	if strings.Contains(tauReq.Path, "webhooks") { //needed because the webhook endpoints are missing this header
		if len(body) == 0 {
			body = []byte("{}") //for spurious empty DELETE responses, bug from api
		}
		body = []byte(`{"success":true,"payload":` + string(body) + "}")
	}
	var respJSON struct {
		Success bool            `json:"success"`
		Message json.RawMessage `json:"msg"`
		Data    json.RawMessage `json:"data"`
		Payload json.RawMessage `json:"payload"`
	}
	if err := json.Unmarshal(body, &respJSON); err != nil {
		return nil, fmt.Errorf("doTauRequest-> Unmarshal error: %v \n resp code: %d, body=%s", err, resp.StatusCode, body)
	}
	if !respJSON.Success {
		msg := string(respJSON.Message)
		if msg == "" {
			msg = string(body)
		}
		if strings.Contains(msg, "Invalid token") {
			msg += "Authorization: Bearer " + t.APIKey
			msg += signatureDebugInfo
		}
		if strings.Contains(msg, "signature") {
			msg += signatureDebugInfo
		}
		return nil, fmt.Errorf(msg)
	}
	if tauReq.Version == 1 {
		return respJSON.Data, err
	}
	return respJSON.Payload, err
}
