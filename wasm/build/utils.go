package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
	"math/big"
	"strconv"
)

var validation = validator.New()

func jsErr(err error, message string) string {
	if message == "" {
		return fmt.Sprintf("{'error': '%s','message': '%s'}", err.Error(), message)
	}
	if err == nil {
		return fmt.Sprintf("{'error': '%s'}", message)
	}
	return fmt.Sprintf("{'error': '%s'}", err.Error())
}

func parserHexToBigInt(arg string) (*big.Int, string) {
	if err := validation.Var(arg, "required,hexadecimal"); err != nil {
		return nil, jsErr(err, "Invalid argument input passed")
	}

	value := common.FromHex(arg)
	return new(big.Int).SetBytes(value), ""
}

func parserStringToInt(arg string) (int, string) {
	i, err := strconv.Atoi(arg)
	if err != nil {
		return -1, jsErr(err, "Invalid argument input passed")
	}
	return i, ""
}
