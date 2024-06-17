package battlemanager

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
func TestRandGopherLocation(t *testing.T) {
	cnt := 8
	res := randGopherLocation(cnt)
	for _, v := range res.Matrix {
		t.Logf("%+v", v)
	}
	t.Logf("%+v", res)
}

func TestRandomSelectNumber(t *testing.T) {
	v := randomSelectNumber()
	t.Logf("%+v", v)
}

func TestGopherLocation(t *testing.T) {
	req := &GopherLocationsRequest{}
	DefaultBattleManager.GopherLocations(s, req)
	resp := newwork.LastResponse().(*GopherLocationsResponse)
	data, _ := json.Marshal(resp)
	t.Logf("%+v", string(data))
}

func TestHammerGopher(t *testing.T) {
	reqs := []*HammerGopherRequest{
		{
			Round: "1168631676695",
			Index: 6,
		},
		{
			Round: "9861189672949",
			Index: 4,
		},
		{
			Round: "1651702278068",
			Index: 5,
		},
	}

	for _, req := range reqs {
		DefaultBattleManager.HammerGopher(s, req)
		resp := newwork.LastResponse().(*HammerGopherResponse)
		data, _ := json.Marshal(resp)
		t.Logf("%+v", string(data))
	}
}

func TestGophers(t *testing.T) {
	req := &GophersRequest{}
	DefaultBattleManager.Gophers(s, req)
	resp := newwork.LastResponse().(*GophersResponse)
	data, _ := json.Marshal(resp)
	t.Logf("%+v", string(data))
}

func TestHammers(t *testing.T) {
	req := &HammersRequest{}
	DefaultBattleManager.Hammers(s, req)
	resp := newwork.LastResponse().(*HammersResponse)
	data, _ := json.Marshal(resp)
	t.Logf("%+v", string(data))
}

func TestLeaderboard(t *testing.T) {
	req := &LeaderboardRequest{
		Type: 1,
	}
	DefaultBattleManager.Leaderboard(s, req)
	resp := newwork.LastResponse().(*LeaderboardResponse)
	data, _ := json.Marshal(resp)
	t.Logf("%+v", string(data))
}

func TestPrizeDraw(t *testing.T) {
	req := &PrizeDrawRequest{
		PrizeType: 1,
	}
	DefaultBattleManager.PrizeDraw(s, req)
	resp := newwork.LastResponse().(*PrizeDrawResponse)
	data, _ := json.Marshal(resp)
	t.Logf("%+v", string(data))
}

func TestTasks(t *testing.T) {
	req := &TasksRequest{}
	DefaultBattleManager.Tasks(s, req)
	resp := newwork.LastResponse().(*TasksResponse)
	data, _ := json.Marshal(resp)
	t.Logf("%+v", string(data))
}

func TestEquipHammer(t *testing.T) {
	req := &EquipHammerRequest{
		HammerId: 4,
	}
	DefaultBattleManager.EquipHammer(s, req)
	resp := newwork.LastResponse().(*EquipHammerResponse)
	data, _ := json.Marshal(resp)
	t.Logf("%+v", string(data))
}
