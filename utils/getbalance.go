package utils

import (
	"fmt"
	"math/big"
	"solanabot/config"

	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// 获取账户余额
func GetBalance() error {
	out, err := config.RPC_CLIENT.GetBalance(config.CTX, config.PUBLIC_KEY, rpc.CommitmentFinalized)
	if err != nil {
		return err
	}
	spew.Dump(out)
	spew.Dump(out.Value) // total lamports on the account; 1 sol = 1000000000 lamports

	var lamportsOnAccount = new(big.Float).SetUint64(uint64(out.Value))
	// Convert lamports to sol:
	var solBalance = new(big.Float).Quo(lamportsOnAccount, new(big.Float).SetUint64(solana.LAMPORTS_PER_SOL))

	// WARNING: this is not a precise conversion.
	fmt.Println("◎", solBalance.Text('f', 10))
	return nil
}
