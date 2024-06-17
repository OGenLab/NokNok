// Package local is an implement of kms with local storage.
package local

import (
	"fmt"

	"lightspeed.2dao3.com/wallet/usercenter/sdk/kms/util"
)

// LocalKms locak kms
type LocalKms struct {
	allRegiteredCryptoClazz map[string]*util.InnerCryptClass
}

// NewLocalKms new a LocalKms instance
func NewLocalKms() *LocalKms {
	return &LocalKms{
		allRegiteredCryptoClazz: map[string]*util.InnerCryptClass{},
	}
}

// Register register a cryptclass
func (lk *LocalKms) Register(c util.InnerCryptClass) {
	lk.allRegiteredCryptoClazz[c.Name] = &c
}

func (lk LocalKms) getClassAndKey(scene, keyid string) (*util.InnerCryptClass, util.Key, error) {
	// Step 1. get class
	c, ok := lk.allRegiteredCryptoClazz[scene]
	if !ok {
		return nil, nil, fmt.Errorf("unregister crypto class:[%s]", scene)
	}

	// Step 2. get key
	k, okk := c.Keys[keyid]
	if !okk {
		return nil, nil, fmt.Errorf("crypto class[%s] keyid[%s] not found", scene, keyid)
	}
	return c, k, nil
}

// Sign sign a data, returns encoded signature or error
//
//	@scene scene for public key, one of {patk, uatk, vatk}
//	@keyid keyid for get privatekey to sign
//	@data  target data
func (lk LocalKms) Sign(scene, keyid string, data []byte) ([]byte, error) {
	// Step 1. get class & key
	_, k, err := lk.getClassAndKey(scene, keyid)
	if err != nil {
		return []byte{}, fmt.Errorf("local kms get key error:%s", err)
	}

	// Step 2. do sign
	sig, err := k.Sign(data)
	if err != nil {
		return []byte{}, fmt.Errorf("local kms sign error:%s", err)
	}
	return sig, nil
}

// Verify verify sign, returns nil if signature is valid
//
//	@scene scene for public key, one of {patk, uatk, vatk}
//	@keyid keyid for get publickey to verify token
//	@data  target data
func (lk LocalKms) Verify(scene, keyid string, data []byte, sign []byte) error {
	// Step 1. get class & key
	_, k, err := lk.getClassAndKey(scene, keyid)
	if err != nil {
		return fmt.Errorf("local kms get key error:%s", err)
	}

	// Step 2. do verify
	err = k.Verify(data, sign)
	return err
}

// GetPublicKey get the verify publickey
//
//	@scene scene for public key, one of {patk, uatk, vatk}
//	@keyid keyid for get publickey to verify token
func (lk LocalKms) GetPublicKey(scene, keyid string) interface{} {
	// Step 1. get class & key
	_, k, err := lk.getClassAndKey(scene, keyid)
	if err != nil {
		return nil
	}
	return k
}

// Encrypt encrypt data by kms
func (lk LocalKms) Encrypt(scene, keyid string, plaintext []byte) (crypttext []byte, err error) {
	// Step 1. get class & key
	_, k, err := lk.getClassAndKey(scene, keyid)
	if err != nil {
		return []byte{}, fmt.Errorf("local kms get key error:%s", err)
	}

	// Step 2. do encrypt
	return k.Encrypt(plaintext)
}

// Decrypt decrypt data by kms
func (lk LocalKms) Decrypt(scene, keyid string, crypttext []byte) (plaintext []byte, err error) {
	// Step 1. get class & key
	_, k, err := lk.getClassAndKey(scene, keyid)
	if err != nil {
		return []byte{}, fmt.Errorf("local kms get key error:%s", err)
	}

	// Step 2. do decrypt
	return k.Decrypt(crypttext)
}
