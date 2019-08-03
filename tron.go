package addressutil

import (
	"bytes"
	"crypto/sha256"
	"github.com/suyhuai/addressutil/base58"
	"golang.org/x/crypto/sha3"
)

type TRONAddress struct {
	addr   string
	pubKey []byte
}

func NewTRONAddress(pubKey []byte) (*TRONAddress, error) {
	if len(pubKey) != 65 {
		return nil, ErrPublicKeyFormat
	}

	address := &TRONAddress{
		pubKey: pubKey[1:],
		addr:   tronAddrFromPub(pubKey),
	}

	return address, nil
}

func (t *TRONAddress) String() string {
	return t.addr
}

func (t *TRONAddress) Url() string {
	return t.String()
}

func tronAddrFromPub(pub []byte) string {
	// #1 取公钥仅包含x，y坐标的64字节的byte数组
	pubBytes := pub[1:]

	// #2
	hash := sha3.NewLegacyKeccak256()
	hash.Write(pubBytes)
	hashed := hash.Sum(nil)
	last20 := hashed[len(hashed)-20:]

	// #3
	addr41 := append([]byte{0x41}, last20...)

	// #4
	hash2561 := sha256.Sum256(addr41)
	hash2562 := sha256.Sum256(hash2561[:])
	checksum := hash2562[:4]

	// #5
	rawAddr := append(addr41, checksum...)
	// #6
	tronAddr := base58.Encode(rawAddr)
	return tronAddr
}

func CheckTRONAddress(base58Addr string) bool {
	rawAddr := base58.Decode(base58Addr)
	if len(rawAddr) != 25 || len(rawAddr) == 0 {
		return false
	}
	if rawAddr[0] != 0x41 {
		return false
	}

	addr41 := rawAddr[:21]
	hash2561 := sha256.Sum256(addr41)
	hash2562 := sha256.Sum256(hash2561[:])
	if !bytes.Equal(rawAddr[21:25], hash2562[:4]) {
		return false
	}
	return true
}
