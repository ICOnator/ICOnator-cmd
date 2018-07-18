package main

import (
	"crypto/rand"
	"github.com/ogier/pflag"
	"os"
	"encoding/base64"
	"log"
	"math/big"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"io/ioutil"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"encoding/csv"
	"time"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/crypto"
	"context"
)

const (
	GAS_LIMIT uint64 = 21000
)

func main() {

	var (
		privateKeyCSV 		string
		ethereumAddress		string
		ethereumClientURL 	string
		ethereumNet		 	string
		gasCost				int
		dryRun				bool
	)

	flag := pflag.NewFlagSet("consolidate1000-eth", pflag.ExitOnError)

	flag.Usage = func() {
		println("Usage:")
		println("  consolidate1000-eth [-s --private-keys filename] [-d --dryrun boolean] " +
			"[-e --ethereum-to address] [-u --url scheme://url:port] [-g --gascost gwei] [-N --network main|rinkeby|test]")
		println()
		flag.PrintDefaults()
		println()
	}

	flag.StringVarP(&privateKeyCSV, "private-keys", "s", "", "CSV filename to store the private keys")
	flag.BoolVarP(&dryRun, "dryrun", "d", false, "Don't issue the tx, print the tx to stdout")

	flag.StringVarP(&ethereumAddress, "ethereum-to", "e", "", "The ethereum address where to consolidate the transactions")
	flag.StringVarP(&ethereumClientURL, "url", "u", "http://127.0.0.1:8545", "The ethereum URL where a full node is running")
	flag.IntVarP(&gasCost, "gascost", "g", 1, "Gas cost in gwei, check https://ethgasstation.info/ for current price")
	flag.StringVarP(&ethereumNet, "network", "N", "main", "Specify ethereum network: main,test,rinkeby")
	flag.Parse(os.Args[1:])

	if len(ethereumAddress) == 0 || len(privateKeyCSV) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	password := randStr(20)
	dirName, err:=ioutil.TempDir(os.TempDir(), "consolidate-eth")

	if err != nil {
		log.Fatalf("TempDir failed: %v\n", err)
	}

	//we can use keystore only
	//https://ethereum.stackexchange.com/questions/13464/how-to-setup-the-account-manager-type-to-sign-transactions-in-go
	ks := keystore.NewKeyStore(
		dirName,
		keystore.LightScryptN,
		keystore.LightScryptP)

	defer os.RemoveAll(dirName)

	f, err := os.Open(privateKeyCSV)
	if err != nil {
		log.Fatalf("Open file failed: %v\n", err)
	}
	defer f.Close()
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatalf("Read file failed: %v\n", err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 3 * time.Second)
	eth, err := ethclient.DialContext(ctx, ethereumClientURL)
	if err != nil {
		log.Fatalf("Connection failed: %v\n", err)
	}

	total := big.NewInt(0	)
	sendToAddress := common.HexToAddress(ethereumAddress)
	if err != nil {
		log.Fatalf("HexToAddress failed: %v\n", err)
	}

	gasCostWei := big.NewInt(int64(gasCost))
	mul := big.NewInt(1000000000) //10^9
	gasCostWei.Mul(gasCostWei, mul)

	totalGasCost := big.NewInt(0).Mul(gasCostWei, big.NewInt(int64(GAS_LIMIT)))

	var param *params.ChainConfig
	switch ethereumNet {
	case "rinkeby":
		param = params.RinkebyChainConfig
	case "test":
		param = params.TestChainConfig
	default:
		param = params.MainnetChainConfig
	}

	//https://golangcode.com/how-to-read-a-csv-file-into-a-struct/
	for _, line := range lines {
		if line[0] == "Private Key" {
			continue
		}
		privateKey, err := crypto.HexToECDSA(line[0][2:])
		if err != nil {
			log.Fatalf("HexToECDSA failed: %v\n", err)
		}
		acc, err := ks.ImportECDSA(privateKey, password)
		if err != nil {
			log.Fatalf("ImportECDSA failed: %v\n", err)
		}
		//fetch account balance
		balance, err:=eth.BalanceAt(ctx, acc.Address, nil)
		if err != nil {
			log.Fatalf("BalanceAt failed: %v\n", err)
		}

		if balance.Cmp(big.NewInt(0)) == 0 {
			fmt.Printf(".")
			continue
		}
		fmt.Printf("Balance of %v is %v\n", line[2], balance)
		nonce, err:=eth.NonceAt(ctx, acc.Address, nil)
		if err != nil {
			log.Fatalf("NonceAt failed: %v\n", err)
		}

		//outdated, but gives a hint on the API:
		//https://github.com/ethereum/go-ethereum/wiki/Native:-Account-management
		newBalance := balance.Sub(balance, totalGasCost)
		rawTx := types.NewTransaction(nonce, sendToAddress, newBalance, GAS_LIMIT, gasCostWei, nil)
		err = ks.Unlock(acc, password)
		if err != nil {
			log.Fatalf("Unlock failed: %v\n", err)
		}

		signedTx, err := ks.SignTx(acc, rawTx, param.ChainID)

		if err != nil {
			log.Fatalf("SignTx failed: %v\n", err)
		}

		if dryRun {
			fmt.Printf("Create tx: %v\n",signedTx);
		} else {
			err = eth.SendTransaction(ctx, signedTx)
			if err != nil {
				log.Fatalf("SendTransaction failed: %v\n", err)
			}
		}
		total.Add(total,newBalance)
		fmt.Printf("New total balance: %v\n", total)
		//remove key again
		ks.Delete(acc, password)
	}
	fmt.Printf("Ethereum consolidating done\n")
}

//generate some random strings
//https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang#31832326
func randStr(len int) string {
	buff := make([]byte, len)
	rand.Read(buff)
	str := base64.StdEncoding.EncodeToString(buff)
	// Base 64 can be longer than len
	return str[:len]
}