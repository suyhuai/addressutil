package chaincfg

import (
	"math"
	"math/big"
	"time"

	"github.com/suyhuai/addressutil/util"
)

const (
	MainNet  util.BitcoinNet = 0xe8f3e1e3
	TestNet  util.BitcoinNet = 0xfabfb5da
	TestNet3 util.BitcoinNet = 0xf4f3e5f4
	SimNet   util.BitcoinNet = 0x12141c16
)

const (
	DeploymentTestDummy = iota
	DeploymentCSV
	DefinedDeployments
)

var (
	bigOne             = big.NewInt(1)
	mainPowLimit       = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 224), bigOne)
	regressionPowLimit = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 255), bigOne)
	testNet3PowLimit   = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 224), bigOne)
	simNetPowLimit     = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 255), bigOne)
)

var (
	registeredNets      = make(map[util.BitcoinNet]struct{})
	pubKeyHashAddrIDs   = make(map[byte]struct{})
	scriptHashAddrIDs   = make(map[byte]struct{})
	cashAddressPrefixes = make(map[string]struct{})
	hdPrivToPubKeyIDs   = make(map[[4]byte][]byte)
)

func IsPubKeyHashAddrID(id byte) bool {
	_, ok := pubKeyHashAddrIDs[id]
	return ok
}

func IsScriptHashAddrID(id byte) bool {
	_, ok := scriptHashAddrIDs[id]
	return ok
}

var MainNetParams = util.Params{
	Name:        "mainnet",
	Net:         MainNet,
	DefaultPort: "8333",
	DNSSeeds: []util.DNSSeed{
		{"seed.bitcoinabc.org", true},
		{"seed-abc.bitcoinforks.org", true},
		{"btccash-seeder.bitcoinunlimited.info", true},
		{"seed.bitprim.org", true},
		{"seed.deadalnix.me", true},
		{"seeder.criptolayer.net", true},
	},

	// Chain parameters
	GenesisBlock:  &genesisBlock,
	GenesisHash:   &genesisHash,
	PowLimit:      mainPowLimit,
	PowLimitBits:  0x1d00ffff,
	BIP0034Height: 227931, // 000000000000024b89b42a942fe0d9fea3bb44ab7bd1b19115dd6a759c0808b8
	BIP0065Height: 388381, // 000000000000000004c2b624ed5d7756c508d90fd0da2c7c679febfa6c4735f0
	BIP0066Height: 363725, // 00000000000000000379eaa19dce8c9b722d46ae6a57c2f1a988119488b50931

	UahfForkHeight: 478558, // 0000000000000000011865af4122fe3b144e2cbeea86142e8ff2fb4107352d43
	DaaForkHeight:  504031, // 0000000000000000011ebf65b60d0a3de80b8175be709d653b4c1a1beeb6ab9c

	CoinbaseMaturity:         100,
	SubsidyReductionInterval: 210000,
	TargetTimespan:           time.Hour * 24 * 14, // 14 days
	TargetTimePerBlock:       time.Minute * 10,    // 10 minutes
	RetargetAdjustmentFactor: 4,                   // 25% less, 400% more
	ReduceMinDifficulty:      false,
	MinDiffReductionTime:     0,
	GenerateSupported:        false,

	// Checkpoints ordered from oldest to newest.
	Checkpoints: []util.Checkpoint{
		{11111, newHashFromStr("0000000069e244f73d78e8fd29ba2fd2ed618bd6fa2ee92559f542fdb26e7c1d")},
		{33333, newHashFromStr("000000002dd5588a74784eaa7ab0507a18ad16a236e7b1ce69f00d7ddfb5d0a6")},
		{74000, newHashFromStr("0000000000573993a3c9e41ce34471c079dcf5f52a0e824a81e7f953b8661a20")},
		{105000, newHashFromStr("00000000000291ce28027faea320c8d2b054b2e0fe44a773f3eefb151d6bdc97")},
		{134444, newHashFromStr("00000000000005b12ffd4cd315cd34ffd4a594f430ac814c91184a0d42d2b0fe")},
		{168000, newHashFromStr("000000000000099e61ea72015e79632f216fe6cb33d7899acb35b75c8303b763")},
		{193000, newHashFromStr("000000000000059f452a5f7340de6682a977387c17010ff6e6c3bd83ca8b1317")},
		{210000, newHashFromStr("000000000000048b95347e83192f69cf0366076336c639f9b7228e9ba171342e")},
		{216116, newHashFromStr("00000000000001b4f4b433e81ee46494af945cf96014816a4e2370f11b23df4e")},
		{225430, newHashFromStr("00000000000001c108384350f74090433e7fcf79a606b8e797f065b130575932")},
		{250000, newHashFromStr("000000000000003887df1f29024b06fc2200b55f8af8f35453d7be294df2d214")},
		{267300, newHashFromStr("000000000000000a83fbd660e918f218bf37edd92b748ad940483c7c116179ac")},
		{279000, newHashFromStr("0000000000000001ae8c72a0b0c301f67e3afca10e819efa9041e458e9bd7e40")},
		{300255, newHashFromStr("0000000000000000162804527c6e9b9f0563a280525f9d08c12041def0a0f3b2")},
		{319400, newHashFromStr("000000000000000021c6052e9becade189495d1c539aa37c58917305fd15f13b")},
		{343185, newHashFromStr("0000000000000000072b8bf361d01a6ba7d445dd024203fafc78768ed4368554")},
		{352940, newHashFromStr("000000000000000010755df42dba556bb72be6a32f3ce0b6941ce4430152c9ff")},
		{382320, newHashFromStr("00000000000000000a8dc6ed5b133d0eb2fd6af56203e4159789b092defd8ab2")},
	},

	RuleChangeActivationThreshold: 1916, // 95% of MinerConfirmationWindow
	MinerConfirmationWindow:       2016, //
	Deployments: []util.ConsensusDeployment{
		DeploymentTestDummy: {
			BitNumber:  28,
			StartTime:  1199145601, // January 1, 2008 UTC
			ExpireTime: 1230767999, // December 31, 2008 UTC
		},
		DeploymentCSV: {
			BitNumber:  0,
			StartTime:  1462060800, // May 1st, 2016
			ExpireTime: 1493596800, // May 1st, 2017
		},
	},

	// Mempool parameters
	RelayNonStdTxs: false,

	// The prefix for the cashaddress
	CashAddressPrefix: "bitcoincash", // always bitcoincash for mainnet

	// Address encoding magics
	LegacyPubKeyHashAddrID: 0x00, // starts with 1
	LegacyScriptHashAddrID: 0x05, // starts with 3
	PrivateKeyID:           0x80, // starts with 5 (uncompressed) or K (compressed)

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x88, 0xad, 0xe4}, // starts with xprv
	HDPublicKeyID:  [4]byte{0x04, 0x88, 0xb2, 0x1e}, // starts with xpub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 145,
}

