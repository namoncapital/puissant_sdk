package main

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/namoncapital/puissant_sdk/demo"
)

func main() {
	conf, client := demo.GetClient()

	chainID, err := client.General.ChainID(context.Background())
	if err != nil {
		log.Panicln(err.Error())
	}
	log.Printf("chainID: %s", chainID.String())

	countWallet := len(conf.Wallet)
	gasPrice, _ := client.SuggestGasPrice(context.Background())
	value := big.NewInt(1e17)

	var rawTxs []hexutil.Bytes
	// var txs []*types.Transaction
	for i := 0; i < countWallet; i++ {
		pk := conf.Wallet[i]
		privateKey, fromAddress := demo.StrToPK(pk)
		nonce, err := client.General.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			log.Panicln(err.Error())
		}

		gasLimit := uint64(21000)
		toPk := conf.Wallet[0]
		if i < countWallet-1 {
			toPk = conf.Wallet[i+1]
		}
		_, toAddress := demo.StrToPK(toPk)
		// gas sort
		thisGas := big.NewInt(0).Add(big.NewInt(int64(countWallet-i)), gasPrice)
		tx := types.NewTransaction(nonce, toAddress, value, gasLimit, thisGas, nil)
		// change next value
		value.Sub(value, thisGas.Mul(thisGas, big.NewInt(int64(gasLimit))))
		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
		if err != nil {
			log.Panicln(err.Error())
		}
		// txs = append(txs, signedTx)
		rawTxBytes, _ := rlp.EncodeToBytes(signedTx)

		rawTxs = append(rawTxs, rawTxBytes)
	}

	// send puissant tx
	res, err := client.SendPuissant(context.Background(), rawTxs, uint64(time.Now().Unix()+60), nil)
	// res, err := client.SendPuissantTxs(context.Background(), txs, time.Now().Unix()+60, nil)
	if err != nil {
		log.Panicln(err.Error())
	}

	log.Println(res)
}
