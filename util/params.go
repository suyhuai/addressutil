package util

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"
)

type BitcoinNet uint32

const (
	MainNet  BitcoinNet = 0xe8f3e1e3
	TestNet  BitcoinNet = 0xfabfb5da
	TestNet3 BitcoinNet = 0xf4f3e5f4
	SimNet   BitcoinNet = 0x12141c16
)

const (
	DeploymentTestDummy = iota
	DeploymentCSV
	DeploymentSegwit
)

var (
	bigOne             = big.NewInt(1)
	mainPowLimit       = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 224), bigOne)
	regressionPowLimit = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 255), bigOne)
	testNet3PowLimit   = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 224), bigOne)
	simNetPowLimit     = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 255), bigOne)
)

var (
	registeredNets       = make(map[BitcoinNet]struct{})
	pubKeyHashAddrIDs    = make(map[byte]struct{})
	scriptHashAddrIDs    = make(map[byte]struct{})
	cashAddressPrefixes  = make(map[string]struct{})
	hdPrivToPubKeyIDs    = make(map[[4]byte][]byte)
	bech32SegwitPrefixes = make(map[string]struct{})
)

var (
	ErrDuplicateNet = errors.New("duplicate Bitcoin network")
)

func IsPubKeyHashAddrID(id byte) bool {
	_, ok := pubKeyHashAddrIDs[id]
	return ok
}
func IsScriptHashAddrID(id byte) bool {
	_, ok := scriptHashAddrIDs[id]
	return ok
}

func IsBech32SegwitPrefix(prefix string) bool {
	prefix = strings.ToLower(prefix)
	_, ok := bech32SegwitPrefixes[prefix]
	return ok
}

type DNSSeed struct {
	Host         string
	HasFiltering bool
}

type Checkpoint struct {
	Height int32
	Hash   *Hash
}

type ConsensusDeployment struct {
	BitNumber  uint8
	StartTime  uint64
	ExpireTime uint64
}

type MsgBlock struct {
	Header       BlockHeader
	Transactions []*MsgTx
}

type BlockHeader struct {
	Version    int32
	PrevBlock  Hash
	MerkleRoot Hash
	Timestamp  time.Time
	Bits       uint32
	Nonce      uint32
}

type MsgTx struct {
	Version  int32
	TxIn     []*TxIn
	TxOut    []*TxOut
	LockTime uint32
}

type TxIn struct {
	PreviousOutPoint OutPoint
	SignatureScript  []byte
	Sequence         uint32
}

type OutPoint struct {
	Hash  Hash
	Index uint32
}

type TxOut struct {
	Value    int64
	PkScript []byte
}

type Params struct {
	Name         string
	Net          BitcoinNet
	DefaultPort  string
	DNSSeeds     []DNSSeed
	GenesisBlock *MsgBlock
	GenesisHash  *Hash
	PowLimit     *big.Int
	PowLimitBits uint32

	BIP0034Height int32
	BIP0065Height int32
	BIP0066Height int32

	UahfForkHeight int32
	DaaForkHeight  int32

	MagneticAnomalyActivationTime uint64
	GreatWallActivationTime       uint64
	CoinbaseMaturity              uint16
	SubsidyReductionInterval      int32
	TargetTimespan                time.Duration
	TargetTimePerBlock            time.Duration
	RetargetAdjustmentFactor      int64
	ReduceMinDifficulty           bool
	MinDiffReductionTime          time.Duration
	GenerateSupported             bool
	Checkpoints                   []Checkpoint
	RuleChangeActivationThreshold uint32
	MinerConfirmationWindow       uint32
	Deployments                   []ConsensusDeployment

	// Mempool parameters
	RelayNonStdTxs  bool
	Bech32HRPSegwit string

	// Address encoding magics
	PubKeyHashAddrID        byte // First byte of a P2PKH address
	ScriptHashAddrID        byte // First byte of a P2SH address
	WitnessPubKeyHashAddrID byte // First byte of a P2WPKH address
	WitnessScriptHashAddrID byte // First byte of a P2WSH address
	// The prefix used for the cashaddress. This is different for each network.
	CashAddressPrefix string

	// Address encoding magics
	LegacyPubKeyHashAddrID byte // First byte of a P2PKH address
	LegacyScriptHashAddrID byte // First byte of a P2SH address
	PrivateKeyID           byte // First byte of a WIF private key

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID [4]byte
	HDPublicKeyID  [4]byte

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType uint32
}

