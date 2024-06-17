package battlemanager

import (
	"math/rand"
	"time"

	"github.com/lonng/nano/session"
	"lightspeed.2dao3.com/nokserver/common/entities"
	"lightspeed.2dao3.com/nokserver/common/task"
	"lightspeed.2dao3.com/nokserver/consts/errcode"
	"lightspeed.2dao3.com/nokserver/internal/protocol"
	"lightspeed.2dao3.com/nokserver/model"
	"lightspeed.2dao3.com/nokserver/model/t_user_hammer"
	"lightspeed.2dao3.com/nokserver/model/t_user_task"
	"qoobing.com/gomod/log"
)

type PrizeDrawRequest struct {
	PrizeType int `json:"prizeType"`
}

type PrizeDrawResponse struct {
	protocol.BaseResponse
	DrawInfo *protocol.DrawInfo `json:"drawInfo"`
}

func (u *BattleManager) PrizeDraw(s *session.Session, req *PrizeDrawRequest) error {
	var (
		m          = model.NewModelDefault()
		uid        = s.UID()
		allTaskSet = task.TaskSet.AllSet()
	)
	defer m.Close()

	if uid <= 0 {
		return s.Response(&PrizeDrawResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.NOT_LOGIN, nil),
		})
	}
	if req.PrizeType != task.DRAW_HAMMER_ID {
		return s.Response(&PrizeDrawResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.PARAMETER_INVALID, nil),
		})
	}

	m.Begin()

	// 1. 查询是否还有抽奖机会
	userTask, err := t_user_task.CreateOrUpdateScheduleCount(m.DB, uid, task.DRAW_TASK_TYPE, task.DRAW_HAMMER_ID)
	if err != nil {
		return s.Response(&PrizeDrawResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.DATABASE_ERROR, err),
		})
	}
	if userTask.F_schedule_count >= allTaskSet[task.DRAW_TASK_TYPE][task.DRAW_HAMMER_ID].F_threshold {
		return s.Response(&PrizeDrawResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.NOT_DRAW_COUNT, err),
		})
	}

	// 2. 随机选择锤子id
	hammerId := randomSelectNumber()
	// 2.1 查询用户是否持有该类型的锤子, 有递增, 没有创建
	_, err = t_user_hammer.CreateOrUpdateCount(m.DB, uid, hammerId, 1)
	if err != nil {
		return s.Response(&PrizeDrawResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.DATABASE_ERROR, err),
		})
	}

	if err := m.Commit(); err != nil {
		return s.Response(&PrizeDrawResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.DATABASE_ERROR, err),
		})
	}

	return s.Response(&PrizeDrawResponse{
		BaseResponse: protocol.NewBaseResponse(errcode.SUCCESS, nil),
		DrawInfo: &protocol.DrawInfo{
			PrizeType:    req.PrizeType,
			PrizeId:      hammerId,
			Count:        userTask.F_schedule_count,
			Threshold:    allTaskSet[task.DRAW_TASK_TYPE][task.DRAW_HAMMER_ID].F_threshold,
			LastDrawTime: time.Now().Unix(),
		},
	})
}

func randomSelectNumber() uint64 {
	var (
		entSet        = entities.EntSet
		ids           = entSet.HammerIds()
		probabilities = entSet.HammerProbabilities()
	)
	// 确保 ids 和 probabilities 的长度相同
	if len(ids) != len(probabilities) {
		log.Debugf("RandomSelectNumber ids != probabilities, ids: %+v, probabilities: %+v", ids, probabilities)
		return ids[0]
	}

	// 生成随机数种子
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 计算总概率
	totalProbability := 0
	for _, p := range probabilities {
		totalProbability += p
	}

	// 生成随机概率值
	randomValue := r.Intn(totalProbability)

	// 根据随机概率值选择数字
	cumulativeProbability := 0
	for i, p := range probabilities {
		cumulativeProbability += p
		if randomValue <= cumulativeProbability {
			return ids[i]
		}
	}

	return ids[0]
}
