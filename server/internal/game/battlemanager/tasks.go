package battlemanager

import (
	"github.com/lonng/nano/session"
	"lightspeed.2dao3.com/nokserver/consts/errcode"
	"lightspeed.2dao3.com/nokserver/internal/protocol"
	"lightspeed.2dao3.com/nokserver/model"
	"lightspeed.2dao3.com/nokserver/model/t_task"
)

type TasksRequest struct {
	Type int `json:"type"`
}

type TasksResponse struct {
	protocol.BaseResponse
	TaskInfos []*protocol.TaskInfo `json:"taskInfos"`
}

func (u *BattleManager) Tasks(s *session.Session, req *TasksRequest) error {
	var (
		m   = model.NewModelDefault()
		uid = s.UID()
	)
	defer m.Close()

	if uid <= 0 {
		return s.Response(&TasksResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.NOT_LOGIN, nil),
		})
	}

	tasks, err := t_task.SelectTask(m.DB, req.Type)
	if err != nil {
		return s.Response(&TasksResponse{
			BaseResponse: protocol.NewBaseResponse(errcode.DATABASE_ERROR, err),
		})
	}

	taskInfos := make([]*protocol.TaskInfo, len(tasks))
	for i, v := range tasks {
		taskInfos[i] = &protocol.TaskInfo{
			Id:        v.F_task_id,
			Season:    v.F_season,
			Type:      v.F_type,
			Threshold: v.F_threshold,
			Reward:    v.F_reward,
			StartTime: v.F_start_time.Unix(),
			EndTime:   v.F_end_time.Unix(),
		}
	}

	return s.Response(&TasksResponse{
		BaseResponse: protocol.NewBaseResponse(errcode.SUCCESS, nil),
		TaskInfos:    taskInfos,
	})
}
