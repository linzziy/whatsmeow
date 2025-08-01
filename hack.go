package whatsmeow

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"go.mau.fi/util/random"
	waBinary "go.mau.fi/whatsmeow/binary"
	"go.mau.fi/whatsmeow/socket"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/util"
	"strconv"
	"strings"
	"time"
)

type DevicePlatform int

const IosMessageIDPrefix = "3A"
const (
	PlatformWeb DevicePlatform = iota
	PlatformAndroid
	PlatformIos
)

func (d DevicePlatform) String() string {
	return [...]string{"Web", "Android", "Ios"}[d]
}

func (cli *Client) SetSocket(ns *socket.NoiseSocket) {
	cli.socket = ns
}

func (cli *Client) GetSocket() *socket.NoiseSocket {
	return cli.socket
}

func (cli *Client) HandleFrame(data []byte) {
	cli.handleFrame(data)
}

func (cli *Client) ResetExpectedDisconnect() {
	cli.resetExpectedDisconnect()
}

func (cli *Client) AutoReconnect() {
	cli.autoReconnect()
}

func (cli *Client) OnDisconnect(ns *socket.NoiseSocket, remote bool) {
	cli.onDisconnect(ns, remote)
}

func (cli *Client) SendNode(node waBinary.Node) error {
	_, err := cli.sendNodeAndGetData(node)
	return err
}

// GenerateMessageID generates a random string that can be used as a message ID on WhatsApp.
//
//	android：B8B38A4AEB507214B9906377C31E8381、63E28F596B616D5E954624E62E8A3731、386C77891A4713851C82C648311545D8
//	web：3EB0 0D821F02E0A0A3F845
//	ios：3A A79EBD626D54F60DDF、3A 7C2E07FD92B63697CC、3A 870F338797859209B7、3A A79EBD626D54F60DDF
//
//	msgID := cli.GenerateMessageIDPlatform(platform)
//	cli.SendMessage(context.Background(), targetJID, &waE2E.Message{...}, whatsmeow.SendRequestExtra{ID: msgID})
func (cli *Client) GenerateMessageIDPlatform(platform DevicePlatform) types.MessageID {
	if cli != nil && cli.MessengerConfig != nil {
		return types.MessageID(strconv.FormatInt(GenerateFacebookMessageID(), 10))
	}
	data := make([]byte, 8, 8+20+16)
	binary.BigEndian.PutUint64(data, uint64(time.Now().Unix()))
	ownID := cli.getOwnID()
	if !ownID.IsEmpty() {
		data = append(data, []byte(ownID.User)...)
		data = append(data, []byte("@c.us")...)
	}
	data = append(data, random.Bytes(16)...)
	hash := sha256.Sum256(data)
	switch platform {
	case PlatformWeb:
		return WebMessageIDPrefix + strings.ToUpper(hex.EncodeToString(hash[:9]))
	case PlatformAndroid:
		return strings.ToUpper(hex.EncodeToString(hash[:16]))
	case PlatformIos:
		return IosMessageIDPrefix + strings.ToUpper(hex.EncodeToString(hash[:9]))
	default:
		panic("unknown platform")
	}
}

// 是否可以作为安卓端生成规则使用？？
func (cli *Client) GenerateMessageIDHash() types.MessageID {
	data := make([]byte, 8, 8+20+16)
	binary.BigEndian.PutUint64(data, uint64(time.Now().Unix()))
	ownID := cli.getOwnID()
	if !ownID.IsEmpty() {
		data = append(data, []byte(ownID.User)...)
		data = append(data, []byte("@c.us")...)
	}
	data = append(data, random.Bytes(16)...)
	hash := md5.Sum(data)
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

// Test?
func (cli *Client) InitLogin() {
	_, err := cli.sendIQ(infoQuery{
		Namespace: "w",
		Type:      "get",
		To:        types.ServerJID,
		Context:   cli.BackgroundEventCtx,
		Content:   []waBinary.Node{{Tag: "props", Attrs: waBinary.Attrs{"protocol": 2, "hash": ""}}},
	})
	if err != nil {
		cli.Log.Errorf(err.Error())
	}

	_, err = cli.sendIQ(infoQuery{
		Namespace: "w:b",
		Type:      "get",
		To:        types.ServerJID,
		Context:   cli.BackgroundEventCtx,
		Content:   []waBinary.Node{{Tag: "list"}},
	})
	if err != nil {
		cli.Log.Errorf(err.Error())
	}

	_, err = cli.sendIQ(infoQuery{
		Namespace: "md",
		Type:      "set",
		To:        types.ServerJID,
		Context:   cli.BackgroundEventCtx,
		Content:   []waBinary.Node{{Tag: "remove-companion-device", Attrs: waBinary.Attrs{"all": true, "reason": "user_initiated"}}},
	})
	if err != nil {
		cli.Log.Errorf(err.Error())
	}
}

//	func (cli *Client) DeltaSync(jid types.JID) {
//		usync, err := cli.usync(context.TODO(), []types.JID{jid}, "delta", "interactive", []waBinary.Node{
//			{Tag: "business", Content: []waBinary.Node{{Tag: "verified_name"}}},
//			{Tag: "contact"},
//			{Tag: "status"},
//		})
//		if err != nil {
//			return
//		}
//		print(usync)
//	}
func (cli *Client) ContactsSync(jids []types.JID, upload bool) (*waBinary.Node, error) {
	return cli.usync(context.TODO(), jids, util.If(upload, "delta", "query"), "interactive", []waBinary.Node{
		{Tag: "business", Content: []waBinary.Node{{Tag: "verified_name"}}},
		{Tag: "contact"},
		{Tag: "status"},
	})
}
func (cli *Client) ContactSync(jid types.JID, upload bool) (*waBinary.Node, error) {
	return cli.ContactsSync([]types.JID{jid}, upload)
}
