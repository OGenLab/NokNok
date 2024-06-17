package local

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"strings"

	"qoobing.com/gomod/log"
)

type ecdsaKey ecdsa.PrivateKey

// NewEcdsaKeyBySeeds new an ecdsaKey by seeds
func NewEcdsaKeyBySeeds(key string) ecdsaKey {
	// data masking
	length := len(key)
	keyMask := key[:6] + strings.Repeat("*", length-12) + key[length-6:]
	log.Infof("key mask %v", keyMask)
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), strings.NewReader(key))
	if err != nil {
		panic("NewEcdsaKey generate privatekey failed:" + err.Error())
	}
	log.Infof("public key: %v", privateKey.PublicKey)
	return ecdsaKey(*privateKey)
}

// Sign sign a data, returns encoded signature or error
func (ek ecdsaKey) Sign(data []byte) ([]byte, error) {
	hash := sha256.Sum256(data)
	privateKey := (*ecdsa.PrivateKey)(&ek)
	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	if err != nil {
		return []byte{}, err
	}

	return sig, nil
}

// Verify verify sign, returns nil if signature is valid
func (ek ecdsaKey) Verify(data []byte, sig []byte) error {
	hash := sha256.Sum256(data)
	privateKey := (*ecdsa.PrivateKey)(&ek)
	if valid := ecdsa.VerifyASN1(&privateKey.PublicKey, hash[:], sig); !valid {
		return errors.New("invalid sig")
	}
	return nil
}

// GetPublicKey get the verify publickey
func (ek ecdsaKey) GetPublicKey() interface{} {
	privateKey := (*ecdsa.PrivateKey)(&ek)
	return &privateKey.PublicKey
}

// Encrypt encrypt data by kms
func (ek ecdsaKey) Encrypt(plaintext []byte) (crypttext []byte, err error) {
	panic("unsuported method")
}

// Decrypt decrypt data by kms
func (ek ecdsaKey) Decrypt(crypttext []byte) (plaintext []byte, err error) {
	panic("unsuported method")
}
