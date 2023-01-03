.PHONY: build clean-abi-generated solc abigen

all: build abigen go-test

build-bid:
	go run zk/build-bid/main.go
	cd solidity && abigen --bin ./abi/Contract_BidCircuit_sol_Verifier.bin --abi abi/Contract_BidCircuit_sol_Verifier.abi --pkg solidity --out solidity_Contract_BidCircuit.go --type BidCircuit

build:
	go run zk/build/main.go

clean-abi-generated:
	cd solidity && rm -fr ./abi/*

solc: clean-abi-generated
	cd solidity && solc --bin --abi -o ./abi *.sol

abigen: solc
	cd solidity && abigen --bin ./abi/Contract_MerkleCircuit_sol_Verifier.bin --abi abi/Contract_MerkleCircuit_sol_Verifier.abi --pkg solidity --out solidity_Contract_MerkleCircuit.go --type MerkleCircuit
	cd solidity && abigen --bin ./abi/Contract_PrivateValueCircuit_sol_Verifier.bin --abi abi/Contract_PrivateValueCircuit_sol_Verifier.abi --pkg solidity --out solidity_Contract_PrivateValueCircuit.go --type PrivateValueCircuit

abigen-merkle: solc
	cd solidity && abigen --bin ./abi/MerkleProof.bin --abi abi/MerkleProof.abi --pkg solidity --out MerkleProof.go --type MerkleProof

go-test: abigen
	cd solidity && go test

remixd:
	remixd -s ./ -u https://remix.ethereum.org

build-wasm: build
	cd wasm/build && GOOS=js GOARCH=wasm go build -o  ../../assets/json.wasm

