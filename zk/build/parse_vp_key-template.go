package main

import (
	"fmt"
	"gnark-bid/zk"
)

func readJsonVPKey() (map[string]*zk.VPKey, error) {
	jsonData := []byte(`{{.}}`)
	return zk.ParseMapData(jsonData)
}

func GetVPKey(name string) (*zk.VPKey, error) {
	keys, err := readJsonVPKey()
	if err != nil {
		return nil, err
	}
	vp := keys[name]
	if vp == nil {
		return nil, fmt.Errorf("no key for %s", name)
	}
	return vp, nil
}