func newHashFromStr(hexStr string) *Hash {
	hash, err := NewHashFromStr(hexStr)
	if err != nil {
		panic(err)
	}
	return hash
}

const HashSize = 32

type Hash [HashSize]byte

// MaxHashStringSize is the maximum length of a Hash hash string.
const MaxHashStringSize = HashSize * 2

// ErrHashStrSize describes an error that indicates the caller specified a hash
// string that has too many characters.
var ErrHashStrSize = fmt.Errorf("max hash string length is %v bytes", MaxHashStringSize)

func (hash Hash) String() string {
	for i := 0; i < HashSize/2; i++ {
		hash[i], hash[HashSize-1-i] = hash[HashSize-1-i], hash[i]
	}
	return hex.EncodeToString(hash[:])
}

// CloneBytes returns a copy of the bytes which represent the hash as a byte
// slice.
//
// NOTE: It is generally cheaper to just slice the hash directly thereby reusing
// the same bytes rather than calling this method.
func (hash *Hash) CloneBytes() []byte {
	newHash := make([]byte, HashSize)
	copy(newHash, hash[:])

	return newHash
}

// SetBytes sets the bytes which represent the hash.  An error is returned if
// the number of bytes passed in is not HashSize.
func (hash *Hash) SetBytes(newHash []byte) error {
	nhlen := len(newHash)
	if nhlen != HashSize {
		return fmt.Errorf("invalid hash length of %v, want %v", nhlen,
			HashSize)
	}
	copy(hash[:], newHash)

	return nil
}

// IsEqual returns true if target is the same as hash.
func (hash *Hash) IsEqual(target *Hash) bool {
	if hash == nil && target == nil {
		return true
	}
	if hash == nil || target == nil {
		return false
	}
	return *hash == *target
}

// NewHash returns a new Hash from a byte slice.  An error is returned if
// the number of bytes passed in is not HashSize.
func NewHash(newHash []byte) (*Hash, error) {
	var sh Hash
	err := sh.SetBytes(newHash)
	if err != nil {
		return nil, err
	}
	return &sh, err
}

// NewHashFromStr creates a Hash from a hash string.  The string should be
// the hexadecimal string of a byte-reversed hash, but any missing characters
// result in zero padding at the end of the Hash.
func NewHashFromStr(hash string) (*Hash, error) {
	ret := new(Hash)
	err := Decode(ret, hash)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// Decode decodes the byte-reversed hexadecimal string encoding of a Hash to a
// destination.
func Decode(dst *Hash, src string) error {
	// Return error if hash string is too long.
	if len(src) > MaxHashStringSize {
		return ErrHashStrSize
	}

	// Hex decoder expects the hash to be a multiple of two.  When not, pad
	// with a leading zero.
	var srcBytes []byte
	if len(src)%2 == 0 {
		srcBytes = []byte(src)
	} else {
		srcBytes = make([]byte, 1+len(src))
		srcBytes[0] = '0'
		copy(srcBytes[1:], src)
	}

	// Hex decode the source bytes to a temporary destination.
	var reversedHash Hash
	_, err := hex.Decode(reversedHash[HashSize-hex.DecodedLen(len(srcBytes)):], srcBytes)
	if err != nil {
		return err
	}

	// Reverse copy from the temporary hash to destination.  Because the
	// temporary was zeroed, the written result will be correctly padded.
	for i, b := range reversedHash[:HashSize/2] {
		dst[i], dst[HashSize-1-i] = reversedHash[HashSize-1-i], b
	}

	return nil
}
