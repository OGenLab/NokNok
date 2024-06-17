package local

import (
	"errors"
	"fmt"

	"qoobing.com/gomod/log"

	"lightspeed.2dao3.com/wallet/usercenter/sdk/kms/util"
)

type myaesKey struct {
	ver  int      // crypto version
	keys []string // Keys
}

type myaesCryptValue struct {
	Version        int    //首字节高4bit
	KeyIndex       uint16 //(版本1)首字节低4bit + 第二字节，共12bit
	EncryptedValue []byte //(版本1)第三字节开始的所有字节
}

// NewMyaesKey new a myaesKey
func NewMyaesKey(keys []string) myaesKey {
	return myaesKey{
		ver:  1, //TODO
		keys: keys,
	}
}

// Sign sign a data, returns encoded signature or error
func (mak myaesKey) Sign(data []byte) ([]byte, error) {
	panic("unsuported method")
}

// Verify verify sign, returns nil if signature is valid
func (mak myaesKey) Verify(data []byte, sign []byte) error {
	panic("unsuported method")
}

// GetPublicKey get the verify publickey
func (mak myaesKey) GetPublicKey() interface{} {
	panic("unsuported method")
}

// Encrypt encrypt data by kms
func (mak myaesKey) Encrypt(plaintext []byte) (crypttext []byte, err error) {
	cv := myaesCryptValue{
		Version:  mak.ver,
		KeyIndex: uint16(1),
	}

	key := mak.keys[cv.KeyIndex-1]
	cv.EncryptedValue, err = util.AesEncrypt(plaintext, []byte(key))
	if err != nil {
		return []byte{}, errors.New(fmt.Sprintf("aes encrypt error:%s", err.Error()))
	}

	return cv.encode()
}

// Decrypt decrypt data by kms
func (mak myaesKey) Decrypt(crypttext []byte) (plaintext []byte, err error) {
	var cv = new(myaesCryptValue)
	if err := cv.decode(crypttext); err != nil {
		return []byte{}, fmt.Errorf("decode error:'%s'", err.Error())
	}
	if cv.KeyIndex > uint16(len(mak.keys)) {
		return []byte{}, fmt.Errorf(
			"invalid data: 'keyindex(%d) is out of range'", cv.KeyIndex)
	}
	key := mak.keys[cv.KeyIndex-1]
	plain, err := util.AesDecrypt(cv.EncryptedValue, []byte(key))
	if err != nil {
		return []byte{}, fmt.Errorf("aesDecrypt error:'%s'", err.Error())
	}
	return plain, nil
}

// encode encode myaesCryptValue to []byte
func (cv *myaesCryptValue) encode() (crypted []byte, err error) {
	v := byte(byte(cv.Version<<4) | byte(cv.KeyIndex>>8))
	i := byte(cv.KeyIndex & 0x00FF)
	r := []byte{v, i}
	return append(r, cv.EncryptedValue...), nil
}

// decode decode []byte to myaesCryptValue
func (cv *myaesCryptValue) decode(crypted []byte) (err error) {
	if len(crypted) <= 2 {
		log.Debugf("crypted=[%s]", string(crypted))
		return fmt.Errorf("invalid encrypted data: 'len(crypted)=%d<=2", len(crypted))
	}
	cv.Version = int(crypted[0] >> 4)
	if cv.Version != 1 {
		return errors.New(fmt.Sprintf(
			"invalid data: 'version(%d) is not supported'", cv.Version))
	}
	cv.KeyIndex = uint16(crypted[0]&0x0F)<<8 | uint16(crypted[1])
	if cv.KeyIndex <= 0 {
		return errors.New(fmt.Sprintf(
			"invalid data: 'keyindex(%d) is out of range'", cv.KeyIndex))
	}

	cv.EncryptedValue = crypted[2:]
	return nil
}
