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
)

func main() {

	var (
		numberOfKeys     int
		random1     int
		random2     int
		random3     int
		random4     int
	)

	flag := pflag.NewFlagSet("keygen", pflag.ExitOnError)

	flag.Usage = func() {
		println("Usage:")
		println("  keygen [-n num] [-r1 random1] [-r2 random2] [-r3 random3] [-r4 random4]")
		println()
		flag.PrintDefaults()
		println()
	}

	flag.IntVarP(&numberOfKeys, "number", "n", 1, "Number of addresses to generate")
	flag.IntVarP(&random1, "random1", "1", 0, "Optional random number to increase entropy")
	flag.IntVarP(&random2, "random2", "2", 0, "Optional random number to increase entropy")
	flag.IntVarP(&random3, "random3", "3", 0, "Optional random number to increase entropy")
	flag.IntVarP(&random4, "random4", "4", 0, "Optional random number to increase entropy")
	flag.Parse(os.Args[1:])

	fmt.Printf("Private Key,Private Key (WIF),Ethereum Address,Bitcoin Address (SegWit-Bech32),Bitcoin Address (P2PKH-Base58)\n")

	seed := createSeed(random1, random2, random3, random4)
	rand1 := rand.Reader //the crypto random
	rand2 := rand2.New(rand2.NewSource(time.Now().Unix() ^ seed)) //user provided entropy

	rand3 := New(rand1, rand2)

	for i := 0; i < numberOfKeys; i++ {
		key, _ := ecdsa.GenerateKey(secp256k1.S256(), rand3)

		addr := crypto.PubkeyToAddress(key.PublicKey)
		priv, _ := btcec.PrivKeyFromBytes(key.PublicKey.Curve, crypto.FromECDSA(key))
		wif, _ := btcutil.NewWIF(priv, &chaincfg.MainNetParams, true)

		bitcoinKey, _ := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), &chaincfg.MainNetParams);
		bitcoinKeyOld, _ := btcutil.NewAddressPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), &chaincfg.MainNetParams);
		fmt.Printf("%v,%v,%v,%v,%v\n", common.ToHex(crypto.FromECDSA(key)), wif, addr.Hex(), bitcoinKey, bitcoinKeyOld)
	}
}

func createSeed(random1 int, random2 int, random3 int, random4 int) int64 {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, random1)
	binary.Write(buf, binary.LittleEndian, random2)
	binary.Write(buf, binary.LittleEndian, random3)
	binary.Write(buf, binary.LittleEndian, random4)
	seed := crypto.Sha3Hash(buf.Bytes()).Big().Int64()
	return seed
}