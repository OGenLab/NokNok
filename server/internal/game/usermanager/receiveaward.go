package usermanager

import (
	"errors"

	"github.com/lonng/nano/session"
	"lightspeed.2dao3.com/nokserver/common/task"
	"lightspeed.2dao3.com/nokserver/consts/errcode"
	"lightspeed.2dao3.com/nokserver/internal/protocol"
	"lightspeed.2dao3.com/nokserver/model"
	"lightspeed.2dao3.com/nokserver/model/t_coins_flow"
	"lightspeed.2dao3.com/nokserver/model/t_player"
	"lightspeed.2dao3.com/nokserver/model/t_user_task"
)

type ReceiveAwardRequest struct {
	TaskType int    `json:"taskType"`
	TaskId   uint64 `json:"taskId"`
}

type ReceiveAwardResponse struct {
	protocol.BaseResponse
	Coins       uint64                `json:"coins"` // coins奖励
	AccountInfo *protocol.AccountInfo `json:"accountInfo"`
}

func (u *UserManager) ReceiveAward(s *session.Session, req *ReceiveAwardRequest) error {
	var (
		m   = model.NewModelDefault()
		uid = s.UID()
	)
	m.Begin()
	defer m.Close()
	if uid <= 0 {
		return s.Response(&ReceiveAwardResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.NOT_LOGIN, nil),
		})
	}

	// 1. 查询用户任务完成详情
	userTask, err := t_user_task.SelectTask(m.DB, uid, req.TaskType, req.TaskId)
	if err != nil && !errors.Is(err, t_user_task.UnCompleteTaskErr) {
		return s.Response(&ReceiveAwardResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.DATABASE_ERROR, err),
		})
	} else if err != nil {
		return s.Response(&ReceiveAwardResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.UNCOMPLTE_TASK_ERR, nil),
		})
	}

	if userTask.F_task_status == t_user_task.NOT_COMPLETED_STATUS {
		return s.Response(&ReceiveAwardResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.UNCOMPLTE_TASK_ERR, nil),
		})
	} else if userTask.F_task_status == t_user_task.RECEIVED_REWARD_STATUS {
		return s.Response(&ReceiveAwardResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.RECEIVED_REWARD_ERR, nil),
		})
	}

	// 2. 发放任务奖励
	taskSet := task.TaskSet.AllSet()
	if v, ok := taskSet[req.TaskType]; !ok || v == nil {
		return s.Response(&ReceiveAwardResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.PARAMETER_INVALID, nil),
		})
	}

	coins := taskSet[req.TaskType][req.TaskId].F_reward
	if err := t_player.AddCoins(m.DB, uid, coins); err != nil {
		return s.Response(&ReceiveAwardResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.DATABASE_ERROR, err),
		})
	}

	// 3. 更新任务状态
	if err := t_user_task.UpdateStatus(m.DB, uid, req.TaskType, req.TaskId, userTask.F_task_status, t_user_task.RECEIVED_REWARD_STATUS); err != nil {
		return s.Response(&ReceiveAwardResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.DATABASE_ERROR, err),
		})
	}

	player, err := t_player.SelectPlayer(m.DB, uid)
	if err != nil {
		return s.Response(&ReceiveAwardResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.DATABASE_ERROR, err),
		})
	}

	// 4. 提交流水表
	if err := t_coins_flow.CreateCoinFlow(m.DB, uid, req.TaskType, int(req.TaskId), coins); err != nil {
		return s.Response(&ReceiveAwardResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.DATABASE_ERROR, err),
		})
	}

	// 5. 提交事务
	if err := m.Commit(); err != nil {
		return s.Response(&ReceiveAwardResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.DATABASE_ERROR, err),
		})
	}

	return s.Response(&ReceiveAwardResponse{
		BaseResponse: protocol.NewBaseResponse(errcode.SUCCESS, nil),
		Coins:        coins,
		AccountInfo: &protocol.AccountInfo{
			Coins: player.F_coins,
		},
	})
}
