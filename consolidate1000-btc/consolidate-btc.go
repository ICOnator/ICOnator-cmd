package main

import (
	"github.com/ogier/pflag"
	"os"
	"log"
	"encoding/csv"
	"github.com/ethereum/go-ethereum/crypto"
	"fmt"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/btcec"
	"errors"
)

const (
	GAS_LIMIT uint64 = 21000
)

func main() {

	var (
		privateKeyCSV 		string
		bitcoinAddress  	string
		bitcoinClientURL 	string
		bitcoinUsername 	string
		bitcoinNet		 	string
		bitcoinPassword 	string
		dryRun				bool
	)

	flag := pflag.NewFlagSet("consolidate1000-btc", pflag.ExitOnError)

	flag.Usage = func() {
		println("Usage:")
		println("  consolidate1000-btc [-s --private-keys filename] [-d --dryrun boolean] [-b --bitcoin address] " +
			"[-u --bitcoin-url url:port] [-U --bitcoin-user username] [-p --bitcoin-pass passsword] [-n --bictoin-net main|reg|test]")
		println()
		flag.PrintDefaults()
		println()
	}

	flag.StringVarP(&privateKeyCSV, "private-keys", "s", "", "CSV filename to store the private keys")
	flag.BoolVarP(&dryRun, "dryrun", "d", false, "Don't issue the tx, print the tx to stdout")

	flag.StringVarP(&bitcoinAddress, "bitcoin", "b", "", "The bitcoin address where to consolidate the transactions")
	flag.StringVarP(&bitcoinClientURL, "bitcoin-url", "u", "http://127.0.0.1:18332", "The bitcoin URL where a full node is running")
	flag.StringVarP(&bitcoinUsername, "bitcoin-user", "U", "", "The bitcoin username for the full node")
	flag.StringVarP(&bitcoinPassword, "bitcoin-pass", "p", "", "The bitcoin password for the full node")
	flag.StringVarP(&bitcoinNet, "bitcoin-net", "n", "main", "Specify bitcoin network: main,test,reg")

	flag.Parse(os.Args[1:])

	if len(bitcoinAddress) == 0 || len(privateKeyCSV) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	consolidateBtc(privateKeyCSV, bitcoinAddress, bitcoinClientURL, bitcoinUsername, bitcoinPassword, bitcoinNet, dryRun)
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

	client.CreateNewAccount("tmp-account");

	for _, line := range lines {
		if line[0] == "Private Key" {
			continue
		}
		privateKey, err := crypto.HexToECDSA(line[0][2:])
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
			log.Printf("less than 5000: %v", amount)
			continue
		}
		total += int64(amount)
		client.ImportPrivKey(wif)

	}
	sendTo, err := btcutil.DecodeAddress(bitcoinAddress, param)
	if err != nil {
		log.Fatal(err)
	}

	if dryRun {
		fmt.Printf("send to: %v, amount: %v BTC\n", sendTo, btcutil.Amount(total).ToBTC());
	} else {
		client.SendToAddress(sendTo, btcutil.Amount(total))
	}

}