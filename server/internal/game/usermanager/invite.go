package usermanager

import (
	"github.com/lonng/nano/session"
	"lightspeed.2dao3.com/nokserver/consts/errcode"
	"lightspeed.2dao3.com/nokserver/internal/protocol"
	"lightspeed.2dao3.com/nokserver/model"
	"lightspeed.2dao3.com/nokserver/model/t_user"
)

type InviteRequest struct {
}

type InviteResponse struct {
	protocol.BaseResponse
	ReferralCode string `json:"referralCode"`
}

func (u *UserManager) Invite(s *session.Session, req *InviteRequest) error {
	var (
		m   = model.NewModelDefault()
		uid = s.UID()
	)
	defer m.Close()

	if uid <= 0 {
		return s.Response(&InviteResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.NOT_LOGIN, nil),
		})
	}

	response := &InviteResponse{
		BaseResponse: protocol.NewBaseResponse(errcode.SUCCESS, nil),
		ReferralCode: t_user.ID2ReferralCode(uid),
	}

	return s.Response(response)
}
