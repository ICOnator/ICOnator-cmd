package main

import (
	"github.com/ethereum/go-ethereum/crypto"
	"fmt"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"crypto/rand"
	"github.com/ethereum/go-ethereum/common"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcd/btcec"
	rand2 "math/rand"
	"github.com/ogier/pflag"
	"os"
	"time"
	"encoding/binary"
	"bytes"
	"io/ioutil"
	"log"
)

func main() {

	var (
		numberOfKeys    int
		random1     	int
		random2     	int
		random3     	int
		random4     	int
		privateKeyCSV 	string
		publicKeyCSV 	string
	)

	flag := pflag.NewFlagSet("keygen", pflag.ExitOnError)

	flag.Usage = func() {
		println("Usage:")
		println("  keygen [-n num] [-s private-key CSV] [-p public-key CSV] [-r1 random1] [-r2 random2] [-r3 random3] [-r4 random4]")
		println()
		flag.PrintDefaults()
		println()
	}

	flag.IntVarP(&numberOfKeys, "number", "n", 1, "Number of addresses to generate")
	flag.IntVarP(&random1, "random1", "1", 0, "Optional random number to increase entropy")
	flag.IntVarP(&random2, "random2", "2", 0, "Optional random number to increase entropy")
	flag.IntVarP(&random3, "random3", "3", 0, "Optional random number to increase entropy")
	flag.IntVarP(&random4, "random4", "4", 0, "Optional random number to increase entropy")
	flag.StringVarP(&privateKeyCSV, "private-keys", "s", "", "Specify a file name to store the private keys in a CSV")
	flag.StringVarP(&publicKeyCSV, "public-keys", "p", "", "Specify a file name to store the public keys only in a CSV")
	flag.Parse(os.Args[1:])

	line:= fmt.Sprintf("Private Key,Private Key (WIF),Ethereum Address,Bitcoin Address (SegWit-Bech32),Bitcoin Address (P2PKH-Base58)\n")
	if len(privateKeyCSV) > 0 {
		ioutil.WriteFile(privateKeyCSV, []byte(line), 0644)
	} else {
		print(line)
	}

	if len(publicKeyCSV) > 0 {
		linePublic:= fmt.Sprintf("Ethereum Address,Bitcoin Address (SegWit-Bech32),Bitcoin Address (P2PKH-Base58)\n")
		ioutil.WriteFile(publicKeyCSV, []byte(linePublic), 0644)
	}

	seed := createSeed(random1, random2, random3, random4)
	rand1 := rand.Reader //the crypto random
	rand2 := rand2.New(rand2.NewSource(time.Now().Unix() ^ seed)) //user provided entropy
	rand3 := New(rand1, rand2)


	var fpriv *os.File;
	if len(privateKeyCSV) > 0 {
		var err error
		fpriv, err = os.OpenFile(privateKeyCSV, os.O_APPEND|os.O_WRONLY, 0644)
		defer fpriv.Close()
		if err != nil {
			log.Fatal(err)
		}
	}

	var fpub *os.File;
	if len(publicKeyCSV) > 0 {
		var err error
		fpub, err = os.OpenFile(publicKeyCSV, os.O_APPEND|os.O_WRONLY, 0644)
		defer fpub.Close()
		if err != nil {
			log.Fatal(err)
		}
	}

	for i := 0; i < numberOfKeys; i++ {
		key, _ := ecdsa.GenerateKey(secp256k1.S256(), rand3)

		ethAddr := crypto.PubkeyToAddress(key.PublicKey)
		btcPriv, _ := btcec.PrivKeyFromBytes(key.PublicKey.Curve, crypto.FromECDSA(key))
		wif, _ := btcutil.NewWIF(btcPriv, &chaincfg.MainNetParams, true)

		btcWitAddr, _ := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), &chaincfg.MainNetParams);
		btcAddr, _ := btcutil.NewAddressPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), &chaincfg.MainNetParams);
		line := fmt.Sprintf("%v,%v,%v,%v,%v\n", common.ToHex(crypto.FromECDSA(key)), wif, ethAddr.Hex(), btcWitAddr, btcAddr)
		if len(privateKeyCSV) > 0 {
			fpriv.WriteString(line)
		} else {
			print(line)
		}

		if len(publicKeyCSV) > 0 {
			linePublic := fmt.Sprintf("%v,%v,%v\n", ethAddr.Hex(), btcWitAddr, btcAddr)
			fpub.WriteString(linePublic)
		}
	}
}

func createSeed(random1 int, random2 int, random3 int, random4 int) int64 {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, random1)
	binary.Write(buf, binary.LittleEndian, random2)
	binary.Write(buf, binary.LittleEndian, random3)
	binary.Write(buf, binary.LittleEndian, random4)
	seed := crypto.Keccak256Hash(buf.Bytes()).Big().Int64()
	return seed
}