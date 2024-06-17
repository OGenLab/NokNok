package protocol

import (
	"testing"
)

func TestParseInitData(t *testing.T) {
	initData := `user=%7B%22id%22%3A6986244361%2C%22first_name%22%3A%22Qoo%22%2C%22last_name%22%3A%22Bryan%22%2C%22username%22%3A%22QooBryan%22%2C%22language_code%22%3A%22zh-hans%22%2C%22allows_write_to_pm%22%3Atrue%7D&chat_instance=-880176127099621663&chat_type=sender&start_param=TTTTTTREFERCODETTTTTTTT&auth_date=1718605698&hash=3d6c0d79bcf47431efd8206b82a5cff572f506791e277d3cf465b84cec37daf2`
	// res, err := telegoutil.ValidateWebAppData(config.Instance().TgConfig.BotToken, initData)
	res, err := ParseInitData(initData)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", res.User)
}
