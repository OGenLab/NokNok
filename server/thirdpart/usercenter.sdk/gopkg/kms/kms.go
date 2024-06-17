// Package kms is a kms(Key Management System) for web3wallet.
package kms

import (
	"encoding/hex"
	"errors"
	"fmt"

	"lightspeed.2dao3.com/wallet/usercenter/sdk/kms/mycrypt/local"
	"lightspeed.2dao3.com/wallet/usercenter/sdk/kms/mycrypt/outer"
	"lightspeed.2dao3.com/wallet/usercenter/sdk/kms/util"
	"qoobing.com/gomod/log"
)

// Kms defined Kms(Key Management System) interface
type Kms interface {
	// Sign sign a data, returns encoded signature or error
	//   @scene scene for public key, one of {patk, uatk, vatk}
	//   @keyid keyid for get privatekey to sign
	//   @data  target data
	Sign(scene, keyid string, data []byte) ([]byte, error)
	// Verify verify sign, returns nil if signature is valid
	//   @scene scene for public key, one of {patk, uatk, vatk}
	//   @keyid keyid for get publickey to verify token
	//   @data  target data
	Verify(scene, keyid string, data []byte, sign []byte) error
	// GetPublicKey get the verify publickey
	//   @scene scene for public key, one of {patk, uatk, vatk}
	//   @keyid keyid for get publickey to verify token
	GetPublicKey(scene, keyid string) interface{}

	// Encrypt encrypt data by kms
	Encrypt(scene, keyid string, plaintext []byte) (crypttext []byte, err error)
	// Decrypt decrypt data by kms
	Decrypt(scene, keyid string, crypttext []byte) (plaintext []byte, err error)
}

// myKms
type myKms struct {
	awsKms             *local.LocalKms // aws kms intance
	localKms           *local.LocalKms // local kms intance
	sceneToKmsInstance map[string]Kms  // scene to kms instance
}

// gkms global kms in this package
var gkms = myKms{
	awsKms:             local.NewLocalKms(), // aws kms intance
	localKms:           local.NewLocalKms(), // local kms intance
	sceneToKmsInstance: map[string]Kms{},    // sceneToKmsInstance scene to kms instance
}

// CryptClass crypt class
type CryptClass struct {
	Type    string
	Name    string
	Version int
	Keys    []outer.KeysClass
}

const (
	// MYCRYPT_TYPE
	MYCRYPT_TYPE_AES_KEYS          = "TYPE_AES_KEYS"
	MYCRYPT_TYPE_MYAES_KEYS        = "TYPE_MYAES_KEYS"
	MYCRYPT_TYPE_ECDSA_AWS         = "TYPE_ECDSA_AWS"
	MYCRYPT_TYPE_ECDSA_SEEDS       = "TYPE_ECDSA_SEEDS"
	MYCRYPT_TYPE_ECDSA_PUBLIC_KEYS = "TYPE_ECDSA_PUBLIC_KEYS"
)

// Init init package
func Init(crypts []CryptClass) {
	gkms.Init(crypts)
}

// InitPublicKeys init kms publickeys
func InitPublicKeys(ip string, name string) (err error) {
	return gkms.InitPublicKeys(ip, name)
}

// GetGlobalKms get global kms instance
func GetGlobalKms() Kms {
	return &gkms
}

// Sign sign a data, returns encoded signature or error
func Sign(scene, keyid string, data []byte) ([]byte, error) {
	return gkms.Sign(scene, keyid, data)
}

// Verify verify sign, returns nil if signature is valid
func Verify(scene, keyid string, data []byte, sign []byte) error {
	return gkms.Verify(scene, keyid, data, sign)
}

// GetPublicKey get the verify publickey
func GetPublicKey(scene, keyid string) interface{} {
	return gkms.GetPublicKey(scene, keyid)
}

// Encrypt encrypt data by kms
func Encrypt(scene, keyid string, plaintext []byte) (crypttext []byte, err error) {
	return gkms.Encrypt(scene, keyid, plaintext)
}

// Decrypt decrypt data by kms
func Decrypt(scene, keyid string, crypttext []byte) (plaintext []byte, err error) {
	return gkms.Decrypt(scene, keyid, crypttext)
}

