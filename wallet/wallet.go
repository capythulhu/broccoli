package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"github.com/thzoid/broccoli/hash"
	"golang.org/x/crypto/ripemd160"
)

// Address version
const VERSION = 0x04

// Network (0x00 = Main Net)
const NETWORK = 0x00

// Wallet address
type Address [25]byte

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

var coinbaseAddress = Address{}

func CoinbaseAddress() Address {
	address := Address{}
	copy(address[:], coinbaseAddress[:])
	return address
}

func (a Address) IsCoinbase() bool {
	return a == coinbaseAddress
}

func (a Address) String() string {
	return string(a[:])
}

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	public := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, public
}

func NewWallet() Wallet {
	privateKey, publicKey := NewKeyPair()
	wallet := Wallet{privateKey, publicKey}
	return wallet
}

func (w Wallet) Hash() [20]byte {
	versioned := append([]byte{VERSION}, w.PublicKey[:]...)
	hash := hash.Hash{}
	hash.SHA256(versioned)

	hasher := ripemd160.New()
	_, err := hasher.Write(hash[:])
	if err != nil {
		log.Panic(err)
	}

	publicRipeMd := [20]byte{}
	copy(publicRipeMd[:], hasher.Sum(nil))

	return publicRipeMd
}

func Checksum(hash []byte) [4]byte {
	firstHash := sha256.Sum256(hash)
	secondHash := sha256.Sum256(firstHash[:])

	checksum := [4]byte{}
	copy(checksum[:], secondHash[:4])

	return checksum
}

func (w Wallet) Address() Address {
	hash := w.Hash()
	versioned := append([]byte{NETWORK}, hash[:]...)
	checksum := Checksum(versioned)

	finalHash := append(versioned[:], checksum[:]...)
	address := Address{}
	copy(address[:], base58Encode(finalHash))

	return address
}
