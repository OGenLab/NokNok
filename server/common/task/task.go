package task

import (
	"sync"
	"time"

	"lightspeed.2dao3.com/nokserver/config"
	"lightspeed.2dao3.com/nokserver/model"
	"lightspeed.2dao3.com/nokserver/model/t_task"
)

const (
	DAILY_TASK_TYPE  = 1
	BASIC_TASK_TYPE  = 2
	DRAW_TASK_TYPE   = 3
	INVITE_TASK_TYPE = 4
)

type task struct {
	allSet    map[int]map[uint64]*t_task.T_task
	ladInvSet map[uint64]*t_task.T_task
	sync.RWMutex
}

var TaskSet = &task{}

func init() {
	refreshTask()
	go refresh()
}

func (t *task) AllSet() map[int]map[uint64]*t_task.T_task {
	t.RLock()
	defer t.RUnlock()
	return t.allSet
}

func (t *task) LadInvSet() map[uint64]*t_task.T_task {
	t.RLock()
	defer t.RUnlock()
	return t.ladInvSet
}

func refreshTask() {
	var (
		m         = model.NewModel()
		allSet    = make(map[int]map[uint64]*t_task.T_task)
		ladInvSet = make(map[uint64]*t_task.T_task)
	)
	defer m.Close()

	tasks, err := t_task.SelectTask(m.DB, 0)
	if err != nil {
		panic(err)
	}

	for i := range tasks {
		if _, ok := allSet[tasks[i].F_type]; !ok {
			allSet[tasks[i].F_type] = make(map[uint64]*t_task.T_task)
		}
		allSet[tasks[i].F_type][tasks[i].F_task_id] = tasks[i]
		if tasks[i].F_type == INVITE_TASK_TYPE {
			ladInvSet[tasks[i].F_task_id] = tasks[i]
		}
	}

	TaskSet.Lock()
	defer TaskSet.Unlock()

	TaskSet.allSet = allSet
	TaskSet.ladInvSet = ladInvSet
}

func refresh() {
	for {
		time.Sleep(time.Duration(config.Instance().NokConfig.RefreshInterval) * time.Minute)
		refreshTask()
	}
}
