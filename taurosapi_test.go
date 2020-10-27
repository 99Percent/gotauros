package taurosapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

var tauros TauAPI
var webhookID int64
var coins []Coin
var balances []Balance
var markets []Market
var order Order
var err error

func init() {
	in, err := ioutil.ReadFile("tokens.json")
	if err != nil {
		log.Fatalf("Unable to load tokens file tokens.json: %v", err)
	}
	if err := json.Unmarshal(in, &tauros); err != nil {
		log.Fatalf("Unable to unmarshall tokens file: %v", err)
	}
}

func TestGetCoins(t *testing.T) {
	coins, err = tauros.GetCoins()
	if err != nil {
		t.Errorf("%v", err)
	}
	if len(coins) == 0 {
		t.Error("no coins available returned")
	}
	//todo: show all coins
}

func TestGetMarkets(t *testing.T) {
	markets, err = tauros.GetMarkets()
	if err != nil {
		t.Errorf("%v", err)
	}
	if len(markets) == 0 {
		t.Error("markets available is zero")
	}
	//todo: show all markets
}

func TestDeleteWebhooks(t *testing.T) {
	if err := tauros.DeleteWebhooks(); err != nil {
		t.Errorf("%v", err)
	}
}

func TestCreateWebhook(t *testing.T) {
	webhookID, err := tauros.CreateWebhook(Webhook{
		Name:              "MyWebhook",
		Endpoint:          "https://somendpoint.com",
		NotifyDeposit:     true,
		NotifyWithdrawal:  true,
		NotifyOrderPlaced: false,
		NotifyOrderFilled: true,
		NotifyTrade:       true,
		IsActive:          true,
	})
	if err != nil {
		t.Errorf("%v", err)
		t.SkipNow()
	}
	if webhookID == 0 {
		t.Errorf("expected webhook id to be not zero")
	}
}

func TestDeleteWebhook(t *testing.T) {
	if err := tauros.DeleteWebhook(webhookID); err != nil {
		t.Errorf("%v", err)
	}
}

func TestCloseAllOrders(t *testing.T) {
	if err = tauros.CloseAllOrders(); err != nil {
		t.Errorf("%v", err)
	}
}

func TestGetBalances(t *testing.T) {
	balances, err = tauros.GetBalances()
	if err != nil {
		t.Errorf("%v", err)
	}
	if len(balances) == 0 {
		t.Error("no balances available returned")
	}
	//todo: show non zero balances
}

func TestTransfer(t *testing.T) {
	err = tauros.Transfer(TransferMsg{
		Recipient: "david@montebit.com",
		Coin:      "MXN",
		Nip:       "119744",
		Amount:    1.99,
	})
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestPlaceOrder(t *testing.T) {
	var available float64
	for _, b := range balances {
		if b.Coin == "BTC" {
			available, _ = b.Balances.Available.Float64()
			break
		}
	}
	if !(available > 0.001) {
		t.Log("no BTC balance available to test placeorder func")
		t.SkipNow()
	}
	if order, err = tauros.PlaceOrder(NewOrder{
		Market: "BTC-MXN",
		Side:   "sell",
		Amount: fmt.Sprintf("%.8f", available*0.1),
		Price:  "250000.0",
		Type:   "limit",
	}); err != nil {
		t.Errorf("Unable to place order: %v", err)
	}
	if (order.ID == 0) || (order.ID < 0) {
		t.Errorf("place Order ID returned zero or negative")
	}
	t.Logf("New Order placed: %+v", order)
}

func TestCloseOrder(t *testing.T) {
	if order.ID == 0 {
		t.Log("order not placed, skipping test")
		t.SkipNow()
	}
	if err = tauros.CloseOrder(order.ID); err != nil {
		t.Errorf("%v", err)
	}
}
