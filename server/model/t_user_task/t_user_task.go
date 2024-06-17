package t_user_task

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"lightspeed.2dao3.com/nokserver/common/task"
	"lightspeed.2dao3.com/nokserver/config"

	"qoobing.com/gomod/log"
)

const (
	NOT_COMPLETED_STATUS      = 1
	NOT_RECEIVE_REWARD_STATUS = 2
	RECEIVED_REWARD_STATUS    = 3
)

var UnCompleteTaskErr = errors.New("uncomplete task")

type T_user_task struct {
	F_uid            int64     `gorm:"column:F_uid"`
	F_task_type      int       `gorm:"column:F_task_type"`
	F_task_id        uint64    `gorm:"column:F_task_id"`
	F_task_status    int       `gorm:"column:F_task_status"`
	F_schedule_count int       `gorm:"column:F_schedule_count"`
	F_create_time    time.Time `gorm:"column:F_create_time;autoCreateTime"` // 首次创建时间
	F_modify_time    time.Time `gorm:"column:F_modify_time;autoUpdateTime"` // 最后更新时间
}

func NewUserTask(uid int64, taskType int, taskId uint64, scheduleCount int) *T_user_task {
	status := NOT_COMPLETED_STATUS
	if task.TaskSet.AllSet()[taskType][taskId].F_threshold == scheduleCount {
		status = NOT_RECEIVE_REWARD_STATUS
	}
	return &T_user_task{
		F_uid:            uid,
		F_task_id:        taskId,
		F_task_type:      taskType,
		F_task_status:    status,
		F_schedule_count: scheduleCount,
	}
}

func (t *T_user_task) TableName() string {
	return "t_user_task"
}

func Create(db *gorm.DB, tx *T_user_task) (err error) {
	if err := db.Table("t_user_task").Create(tx).Error; err != nil {
		log.Errorf("t_user_task create err, tx: %+v, err: %v", tx, err)
		return err
	}
	return nil
}

func CreateUserTask(db *gorm.DB, uid int64, taskType int, taskId uint64, scheduleCount int) (*T_user_task, error) {
	userTask := NewUserTask(uid, taskType, taskId, scheduleCount)
	if err := Create(db, userTask); err != nil {
		return nil, err
	}
	return userTask, nil
}

func Select(db *gorm.DB, fields map[string]interface{}) ([]*T_user_task, error) {
	var (
		txs = []*T_user_task{}
	)

	result := db.Table("t_user_task").Where(fields).Find(&txs)
	if err := result.Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("t_user_task select, err[%s], fields[%v]", err, fields)
		return nil, err
	}
	return txs, nil
}

func Update(db *gorm.DB, sFields map[string]interface{}, uFields map[string]interface{}) error {
	result := db.Table("t_user_task").Where(sFields).Updates(uFields)
	if err := result.Error; err != nil {
		log.Errorf("t_user_task SelectAndUpdate err:%v, sFields:%+v, uFields:%+v", err, sFields, uFields)
	}
	return nil
}

func UpdateStatus(db *gorm.DB, uid int64, taskType int, taskId uint64, originStatus int, newStatus int) error {
	sFields := map[string]interface{}{
		"F_uid":       uid,
		"F_task_type": taskType,
		"F_task_id":   taskId,
	}
	uFields := map[string]interface{}{
		"F_task_status": newStatus,
	}
	return Update(db, sFields, uFields)
}

func UpdateScheduleCount(db *gorm.DB, uid int64, taskType int, taskId uint64, scheduleCount int) error {
	sFields := map[string]interface{}{
		"F_uid":     uid,
		"F_task_id": taskId,
	}

	uFields := map[string]interface{}{
		"F_schedule_count": scheduleCount,
	}

	if task.TaskSet.AllSet()[taskType][taskId].F_threshold == scheduleCount {
		uFields["F_task_status"] = NOT_RECEIVE_REWARD_STATUS
	}

	return Update(db, sFields, uFields)
}

func SelectAllTask(db *gorm.DB, uid int64) ([]*T_user_task, error) {
	var tasks []*T_user_task

	// 获取当前时间
	currentTime := time.Now()

	// 获取俄罗斯时区
	russianTimeZone, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return nil, err
	}

	// 获取俄罗斯早上 8 点的时间
	russianEightAM := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), config.Instance().NokConfig.RefreshTime, 0, 0, 0, russianTimeZone)

	// 查询符合条件的记录
	err = db.Where(`"F_uid" = ?  AND "F_create_time" > ?`, uid, russianEightAM).Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func SelectTask(db *gorm.DB, uid int64, taskType int, taskId uint64) (*T_user_task, error) {
	sFields := map[string]interface{}{
		"F_uid":       uid,
		"F_task_type": taskType,
		"F_task_id":   taskId,
	}

	userTasks, err := Select(db, sFields)
	if err != nil {
		return nil, err
	} else if len(userTasks) == 0 {
		return nil, UnCompleteTaskErr
	}

	return userTasks[0], nil
}

func SelectDailyTask(db *gorm.DB, uid int64, taskType int, taskId uint64) (*T_user_task, error) {
	var tasks []*T_user_task

	// 获取当前时间
	currentTime := time.Now()

	// 获取俄罗斯时区
	russianTimeZone, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return nil, err
	}

	// 获取俄罗斯早上 8 点的时间
	russianEightAM := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), config.Instance().NokConfig.RefreshTime, 0, 0, 0, russianTimeZone)

	// 查询符合条件的记录
	err = db.Where(`"F_uid" = ? AND "F_task_type" = ? AND "F_task_id" = ? AND "F_create_time" > ?`, uid, taskType, taskId, russianEightAM).Find(&tasks).Error
	if err != nil {
		return nil, err
	} else if len(tasks) == 0 {
		return nil, nil
	}

	return tasks[0], nil
}

func SelectBaseTask(db *gorm.DB, uid int64, taskType int, taskId uint64) (*T_user_task, error) {
	var tasks []*T_user_task
	// 查询符合条件的记录
	res := db.Where(`"F_uid" = ? AND "F_task_type" = ? AND "F_task_id" = ?`, uid, taskType, taskId).Find(&tasks)
	if err := res.Error; err != nil {
		return nil, err
	} else if len(tasks) == 0 {
		return nil, nil
	}

	return tasks[0], nil
}

// 查询用户指定taskId的userTask记录
// 1. 如果记录不存在则创建记录
// 2. 如果记录存在则更新完成次数
// 3. 返回的userTask对象中ScheduleCount是最新的
func CreateOrUpdateScheduleCount(db *gorm.DB, uid int64, taskType int, taskId uint64) (userTask *T_user_task, err error) {
	if taskType == task.DAILY_TASK_TYPE || taskType == task.DRAW_TASK_TYPE {
		userTask, err = SelectDailyTask(db, uid, taskType, taskId)
	} else if taskType == task.BASIC_TASK_TYPE || taskType == task.INVITE_TASK_TYPE {
		userTask, err = SelectBaseTask(db, uid, taskType, taskId)
	}
	if err != nil {
		return nil, err
	}

	if userTask == nil {
		userTask, err = CreateUserTask(db, uid, taskType, taskId, 1)
		if err != nil {
			return nil, err
		}
		return userTask, nil
	}

	if userTask.F_schedule_count+1 > task.TaskSet.AllSet()[taskType][taskId].F_threshold {
		return userTask, nil
	}

	userTask.F_schedule_count++
	if err := UpdateScheduleCount(db, uid, taskType, taskId, userTask.F_schedule_count); err != nil {
		return nil, err
	}

	return userTask, nil
}
