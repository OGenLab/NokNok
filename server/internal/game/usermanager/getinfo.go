package usermanager

import (
	"github.com/lonng/nano/session"
	"lightspeed.2dao3.com/nokserver/common/boost"
	"lightspeed.2dao3.com/nokserver/common/entities"
	"lightspeed.2dao3.com/nokserver/common/task"
	"lightspeed.2dao3.com/nokserver/config"
	"lightspeed.2dao3.com/nokserver/consts/errcode"
	"lightspeed.2dao3.com/nokserver/internal/protocol"
	"lightspeed.2dao3.com/nokserver/model"
	"lightspeed.2dao3.com/nokserver/model/t_player"
	"lightspeed.2dao3.com/nokserver/model/t_user"
	"lightspeed.2dao3.com/nokserver/model/t_user_task"
)

type GetInfoRequest struct {
}

type GetInfoResponse struct {
	protocol.BaseResponse
	AccountInfo *protocol.AccountInfo `json:"accountInfo"`
	DrawInfo    []*protocol.DrawInfo  `json:"drawInfo"`
	NokInfo     *protocol.NokInfo     `json:"nokInfo"`
	TaskInfos   []*protocol.TaskInfo  `json:"taskInfo"`
	MetaInfo    *protocol.TaskInfo    `json:"metaInfo"`
}

func (u *UserManager) GetInfo(s *session.Session, req *GetInfoRequest) error {
	var (
		m          = model.NewModelDefault()
		uid        = s.UID()
		allTaskSet = task.TaskSet.AllSet()
	)
	defer m.Close()

	if uid <= 0 {
		return s.Response(GetInfoResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.NOT_LOGIN, nil),
		})
	}

	// 1. 查询t_user表
	user, err := t_user.SelectUser(m.DB, uid)
	if err != nil {
		return s.Response(&GetInfoResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.DATABASE_ERROR, err),
		})
	}

	// 2. 查询t_player表
	player, err := t_player.SelectPlayer(m.DB, uid)
	if err != nil {
		return s.Response(&GetInfoResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.DATABASE_ERROR, err),
		})
	}

	// 3. 查询用户完成任务表, 并筛选出抽奖详情
	var (
		drawInfos []*protocol.DrawInfo
		taskInfos []*protocol.TaskInfo
	)

	tasks, err := t_user_task.SelectAllTask(m.DB, uid)
	if err != nil {
		return s.Response(&GetInfoResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.DATABASE_ERROR, err),
		})
	}

	for _, v := range tasks {
		switch v.F_task_type {
		case task.DRAW_TASK_TYPE:
			drawInfos = append(drawInfos, &protocol.DrawInfo{
				PrizeType:    allTaskSet[v.F_task_type][v.F_task_id].F_type,
				Count:        v.F_schedule_count,
				Threshold:    allTaskSet[v.F_task_type][v.F_task_id].F_threshold,
				LastDrawTime: v.F_modify_time.Unix(),
			})
		default:
			taskInfos = append(taskInfos, &protocol.TaskInfo{
				Id:       v.F_task_id,
				Type:     v.F_task_type,
				Schedule: v.F_schedule_count,
				Status:   v.F_task_status,
			})

		}
	}

	// 返回是否可以继续升级
	nextBoost := player.F_boost
	if _, ok := boost.BoostSet[player.F_boost+1]; ok {
		nextBoost = player.F_boost + 1
	}

	if len(drawInfos) == 0 {
		drawInfos = append(drawInfos, &protocol.DrawInfo{
			PrizeType: entities.HAMMER_TYPE,
			Threshold: task.TaskSet.AllSet()[task.DRAW_TASK_TYPE][task.DRAW_HAMMER_ID].F_threshold,
		})
	}

	var response = &GetInfoResponse{
		BaseResponse: protocol.NewBaseResponse(errcode.SUCCESS, nil),
		AccountInfo: &protocol.AccountInfo{
			IsPremium:           player.F_is_premium,
			Boost:               player.F_boost,
			NextBoost:           nextBoost,
			Coins:               player.F_coins,
			InvitedCount:        user.F_invited_count,
			InvitedPremiumCount: user.F_invited_premium_count,
		},
		DrawInfo: drawInfos,
		NokInfo: &protocol.NokInfo{
			Stamina:     config.Instance().NokConfig.MaxStamina - player.F_stamina,
			HitCount:    player.F_hit_count,
			LastHitTime: player.F_last_hit_time.Unix(),
		},
		TaskInfos: taskInfos,
	}

	return s.Response(response)
}
