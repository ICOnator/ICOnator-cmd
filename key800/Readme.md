# Key Generation Tool

This tool will generate Bitcoin and Ethereum address. You can specify how many addresses you want to create and also provide additional entropy

## Build

```
git clone https://github.com/ICOnator/ICOnator-cmd.git
cd ICOnator-cmd/key800
go get && go build
```
 
This will generate a the binary key800. 

## Example usage

The parameters for this tools are:

```
Usage:
  keygen [-n num] [-s private-key CSV] [-p public-key CSV] [-1 random1] [-2 random2] [-3 random3] [-4 random4]
```

Example:

```
./key800 -n 10 -1 9053294 -2 26906534
```

This will generate:

```
Private Key,Private Key (WIF),Ethereum Address,Bitcoin Address (SegWit-Bech32),Bitcoin Address (P2PKH-Base58)
0x3fa406f363d4cb1a60c3fd8647bb216c982fb08db449fba24e7823898ed5f666,KyMRKdENUKMMqExjPBXFrAhQKRc72pvhDT8AH2TNmAY8TLhqTd45,0x70d2bfc771faedb75618ac3992b1922dfd43dbea,bc1qmh5n6ektslwky8xr7ajr049s34jv2jf30taaku,1MEMkBdvxMCHb12T75PN6TT3WkBaMKX3pT
0x5ac58c652a5bf78566a08c2cbd3f8f3ba282371ffdf33b38d21630d0740f5ba5,KzGABeKwmsmVXSpmbzvYBDm9Dze6FhPB5cFhd33hUqTMqq8U61va,0x6392cbc927c8a651bc915d3b592d21a031807a8c,bc1q0fqywy8kaa3z0dp7k5jeqtt5htzemhnjlyfc9r,1C9QSVPTt5xN2NqVVzop9FbUprGy7DbpMR
0xe1c9640184e1e515c5d7d004fe20e10e900dd3e6bee7429816da6642cee516cc,L4ncNc1xY5sA6dnJBDr6xW9kBy3G1JDJuHA7Ps6dx73puEU4NnDW,0x62650fca251b310d42442a4e5142f3f53dda5f35,bc1qk2mhvtpjwu63cpp7yrzcj2lvmcjz4842f3xp8l,1HHy4HMvUbZrpycUtmnDbD43haA28TFdbr
0xb210664f9ca78ae09610c702eb50d4cb200719112e27f0b2b73e1fae9a57f247,L3Bquy1YXg1oyimGFZfMnPNVhUPoGnKUQ5StkRZzJtEimrPGdGZ9,0xb886aa96016e304133fd55418deb9b1aa35be406,bc1quyuulmlmj27sjq9tulns53vefhfjktzgcv3kxa,1MXtJ2cTupyWpERADSR6GtEhJz95145V7y
0xac1c189d3b484f436a6c4f427537be7352bef43edd22c50b611e7ad205074ef4,L2zGbgMZ3qc96swgW6sUJ1KwEHBr3xYsCKvsWsViKY1DGb42y88p,0x0b403b56af806bf8cbbabd8c893a7097373e735a,bc1q67d6fr0lqxwtk9dxe6ckq0g4faw97ezeqn2ths,1Lf2jVfoZHW9jAcRdyS8wpmV51Xhc7cpiG
0xb64d2cd45408c41aa554b8d2e9caf84b8e39d0334e06b83021f98ebda902051e,L3L5f6jrxEwNDfzYAKCQmE6yZXCZhT5u4hFsAb2M6vkucPitDgjB,0xed4d1ed5b9ad3edf9caea5b417a3f29565ac085e,bc1q0y5eyn9ndfhln6v0xx06kz4lmwl663hht38qhn,1C3eZedPU8StN2GHJ4JpkGettcbtZ3j51x
0x3f5ec36861d91f2049695131c8942fffc6a68184e451a94ccc8cd66ed9796533,KyLtpNdE3HkRVotRfq4upFdiz83zoFpKvYACbX7wy8HV7SZHv5bt,0xad20d4cc81a982e4ff5b293162328f1a034f04b3,bc1qafntcpn7nfzagmq8yu4ctnx3u3dqpzllvlc39a,1NNQC1Gsa8xPjKEiiwKFYufkz5byw66EoN
0x7ec47f238d82744cb7091972e4bed070853b067959a5d30b1c9be32b20585473,L1U8X8RdghUp6Yoy9z8xnTP1j4tXTAyvixJJgiG2QZbzZ1WaxZHy,0x60596256af101b94c0511c144a6e279aa83536da,bc1q76gf80058qsy672fkl24vgsa7vvzkfzjgy4vdt,1PUiRauF1ZbRFQrBHzwqJTU7Z97k6mJub4
0x9f107fa99d3442d1a15270045fdb52988085115258b61929ecfc5248ecfec5f8,L2Yup9KWriM3GwqERR5gj1DmMac4LrXGhMwFBjFZqj6idhCZrN9e,0x071b36315fc82bef802b4df8aeb6202086db8ec9,bc1qgpa0mpgr7tmwqdhupsk4tl36j4vgfrtz408guu,16swao3DPNwFRLuHjpC4uhwpMPqXukYd5a
0x6fff8c07bfa5d2806275d9600bd6862189d8009bc8b6b870e6d139a311f9961e,KzyRMqx4GVVyEhaZH4obTCP4xs5MCKmPLUSEjgwSNucUoxGAUHGH,0x066e8f381c9d3107b11ee3494ccabd8893f15200,bc1qmrkfdk78wptspftcg92atrqfaasd42za80ke80,1LmzNwxPeNoJ1jaY8ekF72y7URPqjDrRpK
```

