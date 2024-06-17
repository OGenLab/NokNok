package usermanager

import (
	"encoding/json"
	"testing"

	"github.com/lonng/nano/mock"
	"github.com/lonng/nano/session"
)

var newwork = mock.NewNetworkEntity()
var s = session.New(newwork)

func init() {
	s.Bind(6986244361)
}
func TestLogin(t *testing.T) {
	req := &LoginToGameServerRequest{
		Token:        "user%3D%257B%2522id%2522%253A6986244361%252C%2522first_name%2522%253A%2522Qoo%2522%252C%2522last_name%2522%253A%2522Bryan%2522%252C%2522username%2522%253A%2522QooBryan%2522%252C%2522language_code%2522%253A%2522zh-hans%2522%252C%2522allows_write_to_pm%2522%253Atrue%257D%26chat_instance%3D-880176127099621663%26chat_type%3Dsender%26start_param%3DTTTTTTREFERCODETTTTTTTT%26auth_date%3D1718589730%26hash%3D05f4e7c191a3b2634ae02368d565fe8afc207bb91f02b4a1641b65f6e9cb78b4&tgWebAppVersion=7.0&tgWebAppPlatform=macos&tgWebAppBotInline=1&tgWebAppThemeParams=%7B%22button_color%22%3A%22%232481cc%22%2C%22button_text_color%22%3A%22%23ffffff%22%2C%22text_color%22%3A%22%23000000%22%2C%22accent_text_color%22%3A%22%232481cc%22%2C%22link_color%22%3A%22%232481cc%22%2C%22secondary_bg_color%22%3A%22%23efeff3%22%2C%22bg_color%22%3A%22%23ffffff%22%2C%22section_bg_color%22%3A%22%23ffffff%22%2C%22destructive_text_color%22%3A%22%23ff3b30%22%2C%22hint_color%22%3A%22%23999999%22%2C%22subtitle_text_color%22%3A%22%23999999%22%2C%22section_header_text_color%22%3A%22%236d6d71%22%2C%22header_bg_color%22%3A%22%23efeff3%22%7D",
		Channel:      "test",
		ReferralCode: "DMEBKPC",
	}

	if err := DefaultUserManger.Login(s, req); err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", newwork.LastResponse().(*LoginToGameServerResponse))
}

func TestGetInfo(t *testing.T) {
	req := &GetInfoRequest{}
	if err := DefaultUserManger.GetInfo(s, req); err != nil {
		t.Fatal(err)
	}
	resp := newwork.LastResponse().(*GetInfoResponse)
	data, _ := json.Marshal(resp)
	t.Logf("%+v", string(data))
}

func TestBoost(t *testing.T) {
	req := &BoostRequest{}
	if err := DefaultUserManger.Boost(s, req); err != nil {
		t.Fatal(err)
	}
	resp := newwork.LastResponse().(*BoostResponse)
	data, _ := json.Marshal(resp)
	t.Logf("%+v", string(data))
}

func TestGetBoost(t *testing.T) {
	req := &GetBoostRequest{}
	if err := DefaultUserManger.GetBoost(s, req); err != nil {
		t.Fatal(err)
	}
	resp := newwork.LastResponse().(*GetBoostResponse)
	data, _ := json.Marshal(resp)
	t.Logf("%+v", string(data))
}

func TestGetHammers(t *testing.T) {
	req := &GetHammerRequest{}
	if err := DefaultUserManger.GetHammer(s, req); err != nil {
		t.Fatal(err)
	}
	resp := newwork.LastResponse().(*GetHammerResponse)
	data, _ := json.Marshal(resp)
	t.Logf("%+v", string(data))
}

func TestInvite(t *testing.T) {
	req := &InviteRequest{}
	if err := DefaultUserManger.Invite(s, req); err != nil {
		t.Fatal(err)
	}
	resp := newwork.LastResponse().(*InviteResponse)
	data, _ := json.Marshal(resp)
	t.Logf("%+v", string(data))
}

func TestReceiveAward(t *testing.T) {
	req := &ReceiveAwardRequest{
		TaskType: 1,
		TaskId:   1,
	}
	if err := DefaultUserManger.ReceiveAward(s, req); err != nil {
		t.Fatal(err)
	}
	resp := newwork.LastResponse().(*ReceiveAwardResponse)
	data, _ := json.Marshal(resp)
	t.Logf("%+v", string(data))
}
