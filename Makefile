.PHONY: build clean-abi-generated solc abigen

all: build abigen go-test

build:
	go run zk/build/main.go

clean-abi-generated:
	cd solidity && rm -fr ./abi/*

solc: clean-abi-generated
	cd solidity && solc --bin --abi -o ./abi contract_g16.sol

abigen: solc
	cd solidity && abigen --bin ./abi/Verifier.bin --abi abi/Verifier.abi --pkg solidity --out solidity_groth16.go --type Verifier

go-test: abigen
	cd solidity && go test

remixd:
	remixd -s ./ -u https://remix.ethereum.org

build-wasm: build
	cd wasm && GOOS=js GOARCH=wasm go build -o  ../assets/json.wasm