// EncryptToHexString encrypt to hex string
func EncryptToHexString(scene, keyid string, value []byte) (hexcrypted string, err error) {
	var b []byte
	if b, err = Encrypt(scene, keyid, value); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// DecryptFromHexString decrypt from hex string
func DecryptFromHexString(scene, keyid string, hexcrypted string) (decrypted []byte, err error) {
	if value, err := hex.DecodeString(hexcrypted); err != nil {
		return nil, errors.New("decrypt failed: decode string from hex failed")
	} else if decrypted, err = Decrypt(scene, keyid, value); err != nil {
		return nil, errors.New(fmt.Sprintf("decrypt failed: '%s'", err.Error()))
	}
	return decrypted, nil
}

// #################################################################
// Init initialize kms by config
//
//	@crypts the crypt class array. example:
//
// #################################################################
//
//	var crypts = []cryptClass{
//		0: {
//			Type:    MYCRYPT_TYPE_ECDSA_AWS,
//			Name:    "uatk",
//			Version: 1,
//			Keys:    []string{"dfak#8237*&^23@1sdjfiweuwrioeakdjfal@#kjfadks"},
//		},
//		1: {
//			Type:    MYCRYPT_TYPE_ECDSA_SEEDS,
//			Name:    "patk",
//			Version: 1,
//			Keys:    []string{"jalksjfou2422394289#@#1rfq24)(*##!Dfkhaoioj#@"},
//		},
//		2: {
//			Type:    MYCRYPT_TYPE_AES_KEYS,
//			Name:    "innervcode",
//			Version: 1,
//			Keys:    []string{"1234567890abcdef", "mhxzkhl2dao3kfsk"},
//		},
//		3: {
//			Type:    MYCRYPT_TYPE_AES_KEYS,
//			Name:    "rediscache",
//			Version: 1,
//			Keys:    []string{"1234567890abcdef"},
//		},
//	}
//
// #################################################################
func (mk myKms) Init(crypts []CryptClass) {
	for _, c := range crypts {
		if _, ok := mk.sceneToKmsInstance[c.Name]; !ok {
			log.Warningf("duplicate cryto calss '%s', will cover the previos one", c.Name)
		}
		newclass := util.InnerCryptClass{
			Name:    c.Name,
			Type:    c.Type,
			Version: c.Version,
			Keys:    map[string]util.Key{},
		}
		switch c.Type {
		case MYCRYPT_TYPE_MYAES_KEYS:
			var keysArr []string
			for _, key := range c.Keys {
				keysArr = append(keysArr, key.Key)
			}
			newclass.Keys[""] = local.NewMyaesKey(keysArr)
			mk.localKms.Register(newclass)
			mk.sceneToKmsInstance[c.Name] = mk.localKms
		case MYCRYPT_TYPE_ECDSA_SEEDS:
			for _, key := range c.Keys {
				log.Infof("ECDSA_SEEDSid %v", key.Id)
				newclass.Keys[key.Id] = local.NewEcdsaKeyBySeeds(key.Key)
			}
			mk.localKms.Register(newclass)
			log.Infof("name:%v", c.Name)
			mk.sceneToKmsInstance[c.Name] = mk.localKms
		case MYCRYPT_TYPE_AES_KEYS:
			for _, key := range c.Keys {
				newclass.Keys[key.Id] = local.NewAesKey(key.Key)
			}
			mk.localKms.Register(newclass)
			mk.sceneToKmsInstance[c.Name] = mk.localKms
		case MYCRYPT_TYPE_ECDSA_AWS:
			panic("TODO: unsupported type:" + c.Type)
		default:
			panic("unsupported type:" + c.Type)
		}
	}
}

func (mk myKms) InitPublicKeys(rootURL string, scene string) (err error) {
	if _, ok := mk.sceneToKmsInstance[scene]; !ok {
		log.Warningf("duplicate cryto calss '%s', will cover the previos one", scene)
	}
	newclass := util.InnerCryptClass{
		Name:    scene,
		Type:    MYCRYPT_TYPE_ECDSA_PUBLIC_KEYS,
		Version: 1,
		Keys:    map[string]util.Key{},
	}
	keysId, ecdsaKeys, err := outer.NewEcdsaKey(rootURL, scene)
	if err != nil {
		log.Warningf("get publickey error: %v", err)
		return err
	}
	for idx, keyid := range keysId {
		log.Infof("ECDSA_PUBLIC_KEYid %v, value %v", keyid, ecdsaKeys[idx])
		newclass.Keys[keyid] = ecdsaKeys[idx]
	}
	mk.localKms.Register(newclass)
	log.Infof("name: %v", scene)
	mk.sceneToKmsInstance[scene] = mk.localKms
	return nil
}

// Sign sign a data, returns encoded signature or error
//
//	@scene scene for public key, one of {patk, uatk, vatk}
//	@keyid keyid for get privatekey to sign
//	@data  target data
func (mk myKms) Sign(scene, keyid string, data []byte) ([]byte, error) {
	kms, ok := mk.sceneToKmsInstance[scene]
	if !ok {
		panic("unknown scene to call kms: scene=" + scene)
	}
	return kms.Sign(scene, keyid, data)
}

// Verify verify sign, returns nil if signature is valid
//
//	@scene scene for public key, one of {patk, uatk, vatk}
//	@keyid keyid for get publickey to verify token
//	@data  target data
func (mk myKms) Verify(scene, keyid string, data []byte, sign []byte) error {
	kms, ok := mk.sceneToKmsInstance[scene]
	if !ok {
		panic("unknown scene to call kms")
	}
	return kms.Verify(scene, keyid, data, sign)
}

// GetPublicKey get the verify publickey
//
//	@scene scene for public key, one of {patk, uatk, vatk}
//	@keyid keyid for get publickey to verify token
func (mk myKms) GetPublicKey(scene, keyid string) interface{} {
	kms, ok := mk.sceneToKmsInstance[scene]
	if !ok {
		panic("unknown scene to call kms")
	}
	return kms.GetPublicKey(scene, keyid)
}

// Encrypt encrypt data by kms
func (mk myKms) Encrypt(scene, keyid string, plaintext []byte) (crypttext []byte, err error) {
	kms, ok := mk.sceneToKmsInstance[scene]
	if !ok {
		panic("unknown scene to call kms")
	}
	return kms.Encrypt(scene, keyid, plaintext)
}

// Decrypt decrypt data by kms
func (mk myKms) Decrypt(scene, keyid string, crypttext []byte) (plaintext []byte, err error) {
	kms, ok := mk.sceneToKmsInstance[scene]
	if !ok {
		panic("unknown scene to call kms")
	}
	return kms.Decrypt(scene, keyid, crypttext)
}
