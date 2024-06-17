package t_task

import (
	"time"

	"gorm.io/gorm"
	"lightspeed.2dao3.com/nokserver/config"

	"qoobing.com/gomod/log"
)

type T_task struct {
	F_task_id     uint64    `gorm:"column:F_task_id"`
	F_season      int       `gorm:"column:F_season"`
	F_type        int       `gorm:"column:F_type"`
	F_type_ext    string    `gorm:"column:F_type_ext"`
	F_threshold   int       `gorm:"column:F_threshold"`
	F_reward      uint64    `gorm:"column:F_reward"`
	F_period      uint64    `gorm:"column:F_period"`
	F_start_time  time.Time `gorm:"column:F_start_time"`
	F_end_time    time.Time `gorm:"column:F_end_time"`
	F_create_time time.Time `gorm:"column:F_create_time;autoCreateTime"` // 首次创建时间
	F_modify_time time.Time `gorm:"column:F_modify_time;autoUpdateTime"` // 最后更新时间
}

func (t *T_task) TableName() string {
	return "t_task"
}

func Create(db *gorm.DB, tx *T_task) (err error) {
	if err := db.Table("t_task").Create(tx).Error; err != nil {
		log.Errorf("t_task create err, tx: %+v, err: %v", tx, err)
		return err
	}
	return nil
}

func Select(db *gorm.DB, fields map[string]interface{}) ([]*T_task, error) {
	var (
		txs = []*T_task{}
	)

	result := db.Table("t_task").Where(fields).Find(&txs)
	if err := result.Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("t_task select, err[%s], fields[%v]", err, fields)
		return nil, err
	}
	return txs, nil
}

func Update(db *gorm.DB, sFields map[string]interface{}, uFields map[string]interface{}) error {
	result := db.Table("t_task").Where(sFields).Updates(uFields)
	if err := result.Error; err != nil {
		log.Errorf("t_task SelectAndUpdate err:%v, sFields:%+v, uFields:%+v", err, sFields, uFields)
	}
	return nil
}

func SelectTask(db *gorm.DB, typ int) ([]*T_task, error) {
	var tasks []*T_task

	// 获取当前时间
	currentTime := time.Now()
	// 获取俄罗斯时区
	russianTimeZone, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return nil, err
	}

	// 获取俄罗斯早上 8 点的时间
	russianEightAM := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), config.Instance().NokConfig.RefreshTime, 0, 0, 0, russianTimeZone)

	db = db.Where(`"F_start_time" < ? AND "F_end_time" > ?`, russianEightAM, russianEightAM)
	if typ != 0 {
		db = db.Where(`"F_type" = ?`, typ)
	}

	res := db.Find(&tasks)
	if err := res.Error; err != nil {
		return nil, err
	}
	return tasks, nil
}
