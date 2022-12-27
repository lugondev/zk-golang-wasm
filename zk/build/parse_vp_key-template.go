package main

import (
	"gnark-bid/zk"
)

func ReadJsonVPKey() (*zk.VPKey, error) {
	jsonData := []byte(`{{.}}`)
	return zk.GetVPKey(jsonData)
}
