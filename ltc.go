package addressutil

import (
	"github.com/suyhuai/addressutil/util"
	"github.com/suyhuai/addressutil/util/ltcutil"
	"github.com/suyhuai/addressutil/util/ltcutil/chaincfg"
)

type LTCNet uint8

const LTC_MAIN_NET LTCNet = 0x30
const LTC_TEST_NET LTCNet = 0x6f

type LTCAddress struct {
	Address

	net    LTCNet
	addr   string
	pubKey []byte
}

func NewLTCAddress(pubKey []byte, main bool) (*LTCAddress, error) {
	var net LTCNet
	if main {
		net = LTC_MAIN_NET
	} else {
		net = LTC_TEST_NET
	}

	return &LTCAddress{
		net:    net,
		pubKey: pubKey,
	}, nil
}

func (a *LTCAddress) String() string {
	if a.addr != "" {
		return a.addr
	}

	ba := &BTCAddress{
		net:    BTCNet(a.net),
		pubKey: a.pubKey,
	}

	a.addr = ba.String()

	return a.addr
}

func (a *LTCAddress) Url() string {
	return a.String()
}

func CheckLTCAddress(address string, main bool) bool {
	var netParam *util.Params
	if main {
		netParam = &chaincfg.MainNetParams
	} else {
		netParam = &chaincfg.TestNet4Params
	}
	addr, err := ltcutil.DecodeAddress(address, netParam)
	if err != nil {
		return false
	}
	return addr.IsForNet(netParam)
}
