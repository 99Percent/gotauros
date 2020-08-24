package taurosapi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"

	"github.com/99percent/gotauros/taurosapi"
)

var tauros taurosapi.TauAPI

func init() {
	in, err := ioutil.ReadFile("tokens.json")
	if err != nil {
		log.Fatalf("Unable to load tokens file tokens.json: %v", err)
	}
	if err := json.Unmarshal(in, &tauros); err != nil {
		log.Fatalf("Unable to unmarshall tokens file: %v", err)
	}
}

func TestCreateWebHook(t *testing.T) {
	id, err := tauros.test_GetWebHooks()
	if err != nil {
		t.Errorf("%v", err)
	}
	if id == 0 {
		t.Error("expected webhookid not zero")
	}
}
