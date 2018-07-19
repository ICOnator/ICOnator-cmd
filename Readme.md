# Key Generation

An important process for an ICO with ICOnator is to generate the Bitcoin and 
Ethereum public and private keys for investors to payin. Its important to mention
that it is very important to keep the private keys secure, as they will be
potentially worth millions. For the previous ICOs we came up with the following
best practice process. This process requires at least 2 people (Alice, Bob, Carol). 
Carol can be one or more person. 

## Preparations
* For the key generation, Alice needs to have git and go installed on their 
machine. 
* As well as one Laptop, new 4 USB sticks from different brands, and a hardware
wallet.
* For the consolidation, a Bitcoin core full client and a Geth full 
client are necessary. Syncing may take a few days, so this needs to be done 
beforehand.
* **Important**: sanity test and getting familiar with the tools. Download and build all command
line tools. Generate 5 keys and payin small amounts of Bitcoin and Ethers in the
mainnet. Run the consolidation script that will transfer these Bitcoins and 
Ethers back to your account. Make sure the coins arrive and can be further used 
(e.g., transfer it to another account). Once you have verified that the funds 
are all there and transferable (minus fees), you are familiar with the process.

## Generate private and public keys
1. Buy a laptop that was never connected to the Internet, 
don't connect it to the Internet! Alice buys the laptop, while
Bob unseals it. Unsealing is supervised by Alice (and Carol).
 
1. Setup the laptop. Setup is done by Alice, supervised by Bob (and Carol).
A password for the laptop is only known by Alice and Bob. If a bootable stick 
with Linux is used, creating the stick is done by Alice, while Bob (and Carol) 
checks if the image on the stick matches the Linux image. For the 
comparison, ```cmp /tmp/ubuntu.img /dev/sdX -n 500GB``` 
can be used.

1. Alice checks out this repository and builds keygen with. This process is 
supervised by Bob (and Carol): 
   ``` 
   git clone https://github.com/ICOnator/ICOnator-cmd.git
   cd ICOnator-cmd/key800
   go get && go build
   ```

1. If the binaries support reproducible builds, the checksum of the binaries are 
compared. If not, the source code has to be checked by Bob (and Carol), and Bob 
builds the binaries supervised by Alice (and Carol).

1. Alice copies the binaries to a USB stick and copies it to the new Laptop.

1. Bob verifies the checksum of the binary on the new laptop matches with 
the newly build binaries. This process is supervised by Alice (and Carol).

1. Alice runs the key generation script with 100'000 keys and specifies a file
to store the private and public keys.

1. Bob copies the private keys to the 2 USB sticks, and public to the 
remaining stick. The public key stick has to be visibly marked as the one with
the public keys to avoid mixing them up.

1. Alice, Bob (and Carol) go to the bank vault and safely store the laptop and the
3 USB keys with the private keys in the bank vault. To access the bank vault, 
Alice and Bob are required to go there. 

1. The stick with the public key can now be used for the payin in 
ICOnator.

##Consolidation
1. Make sure your laptop has the Bitcoin core full client and Ethereum full client
installed. Since this laptop is connected to the Internet, make sure its not compromised.
As no private keys are stored withing the Ethereum full client, the Ethereum full
client can be installed on one machine, while the consolidation script runs on an other
machine. However, private keys are stored in the Bitcoin core client, thus it is
important to be sure that the machine is not compromised.

1. Alice checks out this repository and builds consolidate-eth and consolidate-btc
 with. This process is supervised by Bob (and Carol): 
   ``` 
   git clone https://github.com/ICOnator/ICOnator-cmd.git
   cd ICOnator-cmd/consolidate-eth
   go get && go build
   cd ../consolidate-btc
   go get && go build
   ```
1. If the binaries support reproducible builds, the checksum of the binaries are 
compared. If not, the source code has to be checked by Bob (and Carol), and Bob 
builds the binaries supervised by Alice (and Carol).

1. Alice, Bob (and Carol) go to the bank and get one USB stick, the rest stays in
the vault.

1. The private keys are copied to the laptop with the consolidation binaries by 
Bob. The scripts are run sequentially by Alice. Alice and Bob (and Carol) will 
supervise the process. 

1. The consolidation address is generated from the hardware wallet. Setup first the
hardware wallet. 

1. Read out seed and write to paper (mnemonic 24 word seed). Define a PIN or 
roll dice and write down to paper. Copy addresses from Trezor and save somewhere 
(later needed for contract owner, consolidation addresses). Validate address's 
checksum by validating it with e.g., MyetherWallet.

1. Run the consolidation scripts with the consolidation address from the hardware
wallet. This will take some time. After consolidation check if all funds have
been transferred, watch out for error messages.  

1. Laminate seed on paper and PIN and bring to bank vault, bring the hardware
wallet and the private keys back to the vault as well. There may be Airdrops or coins due to forks
on these accounts. To access the funds, Bob and Alice need to go the bank vault,
get the trezor and bring it back again. 

##Miting

tbd.

* Always 3 people present
* Deploy contract and mint tokens from online laptop
* Set mintDone to true
* Check and validate total amount of tokens
* Save and write down contract address
* Change owner of contract to the address of step 1