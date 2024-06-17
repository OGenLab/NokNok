package battlemanager

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/lonng/nano/session"
	"lightspeed.2dao3.com/nokserver/common/entities"
	"lightspeed.2dao3.com/nokserver/config"
	"lightspeed.2dao3.com/nokserver/consts/errcode"
	"lightspeed.2dao3.com/nokserver/internal/protocol"
	"lightspeed.2dao3.com/nokserver/model"
	"lightspeed.2dao3.com/nokserver/model/cache"
	"qoobing.com/gomod/str"
)

type GopherLocationsRequest struct {
}

type GopherLocationsResponse struct {
	protocol.BaseResponse
	GopherLocations []*GopherLocation `json:"gopherLocations"`
}

type GopherLocation struct {
	Matrix []*location `json:"matrix"`
	Round  string      `json:"round"`
}

type location struct {
	Index int    `json:"index"`
	Id    uint64 `json:"id"`
}

func (u *BattleManager) GopherLocations(s *session.Session, req *GopherLocationsRequest) error {
	var (
		m   = model.NewModelDefault()
		uid = s.UID()
	)
	defer m.Close()

	if uid <= 0 {
		return s.Response(&GopherLocationsResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.NOT_LOGIN, nil),
		})
	}

	resp := &GopherLocationsResponse{
		BaseResponse: protocol.NewBaseResponse(errcode.SUCCESS, nil),
	}

	var (
		round = config.Instance().NokConfig.GopherRound
		cnt   = config.Instance().NokConfig.GopherCnt
	)

	for i := 0; i < round; i++ {
		gl := randGopherLocation(cnt)
		resp.GopherLocations = append(resp.GopherLocations, gl)
		// 写入redis
		if err := setGopherLocation(m.Redis, gl); err != nil {
			return s.Response(&GopherLocationsResponse{
				BaseResponse: protocol.NewBaseResponse(errcode.REDIS_ERROR, err),
			})
		}
	}

	return s.Response(resp)
}

func randGopherLocation(gopherCnt int) *GopherLocation {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	cnt := r.Intn(gopherCnt) + 1

	// 根据cnt在从[0,gopherCnt)中进行cnt轮随机选择得到cnt个位置保存为一个数组
	matrix := make([]*location, 0, cnt)
	set := make(map[int]bool)
	for j := 0; j < cnt; j++ {
		index := r.Intn(gopherCnt)
		if set[index] {
			continue
		}

		set[index] = true
		gopherIds := entities.EntSet.GopherIds()
		id := gopherIds[r.Intn(len(gopherIds))]
		matrix = append(matrix, &location{
			Index: index,
			Id:    id,
		})
	}

	return &GopherLocation{
		Matrix: matrix,
		Round:  str.GetRandomNumString(13),
	}
}

func setGopherLocation(rds redis.Conn, gl *GopherLocation) error {
	data, _ := json.Marshal(gl)
	expireTime := config.Instance().NokConfig.RefreshInterval
	return cache.SetKey(rds, []interface{}{gl.Round, string(data), "EX", expireTime})
}
