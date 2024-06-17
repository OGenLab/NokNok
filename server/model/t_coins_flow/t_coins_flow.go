package t_coins_flow

import (
	"time"

	"gorm.io/gorm"

	"qoobing.com/gomod/log"
)

type T_coins_flow struct {
	F_uid         int64     `gorm:"column:F_uid"`
	F_coins       uint64    `gorm:"column:F_coins"`
	F_task_type   int       `gorm:"column:F_task_type"`
	F_task_id     int       `gorm:"column:F_task_id"`
	F_create_time time.Time `gorm:"column:F_create_time;autoCreateTime"` // 首次创建时间
	F_modify_time time.Time `gorm:"column:F_modify_time;autoUpdateTime"` // 最后更新时间
}

func (t *T_coins_flow) TableName() string {
	return "t_coins_flow"
}

func Create(db *gorm.DB, tx *T_coins_flow) (err error) {
	if err := db.Table("t_coins_flow").Create(tx).Error; err != nil {
		log.Errorf("t_coins_flow create err, tx: %+v, err: %v", tx, err)
		return err
	}
	return nil
}

func Select(db *gorm.DB, fields map[string]interface{}) ([]*T_coins_flow, error) {
	var (
		txs = []*T_coins_flow{}
	)

	result := db.Table("t_coins_flow").Where(fields).Find(&txs)
	if err := result.Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("t_coins_flow select, err[%s], fields[%v]", err, fields)
		return nil, err
	}
	return txs, nil
}

func Update(db *gorm.DB, sFields map[string]interface{}, uFields map[string]interface{}) error {
	result := db.Table("t_coins_flow").Where(sFields).Updates(uFields)
	if err := result.Error; err != nil {
		log.Errorf("t_coins_flow SelectAndUpdate err:%v, sFields:%+v, uFields:%+v", err, sFields, uFields)
	}
	return nil
}

func CreateCoinFlow(db *gorm.DB, uid int64, taskType int, taskId int, coins uint64) error {
	coinFlow := &T_coins_flow{
		F_uid:       uid,
		F_task_type: taskType,
		F_task_id:   taskId,
		F_coins:     coins,
	}
	return Create(db, coinFlow)
}