## Testing with 2 
Tests conducted with: ```./key800 -n 2 -s secret-test.csv -p public-test.csv -N test -1 43543 -2 24566 -3 645737 -4 67843```

Payin to

```
Ethereum Address,Bitcoin Address (P2PKH-Base58)
0x93d499b82c30C27E7265FfDdBbCb0421767bdd64,mopQvr2abuWAynp1hnMxcNddk8n2v5wkVY
0x586efC5E9463C077D492a449f00a8eDa4675Da36,mrp1xXbaWbRYDN81ds47tWJrzWwX7k5XpN
```

### Ethereum Rinkeby
For testing, 0.25 ETH was sent to the first address, and 0.33 ETH was sent to the second address:
https://rinkeby.etherscan.io/address/0x93d499b82c30C27E7265FfDdBbCb0421767bdd64
https://rinkeby.etherscan.io/address/0x586efC5E9463C077D492a449f00a8eDa4675Da36

Consolidate Ethereum to: 0x0CbdF5B0c4E117619631bA4b97dC0d439ADAbD88, current balance 0.999 ETH

```
./consolidate1000-eth -N "rinkeby" -u "http://127.0.0.1:8545" -s ../key800/secret-test.csv -e "0x0CbdF5B0c4E117619631bA4b97dC0d439ADAbD88"
```
Output:
```
INFO [07-17|16:52:32.584] Submitted transaction                    fullhash=0x41511f1401c32d1e66b20500e6900f428354c89cc3f20ea6f83c129954384551 recipient=0x0CbdF5B0c4E117619631bA4b97dC0d439ADAbD88
INFO [07-17|16:52:32.736] Submitted transaction                    fullhash=0x78a88badb921c2778949e3423ad2577ac4a4ec752a12d96e3140543c02b3b74b recipient=0x0CbdF5B0c4E117619631bA4b97dC0d439ADAbD88
```
https://rinkeby.etherscan.io/address/0x0CbdF5B0c4E117619631bA4b97dC0d439ADAbD88, new balance 1.579114416 Ether

The the other balances are:

balance of 0x93d499b82c30C27E7265FfDdBbCb0421767bdd64 is 0
balance of 0x586efC5E9463C077D492a449f00a8eDa4675Da36 is 0

Geth was run with the following settings

```
geth --rinkeby --fast --rpc --rpcapi "eth"
```

### Bitcoin Testnet3
For testing, 1.3 BTC was sent to the first address, 0.865 BTC was send to the second address:
https://live.blockcypher.com/btc-testnet/address/mrp1xXbaWbRYDN81ds47tWJrzWwX7k5XpN/
https://live.blockcypher.com/btc-testnet/address/mopQvr2abuWAynp1hnMxcNddk8n2v5wkVY/

Consolidate Bitcoin to mjdqkWc34TYbFxqHaJ1mGsjuRbPC57Bqsq, current balance 0 BTC

```
./consolidate1000-btc -s ../key800/secret-test.csv -b mjdqkWc34TYbFxqHaJ1mGsjuRbPC57Bqsq -N test -U test -p me -i
```

After consolidation, both addresses are empty and mjdqkWc34TYbFxqHaJ1mGsjuRbPC57Bqsq has 2.16541752 BTC:
https://live.blockcypher.com/btc-testnet/address/mjdqkWc34TYbFxqHaJ1mGsjuRbPC57Bqsq/

Bitcoin Core was run with the following settings

```
bitcoind -testnet -rpcuser=user -rpcpassword=pass
```

## Testing with 5 and gaps (3+2)
Tests conducted with: ```./key800 -n 5 -s secret-test.csv -p public-test.csv -N test -1 43543 -2 24566 -3 645737 -4 67843```

Payin to

