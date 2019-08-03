package bchutil

import (
	"encoding/hex"
	"errors"
	"github.com/suyhuai/addressutil/base58"
	bchec "github.com/suyhuai/addressutil/ecc"
	"github.com/suyhuai/addressutil/ripemd160"
	"github.com/suyhuai/addressutil/util"
	"github.com/suyhuai/addressutil/util/bchutil/chaincfg"
)

var (
	ErrChecksumMismatch   = errors.New("checksum mismatch")
	ErrUnknownAddressType = errors.New("unknown address type")
	ErrAddressCollision   = errors.New("address collision")
	ErrInvalidFormat      = errors.New("invalid format: version and/or checksum bytes missing")
	Prefixes              map[*util.Params]string
)

type AddressType int
type PubKeyFormat int

const (
	AddrTypePayToPubKeyHash AddressType  = 0
	AddrTypePayToScriptHash AddressType  = 1
	PKFUncompressed         PubKeyFormat = iota
	PKFCompressed
	PKFHybrid
)

func init() {
	Prefixes = make(map[*util.Params]string)
	Prefixes[&chaincfg.MainNetParams] = "bitcoincash"
	Prefixes[&chaincfg.TestNet3Params] = "bchtest"
	Prefixes[&chaincfg.RegressionNetParams] = "bchreg"
	Prefixes[&chaincfg.SimNetParams] = "bchsim"
}

type Address interface {
	IsForNet(*util.Params) bool
}

func DecodeAddress(addr string, defaultNet *util.Params) (Address, error) {
	pre, ok := Prefixes[defaultNet]
	if !ok {
		return nil, errors.New("unknown network parameters")
	}
	if len(addr) < len(pre)+2 {
		return nil, errors.New("invalid length address")
	}

	// Add prefix if it does not exist
	addrWithPrefix := addr
	if addr[:len(pre)+1] != pre+":" {
		addrWithPrefix = pre + ":" + addr
	}

	// Switch on decoded length to determine the type.
	decoded, _, typ, err := checkDecodeCashAddress(addrWithPrefix)
	if err == nil {
		switch len(decoded) {
		case ripemd160.Size: // P2PKH or P2SH
			switch typ {
			case AddrTypePayToPubKeyHash:
				return newAddressPubKeyHash(decoded, defaultNet)
			default:
				return nil, ErrUnknownAddressType
			}
		default:
			return nil, errors.New("decoded address is of unknown size")
		}
	} else if err == ErrChecksumMismatch {
		return nil, ErrChecksumMismatch
	}

	// Serialized public keys are either 65 bytes (130 hex chars) if
	// uncompressed/hybrid or 33 bytes (66 hex chars) if compressed.
	if len(addr) == 130 || len(addr) == 66 {
		serializedPubKey, err := hex.DecodeString(addr)
		if err != nil {
			return nil, err
		}
		return NewAddressPubKey(serializedPubKey, defaultNet)
	}

	// Switch on decoded length to determine the type.
	decoded, netID, err := base58.CheckDecode(addr)
	if err != nil {
		if err == base58.ErrChecksum {
			return nil, ErrChecksumMismatch
		}
		return nil, errors.New("decoded address is of unknown format")
	}
	switch len(decoded) {
	case ripemd160.Size: // P2PKH or P2SH
		isP2PKH := chaincfg.IsPubKeyHashAddrID(netID)
		isP2SH := chaincfg.IsScriptHashAddrID(netID)
		switch hash160 := decoded; {
		case isP2PKH && isP2SH:
			return nil, ErrAddressCollision
		case isP2PKH:
			return newLegacyAddressPubKeyHash(hash160, netID)
		case isP2SH:
			return newLegacyAddressScriptHashFromHash(hash160, netID)
		default:
			return nil, ErrUnknownAddressType
		}

	default:
		return nil, errors.New("decoded address is of unknown size")
	}
}

