package meshtastic

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/brianshea2/meshmap.net/internal/meshtastic/generated"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"google.golang.org/protobuf/proto"
)

var DefaultKey = []byte{
	0xd4, 0xf1, 0xbb, 0x3a,
	0x20, 0x29, 0x07, 0x59,
	0xf0, 0xbc, 0xff, 0xab,
	0xcf, 0x4e, 0x69, 0x01,
}

func NewBlockCipher(key []byte) cipher.Block {
	c, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return c
}

type MQTTClient struct {
	Topics         []string
	TopicRegex     *regexp.Regexp
	Accept         func(from uint32) bool
	BlockCipher    cipher.Block
	MessageHandler func(from uint32, topic string, portNum generated.PortNum, payload []byte)
	mqtt.Client
}

func (c *MQTTClient) Connect() error {
	randomId := make([]byte, 4)
	rand.Read(randomId)
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://mqtt.meshtastic.org:1883")
	opts.SetClientID(fmt.Sprintf("meshobserv-%x", randomId))
	opts.SetUsername("meshdev")
	opts.SetPassword("large4cats")
	opts.SetOrderMatters(false)
	opts.SetDefaultPublishHandler(c.handleMessage)
	c.Client = mqtt.NewClient(opts)
	token := c.Client.Connect()
	<-token.Done()
	if err := token.Error(); err != nil {
		return err
	}
	log.Print("[info] connected")
	topics := make(map[string]byte)
	for _, topic := range c.Topics {
		topics[topic] = 0
	}
	token = c.SubscribeMultiple(topics, nil)
	<-token.Done()
	if err := token.Error(); err != nil {
		return err
	}
	log.Print("[info] subscribed")
	return nil
}

func (c *MQTTClient) Disconnect() {
	if c.IsConnected() {
		c.Client.Disconnect(1000)
	}
}

func (c *MQTTClient) handleMessage(_ mqtt.Client, msg mqtt.Message) {
	// filter topic
	topic := msg.Topic()
	if !c.TopicRegex.MatchString(topic) {
		return
	}
	// parse ServiceEnvelope
	var envelope generated.ServiceEnvelope
	if err := proto.Unmarshal(msg.Payload(), &envelope); err != nil {
		log.Printf("[warn] could not parse ServiceEnvelope on %v: %v", topic, err)
		return
	}
	// get MeshPacket
	packet := envelope.GetPacket()
	if packet == nil {
		log.Printf("[warn] skipping ServiceEnvelope with no MeshPacket on %v", topic)
		return
	}
	// no anonymous packets
	from := packet.GetFrom()
	if from == 0 {
		log.Printf("[warn] skipping MeshPacket from unknown on %v", topic)
		return
	}
	// check sender
	if c.Accept != nil && !c.Accept(from) {
		return
	}
	// get Data, try decoded first
	data := packet.GetDecoded()
	if data == nil {
		// data must be (probably) encrypted
		encrypted := packet.GetEncrypted()
		if encrypted == nil {
			log.Printf("[warn] skipping MeshPacket from %v with no data on %v", from, topic)
			return
		}
		// decrypt
		nonce := make([]byte, 16)
		binary.LittleEndian.PutUint32(nonce[0:], packet.GetId())
		binary.LittleEndian.PutUint32(nonce[8:], from)
		decrypted := make([]byte, len(encrypted))
		cipher.NewCTR(c.BlockCipher, nonce).XORKeyStream(decrypted, encrypted)
		// parse Data
		data = new(generated.Data)
		if err := proto.Unmarshal(decrypted, data); err != nil {
			// ignore, probably encrypted with other psk
			return
		}
	}
	c.MessageHandler(from, topic, data.GetPortnum(), data.GetPayload())
}

func init() {
	mqtt.ERROR = log.New(os.Stderr, "[error] mqtt: ", log.Flags()|log.Lmsgprefix)
	mqtt.CRITICAL = log.New(os.Stderr, "[crit] mqtt: ", log.Flags()|log.Lmsgprefix)
	mqtt.WARN = log.New(os.Stderr, "[warn] mqtt: ", log.Flags()|log.Lmsgprefix)
}
