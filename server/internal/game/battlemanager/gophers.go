package battlemanager

import (
	"github.com/lonng/nano/session"
	"lightspeed.2dao3.com/nokserver/common/entities"
	"lightspeed.2dao3.com/nokserver/consts/errcode"
	"lightspeed.2dao3.com/nokserver/internal/protocol"
	"lightspeed.2dao3.com/nokserver/model"
)

type GophersRequest struct {
}

type GophersResponse struct {
	protocol.BaseResponse
	Gophers []*protocol.EntitiesInfo `json:"gophers"`
}

func (u *BattleManager) Gophers(s *session.Session, req *GophersRequest) error {
	var (
		m   = model.NewModelDefault()
		uid = s.UID()
	)
	defer m.Close()

	if uid <= 0 {
		return s.Response(&GophersResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.NOT_LOGIN, nil),
		})
	}

	gopherSet := entities.EntSet.GopherSet()
	gophers := make([]*protocol.EntitiesInfo, 0, len(gopherSet))
	for _, gopher := range gopherSet {
		gophers = append(gophers, &protocol.EntitiesInfo{
			Id:          gopher.F_entities_id,
			Season:      gopher.F_season,
			Name:        gopher.F_name,
			ImageUrl:    gopher.F_image_url,
			Probability: gopher.F_probability,
		})
	}

	return s.Response(&GophersResponse{
		BaseResponse: protocol.NewBaseResponse(errcode.SUCCESS, nil),
		Gophers:      gophers,
	})
}