func checkDecodeCashAddress(input string) (result []byte, prefix string, t AddressType, err error) {
	prefix, data, err := DecodeCashAddress(input)
	if err != nil {
		return data, prefix, AddrTypePayToPubKeyHash, err
	}
	data, err = convertBits(data, 5, 8, false)
	if err != nil {
		return data, prefix, AddrTypePayToPubKeyHash, err
	}
	if len(data) != 21 {
		return data, prefix, AddrTypePayToPubKeyHash, errors.New("incorrect data length")
	}
	switch data[0] {
	case 0x00:
		t = AddrTypePayToPubKeyHash
	case 0x08:
		t = AddrTypePayToScriptHash
	}
	return data[1:21], prefix, t, nil
}

type AddressPubKeyHash struct {
	hash   [ripemd160.Size]byte
	prefix string
}

func newAddressPubKeyHash(pkHash []byte, net *util.Params) (*AddressPubKeyHash, error) {
	// Check for a valid pubkey hash length.
	if len(pkHash) != ripemd160.Size {
		return nil, errors.New("pkHash must be 20 bytes")
	}

	prefix, ok := Prefixes[net]
	if !ok {
		return nil, errors.New("unknown network parameters")
	}

	addr := &AddressPubKeyHash{prefix: prefix}
	copy(addr.hash[:], pkHash)
	return addr, nil
}

func (a *AddressPubKeyHash) IsForNet(net *util.Params) bool {
	checkPre, ok := Prefixes[net]
	if !ok {
		return false
	}
	return a.prefix == checkPre
}

type AddressPubKey struct {
	pubKeyFormat PubKeyFormat
	pubKey       *bchec.PublicKey
	pubKeyHashID byte
}

func NewAddressPubKey(serializedPubKey []byte, net *util.Params) (*AddressPubKey, error) {
	pubKey, err := bchec.ParsePubKey(serializedPubKey, bchec.S256())
	if err != nil {
		return nil, err
	}

	// Set the format of the pubkey.  This probably should be returned
	// from bchec, but do it here to avoid API churn.  We already know the
	// pubkey is valid since it parsed above, so it's safe to simply examine
	// the leading byte to get the format.
	pkFormat := PKFUncompressed
	switch serializedPubKey[0] {
	case 0x02, 0x03:
		pkFormat = PKFCompressed
	case 0x06, 0x07:
		pkFormat = PKFHybrid
	}

	return &AddressPubKey{
		pubKeyFormat: pkFormat,
		pubKey:       pubKey,
		pubKeyHashID: net.LegacyPubKeyHashAddrID,
	}, nil
}

func (a *AddressPubKey) IsForNet(net *util.Params) bool {
	return a.pubKeyHashID == net.LegacyPubKeyHashAddrID
}

type LegacyAddressPubKeyHash struct {
	hash  [ripemd160.Size]byte
	netID byte
}

func newLegacyAddressPubKeyHash(pkHash []byte, netID byte) (*LegacyAddressPubKeyHash, error) {
	// Check for a valid pubkey hash length.
	if len(pkHash) != ripemd160.Size {
		return nil, errors.New("pkHash must be 20 bytes")
	}

	addr := &LegacyAddressPubKeyHash{netID: netID}
	copy(addr.hash[:], pkHash)
	return addr, nil
}

func (a *LegacyAddressPubKeyHash) IsForNet(net *util.Params) bool {
	return a.netID == net.LegacyPubKeyHashAddrID
}

type LegacyAddressScriptHash struct {
	hash  [ripemd160.Size]byte
	netID byte
}

func newLegacyAddressScriptHashFromHash(scriptHash []byte, netID byte) (*LegacyAddressScriptHash, error) {
	// Check for a valid script hash length.
	if len(scriptHash) != ripemd160.Size {
		return nil, errors.New("scriptHash must be 20 bytes")
	}

	addr := &LegacyAddressScriptHash{netID: netID}
	copy(addr.hash[:], scriptHash)
	return addr, nil
}

func (a *LegacyAddressScriptHash) IsForNet(net *util.Params) bool {
	return a.netID == net.LegacyScriptHashAddrID
}

