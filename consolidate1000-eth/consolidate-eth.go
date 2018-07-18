package main

import (
	"crypto/rand"
	"github.com/ogier/pflag"
	"os"
	"encoding/base64"
	"io/ioutil"
	"log"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"encoding/csv"
	"time"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"fmt"
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
			"[-e --ethereum address] [-u --ethereum-url url:port] [-g --gascost gwei] [-n --ethereum-net main|rinkeby|test]")
		println()
		flag.PrintDefaults()
		println()
	}

	flag.StringVarP(&privateKeyCSV, "private-keys", "s", "", "CSV filename to store the private keys")
	flag.BoolVarP(&dryRun, "dryrun", "d", false, "Don't issue the tx, print the tx to stdout")

	flag.StringVarP(&ethereumAddress, "ethereum", "e", "", "The ethereum address where to consolidate the transactions")
	flag.StringVarP(&ethereumClientURL, "ethereum-url", "u", "http://127.0.0.1:8545", "The ethereum URL where a full node is running")
	flag.IntVarP(&gasCost, "gascost", "g", 1, "Gas cost in gwei, check https://ethgasstation.info/ for current price")
	flag.StringVarP(&ethereumNet, "ethereum-net", "n", "main", "Specify ethereum network: main,test,rinkeby")
	flag.Parse(os.Args[1:])

	if len(ethereumAddress) == 0 || len(privateKeyCSV) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	consolidateEth(privateKeyCSV, ethereumAddress, ethereumClientURL, gasCost, ethereumNet, dryRun)
}

func consolidateEth(privateKeyCSV string, ethereumAddress string, ethereumClientURL string, gasCost int, ethereumNet string, dryRun bool) {
	password := randStr(20)
	dirName, err:=ioutil.TempDir(os.TempDir(), "consolidate-eth")

	if err != nil {
		log.Fatal(err)
	}

	//we can use keystore only
	//https://ethereum.stackexchange.com/questions/13464/how-to-setup-the-account-manager-type-to-sign-transactions-in-go
	ks := keystore.NewKeyStore(
		dirName,
		keystore.LightScryptN,
		keystore.LightScryptP)

	f, err := os.Open(privateKeyCSV)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 3 * time.Second)
	eth, err := ethclient.DialContext(ctx, ethereumClientURL)
	if err != nil {
		log.Fatal(err)
	}

	total := big.NewInt(0	)
	sendToAddress := common.HexToAddress(ethereumAddress)
	if err != nil {
		log.Fatal(err)
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
			log.Fatal(err)
		}
		acc, err := ks.ImportECDSA(privateKey, password)
		if err != nil {
			log.Fatal(err)
		}
		//fetch account balance
		balance, err:=eth.BalanceAt(ctx, acc.Address, nil)
		if err != nil {
			log.Fatal(err)
		}

		if balance.Cmp(big.NewInt(0)) == 0 {
			fmt.Printf(".")
			continue
		}
		fmt.Printf("balance of %v is %v\n", line[2], balance)
		nonce, err:=eth.NonceAt(ctx, acc.Address, nil)
		if err != nil {
			log.Fatal(err)
		}

		//outdated, but gives a hint on the API:
		//https://github.com/ethereum/go-ethereum/wiki/Native:-Account-management
		newBalance := balance.Sub(balance, totalGasCost)
		rawTx := types.NewTransaction(nonce, sendToAddress, newBalance, GAS_LIMIT, gasCostWei, nil)
		err = ks.Unlock(acc, password)
		if err != nil {
			log.Fatal(err)
		}

		signedTx, err := ks.SignTx(acc, rawTx, param.ChainID)

		if err != nil {
			log.Fatal(err)
		}

		if dryRun {
			fmt.Printf("create tx: %v\n",signedTx);
		} else {
			err = eth.SendTransaction(ctx, signedTx)
			if err != nil {
				log.Fatal(err)
			}
		}
		total.Add(total,newBalance)
		fmt.Printf("New total: %v\n", total)
	}
	fmt.Printf("\n")
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