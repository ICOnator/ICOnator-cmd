package main

import (
	"github.com/ogier/pflag"
	"os"
	"log"
	"fmt"
	"encoding/json"
	"net/http"
	"strings"
	"io/ioutil"
	"errors"
	"encoding/csv"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/btcsuite/btcd/btcec/v2"
)

func main() {

	var (
		privateKeyCSV 		string
		bitcoinAddress  	string
		bitcoinClientURL 	string
		bitcoinUsername 	string
		bitcoinPassword 	string
		bitcoinNet		 	string
		dryRun				bool
		ignoreNotEmpty		bool
	)

	flag := pflag.NewFlagSet("consolidate1000-btc", pflag.ExitOnError)

	flag.Usage = func() {
		println("Usage:")
		println("  consolidate1000-btc [-s --private-keys filename] [-d --dryrun boolean] [-i --ignore boolean] [-b --bitcoin-to address] " +
			"[-u --url url:port] [-U --username username] [-p --password password] [-N --network main|reg|test]")
		println()
		flag.PrintDefaults()
		println()
	}

	flag.StringVarP(&privateKeyCSV, "private-keys", "s", "", "CSV filename to store the private keys")
	flag.BoolVarP(&dryRun, "dryrun", "d", false, "Don't issue the tx, print the tx to stdout")
	flag.BoolVarP(&ignoreNotEmpty, "ignore", "i", false, "Ignore if wallet is not empty")

	flag.StringVarP(&bitcoinAddress, "bitcoin-to", "b", "", "The bitcoin address where to consolidate the transactions")
	flag.StringVarP(&bitcoinClientURL, "url", "u", "127.0.0.1:18332", "The bitcoin URL where a full node is running")
	flag.StringVarP(&bitcoinUsername, "username", "U", "", "The username for the bitcoin RPC client")
	flag.StringVarP(&bitcoinPassword, "password", "p", "", "The password for theh bitcoin RPC client")
	flag.StringVarP(&bitcoinNet, "network", "N", "main", "Specify bitcoin network: main,test,reg")

	flag.Parse(os.Args[1:])

	if len(bitcoinAddress) == 0 || len(privateKeyCSV) == 0 || len(bitcoinUsername) == 0 || len(bitcoinPassword) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	f, err := os.Open(privateKeyCSV)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	//start bitcoind with e.g. bitcoind -testnet -rpcuser=test -rpcpassword=me
	c := Client {"http://"+bitcoinUsername+":"+bitcoinPassword+"@"+bitcoinClientURL, 1};

	var param *chaincfg.Params
	switch bitcoinNet {
	case "reg":
		param = &chaincfg.RegressionNetParams
	case "test":
		param = &chaincfg.TestNet3Params
	default:
		param = &chaincfg.MainNetParams
	}

	balance, err := c.getBalance();
	if err != nil {
		log.Fatalf("GetBalance failed: %v\n", err)
	}
	if !ignoreNotEmpty && balance != 0 {
		log.Fatalf("The wallet needs to be empty. It has a current balance of: %v BTC\n", balance)
	}

	var last *btcutil.WIF;

	for _, line := range lines {
		if line[0] == "Private Key" {
			continue
		}

		if last != nil {
			err = c.importPrivKeyRescan(last.String(),  false)
			if err != nil {
				log.Fatalf("not importing... %v: %v\n", last, err)
			}
		}

		privateKey, err := crypto.HexToECDSA(line[0][2:])
		btcPriv, btcPub := btcec.PrivKeyFromBytes(privateKey.PublicKey.Curve, crypto.FromECDSA(privateKey))
		if err != nil {
			log.Fatal(err)
		}
		//sanity check
		if btcPub.X.Cmp(privateKey.X) != 0 || btcPub.Y.Cmp(privateKey.Y) != 0 || *btcPub.Curve.Params() != *privateKey.Curve.Params() {
			log.Fatal("not the same public keys\n")
		}

		wif, err := btcutil.NewWIF(btcPriv, param, true)
		if err != nil {
			log.Fatalf("Unable to create WIF: %v\n", err)
		}
		last = wif;
		print(".");
	}
	print("\n");
	log.Printf("rescanning... (this may take a while)\n")
	err = c.importPrivKeyRescan(last.String(),  true)
	if err != nil {
		log.Fatalf("not importing... %v: %v\n", last, err)
	}

	balance, err = c.getBalance();
	if err != nil {
		log.Fatalf("GetBalance failed: %v\n", err)
	}
	fmt.Printf("Consolidating %9.8f BTC\n", balance)

	sendTo, err := btcutil.DecodeAddress(bitcoinAddress, param)
	if err != nil {
		log.Fatalf("DecodeAddress failed: %v\n", err)
	}

	//set fees
	fee, err := c.estimateSmartFee(6)
	if err != nil {
		log.Fatalf("EstimateFee failed: %v\n", err)
	}

	fmt.Printf("Set fee per KB of %9.8f BTC\n", fee)
	err = c.setTxFee(fee)
	if err != nil {
		log.Fatalf("SetTxFee failed: %v\n", err)
	}

	if dryRun {
		fmt.Printf("send to: %v, amount: %9.8f BTC\n", sendTo, balance);
	} else {
		txid, err := c.sendToAddress(sendTo.String(), balance, true)
		if err != nil {
			log.Fatalf("SendToAddress failed: %v\n", err)
		}
		fmt.Printf("send %9.8f BTC to: %v, txid: %v\n", balance, sendTo, txid);
	}

	balance, err = c.getBalanceMinConf( 0)
	if err != nil {
		log.Fatalf("GetBalance failed: %v\n", err)
	}
	if balance != 0 {
		log.Fatalf("The wallet needs to be empty. It has a current balance of: %9.8f BTC\n", balance)
	}

	fmt.Printf("Bitcoin consolidating done\n")
}

