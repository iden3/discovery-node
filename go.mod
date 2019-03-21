module discovery-node

go 1.12

require (
	github.com/ethereum/go-ethereum v1.8.23
	github.com/fatih/color v1.7.0
	github.com/spf13/viper v1.2.1
	github.com/urfave/cli v1.20.0
	github.com/vocdoni/go-dvote v0.0.0-20190318130547-148b652c8a49
)

replace github.com/ethereum/go-ethereum => ../../go-ethereum

replace github.com/vocdoni/go-dvote => ../../go-dvote
