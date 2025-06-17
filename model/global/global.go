package global

import (
	"solanabot/utils"

	"github.com/gagliardetto/solana-go"
)

var (
	PRIVATE_KEY solana.PrivateKey
)

func init() {
	PRIVATE_KEY = utils.GetKeypair()
}
