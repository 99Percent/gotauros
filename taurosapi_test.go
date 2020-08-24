package taurosapi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"
)

var tauros TauAPI
var webhookID int64

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
	coins, err := tauros.GetCoins()
	if err != nil {
		t.Errorf("%v", err)
	}
	if len(coins) == 0 {
		t.Error("coins available is zero")
	}
}

func TestGetMarkets(t *testing.T) {
	markets, err := tauros.GetMarkets()
	if err != nil {
		t.Errorf("%v", err)
	}
	if len(markets) == 0 {
		t.Error("markets available is zero")
	}
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
