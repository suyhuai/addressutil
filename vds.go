package addressutil

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/suyhuai/addressutil/base58"
	"github.com/suyhuai/addressutil/ecc"
	"github.com/suyhuai/addressutil/hash160"
	"github.com/suyhuai/addressutil/ripemd160"
	"github.com/suyhuai/addressutil/util"
	"github.com/suyhuai/addressutil/util/btcutil"
)

var (
	P2PKHAddrId = []byte{0x10, 0x1C}
	MainAddrId  = "101c"
)

type VDSAddress struct {
	addr   string
	pubKey []byte
}

func NewVDSAddress(pubKey []byte) (*VDSAddress, error) {
	if len(pubKey) != 65 || pubKey[0] != 0x04 {
		return nil, ErrPublicKeyFormat
	}

	addr, err := VdsAddrFromPub(pubKey)
	if err != nil {
		return nil, err
	}

	address := &VDSAddress{
		pubKey: pubKey,
		addr:   addr,
	}

	return address, nil
}

func (t *VDSAddress) String() string {
	return t.addr
}

func (t *VDSAddress) Url() string {
	return t.String()
}

func encodeAddr(addrHash []byte, prefix []byte) (string, error) {
	if len(addrHash) != ripemd160.Size {
		return "", errors.New("incorrect hash length")
	}

	body := append(prefix, addrHash[:ripemd160.Size]...)
	chk := addrChecksum(body)

	var checksum [4]byte
	copy(checksum[:], chk[:4])

	return base58.Encode(append(body, checksum[:]...)), nil
}

func addrChecksum(input []byte) []byte {
	first := sha256.Sum256(input)
	second := sha256.Sum256(first[:])
	return second[:4]
}

func VdsAddrFromPub(pub []byte) (string, error) {
	pubKey, err := ecc.ParsePubKey(pub, ecc.S256())
	if err != nil {
		return "", err
	}

	addrPubKey, err := btcutil.NewAddressPubKey(pubKey.SerializeCompressed(), &util.Params{})
	if err != nil {
		return "", err
	}
	address, err := encodeAddr(hash160.Hash160(addrPubKey.ScriptAddress())[:ripemd160.Size], P2PKHAddrId)
	return address, err
}

func CheckVDSAddress(address string) bool {
	addr := base58.Decode(address)
	if len(addr) != 26 {
		return false
	}
	body := addr[:len(addr)-4]

	prefix := hex.EncodeToString(body[:len(body)-ripemd160.Size])
	if prefix != MainAddrId {
		return false
	}

	checksum := hex.EncodeToString(addr[len(addr)-4:])
	checksum2 := hex.EncodeToString(addrChecksum(body))
	if checksum != checksum2 {
		return false
	}

	return true
}