var RegressionNetParams = util.Params{
	Name:        "regtest",
	Net:         TestNet,
	DefaultPort: "18444",
	DNSSeeds:    []util.DNSSeed{},

	// Chain parameters
	GenesisBlock:     &regTestGenesisBlock,
	GenesisHash:      &regTestGenesisHash,
	PowLimit:         regressionPowLimit,
	PowLimitBits:     0x207fffff,
	CoinbaseMaturity: 100,
	BIP0034Height:    100000000, // Not active - Permit ver 1 blocks
	BIP0065Height:    1351,      // Used by regression tests
	BIP0066Height:    1251,      // Used by regression tests

	UahfForkHeight: 0, // Always active on regtest
	DaaForkHeight:  0, // Always active on regtest

	SubsidyReductionInterval: 150,
	TargetTimespan:           time.Hour * 24 * 14, // 14 days
	TargetTimePerBlock:       time.Minute * 10,    // 10 minutes
	RetargetAdjustmentFactor: 4,                   // 25% less, 400% more
	ReduceMinDifficulty:      true,
	MinDiffReductionTime:     time.Minute * 20, // TargetTimePerBlock * 2
	GenerateSupported:        true,

	// Checkpoints ordered from oldest to newest.
	Checkpoints: nil,

	RuleChangeActivationThreshold: 108, // 75%  of MinerConfirmationWindow
	MinerConfirmationWindow:       144,
	Deployments: []util.ConsensusDeployment{
		DeploymentTestDummy: {
			BitNumber:  28,
			StartTime:  0,             // Always available for vote
			ExpireTime: math.MaxInt64, // Never expires
		},
		DeploymentCSV: {
			BitNumber:  0,
			StartTime:  0,             // Always available for vote
			ExpireTime: math.MaxInt64, // Never expires
		},
	},

	// Mempool parameters
	RelayNonStdTxs: true,

	// The prefix for the cashaddress
	CashAddressPrefix: "bchreg", // always bchreg for reg testnet

	// Address encoding magics
	LegacyPubKeyHashAddrID: 0x6f, // starts with m or n
	LegacyScriptHashAddrID: 0xc4, // starts with 2
	PrivateKeyID:           0xef, // starts with 9 (uncompressed) or c (compressed)

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with tprv
	HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with tpub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 1, // all coins use 1
}

