package addressutil

import (
	"crypto/sha256"
	"github.com/suyhuai/addressutil/base58"
	"github.com/suyhuai/addressutil/util"
	"github.com/suyhuai/addressutil/util/btcutil"
	"github.com/suyhuai/addressutil/util/btcutil/chaincfg"
	"golang.org/x/crypto/ripemd160"
)

type BTCNet uint8

const BTC_MAIN_NET BTCNet = 0x00
const BTC_TEST_NET BTCNet = 0x6f

type BTCAddress struct {
	Address

	net    BTCNet
	addr   string
	pubKey []byte
}

func NewBTCAddress(pubKey []byte, main bool) (*BTCAddress, error) {
	if len(pubKey) != 65 || pubKey[0] != 0x04 {
		return nil, ErrPublicKeyFormat
	}

	var net BTCNet
	if main {
		net = BTC_MAIN_NET
	} else {
		net = BTC_TEST_NET
	}

	return &BTCAddress{
		net:    net,
		pubKey: pubKey,
	}, nil
}

func (a *BTCAddress) String() string {
	if a.addr != "" {
		return a.addr
	}

	h1 := sha256.Sum256(a.pubKey)
	hash := ripemd160.New()
	hash.Write(h1[:])
	h2 := hash.Sum(nil)

	a.addr = base58.CheckEncode(h2[:], byte(a.net))
	return a.addr
}

func (a *BTCAddress) Url() string {
	return a.String()
}

func CheckBTCAddress(address string, main bool) bool {
	var netParam *util.Params
	if main {
		netParam = &chaincfg.MainNetParams
	} else {
		netParam = &chaincfg.TestNet3Params
	}
	addr, err := btcutil.DecodeAddress(address, netParam)
	if err != nil {
		return false
	}
	return addr.IsForNet(netParam)
}
