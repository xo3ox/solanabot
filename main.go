package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"solanabot/config"

	bi "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

const (
	JupiterAPIQuoteURL = "https://quote-api.jup.ag/v6/quote"
	JupiterAPISwapURL  = "https://quote-api.jup.ag/v6/swap"
)

type QuoteResponse struct {
	InAmount   string `json:"inAmount"`
	OutAmount  string `json:"outAmount"`
	InputMint  string `json:"inputMint"`
	OutputMint string `json:"outputMint"`
	Amount     string `json:"amount"`
}

type SwapResponse struct {
	SwapTransaction string `json:"swapTransaction"`
}

func getQuote(inputMint, outputMint string, amount uint64) (*QuoteResponse, []byte, error) {
	u, _ := url.Parse(JupiterAPIQuoteURL)
	q := u.Query()
	q.Set("inputMint", inputMint)
	q.Set("outputMint", outputMint)
	q.Set("amount", fmt.Sprintf("%d", amount))
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != 200 {
		return nil, body, fmt.Errorf("quote failed: status %d: %s", resp.StatusCode, string(body))
	}

	var quote QuoteResponse
	if err := json.Unmarshal(body, &quote); err != nil {
		return nil, body, err
	}
	return &quote, body, nil
}

func requestSwapTransaction(pubkey string, quoteRaw []byte) (*SwapResponse, error) {
	bodyMap := map[string]json.RawMessage{
		"userPublicKey": json.RawMessage(fmt.Sprintf("\"%s\"", pubkey)),
		"wrapUnwrapSOL": json.RawMessage("true"),
		"slippageBps":   json.RawMessage("50"),
		"quoteResponse": quoteRaw,
	}
	reqBody, _ := json.Marshal(bodyMap)

	resp, err := http.Post(JupiterAPISwapURL, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status: %d, body: %s", resp.StatusCode, string(body))
	}

	var swapResp SwapResponse
	if err := json.Unmarshal(body, &swapResp); err != nil {
		return nil, err
	}
	return &swapResp, nil
}

func main() {
	ctx := context.Background()
	rpcClient := rpc.New(rpc.DevNet_RPC)

	// token mint address
	inputMint := "So11111111111111111111111111111111111111112" // SOL

	outputMint := "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v" // USDC
	amount := uint64(1_000_000_000)                              // 1 SOL

	quote1, rawQuote1, err := getQuote(inputMint, outputMint, amount)
	if err != nil {
		log.Fatalf("GetQuote1 error: %v", err)
	}

	outAmount1, ok := new(big.Int).SetString(quote1.OutAmount, 10)
	if !ok {
		log.Fatalf("Failed to parse quote1.OutAmount")
	}
	fmt.Printf("1 SOL -> ~%s USDC\n", quote1.OutAmount)

	quote2, _, err := getQuote(outputMint, inputMint, outAmount1.Uint64())
	if err != nil {
		log.Fatalf("GetQuote2 error: %v", err)
	}
	backAmount, _ := new(big.Int).SetString(quote2.OutAmount, 10)
	fmt.Printf("%s USDC -> ~%s SOL\n", quote1.OutAmount, quote2.OutAmount)

	if backAmount.Cmp(big.NewInt(int64(amount))) <= 0 {
		fmt.Println("没有套利空间.")
		return
	}
	fmt.Println("检测到套利机会!")

	swapResp, err := requestSwapTransaction(config.PUBLIC_KEY.String(), rawQuote1)
	if err != nil {
		log.Fatalf("Swap tx request failed: %v", err)
	}

	txBytes, err := base64.StdEncoding.DecodeString(swapResp.SwapTransaction)
	if err != nil {
		log.Fatalf("Failed to decode transaction: %v", err)
	}

	decoder := bi.NewBinDecoder(txBytes)
	tx, err := solana.TransactionFromDecoder(decoder)
	if err != nil {
		log.Fatalf("Failed to decode transaction: %v", err)
	}

	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(config.PUBLIC_KEY) {
			return &config.PRIVATE_KEY
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	sig, err := rpcClient.SendTransactionWithOpts(ctx, tx, rpc.TransactionOpts{
		SkipPreflight:       false,
		PreflightCommitment: rpc.CommitmentFinalized,
	})
	if err != nil {
		log.Fatalf("Transaction send failed: %v", err)
	}
	fmt.Printf("Transaction sent: %s\n", sig.String())
}