var TestNet3Params = util.Params{
	Name:        "testnet3",
	Net:         TestNet3,
	DefaultPort: "18333",
	DNSSeeds: []util.DNSSeed{
		{"testnet-seed.bitcoinabc.org", true},
		{"testnet-seed-abc.bitcoinforks.org", true},
		{"testnet-seed.bitprim.org", true},
		{"testnet-seed.deadalnix.me", true},
		{"testnet-seeder.criptolayer.net", true},
	},

	// Chain parameters
	GenesisBlock:  &testNet3GenesisBlock,
	GenesisHash:   &testNet3GenesisHash,
	PowLimit:      testNet3PowLimit,
	PowLimitBits:  0x1d00ffff,
	BIP0034Height: 21111,  // 0000000023b3a96d3484e5abb3755c413e7d41500f8e2a5c3f0dd01299cd8ef8
	BIP0065Height: 581885, // 00000000007f6655f22f98e72ed80d8b06dc761d5da09df0fa1dc4be4f861eb6
	BIP0066Height: 330776, // 000000002104c8c45e99a8853285a3b592602a3ccde2b832481da85e9e4ba182

	UahfForkHeight: 1155875, // 00000000f17c850672894b9a75b63a1e72830bbd5f4c8889b5c1a80e7faef138
	DaaForkHeight:  1188697, // 0000000000170ed0918077bde7b4d36cc4c91be69fa09211f748240dabe047fb

	CoinbaseMaturity:         100,
	SubsidyReductionInterval: 210000,
	TargetTimespan:           time.Hour * 24 * 14, // 14 days
	TargetTimePerBlock:       time.Minute * 10,    // 10 minutes
	RetargetAdjustmentFactor: 4,                   // 25% less, 400% more
	ReduceMinDifficulty:      true,
	MinDiffReductionTime:     time.Minute * 20, // TargetTimePerBlock * 2
	GenerateSupported:        false,

	// Checkpoints ordered from oldest to newest.
	Checkpoints: []util.Checkpoint{
		{546, newHashFromStr("000000002a936ca763904c3c35fce2f3556c559c0214345d31b1bcebf76acb70")},
		{100000, newHashFromStr("00000000009e2958c15ff9290d571bf9459e93b19765c6801ddeccadbb160a1e")},
		{200000, newHashFromStr("0000000000287bffd321963ef05feab753ebe274e1d78b2fd4e2bfe9ad3aa6f2")},
		{300001, newHashFromStr("0000000000004829474748f3d1bc8fcf893c88be255e6d7f571c548aff57abf4")},
		{400002, newHashFromStr("0000000005e2c73b8ecb82ae2dbc2e8274614ebad7172b53528aba7501f5a089")},
		{500011, newHashFromStr("00000000000929f63977fbac92ff570a9bd9e7715401ee96f2848f7b07750b02")},
		{600002, newHashFromStr("000000000001f471389afd6ee94dcace5ccc44adc18e8bff402443f034b07240")},
		{700000, newHashFromStr("000000000000406178b12a4dea3b27e13b3c4fe4510994fd667d7c1e6a3f4dc1")},
		{800010, newHashFromStr("000000000017ed35296433190b6829db01e657d80631d43f5983fa403bfdb4c1")},
		{900000, newHashFromStr("0000000000356f8d8924556e765b7a94aaebc6b5c8685dcfa2b1ee8b41acd89b")},
		{1000007, newHashFromStr("00000000001ccb893d8a1f25b70ad173ce955e5f50124261bbbc50379a612ddf")},
	},

	// Consensus rule change deployments.
	//
	// The miner confirmation window is defined as:
	//   target proof of work timespan / target proof of work spacing
	RuleChangeActivationThreshold: 1512, // 75% of MinerConfirmationWindow
	MinerConfirmationWindow:       2016,
	Deployments: []util.ConsensusDeployment{
		DeploymentTestDummy: {
			BitNumber:  28,
			StartTime:  1199145601, // January 1, 2008 UTC
			ExpireTime: 1230767999, // December 31, 2008 UTC
		},
		DeploymentCSV: {
			BitNumber:  0,
			StartTime:  1456790400, // March 1st, 2016
			ExpireTime: 1493596800, // May 1st, 2017
		},
	},

	// Mempool parameters
	RelayNonStdTxs: true,

	// The prefix for the cashaddress
	CashAddressPrefix: "bchtest", // always bchtest for testnet

	// Address encoding magics
	LegacyPubKeyHashAddrID: 0x6f, // starts with m or n
	LegacyScriptHashAddrID: 0xc4, // starts with 2
	PrivateKeyID:           0xef, // starts with 9 (uncompressed) or c (compressed)

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with tprv
	HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with tpub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 1, // all coins use 1
}

