package t_boost

import (
	"time"

	"gorm.io/gorm"

	"qoobing.com/gomod/log"
)

type T_boost struct {
	F_level          int       `gorm:"column:F_level"`
	F_needed_coins   uint64    `gorm:"column:F_needed_coins"`
	F_consume_samina int       `gorm:"column:F_consume_samina"`
	F_coins_rate     int       `gorm:"column:F_coins_rate"`
	F_create_time    time.Time `gorm:"column:F_create_time;autoCreateTime"` // 首次创建时间
	F_modify_time    time.Time `gorm:"column:F_modify_time;autoUpdateTime"` // 最后更新时间
}

func (t *T_boost) TableName() string {
	return "t_boost"
}

func Create(db *gorm.DB, tx *T_boost) (err error) {
	if err := db.Table("t_boost").Create(tx).Error; err != nil {
		log.Errorf("t_boost create err, tx: %+v, err: %v", tx, err)
		return err
	}
	return nil
}

func Select(db *gorm.DB, fields map[string]interface{}) ([]*T_boost, error) {
	var (
		txs = []*T_boost{}
	)

	result := db.Table("t_boost").Where(fields).Find(&txs)
	if err := result.Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("t_boost select, err[%s], fields[%v]", err, fields)
		return nil, err
	}
	return txs, nil
}

func Update(db *gorm.DB, sFields map[string]interface{}, uFields map[string]interface{}) error {
	result := db.Table("t_boost").Where(sFields).Updates(uFields)
	if err := result.Error; err != nil {
		log.Errorf("t_boost SelectAndUpdate err:%v, sFields:%+v, uFields:%+v", err, sFields, uFields)
	}
	return nil
}
