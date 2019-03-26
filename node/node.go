package node

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fatih/color"
	"github.com/iden3/discovery-node/config"
	"github.com/iden3/discovery-node/db"
	"github.com/iden3/discovery-node/discovery"
	"github.com/iden3/discovery-node/utils"
	"github.com/syndtr/goleveldb/leveldb/errors"
	swarm "github.com/vocdoni/go-dvote/net/swarm"
)

type NodeSrv struct {
	discsrv     discovery.DiscoveryService
	db          *db.Db
	dbOwnIds    *db.Db
	dbAnswCache *db.Db
	sn          *swarm.SimplePss
}

func RunNode() (*NodeSrv, error) {
	fmt.Println("initializing node")

	sto, err := db.New("./tmp/iddb")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	stoIdentities := sto.WithPrefix([]byte("identities"))
	stoAnswers := sto.WithPrefix([]byte("answers"))

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
		discsrv:     dscsrv,
		db:          sto,
		dbOwnIds:    stoIdentities,
		dbAnswCache: stoAnswers,
		sn:          sn,
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
	node.dbOwnIds.Put(id.IdAddr.Bytes(), idBytes)
	return nil
}

// DiscoverId checks if the nade has a fresh data from the id, if not, asks to the network for an idenity address
func (node *NodeSrv) DiscoverId(idAddr common.Address) (*discovery.Id, error) {
	// check if is an own identity that this node holds
	idBytes, err := node.dbOwnIds.Get(idAddr.Bytes())
	if err != errors.ErrNotFound {
		id, err := discovery.IdFromBytes(idBytes)
		return id, err
	}

	// check if this node has already a fresh copy of the packet of idAddr
	answerBytes, err := node.dbAnswCache.Get(idAddr.Bytes())
	if err != errors.ErrNotFound {
		// the node has the packet
		answer, err := discovery.AnswerFromBytes(answerBytes)
		if err != nil {
			return nil, err
		}
		if answer.Timestamp < time.Now().Unix()-1000 {
			// the data is a fresh copy
			// set id data structure from answer
			return answer.Id(), nil
		}
	}

	// if answer not found, ask to the network for it

	query, err := node.discsrv.NewQueryPacket(idAddr)
	if err != nil {
		return nil, err
	}

	// send the packet over Pss Swarm
	qBytes, err := query.Bytes()
	if err != nil {
		return nil, err
	}
	msg := hexutil.Encode(qBytes)
	err = node.sn.PssPub(config.C.Pss.Kind, config.C.Pss.Key, config.C.Pss.Topic, msg, "")

	return nil, err
}

// ResolveId checks if the node knows the idAddress data, if it knows, returns the data
func (node *NodeSrv) ResolveId(idAddr common.Address) (*discovery.Id, error) {
	idBytes, err := node.dbOwnIds.Get(idAddr.Bytes())
	if err != nil {
		return nil, err
	}
	id, err := discovery.IdFromBytes(idBytes)
	return id, err
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

		id, err := node.ResolveId(query.About)
		if err != nil {
			return err
		}
		fmt.Println(id)
		// TODO return id to the requester
		return nil
	case discovery.ANSWERMSG:
		answer, err := discovery.AnswerFromBytes(msg)
		if err != nil {
			return err
		}
		// TODO check query packet (PoW, Signature, etc)

		// TODO store data
		fmt.Println(answer)
		answerBytes, err := answer.Bytes()
		if err != nil {
			return err
		}
		node.dbAnswCache.Put(answer.About.Bytes(), answerBytes)
		return nil
	}

	return nil

}
