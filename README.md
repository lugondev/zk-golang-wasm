# gnark-bid

This repo contains tests (interop or integration) that may drag some extra dependencies, for the following projects:

* [`gnark (v0.7.1)`: a framework to execute (and verify) algorithms in zero-knowledge](https://github.com/consensys/gnark) 
* [`gnark-crypto (v0.7.0)`: provides efficient cryptographic primitives](https://github.com/consensys/gnark-crypto)

## Solidity verifier

```bash
go generate
go test
```
or
```bash
make
```

It needs `solc` and `abigen` (1.10.17-stable).