var SimNetParams = util.Params{
	Name:        "simnet",
	Net:         SimNet,
	DefaultPort: "18555",
	DNSSeeds:    []util.DNSSeed{}, // NOTE: There must NOT be any seeds.

	// Chain parameters
	GenesisBlock:             &simNetGenesisBlock,
	GenesisHash:              &simNetGenesisHash,
	PowLimit:                 simNetPowLimit,
	PowLimitBits:             0x207fffff,
	BIP0034Height:            0, // Always active on simnet
	BIP0065Height:            0, // Always active on simnet
	BIP0066Height:            0, // Always active on simnet
	UahfForkHeight:           0, // Always active on simnet
	DaaForkHeight:            0, // Always active on simnet
	CoinbaseMaturity:         100,
	SubsidyReductionInterval: 210000,
	TargetTimespan:           time.Hour * 24 * 14, // 14 days
	TargetTimePerBlock:       time.Minute * 10,    // 10 minutes
	RetargetAdjustmentFactor: 4,                   // 25% less, 400% more
	ReduceMinDifficulty:      true,
	MinDiffReductionTime:     time.Minute * 20, // TargetTimePerBlock * 2
	GenerateSupported:        true,

	// Checkpoints ordered from oldest to newest.
	Checkpoints: nil,

	// Consensus rule change deployments.
	//
	// The miner confirmation window is defined as:
	//   target proof of work timespan / target proof of work spacing
	RuleChangeActivationThreshold: 75, // 75% of MinerConfirmationWindow
	MinerConfirmationWindow:       100,
	Deployments: []util.ConsensusDeployment{
		DeploymentTestDummy: {
			BitNumber:  28,
			StartTime:  0,             // Always available for vote
			ExpireTime: math.MaxInt64, // Never expires
		},
		DeploymentCSV: {
			BitNumber:  0,
			StartTime:  0,             // Always available for vote
			ExpireTime: math.MaxInt64, // Never expires
		},
	},

	RelayNonStdTxs: true,

	CashAddressPrefix: "bchsim", // always bchsim for simnet

	LegacyPubKeyHashAddrID: 0x3f, // starts with S
	LegacyScriptHashAddrID: 0x7b, // starts with s
	PrivateKeyID:           0x64, // starts with 4 (uncompressed) or F (compressed)

	HDPrivateKeyID: [4]byte{0x04, 0x20, 0xb9, 0x00}, // starts with sprv
	HDPublicKeyID:  [4]byte{0x04, 0x20, 0xbd, 0x3a}, // starts with spub

	HDCoinType: 115, // ASCII for s
}

func Register(params *util.Params) error {
	if _, ok := registeredNets[params.Net]; ok {
		return util.ErrDuplicateNet
	}
	registeredNets[params.Net] = struct{}{}
	pubKeyHashAddrIDs[params.LegacyPubKeyHashAddrID] = struct{}{}
	scriptHashAddrIDs[params.LegacyScriptHashAddrID] = struct{}{}
	hdPrivToPubKeyIDs[params.HDPrivateKeyID] = params.HDPublicKeyID[:]

	cashAddressPrefixes[params.CashAddressPrefix+":"] = struct{}{}
	return nil
}

func newHashFromStr(hexStr string) *util.Hash {
	hash, err := util.NewHashFromStr(hexStr)
	if err != nil {
		panic(err)
	}
	return hash
}

func mustRegister(params *util.Params) {
	if err := Register(params); err != nil {
		panic("failed to register network: " + err.Error())
	}
}

func init() {
	mustRegister(&MainNetParams)
	mustRegister(&TestNet3Params)
	mustRegister(&RegressionNetParams)
	mustRegister(&SimNetParams)
}
