package addressutil

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/suyhuai/addressutil/util/eosutil"
)

var ErrPublicKeyFormat = errors.New("public key format error")

type Address interface {
	String() string
	Url() string
}

func NewAddress(chain string, pubKey []byte, main bool) (addr Address, err error) {
	switch chain {
	case "BTC":
		addr, err = NewBTCAddress(pubKey, main)
	case "ETH":
		addr, err = NewETHAddress(pubKey)
	case "LTC":
		addr, err = NewLTCAddress(pubKey, main)
	case "BCH":
		addr, err = NewBCHAddress(pubKey, main)
	case "ETC":
		addr, err = NewETHAddress(pubKey)
	case "OMNI":
		addr, err = NewBTCAddress(pubKey, main)
	case "TRON":
		addr, err = NewTRONAddress(pubKey)
	case "VDS":
		addr, err = NewVDSAddress(pubKey)
	default:
		err = fmt.Errorf("unsupport chain type %s", chain)
	}

	return
}

func CheckAddress(address, chain string, main bool) bool {
	switch chain {
	case "BTC", "OMNI":
		return CheckBTCAddress(address, main)
	case "BCH":
		return CheckBCHAddress(address, main)
	case "LTC":
		return CheckLTCAddress(address, main)
	case "ETH", "ETC":
		return CheckETHAddress(address)
	case "EOS":
		return eosutil.CheckEOSAccount(address)
	case "IOST":
		m, _ := regexp.MatchString(`^([a-z0-9_]{5,11})$`, address)
		return m
	case "TRON":
		return CheckTRONAddress(address)
	case "VDS":
		return CheckVDSAddress(address)
	default:
		return true
	}
}

func AddressUrl(address, _chain string) string {
	return address
}
