package taurosapi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"
)

var tauros TauApi

func init() {
	in, err := ioutil.ReadFile("tokens.json")
	if err != nil {
		log.Fatalf("Unable to load tokens file tokens.json: %v", err)
	}
	if err := json.Unmarshal(in, &tauros); err != nil {
		log.Fatalf("Unable to unmarshall tokens file: %v", err)
	}
}

func TestGetWebHooks(t *testing.T) {
	w, err := tauros.test_GetWebHooks()
	if err != nil {
		t.Errorf("%v", err)
	}
}
