package local

import "lightspeed.2dao3.com/wallet/usercenter/sdk/kms/util"

type aesKey string

// NewAesKey new an aesKey
func NewAesKey(key string) aesKey {
	return aesKey(key)
}

// Sign sign a data, returns encoded signature or error
func (ak aesKey) Sign(data []byte) ([]byte, error) {
	panic("unsuported method")
}

// Verify verify sign, returns nil if signature is valid
func (ak aesKey) Verify(data []byte, sign []byte) error {
	panic("unsuported method")
}

// GetPublicKey get the verify publickey
func (ak aesKey) GetPublicKey() interface{} {
	return string(ak)
}

// Encrypt encrypt data by kms
func (ak aesKey) Encrypt(plaintext []byte) (crypttext []byte, err error) {
	key := []byte(ak)
	return util.AesEncrypt(plaintext, key)
}

// Decrypt decrypt data by kms
func (ak aesKey) Decrypt(crypttext []byte) (plaintext []byte, err error) {
	key := []byte(ak)
	return util.AesDecrypt(plaintext, key)
}
