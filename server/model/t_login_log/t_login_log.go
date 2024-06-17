package t_login_log

import (
	"time"

	"gorm.io/gorm"

	"qoobing.com/gomod/log"
)

type T_login_log struct {
	F_uid         int64     `gorm:"column:F_uid"`
	F_is_premium  bool      `gorm:"column:F_is_premium"`
	F_create_time time.Time `gorm:"column:F_create_time;autoCreateTime"` // 首次创建时间
}

func NewLoginLog(uid int64, isPremium bool) *T_login_log {
	return &T_login_log{
		F_uid:        uid,
		F_is_premium: isPremium,
	}
}

func (t *T_login_log) TableName() string {
	return "t_login_log"
}

func Create(db *gorm.DB, tx *T_login_log) (err error) {
	if err := db.Table("t_login_log").Create(tx).Error; err != nil {
		log.Errorf("t_login_log create err, tx: %+v, err: %v", tx, err)
		return err
	}
	return nil
}

func CreateLog(db *gorm.DB, uid int64, isPremium bool) (*T_login_log, error) {
	loginLog := NewLoginLog(uid, isPremium)
	if err := Create(db, loginLog); err != nil {
		return nil, err
	}
	return loginLog, nil
}

func Select(db *gorm.DB, fields map[string]interface{}) ([]*T_login_log, error) {
	var (
		txs = []*T_login_log{}
	)

	result := db.Table("t_login_log").Where(fields).Find(&txs)
	if err := result.Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("t_login_log select, err[%s], fields[%v]", err, fields)
		return nil, err
	}
	return txs, nil
}

func Update(db *gorm.DB, sFields map[string]interface{}, uFields map[string]interface{}) error {
	result := db.Table("t_login_log").Where(sFields).Updates(uFields)
	if err := result.Error; err != nil {
		log.Errorf("t_login_log SelectAndUpdate err:%v, sFields:%+v, uFields:%+v", err, sFields, uFields)
	}
	return nil
}
