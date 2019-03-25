package node

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fatih/color"
	"github.com/iden3/discovery-node/config"
	"github.com/iden3/discovery-node/db"
	"github.com/iden3/discovery-node/discovery"
	"github.com/iden3/discovery-node/utils"
	swarm "github.com/vocdoni/go-dvote/net/swarm"
)

type NodeSrv struct {
	discsrv discovery.DiscoveryService
	db      *db.Db
	sn      *swarm.SimplePss
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

	err = sn.Init()
	if err != nil {
		color.Red(err.Error())
		os.Exit(0)
	}
	sn.PssSub(config.C.Pss.Kind, config.C.Pss.Key, config.C.Pss.Topic, "")
	defer sn.PssTopics[config.C.Pss.Topic].Unregister()

	fmt.Print("pubK: ")
	color.Cyan(utils.PublicKeyToString(publicKey))
	// fmt.Println("pubK", sn.PssPubKey) // is the same than publicKey, with 0x04 at the begining

	publicKeyECDSA := *publicKey.(*ecdsa.PublicKey)
	addr := crypto.PubkeyToAddress(publicKeyECDSA)

	dscsrv, err := discovery.NewDiscoveryService(addr, &publicKeyECDSA, "url", "Active", []byte{})
	if err != nil {
		color.Red(err.Error())
		os.Exit(0)
	}

	node := &NodeSrv{
		discsrv: dscsrv,
		db:      sto,
		sn:      sn,
	}

	go func() {
		for {
			pmsg := <-sn.PssTopics[config.C.Pss.Topic].Delivery
			fmt.Print("[MSG RECEIVED]: ")
			color.Yellow(string(pmsg.Msg))

			node.HandleMsg(pmsg.Msg)
		}
	}()

	return node, nil
}

func (node *NodeSrv) StoreId(id discovery.Id) error {
	idBytes, err := id.Bytes()
	if err != nil {
		return err
	}
	node.db.Put(id.IdAddr.Bytes(), idBytes)
	return nil
}

// DiscoverId checks if the nade has a fresh data from the id, if not, asks to the network for an idenity address
func (node *NodeSrv) DiscoverId(id discovery.Id) error {
	query, err := node.discsrv.NewQueryPacket(id.IdAddr)
	if err != nil {
		return err
	}
	fmt.Println(query)

	// TODO send the packet over Pss Swarm

	return nil
}

// ResolveId checks if the node knows the idAddress data, if it knows, returns the data
func (node *NodeSrv) ResolveId(idAddr common.Address) error {

	return nil
}

// HandleMsg
func (node *NodeSrv) HandleMsg(msg []byte) error {
	switch msg[0] {
	case discovery.QUERYMSG:
		query, err := discovery.QueryFromBytes(msg)
		if err != nil {
			return err
		}
		// TODO check query packet (PoW, Signature, etc)

		err = node.ResolveId(query.About)

		return nil
	case discovery.ANSWERMSG:
		answer, err := discovery.AnswerFromBytes(msg)
		if err != nil {
			return err
		}
		// TODO check query packet (PoW, Signature, etc)

		// TODO store data
		fmt.Println(answer)
		return nil
	}

	return nil

}
