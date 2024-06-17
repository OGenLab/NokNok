package usermanager

import (
	"errors"
	"fmt"
	"time"

	"github.com/lonng/nano/session"
	"lightspeed.2dao3.com/nokserver/common/task"
	"lightspeed.2dao3.com/nokserver/consts/errcode"
	"lightspeed.2dao3.com/nokserver/internal/protocol"
	"lightspeed.2dao3.com/nokserver/model"
	"lightspeed.2dao3.com/nokserver/model/t_login_log"
	"lightspeed.2dao3.com/nokserver/model/t_player"
	"lightspeed.2dao3.com/nokserver/model/t_user"
	"lightspeed.2dao3.com/nokserver/model/t_user_task"
	"qoobing.com/gomod/log"
)

type LoginToGameServerRequest struct {
	Token        string `json:"token"`
	ReferralCode string `json:"referralCode"`
	Channel      string `json:"channel"`
}

type LoginToGameServerResponse struct {
	protocol.BaseResponse
}

func (u *UserManager) Login(s *session.Session, req *LoginToGameServerRequest) error {
	var (
		m          = model.NewModelDefault()
		allTaskSet = task.TaskSet.AllSet()
		userPush   = false
	)
	m.Begin()
	defer m.Close()

	// 解析token
	webAppData, err := protocol.MockParseInitData(req.Token) // TODO: daewoo, mock
	if err != nil {
		return s.Response(&LoginToGameServerResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.PARAMETER_INVALID, err),
		})
	}

	var (
		tgUser = webAppData.User
		uid    = int64(tgUser.ID)
	)

	s.Bind(uid)

	log.Debugf("player: %d, webAppData: %+v, referralCode: %+v, channel: %v", uid, webAppData, req.ReferralCode, req.Channel)

	// 插入登录表
	if _, err = t_login_log.CreateLog(m.DB, uid, tgUser.IsPremium); err != nil {
		return s.Response(&LoginToGameServerResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.DATABASE_ERROR, err),
		})
	}

	var invitingUser = &t_user.T_user{}
	if len(req.ReferralCode) != 0 {
		invitingUser, err = t_user.SelectUserByRefferCode(m.DB, req.ReferralCode)
		if err != nil && !errors.Is(err, t_user.ErrAccountNotExist) {
			return s.Response(&LoginToGameServerResponse{
				BaseResponse: protocol.NewBaseResponse(errcode.DATABASE_ERROR, err),
			})
		}
	}

	// 查询t_user表用户是否存在, 不存在则新增
	user, userErr := t_user.SelectOrCreate(m.DB, uid, invitingUser.F_uid, req.Channel, tgUser)
	if userErr != nil && !errors.Is(userErr, t_user.ErrAccountNotExist) {
		return s.Response(&LoginToGameServerResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.PARAMETER_INVALID, err),
		})
	}

	// 查询t_player表用户是否存在, 不存在则新增
	_, playerErr := t_player.SelectOrCreate(m.DB, uid, tgUser)
	if playerErr != nil && !errors.Is(playerErr, t_player.ErrAccountNotExist) {
		return s.Response(&LoginToGameServerResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.DATABASE_ERROR, err),
		})
	}

	// 新用户处理
	if userErr != nil && playerErr != nil && invitingUser != nil {
		if err := hanlderInvite(m, user, invitingUser); err != nil {
			return s.Response(&LoginToGameServerResponse{
				BaseResponse: protocol.NewBaseResponse(errcode.DATABASE_ERROR, err),
			})
		}
	}

	// 查询日常任务-游戏登录 是否已完成, 未完成则插入用户完成任务表, 同时奖励金币
	userTask, err := t_user_task.CreateOrUpdateScheduleCount(m.DB, user.F_uid, task.DAILY_TASK_TYPE, task.LOGIN_GAME_ID)
	if err != nil {
		return s.Response(&LoginToGameServerResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.DATABASE_ERROR, err),
		})
	}
	if userTask.F_schedule_count == allTaskSet[task.DAILY_TASK_TYPE][task.LOGIN_GAME_ID].F_threshold {
		userPush = true
	}

	if err := m.Commit(); err != nil {
		return s.Response(&LoginToGameServerResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.DATABASE_ERROR, err),
		})
	}

	if userPush {
		s.Push("onTaskComplete", &protocol.OnTaskCompleteRequest{})
	}

	if p, ok := u.player(uid); !ok {
		log.Infof("players: %d not online, create new players", uid)
		p = newPlayer(s, uid)
		u.setPlayer(uid, p)
	} else {
		p.bindSession(s)
	}

	return s.Response(&LoginToGameServerResponse{
		BaseResponse: protocol.NewBaseResponse(errcode.SUCCESS, nil),
	})
}

