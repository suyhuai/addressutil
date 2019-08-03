package addressutil

import (
	"encoding/hex"
	"testing"
)

func TestLTCAddress(t *testing.T) {
	keys := map[string]string{
		"LdNQoxcHSqEX6jLpRH6V5op12uF9KE5KYY": "042303bd28496c2f9ca2c8530737fa9059b9277dfafe66a35c321a439490e05cb0b49c04207f9d61124c1f7b30ae578900df60ce778806fb358d4041ba28411d73",
	}

	for addr, pubKey := range keys {
		pub, _ := hex.DecodeString(pubKey)
		if a, err := NewLTCAddress(pub, true); err != nil {
			t.Log(err)
			t.Fail()
		} else if a.String() != addr {
			t.Log("Address mismatch", a, addr)
			t.Fail()
		}
	}

}

func TestCheckLTCAddress(t *testing.T) {
	validAddresses := []string{
		"LiPgsBxvBBDR6TaTYnUvwML7rb3dTbMECK", // 普通地址
		"MLro9kXkvYRe4nHfCnWmaW94gHr6GgrqnL", // 从3开头的编码而得，可以和3开头的相互转换
	}
	invalidAddresses := []string{
		"LiPgsBxvBBDR6TaTYnUvwML7rb3dTbMECL",
		"MLro9kXkvYRe4nHfCnWmaW94gHr6GgrqnM",
	}

	for _, address := range validAddresses {
		if !CheckLTCAddress(address, true) {
			t.Fail()
		}
	}

	for _, address := range invalidAddresses {
		if CheckLTCAddress(address, true) {
			t.Fail()
		}
	}
}
