package main

import "testing"

func TestReadJsonVPKey(t *testing.T) {
	keys, err := ReadJsonVPKey()
	if err != nil {
		t.Error(err)
	}
	for k := range keys {
		t.Log(k)
	}
}
