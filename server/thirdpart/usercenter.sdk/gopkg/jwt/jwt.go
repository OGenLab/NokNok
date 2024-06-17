// Package jwt is jwt(json web token) helper for web3wallet.
package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"lightspeed.2dao3.com/wallet/usercenter/sdk/kms"
	"qoobing.com/gomod/log"
)

var (
	jwtkms            kms.Kms = nil
	kmsSignMethodUatk         = NewKmsSignMethod("ES256UATK", "uatk")
	kmsSignMethodVatk         = NewKmsSignMethod("ES256VATK", "vatk")
	kmsSignMethodPatk         = NewKmsSignMethod("ES256PATK", "patk")
)

// LoginClaims usercenter login claims in jwt.
// @UserId: the usercenter userid
type LoginClaims struct {
	jwt.StandardClaims
	AppId  string `json:"appid"` // AppId
	UserId uint64 `json:"uid"`   // UserId
	KeyId  string `json:"keyid"`
}

// GenerateUatk generate user access token
//
//	@userid the user's id to generate token
func GenerateUatk(userid uint64, stdclaims jwt.StandardClaims, keyid string) (uatk string, err error) {
	var claims = LoginClaims{
		UserId:         userid,
		StandardClaims: stdclaims,
		KeyId:          keyid,
	}

	token := jwt.NewWithClaims(kmsSignMethodUatk, claims)
	uatk, err = token.SignedString(keyid)
	return uatk, err
}

// VerifyUatk verify the uatk is valid or not,
// return userid if no error.
func VerifyUatk(uatk string) (userid uint64, expireTime int64, err error) {
	var (
		claims  = LoginClaims{}
		keyfunc = func(token *jwt.Token) (interface{}, error) {
			if cl, ok := token.Claims.(*LoginClaims); ok {
				log.Infof("key id: %v", cl.KeyId)
				return cl.KeyId, nil
			}
			panic("unreachable code")
		}
	)
	if _, err = jwt.ParseWithClaims(uatk, &claims, keyfunc); err != nil {
		return 0, 0, err
	}
	return claims.UserId, claims.ExpiresAt, nil
}

// PayClaims pay claims in jwt.
// @To: the verify target address email(or phone number)
// @Vcode: the verify code send to user
// @VcsId: the verify code scene id
type PayClaims struct {
	jwt.StandardClaims
	UserId  uint64 `json:"uid"`     // UserId
	Address string `json:"address"` // Pay accout address
	KeyId   string `json:"keyid"`
}

// GeneratePatk generate user pay access token
//
//	@userid the user's id to generate token
//	@address the user's account address to pay
func GeneratePatk(userid uint64, address string, stdclaims jwt.StandardClaims, keyid string) (patk string, err error) {
	var claims = PayClaims{
		UserId:         userid,
		Address:        address,
		StandardClaims: stdclaims,
		KeyId:          keyid,
	}

	token := jwt.NewWithClaims(kmsSignMethodPatk, claims)
	patk, err = token.SignedString(keyid)
	return patk, err
}

// VerifyPatk verify the patk is valid or not,
// return (userid, address) if no error.
func VerifyPatk(patk string) (userid uint64, address string, err error) {
	var (
		claims  = PayClaims{}
		keyfunc = func(token *jwt.Token) (interface{}, error) {
			if cl, ok := token.Claims.(*PayClaims); ok {
				log.Infof("key id: %v", cl.KeyId)
				return cl.KeyId, nil
			}
			panic("unreachable code")
		}
	)

	_, err = jwt.ParseWithClaims(patk, &claims, keyfunc)
	if err != nil {
		return 0, "", err
	}

	return claims.UserId, claims.Address, nil
}

// VcodeClaims verify code claims in jwt.
// @To: the verify target address email(or phone number)
// @Vcode: the verify code send to user
// @VcsId: the verify code scene id
type VcodeClaims struct {
	jwt.StandardClaims
	To    string `json:"to"`    // Vcode SendTo Address
	Vcode secstr `json:"vcode"` // Vcode (plain in struct & crypted in json)
	VcsId string `json:"vcsid"` // Vcode Scene Id
	KeyId string `json:"keyid"`
}

// GenerateVatk generate verify-code token
//
//	@vcsid the VCodeSceneID for vcode
//	@vcode the vcode of the token
//	@to    the vcode destination
func GenerateVatk(vcsid, vcode, to string, stdclaims jwt.StandardClaims, keyid string) (vatk string, err error) {
	var claims = VcodeClaims{
		To:             to,
		VcsId:          vcsid,
		Vcode:          secstr(vcode),
		StandardClaims: stdclaims,
		KeyId:          keyid,
	}

	token := jwt.NewWithClaims(kmsSignMethodVatk, claims)
	vatk, err = token.SignedString(keyid)
	return vatk, err
}

// VerifyVatk verify the vatk is valid or not,
// return (vscid, vcode, to) if no error.
func VerifyVatk(vatk string) (vscid, vcode, to string, err error) {
	var (
		claims  = VcodeClaims{}
		keyfunc = func(token *jwt.Token) (interface{}, error) {
			if cl, ok := token.Claims.(*VcodeClaims); ok {
				log.Infof("key id: %v", cl.KeyId)
				return cl.KeyId, nil
			}
			panic("unreachable code")
		}
	)

	_, err = jwt.ParseWithClaims(vatk, &claims, keyfunc)
	if err != nil {
		return "", "", "", err
	}

	return claims.VcsId, string(claims.Vcode), claims.To, nil
}

// Init init jwt
func Init(ucjwtkms kms.Kms) (err error) {
	jwtkms = ucjwtkms
	return nil
}
