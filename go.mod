module discovery-node

go 1.12

require (
	github.com/ethereum/go-ethereum v1.8.23
	github.com/fatih/color v1.7.0
	github.com/gin-contrib/cors v0.0.0-20190301062745-f9e10995c85a
	github.com/gin-gonic/gin v1.3.0
	github.com/iden3/discovery-research/discovery-node v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.4.0
	github.com/spf13/viper v1.2.1
	github.com/stretchr/testify v1.2.2
	github.com/syndtr/goleveldb v0.0.0-20190318030020-c3a204f8e965
	github.com/urfave/cli v1.20.0
	github.com/vocdoni/go-dvote v0.0.0-20190318130547-148b652c8a49
)

replace github.com/ethereum/go-ethereum => ../../../go-ethereum

replace github.com/vocdoni/go-dvote => ../../../go-dvote

replace github.com/iden3/discovery-research/discovery-node => ./