```
Ethereum Address,Bitcoin Address (P2PKH-Base58)
0x3BFA6A123b5f57aCE4822eB81a40492d07770AD6,miYzFVaPRJxuF4cJCrBjH1du72BwWqxLVy
0x16b65f95a74172ba17D74cc32595A873f32cfa56,mpzKwaqF12SwZA5HTAwcQLQUnj7GCiQ5ke
0xc938dC49A05FDf46c2078a8E13Fe66e28CE1CcDd,n1LcEGxDDuUu48gj1jhV8pP5UyuWBm5vkN
0xD79F0B989FF359e1E642e46500DEa42880722a64,moG9Pfdrm1VqifdEgJ39x9ghE2kG5YdZfh
0x797f6560A77df0992A4Eec1C3f6A72eE956f370e,mzxuricfQ4jgdjAe8xXcjChjrzXRnjFR6m

```

### Ethereum Rinkeby
The adresses were filled as follows. 
```
https://rinkeby.etherscan.io/address/0x3BFA6A123b5f57aCE4822eB81a40492d07770AD6 0.0
https://rinkeby.etherscan.io/address/0x16b65f95a74172ba17D74cc32595A873f32cfa56 0.1
https://rinkeby.etherscan.io/address/0xc938dC49A05FDf46c2078a8E13Fe66e28CE1CcDd 0.2 
https://rinkeby.etherscan.io/address/0xD79F0B989FF359e1E642e46500DEa42880722a64 0.3
https://rinkeby.etherscan.io/address/0x797f6560A77df0992A4Eec1C3f6A72eE956f370e 0.0
```

Consolidation to: 0x3dcBbCC53C8a9B2f965d89B92b3d0041da24eAE6

```
./consolidate1000-eth -s ../key800/secret-test.csv -N rinkeby -e 0x3dcBbCC53C8a9B2f965d89B92b3d0041da24eAE6
```

Output:

```
https://rinkeby.etherscan.io/address/0x3dcBbCC53C8a9B2f965d89B92b3d0041da24eAE6 0.599937

https://rinkeby.etherscan.io/address/0x3BFA6A123b5f57aCE4822eB81a40492d07770AD6 0.0
https://rinkeby.etherscan.io/address/0x16b65f95a74172ba17D74cc32595A873f32cfa56 0.0
https://rinkeby.etherscan.io/address/0xc938dC49A05FDf46c2078a8E13Fe66e28CE1CcDd 0.0 
https://rinkeby.etherscan.io/address/0xD79F0B989FF359e1E642e46500DEa42880722a64 0.0
https://rinkeby.etherscan.io/address/0x797f6560A77df0992A4Eec1C3f6A72eE956f370e 0.0
```

### Bitcoin Testnet3
The adresses were filled as follows. 
```
https://live.blockcypher.com/btc-testnet/address/miYzFVaPRJxuF4cJCrBjH1du72BwWqxLVy 0.98567784
https://live.blockcypher.com/btc-testnet/address/mpzKwaqF12SwZA5HTAwcQLQUnj7GCiQ5ke 0.0
https://live.blockcypher.com/btc-testnet/address/n1LcEGxDDuUu48gj1jhV8pP5UyuWBm5vkN 1.04248201
https://live.blockcypher.com/btc-testnet/address/moG9Pfdrm1VqifdEgJ39x9ghE2kG5YdZfh 5.37656688
https://live.blockcypher.com/btc-testnet/address/mzxuricfQ4jgdjAe8xXcjChjrzXRnjFR6m 0.0
```

Consolidation to: ms2Dt98uj11zvQ1hdPt3T5bD3QLRZXWGF, executed: 

```
./consolidate1000-btc -s ../key800/secret-test.csv -N test -U test -p me -b ms2Dt98uj11zvQ1hdPt3T5bD3QLRZXWGFC
```

Output: 

```
https://live.blockcypher.com/btc-testnet/address/ms2Dt98uj11zvQ1hdPt3T5bD3QLRZXWGFC 3.32815494

https://live.blockcypher.com/btc-testnet/address/miYzFVaPRJxuF4cJCrBjH1du72BwWqxLVy 0.0
https://live.blockcypher.com/btc-testnet/address/mpzKwaqF12SwZA5HTAwcQLQUnj7GCiQ5ke 0.0
https://live.blockcypher.com/btc-testnet/address/n1LcEGxDDuUu48gj1jhV8pP5UyuWBm5vkN 0.0
https://live.blockcypher.com/btc-testnet/address/moG9Pfdrm1VqifdEgJ39x9ghE2kG5YdZfh 0.0
https://live.blockcypher.com/btc-testnet/address/mzxuricfQ4jgdjAe8xXcjChjrzXRnjFR6m 0.0
```