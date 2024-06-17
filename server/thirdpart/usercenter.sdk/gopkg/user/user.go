// Package user provide an infomation getter to other service 
// who use web3wallet's user system.
package user

import (
	"bytes"
	"encoding/json"
	"net/http"

	"qoobing.com/gomod/log"
)

// User info
type User struct {
	UserId  uint64 `json:"userid"`
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

// request info
type requestArgs struct {
	Uatk string `json:"uatk"`
}

// Usercenter config info
type Usercenter struct {
	RootURL              string
	IsUserInfoValidation string
	AppId                string
	AppSecretKey         string
}

var user Usercenter

// GetUserInfo is to get user information given uatk
func GetUserInfo(uatk string) (usr User, err error) {
	var (
		appIdReq = "appid=" + user.AppId
		addrReq  = user.RootURL + "/innerapi/get_user_info"
	)
	if user.IsUserInfoValidation == "false" {
		return User{}, nil
	}
	requestarg := requestArgs{Uatk: uatk}
	jsonData, err := json.Marshal(requestarg)
	if err != nil {
		log.Warningf("%v, try again", err)
		return User{}, err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", addrReq, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Warningf("request error: %v, try again", err)
		return User{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-APP", appIdReq)
	resp, err := client.Do(req)
	if err != nil {
		log.Warningf("request error: %v, try again", err)
		return User{}, err
	}
	defer resp.Body.Close()
	var userOutput User
	err = json.NewDecoder(resp.Body).Decode(&userOutput)
	if err != nil {
		log.Warningf("decode error: %v, try again", err)
		return User{}, err
	}
	return userOutput, nil
}

// Init get_user_info config
func Init(usr Usercenter) {
	user = usr
}