func hanlderInvite(m *model.Model, invitedUser *t_user.T_user, invitingUser *t_user.T_user) error {
	var (
		allTaskSet       = task.TaskSet.AllSet()
		ladInvSet        = task.TaskSet.LadInvSet()
		invAddCoins      = allTaskSet[task.BASIC_TASK_TYPE][task.INVITE_COMMON_ID].F_reward // 邀请普通用户奖励
		invitingUserPush = false
	)

	if invitedUser.F_is_premium {
		invAddCoins = allTaskSet[task.BASIC_TASK_TYPE][task.INVITE_PREMIUM_ID].F_reward // 邀请premium用户奖励
	}

	// 1. 邀请人和被邀请人都可以获得邀请新用户奖励
	addInvitedUserCoins := invAddCoins
	addInvitingUserCoins := invAddCoins

	// 2. 如果邀请人的日常任务还没有完成, 则完成该任务
	if allTaskSet[task.DAILY_INVITE_ID] != nil && !time.Now().After(allTaskSet[task.DAILY_TASK_TYPE][task.DAILY_INVITE_ID].F_end_time) { // 校验日常任务-邀请1个新用户是否有效
		// 查询邀请人是否完成了该每日任务
		userTask, err := t_user_task.CreateOrUpdateScheduleCount(m.DB, invitingUser.F_uid, task.DAILY_TASK_TYPE, task.DAILY_INVITE_ID) // 邀请人的userTask详情
		if err != nil {
			return fmt.Errorf("failed to select or create and update daily task for daily invite: %w", err)
		}

		if userTask.F_schedule_count == allTaskSet[task.DAILY_TASK_TYPE][task.DAILY_INVITE_ID].F_threshold { // 完成校验日常任务-邀请1个新用户
			invitingUserPush = true
		}
	}

	// 3. 完成累计邀请人数任务
	for _, ladTask := range ladInvSet {
		if time.Now().After(ladTask.F_end_time) {
			continue
		}

		userTask, err := t_user_task.CreateOrUpdateScheduleCount(m.DB, invitingUser.F_uid, task.INVITE_TASK_TYPE, ladTask.F_task_id) // 邀请人的userTask详情
		if err != nil {
			return fmt.Errorf("failed to select or create and update daily task for lad invite: %w", err)
		}
		if userTask.F_schedule_count == ladTask.F_threshold {
			invitingUserPush = true
		}
	}

	// 4. 增加累计积分
	if err := t_player.AddCoins(m.DB, invitingUser.F_uid, addInvitingUserCoins); err != nil {
		return fmt.Errorf("failed to add coins to inviting user: %w", err)
	}
	if err := t_player.AddCoins(m.DB, invitingUser.F_uid, addInvitedUserCoins); err != nil {
		return fmt.Errorf("failed to add coins to invited user: %w", err)
	}

	// 5. 更新邀请人数
	if err := t_user.UpdateInvitedCount(m.DB, invitingUser.F_uid, invitedUser.F_is_premium); err != nil {
		return fmt.Errorf("failed to update invited count: %w", err)
	}

	// 6. 进行推送
	if invitingUserPush {
		player, ok := DefaultUserManger.player(invitingUser.F_uid)
		// 如果玩家在线
		if ok && player.session != nil {
			player.session.Push("onTaskComplete", &protocol.OnTaskCompleteRequest{})
		}
	}

	return nil
}
