package node

import (
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fatih/color"
	"github.com/iden3/discovery-node/config"
	"github.com/iden3/discovery-node/db"
	"github.com/iden3/discovery-node/discovery"
	"github.com/iden3/discovery-node/utils"
	swarm "github.com/vocdoni/go-dvote/net/swarm"
)

type NodeSrv struct {
	db *db.Db
	sn *swarm.SimplePss
}

func RunNode() (*NodeSrv, error) {
	fmt.Println("initializing node")

	sto, err := db.New("./tmp/iddb")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	sn := new(swarm.SimplePss)

	// set a random privK
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	// set private key
	sn.Key = privateKey

	// set swarm.sn configuration
	sn.Ports = swarm.NewSwarmPorts()
	sn.Datadir = config.C.Datadir
	sn.Ports.WebSockets = config.C.Ports.WebSockets
	sn.Ports.HTTPRPC = config.C.Ports.HTTPRPC
	sn.Ports.Bzz = config.C.Ports.Bzz
	sn.Ports.P2P = config.C.Ports.P2P
	sn.LogLevel = config.C.Pss.LogLevel

	publicKey := privateKey.Public()
	fmt.Print("pubK: ")
	color.Cyan(utils.PublicKeyToString(publicKey))

	err = sn.Init()
	if err != nil {
		color.Red(err.Error())
		os.Exit(0)
	}
	sn.PssSub(config.C.Pss.Kind, config.C.Pss.Key, config.C.Pss.Topic, "")
	defer sn.PssTopics[config.C.Pss.Topic].Unregister()

	fmt.Println("pubK", sn.PssPubKey)

	go func() {
		for {
			pmsg := <-sn.PssTopics[config.C.Pss.Topic].Delivery
			fmt.Print("[MSG RECEIVED]: ")
			color.Yellow(string(pmsg.Msg))
		}
	}()

	node := &NodeSrv{
		db: sto,
		sn: sn,
	}

	return node, nil
}

func (node *NodeSrv) StoreId(id discovery.Id) error {
	// node.db.Put(id.IdAddr.Bytes(), )
	return nil
}
