package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/gagliardetto/solana-go"
	"github.com/joho/godotenv"
)

func NewRandomPrivateKey() {
	privateKeyStr, err := solana.NewRandomPrivateKey()
	if err != nil {
		fmt.Println("生成随机私钥失败: %w", err)
	}

	publicKey := privateKeyStr.PublicKey()

	fmt.Println("privateKeyByte", []byte(privateKeyStr))
	fmt.Println("publicKeyStr  ", privateKeyStr.String())
	fmt.Println("publicKey", publicKey)

}

// GetKeypair 从 .env 文件中读取私钥，并返回私钥和公钥（地址）
func GetKeypair(filenames ...string) (string, string) {
	// 加载 .env 文件
	if err := godotenv.Load(filenames...); err != nil {
		log.Fatal("读取 .env 文件失败")
	}

	// 从环境变量读取私钥
	privateKeyStr := os.Getenv("PRIVATE_KEY")
	if privateKeyStr == "" {
		log.Fatal(".env 文件中未设置 PRIVATE_KEY")
	}

	// 将 Base58 字符串转换为私钥对象
	privateKey, err := solana.PrivateKeyFromBase58(privateKeyStr)
	if err != nil {
		log.Fatalf("无效的 private key: %v", err)
	}

	// 获取公钥（地址）
	publicKey := privateKey.PublicKey()
	fmt.Printf("Public Key (Address): %s\n", publicKey)
	fmt.Printf("Private Key ([]byte): %v\n", []byte(privateKey))
	fmt.Printf("Private Key: %s\n", privateKey)

	return privateKey.String(), publicKey.String()
}
