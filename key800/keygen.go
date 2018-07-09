package main

import (
	"github.com/ethereum/go-ethereum/crypto"
	"fmt"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"crypto/rand"
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
	"io"
	"errors"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func main() {

	var (
		numberOfKeys  int
		random1       int
		random2       int
		random3       int
		random4       int
		privateKeyCSV string
		publicKeyCSV  string
		btcNet        string
	)

	flag := pflag.NewFlagSet("keygen", pflag.ExitOnError)

	flag.Usage = func() {
		println("Usage:")
		println("  key800 [-n num] [-s filename] [-p filename] [-b main|reg|test][-r1 random1] [-r2 random2] [-r3 random3] [-r4 random4]")
		println()
		flag.PrintDefaults()
		println()
	}

	flag.IntVarP(&numberOfKeys, "number", "n", 1, "Number of addresses to generate")
	flag.IntVarP(&random1, "random1", "1", 0, "Optional random number to increase entropy")
	flag.IntVarP(&random2, "random2", "2", 0, "Optional random number to increase entropy")
	flag.IntVarP(&random3, "random3", "3", 0, "Optional random number to increase entropy")
	flag.IntVarP(&random4, "random4", "4", 0, "Optional random number to increase entropy")
	flag.StringVarP(&privateKeyCSV, "private-keys", "s", "", "CSV filename to store the private keys")
	flag.StringVarP(&publicKeyCSV, "public-keys", "p", "", "CSV filename to store the public keys only")
	flag.StringVarP(&btcNet, "btcNet", "b", "main", "Specify bitcoin network: main,test,reg")
	flag.Parse(os.Args[1:])

	line := fmt.Sprintf("Private Key,Private Key (WIF),Ethereum Address,Bitcoin Address (SegWit-Bech32),Bitcoin Address (P2PKH-Base58)\n")
	if len(privateKeyCSV) > 0 {
		ioutil.WriteFile(privateKeyCSV, []byte(line), 0644)
	} else {
		print(line)
	}

	if len(publicKeyCSV) > 0 {
		linePublic := fmt.Sprintf("Ethereum Address,Bitcoin Address (SegWit-Bech32),Bitcoin Address (P2PKH-Base58)\n")
		ioutil.WriteFile(publicKeyCSV, []byte(linePublic), 0644)
	}

	seed := createSeed(random1, random2, random3, random4)
	rand1 := rand.Reader                                          //the crypto random
	rand2 := rand2.New(rand2.NewSource(time.Now().Unix() ^ seed)) //user provided entropy
	rand3 := newDualReader(rand1, rand2)

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

	var param *chaincfg.Params
	switch btcNet {
	case "reg":
		param = &chaincfg.RegressionNetParams
	case "test":
		param = &chaincfg.TestNet3Params
	default:
		param = &chaincfg.MainNetParams
	}

	for i := 0; i < numberOfKeys; i++ {
		key, err := ecdsa.GenerateKey(secp256k1.S256(), rand3)
		if err != nil {
			log.Fatal(err)
		}

		ethAddr := crypto.PubkeyToAddress(key.PublicKey)
		btcPriv, btcPub := btcec.PrivKeyFromBytes(key.PublicKey.Curve, crypto.FromECDSA(key))

		//sanity check
		if btcPub.X.Cmp(key.X) != 0 || btcPub.Y.Cmp(key.Y) != 0 || *btcPub.Curve.Params() != *key.Curve.Params() {
			log.Fatal(errors.New("not the same public key"))
		}

		wif, err := btcutil.NewWIF(btcPriv, param, true)
		if err != nil {
			log.Fatal(err)
		}

		btcWitAddr, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), param);
		if err != nil {
			log.Fatal(err)
		}

		btcAddr, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), param);
		if err != nil {
			log.Fatal(err)
		}

		line := fmt.Sprintf("%v,%v,%v,%v,%v\n", hexutil.Encode(crypto.FromECDSA(key)), wif, ethAddr.Hex(), btcWitAddr, btcAddr)
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
	return crypto.Keccak256Hash(buf.Bytes()).Big().Int64()
}

type dualReader struct {
	reader1 io.Reader
	reader2 io.Reader
}

func newDualReader(reader1 io.Reader, reader2 io.Reader) *dualReader {
	return &dualReader{reader1, reader2}
}

func (mr *dualReader) Read(p []byte) (n int, err error) {
	len := len(p)
	tmp1 := make([]byte, len)
	tmp2 := make([]byte, len)

	n1, err1 := mr.reader1.Read(tmp1)
	n2, err2 := mr.reader2.Read(tmp2)

	if err1 != nil {
		return n1, err1;
	}

	if err2 != nil {
		return n2, err2;
	}

	if len != n1 && len != n2 {
		return 0, errors.New("did not read same length")
	}

	for i := 0; i < n1; i++ {
		p[i] = tmp1[i] ^ tmp2[i];
	}

	return n1, nil
}