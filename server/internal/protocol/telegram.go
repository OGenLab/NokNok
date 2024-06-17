package protocol

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"

	"net/url"
	"strconv"
	"strings"

	"github.com/mymmrac/telego/telegoutil"
	"lightspeed.2dao3.com/nokserver/config"
)

var secretKey = generateHMAC([]byte("WebAppData"), []byte(config.Instance().TgConfig.BotToken))

type WebAppInitData struct {
	QueryID      string      `json:"query_id"`
	User         *WebAppUser `json:"user"`
	Receiver     *WebAppUser `json:"receiver"`
	Chat         *WebAppChat `json:"chat"`
	ChatType     string      `json:"chat_type"`
	ChatInstance string      `json:"chat_instance"`
	StartParam   string      `json:"start_param"`
	CanSendAfter int         `json:"can_send_after"`
	AuthDate     int         `json:"auth_date"`
	Hash         string      `json:"hash"`
}
type WebAppUser struct {
	ID                    int    `json:"id"`
	IsBot                 bool   `json:"is_bot"`
	FirstName             string `json:"first_name"`
	LastName              string `json:"last_name"`
	Username              string `json:"username"`
	LanguageCode          string `json:"language_code"`
	IsPremium             bool   `json:"is_premium"`
	AddedToAttachmentMenu bool   `json:"added_to_attachment_menu"`
	AllowsWriteToPM       bool   `json:"allows_write_to_pm"`
	PhotoURL              string `json:"photo_url"`
}

type WebAppChat struct {
	ID       int    `json:"id"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Username string `json:"username"`
	PhotoURL string `json:"photo_url"`
}

func ValidWebAppData(initData string) (url.Values, error) {
	decodedData, err := url.QueryUnescape(initData)
	if err != nil {
		return nil, errors.New("telego: parse query: bad data")
	}

	appData, err := url.ParseQuery(decodedData)
	if err != nil {
		return nil, errors.New("telego: parse query: bad data")
	}

	hash := appData.Get("hash")
	if hash == "" {
		return nil, errors.New("telego: no hash found")
	}

	appData.Del("hash")

	appDataToCheck, _ := url.QueryUnescape(strings.ReplaceAll(appData.Encode(), "&", "\n"))

	secretKey := sha256.Sum256([]byte(config.Instance().TgConfig.BotToken))
	if hex.EncodeToString(hmacHash([]byte(appDataToCheck), secretKey[:])) != hash {
		return nil, errors.New("telego: invalid hash")
	}

	return appData, nil
}

func MockParseInitData(initData string) (*WebAppInitData, error) {
	var data *WebAppInitData
	if err := json.Unmarshal([]byte(initData), &data); err != nil {
		return nil, err
	}
	return data, nil
}

func ParseInitData(initData string) (*WebAppInitData, error) {
	params, err := telegoutil.ValidateWebAppData(config.Instance().TgConfig.BotToken, initData)
	if err != nil {
		return nil, err
	}

	initDataObject := &WebAppInitData{}

	for key, values := range params {
		value := values[0]
		switch key {
		case "query_id":
			initDataObject.QueryID = value
		case "chat_type":
			initDataObject.ChatType = value
		case "chat_instance":
			initDataObject.ChatInstance = value
		case "start_param":
			initDataObject.StartParam = value
		case "hash":
			initDataObject.Hash = value
		case "user", "receiver":
			user := &WebAppUser{}
			json.Unmarshal([]byte(value), user)
			if key == "user" {
				initDataObject.User = user
			} else {
				initDataObject.Receiver = user
			}
		case "chat":
			chat := &WebAppChat{}
			json.Unmarshal([]byte(value), chat)
			initDataObject.Chat = chat
		case "can_send_after", "auth_date":
			intValue, _ := strconv.Atoi(value)
			if key == "can_send_after" {
				initDataObject.CanSendAfter = intValue
			} else {
				initDataObject.AuthDate = intValue
			}
		}
	}

	return initDataObject, nil
}

func generateHMAC(key, data []byte) string {
	hash := hmacHash(key, data)
	return hex.EncodeToString(hash)
}

func hmacHash(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	_, _ = h.Write(data)
	return h.Sum(nil)
}