func DecodeCashAddress(str string) (string, []byte, error) {
	// Go over the string and do some sanity checks.
	lower, upper := false, false
	prefixSize := 0
	for i := 0; i < len(str); i++ {
		c := str[i]
		if c >= 'a' && c <= 'z' {
			lower = true
			continue
		}

		if c >= 'A' && c <= 'Z' {
			upper = true
			continue
		}

		if c >= '0' && c <= '9' {
			// We cannot have numbers in the prefix.
			if prefixSize == 0 {
				return "", nil, errors.New("addresses cannot have numbers in the prefix")
			}

			continue
		}

		if c == ':' {
			// The separator must not be the first character, and there must not
			// be 2 separators.
			if i == 0 || prefixSize != 0 {
				return "", nil, errors.New("the separator must not be the first character")
			}

			prefixSize = i
			continue
		}

		// We have an unexpected character.
		return "", nil, errors.New("unexpected character")
	}

	// We must have a prefix and a data part and we can't have both uppercase
	// and lowercase.
	if prefixSize == 0 {
		return "", nil, errors.New("address must have a prefix")
	}

	if upper && lower {
		return "", nil, errors.New("addresses cannot use both upper and lower case characters")
	}

	// Get the prefix.
	var prefix string
	for i := 0; i < prefixSize; i++ {
		prefix += string(lowerCase(str[i]))
	}

	// Decode values.
	valuesSize := len(str) - 1 - prefixSize
	values := make([]byte, valuesSize)
	for i := 0; i < valuesSize; i++ {
		c := str[i+prefixSize+1]
		// We have an invalid char in there.
		if c > 127 || CharsetRev[c] == -1 {
			return "", nil, errors.New("invalid character")
		}

		values[i] = byte(CharsetRev[c])
	}

	// Verify the checksum.
	if !verifyChecksum(prefix, values) {
		return "", nil, ErrChecksumMismatch
	}

	return prefix, values[:len(values)-8], nil
}

// Base32 conversion contains some licensed code
// https://github.com/sipa/bech32/blob/master/ref/go/src/bech32/bech32.go
// Copyright (c) 2017 Takatoshi Nakagawa
// MIT License
func convertBits(data []byte, fromBits uint, tobits uint, pad bool) ([]byte, error) {
	// General power-of-2 base conversion.
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

func lowerCase(c byte) byte {
	// ASCII black magic.
	return c | 0x20
}

func verifyChecksum(prefix string, payload []byte) bool {
	return polyMod(cat(expandPrefix(prefix), payload)) == 0
}

var CharsetRev = [128]int8{
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 15, -1, 10, 17, 21, 20, 26, 30, 7,
	5, -1, -1, -1, -1, -1, -1, -1, 29, -1, 24, 13, 25, 9, 8, 23, -1, 18, 22,
	31, 27, 19, -1, 1, 0, 3, 16, 11, 28, 12, 14, 6, 4, 2, -1, -1, -1, -1,
	-1, -1, 29, -1, 24, 13, 25, 9, 8, 23, -1, 18, 22, 31, 27, 19, -1, 1, 0,
	3, 16, 11, 28, 12, 14, 6, 4, 2, -1, -1, -1, -1, -1,
}

func polyMod(v []byte) uint64 {
	c := uint64(1)
	for _, d := range v {
		c0 := byte(c >> 35)

		c = ((c & 0x07ffffffff) << 5) ^ uint64(d)

		if c0&0x01 > 0 {
			c ^= 0x98f2bc8e61
		}

		if c0&0x02 > 0 {
			c ^= 0x79b76d99e2
		}

		if c0&0x04 > 0 {
			c ^= 0xf33e5fb3c4
		}

		if c0&0x08 > 0 {
			c ^= 0xae2eabe2a8
		}

		if c0&0x10 > 0 {
			c ^= 0x1e4f43e470
		}
	}
	return c ^ 1
}

func cat(x, y []byte) []byte {
	return append(x, y...)
}

func expandPrefix(prefix string) []byte {
	ret := make([]byte, len(prefix)+1)
	for i := 0; i < len(prefix); i++ {
		ret[i] = prefix[i] & 0x1f
	}

	ret[len(prefix)] = 0
	return ret
}
