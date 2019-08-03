package addressutil

import (
	"errors"
	"fmt"
	"github.com/suyhuai/addressutil/base58"
	"github.com/suyhuai/addressutil/util"
	"github.com/suyhuai/addressutil/util/bchutil"
	"github.com/suyhuai/addressutil/util/bchutil/chaincfg"
)

var bchBase32Encoder = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

type BCHNet uint8
type BCHPrefix string
type AddressType int

const (
	BCH_MAIN_NET BCHNet = 0x00
	BCH_TEST_NET BCHNet = 0x6f

	BCH_MAIN_PREFIX BCHPrefix = "bitcoincash"
	BCH_TEST_PREFIX BCHPrefix = "bchtest"

	AddrTypePayToPubKeyHash AddressType = 0
	AddrTypePayToScriptHash AddressType = 1
)

type BCHAddress struct {
	Address

	prefix BCHPrefix
	net    BCHNet
	addr   string
	pubKey []byte
}

func NewBCHAddress(pubKey []byte, main bool) (*BCHAddress, error) {
	if len(pubKey) != 65 || pubKey[0] != 0x04 {
		return nil, ErrPublicKeyFormat
	}

	var net BCHNet
	var prefix BCHPrefix
	if main {
		net = BCH_MAIN_NET
		prefix = BCH_MAIN_PREFIX
	} else {
		net = BCH_TEST_NET
		prefix = BCH_TEST_PREFIX
	}

	return &BCHAddress{
		net:    net,
		prefix: prefix,
		pubKey: pubKey,
	}, nil
}

func (a *BCHAddress) String() string {
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

func (a *BCHAddress) Url() string {
	return a.String()
}

func CheckBCHAddress(address string, main bool) bool {
	var netParam *util.Params
	if main {
		netParam = &chaincfg.MainNetParams
	} else {
		netParam = &chaincfg.TestNet3Params
	}
	addr, err := bchutil.DecodeAddress(address, netParam)
	if err != nil {
		fmt.Errorf(err.Error())
		return false
	}
	return addr.IsForNet(netParam)
}

func CashAddress(addr string) (string, error) {
	h2, net, err := base58.CheckDecode(addr)
	if err != nil {
		return "", err
	}
	var prefix BCHPrefix
	switch BCHNet(net) {
	case BCH_MAIN_NET:
		prefix = BCH_MAIN_PREFIX
	case BCH_TEST_NET:
		prefix = BCH_TEST_PREFIX
	default:
		errors.New("unsupported address version")
	}

	return string(prefix) + ":" + checkEncodeCashAddress(h2, string(prefix), AddrTypePayToPubKeyHash), nil
}

func checkEncodeCashAddress(input []byte, prefix string, t AddressType) string {
	k, err := packAddressData(t, input)
	if err != nil {
		return ""
	}
	return encode(prefix, k)
}

func packAddressData(addrType AddressType, addrHash []byte) ([]byte, error) {
	if addrType != AddrTypePayToPubKeyHash && addrType != AddrTypePayToScriptHash {
		return nil, errors.New("invalid AddressType")
	}
	versionByte := uint(addrType) << 3
	encodedSize := (uint(len(addrHash)) - 20) / 4
	if (len(addrHash)-20)%4 != 0 {
		return nil, errors.New("invalid address hash size")
	}
	if encodedSize < 0 || encodedSize > 8 {
		return nil, errors.New("encoded size out of valid range")
	}
	versionByte |= encodedSize
	var addrHashUint []byte
	addrHashUint = append(addrHashUint, addrHash...)
	data := append([]byte{byte(versionByte)}, addrHashUint...)
	packedData, err := convertBits(data, 8, 5, true)
	if err != nil {
		return []byte{}, err
	}
	return packedData, nil
}

func encode(prefix string, payload []byte) string {
	checksum := createChecksum(prefix, payload)
	combined := cat(payload, checksum)
	ret := ""

	for _, c := range combined {
		ret += string(bchBase32Encoder[c])
	}

	return ret
}

func convertBits(data []byte, fromBits uint, tobits uint, pad bool) ([]byte, error) {
	var uintArr []uint
	for _, i := range data {
		uintArr = append(uintArr, uint(i))
	}
	acc := uint(0)
	bits := uint(0)
	var ret []uint
	maxv := uint((1 << tobits) - 1)
	maxAcc := uint((1 << (fromBits + tobits - 1)) - 1)
	for _, value := range uintArr {
		acc = ((acc << fromBits) | value) & maxAcc
		bits += fromBits
		for bits >= tobits {
			bits -= tobits
			ret = append(ret, (acc>>bits)&maxv)
		}
	}
	if pad {
		if bits > 0 {
			ret = append(ret, (acc<<(tobits-bits))&maxv)
		}
	} else if bits >= fromBits || ((acc<<(tobits-bits))&maxv) != 0 {
		return []byte{}, errors.New("encoding padding error")
	}
	var dataArr []byte
	for _, i := range ret {
		dataArr = append(dataArr, byte(i))
	}
	return dataArr, nil
}

func createChecksum(prefix string, payload []byte) []byte {
	enc := cat(expandPrefix(prefix), payload)
	enc = cat(enc, []byte{0, 0, 0, 0, 0, 0, 0, 0})
	mod := polyMod(enc)
	ret := make([]byte, 8)
	for i := 0; i < 8; i++ {
		ret[i] = byte((mod >> uint(5*(7-i))) & 0x1f)
	}
	return ret
}

func expandPrefix(prefix string) []byte {
	ret := make([]byte, len(prefix)+1)
	for i := 0; i < len(prefix); i++ {
		ret[i] = prefix[i] & 0x1f
	}

	ret[len(prefix)] = 0
	return ret
}

func cat(x, y []byte) []byte {
	return append(x, y...)
}

func polyMod(data []byte) uint64 {
	c := uint64(1)
	for _, d := range data {
		c0 := uint8(c >> 35)
		c = ((c & 0x07ffffffff) << 5) ^ uint64(d)

		if (c0 & 0x01) > 0 {
			c ^= 0x98f2bc8e61
		}
		if (c0 & 0x02) > 0 {
			c ^= 0x79b76d99e2
		}
		if (c0 & 0x04) > 0 {
			c ^= 0xf33e5fb3c4
		}
		if (c0 & 0x08) > 0 {
			c ^= 0xae2eabe2a8
		}
		if (c0 & 0x10) > 0 {
			c ^= 0x1e4f43e470
		}
	}

	return c ^ 1
}
