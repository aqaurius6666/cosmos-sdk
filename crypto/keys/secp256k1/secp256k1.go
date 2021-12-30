package secp256k1

import (
	"bytes"
	"crypto/sha256"
	"crypto/subtle"
	"io"
	"math/big"

	secp256k1 "github.com/btcsuite/btcd/btcec"
	"github.com/tendermint/tendermint/crypto"
	"golang.org/x/crypto/ripemd160"

	cryptotypes "github.com/aqaurius6666/cosmos-sdk/crypto/types"
)

var (
	_ cryptotypes.PrivKey = &PrivKey{}
)

type PrivKey struct {
	Key []byte
}

type PubKey struct {
	Key []byte
}

const (
	PrivKeySize = 32
	keyType     = "secp256k1"
	PrivKeyName = "tendermint/PrivKeySecp256k1"
	PubKeyName  = "tendermint/PubKeySecp256k1"
)

// Bytes returns the byte representation of the Private Key.
func (privKey *PrivKey) Bytes() []byte {
	return privKey.Key
}

// PubKey performs the point-scalar multiplication from the privKey on the
// generator point to get the pubkey.
func (privKey *PrivKey) PubKey() cryptotypes.PubKey {
	_, pubkeyObject := secp256k1.PrivKeyFromBytes(secp256k1.S256(), privKey.Key)
	pk := pubkeyObject.SerializeCompressed()
	return &PubKey{Key: pk}
}

// Equals - you probably don't need to use this.
// Runs in constant time based on length of the
func (privKey *PrivKey) Equals(other cryptotypes.PrivKey) bool {
	return subtle.ConstantTimeCompare(privKey.Bytes(), other.Bytes()) == 1
}

// GenPrivKey generates a new ECDSA private key on curve secp256k1 private key.
// It uses OS randomness to generate the private key.
func GenPrivKey() *PrivKey {
	return &PrivKey{Key: genPrivKey(crypto.CReader())}
}

// genPrivKey generates a new secp256k1 private key using the provided reader.
func genPrivKey(rand io.Reader) []byte {
	var privKeyBytes [PrivKeySize]byte
	d := new(big.Int)
	for {
		privKeyBytes = [PrivKeySize]byte{}
		_, err := io.ReadFull(rand, privKeyBytes[:])
		if err != nil {
			panic(err)
		}

		d.SetBytes(privKeyBytes[:])
		// break if we found a valid point (i.e. > 0 and < N == curverOrder)
		isValidFieldElement := 0 < d.Sign() && d.Cmp(secp256k1.S256().N) < 0
		if isValidFieldElement {
			break
		}
	}

	return privKeyBytes[:]
}

var one = new(big.Int).SetInt64(1)

func GenPrivKeyFromSecret(secret []byte) *PrivKey {
	secHash := sha256.Sum256(secret)

	fe := new(big.Int).SetBytes(secHash[:])
	n := new(big.Int).Sub(secp256k1.S256().N, one)
	fe.Mod(fe, n)
	fe.Add(fe, one)

	feB := fe.Bytes()
	privKey32 := make([]byte, PrivKeySize)
	// copy feB over to fixed 32 byte privKey32 and pad (if necessary)
	copy(privKey32[32-len(feB):32], feB)

	return &PrivKey{Key: privKey32}
}

//-------------------------------------

var _ cryptotypes.PubKey = &PubKey{}

// PubKeySize is comprised of 32 bytes for one field element
// (the x-coordinate), plus one byte for the parity of the y-coordinate.
const PubKeySize = 33

// Address returns a Bitcoin style addresses: RIPEMD160(SHA256(pubkey))
func (pubKey *PubKey) Address() crypto.Address {
	if len(pubKey.Key) != PubKeySize {
		panic("length of pubkey is incorrect")
	}

	sha := sha256.Sum256(pubKey.Key)
	hasherRIPEMD160 := ripemd160.New()
	hasherRIPEMD160.Write(sha[:]) // does not error
	return crypto.Address(hasherRIPEMD160.Sum(nil))
}

// Bytes returns the pubkey byte format.
func (pubKey *PubKey) Bytes() []byte {
	return pubKey.Key
}

func (pubKey *PubKey) Equals(other cryptotypes.PubKey) bool {
	return bytes.Equal(pubKey.Bytes(), other.Bytes())
}
