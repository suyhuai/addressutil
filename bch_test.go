package addressutil

import (
	"testing"
)

func TestBCHAddress(t *testing.T) {
	keys := map[string][]byte{
		"bitcoincash:qp4qttt9er93g0qwqtemzyw96d7p5jc25q4thz52r4": []byte{0x04, 0x78, 0x14, 0x04, 0x9c, 0xd3, 0x23, 0xb2, 0xf7, 0x07, 0x4c, 0x94, 0xed, 0xc0, 0xf9, 0x61, 0xdb, 0x62, 0xbe, 0x35, 0x68, 0x7c, 0x24, 0x24, 0xb2, 0xad, 0x29, 0xf7, 0xf2, 0x83, 0x7e, 0x03, 0x95, 0xe8, 0xeb, 0xbe, 0xe0, 0x4f, 0x81, 0x98, 0xa9, 0x1c, 0xd9, 0xc8, 0xac, 0xbd, 0x97, 0xaa, 0xd7, 0x68, 0x51, 0x95, 0x7f, 0x3e, 0x63, 0x62, 0xf4, 0xdd, 0x41, 0x92, 0xf3, 0x43, 0x1a, 0x71, 0xfb},
		"bitcoincash:qqtrsdux9q7zuzxn4vpmfqh9tdlm2v9n2vdnkqmpv4": []byte{0x04, 0xe9, 0x16, 0xbd, 0xa7, 0xec, 0x60, 0xfc, 0x10, 0xf1, 0x68, 0x99, 0x32, 0xef, 0x85, 0x92, 0xb7, 0x7d, 0xca, 0x3b, 0x13, 0xd4, 0x8a, 0x82, 0x70, 0x76, 0x2e, 0x07, 0x8f, 0xd9, 0x30, 0xf4, 0x59, 0x8f, 0x53, 0x0b, 0xaf, 0x34, 0x64, 0x68, 0x38, 0x5a, 0x6c, 0xdb, 0x0d, 0x14, 0x3f, 0xce, 0xce, 0x71, 0x8b, 0xf5, 0x9a, 0x8d, 0x83, 0xa4, 0xf9, 0x38, 0x53, 0xb2, 0x95, 0x45, 0x2b, 0xa5, 0x96},
		"bitcoincash:qql9qat3d9qngeeukgyeyv6u22g4pv7qwvk5nm4p5t": []byte{0x04, 0x8d, 0x17, 0x4d, 0xea, 0x66, 0x86, 0x85, 0xad, 0x53, 0x23, 0x81, 0x80, 0xab, 0x36, 0xe4, 0x42, 0xb2, 0x50, 0x21, 0x07, 0xfb, 0xfb, 0xbf, 0x60, 0xc9, 0x26, 0x03, 0xb6, 0x8b, 0x0b, 0x03, 0xd0, 0x79, 0x38, 0xe6, 0x2d, 0x12, 0xfb, 0xff, 0x2d, 0x01, 0x70, 0xd6, 0x16, 0x69, 0x66, 0x5b, 0x37, 0xd7, 0x6e, 0xb1, 0x1b, 0x79, 0x65, 0x1f, 0x2f, 0xb9, 0x64, 0xcd, 0x0b, 0x40, 0x23, 0x28, 0x69},
		"bitcoincash:qqgqcs2794628x94uq4fqjynqswc7kl06syxq83env": []byte{0x04, 0xb2, 0xb6, 0x88, 0xbc, 0x77, 0xec, 0xec, 0xfc, 0x6b, 0x74, 0xb0, 0xf0, 0xd1, 0xda, 0xd7, 0x09, 0xc8, 0xcc, 0x06, 0xd0, 0xe3, 0xad, 0x0e, 0x92, 0x6d, 0x21, 0xca, 0x93, 0x90, 0x74, 0x8a, 0x0b, 0xb3, 0xb0, 0xd0, 0xed, 0xce, 0x71, 0xde, 0x40, 0xe8, 0x38, 0x13, 0x4b, 0xfe, 0x38, 0x6b, 0x8b, 0x1e, 0xf4, 0x81, 0xe0, 0x7e, 0x66, 0x2b, 0x2e, 0x6f, 0x22, 0x95, 0xef, 0x44, 0x53, 0xc9, 0xd0},
	}

	for addr, pubKey := range keys {
		if a, err := NewBCHAddress(pubKey, true); err != nil {
			t.Log(err)
			t.Fail()
		} else if a.String() != addr {
			t.Log("Address mismatch", a, addr)
			t.Fail()
		}
	}
}

func TestCheckBCHAddress(t *testing.T) {
	validAddresses := []string{
		"1BMfnF2h2absXr4JjNMzFeB1XP97NnNXfs",         // 普通地址
		"qpcenuhjnwk0xw4st4x0pyn69vmra29nnvghrpm8jg", // 新版格式
		"37JY6K2gw5rRSkahvCC6maDjRKeMdGSC51",         // 隔离见证或多重签名地址：特殊判断
	}

	invalidAddresses := []string{
		"1BMfnF2h2absXr4JjNMzFeB1XP97NnNXft",
		"qpcenuhjnwk0xw4st4x0pyn69vmra29nnvghrpm8jh",
		"37JY6K2gw5rRSkahvCC6maDjRKeMdGSC52",
	}

	for _, address := range validAddresses {
		if !CheckBCHAddress(address, true) {
			t.Fail()
		}
	}

	for _, address := range invalidAddresses {
		if CheckBCHAddress(address, true) {
			t.Fail()
		}
	}
}
