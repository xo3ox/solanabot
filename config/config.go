package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
)

const (
	SOL_MINT  = "So11111111111111111111111111111111111111112"  // solana测试网
	USDC_MINT = "4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU" // solana测试网 USDC
	USDT_MINT = "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB" // solana测试网 USDT
)

var (
	PRIVATE_KEY solana.PrivateKey
	PUBLIC_KEY  solana.PublicKey

	RPC_CLIENT = rpc.New(rpc.DevNet_RPC) // DevNet_RPC 开发环境
	CTX        = context.Background()
)

func init() {
	PRIVATE_KEY, PUBLIC_KEY = getKeypair("../.env") // 从 .env 文件中读取私钥
}

// getKeypair 从 .env 文件中读取私钥，并返回私钥和公钥（地址）
func getKeypair(filenames ...string) (solana.PrivateKey, solana.PublicKey) {
	// 加载 .env 文件
	if err := godotenv.Load(filenames...); err != nil {
		log.Fatalf("读取 .env 文件失败,err:%v", err)
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

	return privateKey, publicKey
}
