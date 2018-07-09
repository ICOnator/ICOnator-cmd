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
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/btcec"
	"errors"
)

const (
	GAS_LIMIT uint64 = 1000000
)

func main() {

	var (
		privateKeyCSV 		string
		bitcoinAddress  	string
		bitcoinClientURL 	string
		bitcoinUsername 	string
		bitcoinNet		 	string
		bitcoinPassword 	string
		ethereumAddress		string
		ethereumClientURL 	string
		ethereumNet		 	string

		gasCost				int
		dryRun				bool
	)

	flag := pflag.NewFlagSet("keygen", pflag.ExitOnError)

	flag.Usage = func() {
		println("Usage:")
		println("  consolidate1000 [-s filename] [-d boolean] [-b bitcoin address] [--bitcoin-url url:port] [--bitcoin-user] [--bitcoin-pass] [--bictoin-net] [-e ethereum address] [--ethereum-url url:port] [-g gwei] [--ethereum-net]")
		println()
		flag.PrintDefaults()
		println()
	}

	flag.StringVarP(&privateKeyCSV, "private-keys", "s", "", "CSV filename to store the private keys")
	flag.BoolVarP(&dryRun, "dryrun", "d", true, "Don't issue the tx, print the tx to stdout")

	flag.StringVarP(&bitcoinAddress, "bitcoin", "b", "", "The bitcoin address where to consolidate the transactions")
	flag.StringVarP(&bitcoinClientURL, "bitcoin-url", "", "", "The bitcoin URL where a full node is running")
	flag.StringVarP(&bitcoinUsername, "bitcoin-user", "", "", "The bitcoin username for the full node")
	flag.StringVarP(&bitcoinPassword, "bitcoin-pass", "", "", "The bitcoin password for the full node")
	flag.StringVarP(&bitcoinNet, "bitcoin-net", "", "main", "Specify bitcoin network: main,test,reg")

	flag.StringVarP(&ethereumAddress, "ethereum", "e", "", "The ethereum address where to consolidate the transactions")
	flag.StringVarP(&ethereumClientURL, "ethereum-url", "", "", "The ethereum URL where a full node is running")
	flag.IntVarP(&gasCost, "gascost", "g", 10, "Gas cost in gwei, check https://ethgasstation.info/ for current price")
	flag.StringVarP(&ethereumNet, "ethereum-net", "", "main", "Specify ethereum network: main,test,rinkeby")
	flag.Parse(os.Args[1:])

	if len(ethereumAddress) != 0 && len(ethereumClientURL) != 0 {
		consolidateEth(privateKeyCSV, ethereumAddress, ethereumClientURL, gasCost, ethereumNet, dryRun)
	}

	if len(bitcoinAddress) != 0 && len(bitcoinClientURL) != 0 {
		consolidateBtc(privateKeyCSV, bitcoinAddress, bitcoinClientURL, bitcoinUsername, bitcoinPassword, bitcoinNet, dryRun)
	}
}

func consolidateBtc(privateKeyCSV string, bitcoinAddress string, bitcoinClientURL string, bitcoinUsername string, bitcoinPassword string, bitcoinNet string, dryRun bool) {
	f, err := os.Open(privateKeyCSV)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	connCfg := &rpcclient.ConnConfig{
		Host:         bitcoinClientURL,
		User:         bitcoinUsername,
		Pass:         bitcoinPassword,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Shutdown()

	var param *chaincfg.Params
	switch bitcoinNet {
	case "reg":
		param = &chaincfg.RegressionNetParams
	case "test":
		param = &chaincfg.TestNet3Params
	default:
		param = &chaincfg.MainNetParams
	}

	var total int64 = 0;



	for _, line := range lines {
		if line[0] == "Private Key" {
			continue
		}
		privateKey, err := crypto.HexToECDSA(line[0])
		btcPriv, btcPub := btcec.PrivKeyFromBytes(privateKey.PublicKey.Curve, crypto.FromECDSA(privateKey))
		if err != nil {
			log.Fatal(err)
		}
		//sanity check
		if btcPub.X.Cmp(privateKey.X) != 0 || btcPub.Y.Cmp(privateKey.Y) != 0 || *btcPub.Curve.Params() != *privateKey.Curve.Params() {
			log.Fatal(errors.New("not the same public key"))
		}

		wif, err := btcutil.NewWIF(btcPriv, param, true)
		if err != nil {
			log.Fatal(err)
		}

		addr, err := btcutil.NewAddressPubKey(wif.SerializePubKey(), param)

		amount, err := client.GetReceivedByAddress(addr)
		if err != nil {
			log.Fatal(err)
		}
		if int64(amount) < 5000 {
			//ignore addresses with less than 5000 satoshis
			continue
		}
		total += int64(amount)

		//client.New
		//btcjson.TransactionInput{}

	}
	m := make(map[btcutil.Address]btcutil.Amount)

	addr, err := btcutil.DecodeAddress(bitcoinAddress, param)
	if err != nil {
		log.Fatal(err)
	}
	m[addr] = btcutil.Amount(total)

	//client.CreateRawTransaction(m,0) //TODO: locktime?
}

func consolidateEth(privateKeyCSV string, ethereumAddress string, ethereumClientURL string, gasCost int, ethereumNet string, dryRun bool) {
	password := randStr(20)
	file, err:=ioutil.TempFile(os.TempDir(), "consolidate-eth")
	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}

	//we can use keystore only
	//https://ethereum.stackexchange.com/questions/13464/how-to-setup-the-account-manager-type-to-sign-transactions-in-go
	ks := keystore.NewKeyStore(
		file.Name(),
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

	gas := big.NewInt(int64(gasCost))
	mul := big.NewInt(1000000000) //10^9
	gas.Mul(gas, mul)

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
		privateKey, err := crypto.HexToECDSA(line[0])
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
			print(".")
			continue
		}
		nonce, err:=eth.NonceAt(ctx, acc.Address, nil)
		if err != nil {
			log.Fatal(err)
		}

		//outdated, but gives a hint on the API:
		//https://github.com/ethereum/go-ethereum/wiki/Native:-Account-management
		rawTx := types.NewTransaction(nonce, sendToAddress, balance, GAS_LIMIT, gas, nil)
		signedTx, err := ks.SignTx(acc, rawTx, param.ChainID)

		if dryRun {
			fmt.Printf("create tx: %v\n",signedTx);
		} else {
			eth.SendTransaction(ctx, signedTx)
		}
		total.Add(total,balance)
		fmt.Printf("New total: %v\n", total)
	}
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