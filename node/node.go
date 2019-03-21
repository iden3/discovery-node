package core

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fatih/color"
	"github.com/iden3/discovery-node/config"
	"github.com/iden3/discovery-node/utils"
	swarm "github.com/vocdoni/go-dvote/net/swarm"
)

func readInput(sn *swarm.SimplePss) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("enter msg: ")
		msg, _ := reader.ReadString('\n')
		msg = strings.TrimSpace(msg)
		if msg == "" {
			fmt.Println("no message")
			continue
		}
		fmt.Print("[SEND MSG]: ")
		color.Cyan(msg)
		err := sn.PssPub(config.C.Pss.Kind, config.C.Pss.Key, config.C.Pss.Topic, msg, "")
		if err != nil {
			color.Red(err.Error())
		}
	}
}

func RunNode() {
	fmt.Println("initializing node")

	sn := new(swarm.SimplePss)

	// set a random privK
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
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

	// first msg
	currentTime := int64(time.Now().Unix())
	err = sn.PssPub(config.C.Pss.Kind, config.C.Pss.Key, config.C.Pss.Topic, fmt.Sprintf("Hello world from %s at %d", "hostname", currentTime), "")
	if err != nil {
		color.Red(err.Error())
	}

	readInput(sn)
}