type Client struct {
	address string
	id interface{}
}

func (c *Client) getBalance()(float64, error) {
	resp, err:=c.call("getbalance", c.id, nil)
	if err!=nil {
		return 0, err
	}
	result:=resp["result"]
	return result.(float64), err
}

func (c *Client) estimateSmartFee(conf_target uint16)(float64, error) {

	data := make([]interface{}, 1)
	data[0] = conf_target

	resp, err:=c.call("estimatesmartfee", c.id, data)
	if err!=nil {
		return 0, err
	}
	result:=resp["result"]
	return result.(map[string]interface{})["feerate"].(float64), err
}

func (c *Client) setTxFee(newFeeAmount float64) (error) {
	data := make([]interface{}, 1)
	data[0] = newFeeAmount
	_, err:=c.call("settxfee", c.id, data)
	if err!=nil {
		return err
	}

	return nil
}

func (c *Client) getBalanceMinConf(minConf int)(float64, error) {
	data := make([]interface{}, 2)
	data[0] = ""
	data[1] = minConf
	resp, err:=c.call("getbalance", c.id, data)
	if err!=nil {
		return 0, err
	}
	result:=resp["result"]
	return result.(float64), err
}


func (c *Client) sendToAddress(address string, amount float64, feeSubtraction bool) (string, error){
	data := make([]interface{}, 5)
	data[0] = address
	data[1] = amount
	data[2] = ""
	data[3] = ""
	data[4] = feeSubtraction
	resp, err:=c.call("sendtoaddress", c.id, data)
	if err!=nil {
		return "", err
	}

	return resp["result"].(string), nil
}

func (c *Client) importPrivKeyRescan(address string, rescan bool) (error){
	data := make([]interface{}, 3)
	data[0] = address
	data[1] = ""
	data[2] = rescan
	resp, err:=c.call("importprivkey", c.id, data)
	if err!=nil {
		return err
	}
	if (resp["error"] != nil) {
		return errors.New(resp["error"].(string))
	}

	return nil
}

func (c *Client) call(method string, id interface{}, params []interface{})(map[string]interface{}, error){
	data, err := json.Marshal(map[string]interface{}{
		"method": method,
		"id":     id,
		"params": params,
	})
	if err != nil {
		log.Fatalf("Marshal: %v\n", err)
		return nil, err
	}
	resp, err := http.Post(c.address,
		"application/json", strings.NewReader(string(data)))
	if err != nil {
		log.Fatalf("Post: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll: %v\n", err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		log.Fatalf("status code is: %v\n", resp.StatusCode)
	}

	result := make(map[string]interface{})
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalf("Unmarshal: %v\n", body)
		return nil, err
	}
	return result, nil
}
