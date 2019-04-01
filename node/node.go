package node

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/iden3/discovery-node/config"
	"github.com/iden3/discovery-node/db"
	"github.com/iden3/discovery-node/discovery"
	"github.com/iden3/discovery-node/utils"
	"github.com/syndtr/goleveldb/leveldb/errors"
	swarm "github.com/vocdoni/go-dvote/net/swarm"
)

// TIMEOUTQUERYMSG is the maximum amount of time that the node will wait to receive the Answer to a Query packet
const TIMEOUTQUERYMSG = 10000 // 10 seconds
const ACTIVENODETYPE = "ACTIVE"

var conversations map[string]*gin.Context

// NodeSrv contains the services of the node
type NodeSrv struct {
	discsrv     discovery.DiscoveryService
	db          *db.Db
	dbOwnIds    *db.Db
	dbAnswCache *db.Db
	sn          *swarm.SimplePss
	ks          *keystore.KeyStore
	acc         accounts.Account
}

// RunNode starts a new discovery node service
func RunNode() (*NodeSrv, error) {
	fmt.Println("initializing node")

	conversations = make(map[string]*gin.Context)

	if config.C.Mode == ACTIVENODETYPE {
		fmt.Println("starting an " + ACTIVENODETYPE + " node")
	}

	sto, err := db.New(config.C.DbPath)
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
	fmt.Print("pubK: ")
	color.Cyan(utils.PublicKeyToString(publicKey))
	// fmt.Println("pubK", sn.PssPubKey) // is the same than publicKey, with 0x04 at the beginning
	publicKeyECDSA := *publicKey.(*ecdsa.PublicKey)
	addr := crypto.PubkeyToAddress(publicKeyECDSA)
	fmt.Println(addr.Hex())

	err = sn.Init()
	if err != nil {
		color.Red(err.Error())
		os.Exit(0)
	}
	// sn.PssSub(config.C.Pss.Kind, config.C.Pss.Key, config.C.Pss.Topic, "")
	fmt.Println("kad addr: ", string(sn.Pss.Kademlia.BaseAddr()))
	sn.PssSub(config.C.Pss.Kind, config.C.Pss.Key, config.C.Pss.Topic, string(sn.Pss.Kademlia.BaseAddr()))
	// defer sn.PssTopics[config.C.Pss.Topic].Unregister()

	dscsrv, err := discovery.NewDiscoveryService(addr, sn.Pss.Kademlia.BaseAddr(), &publicKeyECDSA, "url", config.C.Mode, []byte{})
	if err != nil {
		color.Red(err.Error())
		os.Exit(0)
	}

	// create new keystore with the privK, and new account
	ks := keystore.NewKeyStore(config.C.KeyStore.Path, keystore.StandardScryptN, keystore.StandardScryptP)
	acc, err := ks.ImportECDSA(privateKey, config.C.KeyStore.Password)
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
		ks:          ks,
		acc:         acc,
	}

	fmt.Println("listening pss swarm, topic: " + config.C.Pss.Topic)
	go func() {
		for {
			pmsg := <-node.sn.PssTopics[config.C.Pss.Topic].Delivery
			// fmt.Print("[MSG RECEIVED]: ")
			// color.Yellow(string(pmsg.Msg))
			msgBytes, err := hex.DecodeString(string(pmsg.Msg))
			if err != nil {
				color.Red(err.Error())
			}
			err = node.HandleMsg(msgBytes)
			if err != nil {
				color.Red(err.Error())
			}
		}
	}()

	return node, nil
}

// StoreId adds the id into the dbOwnIds
func (node *NodeSrv) StoreId(id discovery.Id) error {
	idBytes, err := id.Bytes()
	if err != nil {
		return err
	}
	node.dbOwnIds.Put(id.IdAddr.Bytes(), idBytes)
	return nil
}

