.PHONY: build clean-abi-generated solc abigen

all: build abigen go-test

build-bid-circuit:
	go run zk/build-bid/main.go

build-bid: build-bid-circuit abigen-bid build-wasm

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

abigen-bid: solc
	cd solidity && abigen --bin ./abi/Contract_BiddingCircuit_sol_Verifier.bin --abi abi/Contract_BiddingCircuit_sol_Verifier.abi --pkg solidity --out solidity_Contract_BiddingCircuit.go --type BiddingCircuit

go-test: abigen
	cd solidity && go test

remixd:
	remixd -s ./ -u https://remix.ethereum.org

build-wasm:
	cd wasm/build && GOOS=js GOARCH=wasm go build -o  ../../assets/json.wasm

