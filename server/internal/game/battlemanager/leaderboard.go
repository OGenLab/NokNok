package battlemanager

import (
	"github.com/lonng/nano/session"
	"lightspeed.2dao3.com/nokserver/consts/errcode"
	"lightspeed.2dao3.com/nokserver/internal/protocol"
	"lightspeed.2dao3.com/nokserver/model"
	"lightspeed.2dao3.com/nokserver/model/t_hammer_entities"
)

type LeaderboardRequest struct {
	Type   int `json:"type"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type LeaderboardResponse struct {
	protocol.BaseResponse
	LeaderboardInfo []*protocol.LeaderboardInfo `json:"leaderboardInfo"`
}

func (u *BattleManager) Leaderboard(s *session.Session, req *LeaderboardRequest) error {
	var (
		m      = model.NewModelDefault()
		uid    = s.UID()
		limit  = req.Limit
		offset = req.Offset
	)
	defer m.Close()

	if uid <= 0 {
		return s.Response(&LeaderboardResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.NOT_LOGIN, nil),
		})
	}

	if limit == 0 || limit > 100 {
		limit = 100
	}

	players, err := t_hammer_entities.GetLeaderboard(m.DB, limit, offset)
	if err != nil {
		return s.Response(&LeaderboardResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.DATABASE_ERROR, err),
		})
	}

	leaderboardInfo := make([]*protocol.LeaderboardInfo, len(players))
	for i, v := range players {
		leaderboardInfo[i] = &protocol.LeaderboardInfo{
			Rank:  i + 1,
			Id:    v.F_entities_id,
			Count: v.F_count,
		}
	}

	return s.Response(&LeaderboardResponse{
		BaseResponse:    protocol.NewBaseResponse(errcode.SUCCESS, nil),
		LeaderboardInfo: leaderboardInfo,
	})
}
