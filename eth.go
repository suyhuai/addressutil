package addressutil

import (
	"encoding/hex"
	"github.com/suyhuai/addressutil/util/ethutil"
	"golang.org/x/crypto/sha3"
	"strings"
)

type ETHAddress struct {
	Address

	addr   string
	pubKey []byte
}

func NewETHAddress(pubKey []byte) (*ETHAddress, error) {
	return &ETHAddress{
		pubKey: pubKey,
	}, nil
}

func (a *ETHAddress) String() string {
	if a.addr != "" {
		return a.addr
	}

	hash := sha3.NewLegacyKeccak256()
	hash.Write(a.pubKey[1:])
	b := hash.Sum(nil)
	addrBytes := []byte(hex.EncodeToString(b[12:]))

	hash.Reset()
	hash.Write(addrBytes)
	b = hash.Sum(nil)
	check := hex.EncodeToString(b)

	for i := 0; i < len(addrBytes); i++ {
		if check[i] >= '8' && addrBytes[i] >= 'a' {
			addrBytes[i] -= 0x20
		}
	}

	return "0x" + string(addrBytes)
}

func (a *ETHAddress) Url() string {
	return a.String()
}

func CheckETHAddress(address string) bool {
	addr, err := ethutil.NewMixedcaseAddressFromString(address)
	if err != nil {
		return false
	}

	original := strings.ToLower(addr.Original())
	hex := strings.ToLower(addr.Address().Hex())
	return original == hex
}
