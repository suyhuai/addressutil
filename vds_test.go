package addressutil

import (
	"encoding/hex"
	"testing"
)

var addrPub = map[string]string{
	"VcoCWoY2gJ3QA2NpCR1LgfirvAj6cE48947": "042479618bb2f3835761c880a966f279521856152f009504dfb98825f178865ff576d98a862a8526559ff21be4db0b712967a1e60a9346b615e202c1af1c45cebe",
	"VcSxoFRgkT6ESdPejX7CRSLpdpynvS2zgi5": "047e6a42c4234499c4f5d63beeeaf067aea25456bcf0bb55abb7a4cac164d7841ac3559dc5134aafc304f8ac9725159c9ec10951bcff1262b2a00a0f69c86c87fb", // ef7cc0fcd4acb523decd3c7b0ddc78ff620f6a957e1da51aff0b5f160c84e975
	"Vcodt5VwJGs6wtksfSsg1N9Ah4iiFa27RZt": "04d95abb27882b41a21c2d06f42ecf7e5e84ccd2fd1b1ceab0e981455625f23911fb5fddd681829d74caa7aecf0a7842f8e3f8670b0bced9c830f66455656f7b5e", // d94f077428ea1a78c75c251def821b3829529fdfef3803e81cffb42e8fbf5ff1
	"VcZYPwNbZAVSAPJdtXyNWmEfDd9qfyhuYRk": "04f4f43c002317d56813b869341b07863f3b2517be0ee00d9f01424f3d0638dbe66ba5047fbeee41ef49e95a0919d21812c9bbf291944af6041325117f83094755", // 51788e4716fcb6d4914d3412bbb4a3ce33e8167a0aee6f3ef18de824799435f0
}

func TestNewVdsAddress(t *testing.T) {
	for addr, pubKey := range addrPub {
		pub, err := hex.DecodeString(pubKey)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		if a, err := NewVDSAddress(pub); err != nil {
			t.Log(err)
			t.Fail()
		} else if a.String() != addr {
			t.Log("Address mismatch", a, addr)
			t.Fail()
		}
	}

}

func TestCheckVdsAddress(t *testing.T) {
	for addr := range addrPub {
		if !CheckVDSAddress(addr) {
			t.Log("check address fail: ", addr)
			t.Fail()
		}
	}

}