// DiscoverId checks if the nade has a fresh data from the id, if not, asks to the network for an idenity address
func (node *NodeSrv) DiscoverId(ginContext *gin.Context, idAddr common.Address) (*discovery.Id, error) {
	fmt.Println("start DiscoverId function")
	fmt.Println("checking if id is an own identity")
	// check if is an own identity that this node holds
	idBytes, err := node.dbOwnIds.Get(idAddr.Bytes())
	if err != errors.ErrNotFound {
		id, err := discovery.IdFromBytes(idBytes)
		return id, err
	}

	fmt.Println("check if this node has already a fresh copy of the packet of idAddr")
	// check if this node has already a fresh copy of the packet of idAddr
	answerBytes, err := node.dbAnswCache.Get(idAddr.Bytes())
	if err != errors.ErrNotFound {
		// the node has the packet
		answer, err := discovery.AnswerFromBytes(answerBytes)
		// color.Cyan("node has a copy of the id data")
		if err == nil && answer.Timestamp > time.Now().Unix()-config.C.DiscoverFreshTimeout {
			// the data is a fresh copy
			// set id data structure from answer
			color.Cyan("node has a fresh copy of the id data")

			// send the answer to the client
			ginContext.JSON(200, answer)
			return answer.Id(), nil
		}
	}

	// if answer not found in local databases, ask to the network for it

	fmt.Println("create NewQueryPacket")
	query, err := node.discsrv.NewQueryPacket(idAddr)
	if err != nil {
		return nil, err
	}

	// add query id to the conversations map
	conversations[query.MsgId] = ginContext

	// send the packet over Pss Swarm
	qBytes, err := query.Bytes()
	if err != nil {
		return nil, err
	}
	msg := hex.EncodeToString(qBytes)
	fmt.Println("Send Query over Pss Swarm")
	err = node.sn.PssPub(config.C.Pss.Kind, config.C.Pss.Key, config.C.Pss.Topic, msg, "")

	// wait until the answer is received and sent to the client, or until the timeout is overpassed
	for {
		// check if the conversations[query.MsgId] is already answered
		if _, ok := conversations[query.MsgId]; !ok {
			// once the conversations[query.MsgId] is deleted, this indicates that the answer packet is written in the http connection to the client, so can break the loop and finish the connection
			break
		}

		// check timeout of the msg, if it's too old, remove it
		// as the query.Timestamp is set by the own node, can not be cheated
		if query.Timestamp+TIMEOUTQUERYMSG < time.Now().Unix() {
			fmt.Println("timeout waiting for Answer packet reached, conversations[query.MsgId] deleted. query.MsgId: " + query.MsgId)
			// send error http msg
			conversations[query.MsgId].JSON(500, "error: answer timeout")
			// delete
			delete(conversations, query.MsgId)
			break
		}

		time.Sleep(500 * time.Millisecond)
	}

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

// HandleMsg handles the PSS Swarm messages that the node receives, and depending on the type performs the specified actions
func (node *NodeSrv) HandleMsg(msg []byte) error {
	switch hex.EncodeToString(msg[:discovery.PREFIXLENGTH]) {
	case hex.EncodeToString(discovery.QUERYMSG):
		if config.C.Mode != ACTIVENODETYPE {
			// as a non active node, will not answer QUERY messages
			return nil
		}

		query, err := discovery.QueryFromBytes(msg)
		if err != nil {
			return err
		}
		// TODO check query packet (PoW, Signature, etc)

		id, err := node.ResolveId(query.AboutId)
		if err != nil {
			color.Yellow("received Query msg packet asking for id " + query.AboutId.Hex() + ", and is not in this node")
			return nil
		}
		color.Cyan("-> received QUERY msg packet asking for id " + query.AboutId.Hex() + ", and the id data is in this node")
		fmt.Println("id data found in this node: " + id.IdAddr.Hex())

		fmt.Print("query.MsgId: ")
		color.Cyan(query.MsgId)
		// return id to the requester
		err = node.AnswerAboutId(query, id)
		if err != nil {
			color.Red("error on answering: " + err.Error())
		}

		return nil
	case hex.EncodeToString(discovery.ANSWERMSG):
		color.Green("msg ANSWER received")
		answer, err := discovery.AnswerFromBytes(msg)
		if err != nil {
			return err
		}
		// fmt.Println(answer)
		// TODO check query packet (PoW, Signature, etc)

		fmt.Print("answer.MsgId: ")
		color.Cyan(answer.MsgId)

		// store data in dbAnswCache
		answerBytes, err := answer.Bytes()
		if err != nil {
			return err
		}
		color.Cyan("-> ANSWER received about " + answer.AboutId.Hex() + ", data stored in dbAnswCache")
		node.dbAnswCache.Put(answer.AboutId.Bytes(), answerBytes)

		// check if the conversation[answer.MsgId] is still opened
		if _, ok := conversations[answer.MsgId]; !ok {
			color.Yellow("answer received, but conversation[answer.MsgId] is not open")
			return nil
		}
		// send the answer to the client
		fmt.Println("sending answer through gin context: " + answer.MsgId)
		conversations[answer.MsgId].JSON(200, answer)
		// remove from map
		delete(conversations, answer.MsgId)

		return nil
	default:
		fmt.Println("received pss swarm packet, not recognized type")
	}

	return nil

}

// AnswerAboutId sends (over Pss Swarm directly to the Requester Address) an Answer packet generated from a Query packet and an Id data
func (node *NodeSrv) AnswerAboutId(query *discovery.Query, id *discovery.Id) error {
	answer, err := node.discsrv.NewAnswerPacket(query, id)
	if err != nil {
		return err
	}

	aBytes, err := answer.Bytes()
	if err != nil {
		return err
	}
	msg := hex.EncodeToString(aBytes)
	fmt.Println("Send Answer over Pss Swarm, encrypted with pubK: " + query.RequesterPssPubK.String() + ", and kademlia addr: " + string(query.RequesterKAddr))

	// send a direct message over Pss Swarm
	err = node.sn.PssPub(config.C.Pss.Kind, config.C.Pss.Key, config.C.Pss.Topic, msg, string(query.RequesterKAddr))
	return err
}
