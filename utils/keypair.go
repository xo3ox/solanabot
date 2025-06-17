package utils

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
)

func NewRandomPrivateKey() {
	privateKeyStr, err := solana.NewRandomPrivateKey()
	if err != nil {
		fmt.Println("生成随机私钥失败: %w", err)
	}

	publicKey := privateKeyStr.PublicKey()

	fmt.Println("privateKeyByte", []byte(privateKeyStr))
	fmt.Println("privateKeyStr  ", privateKeyStr.String())
	fmt.Println("publicKey", publicKey)

}
