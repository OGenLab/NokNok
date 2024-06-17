package t_user_hammer

import (
	"time"

	"gorm.io/gorm"

	"qoobing.com/gomod/log"
)

type T_user_hammer struct {
	F_uid             int64     `gorm:"column:F_uid"`
	F_hammer_id       uint64    `gorm:"column:F_hammer_id"`
	F_equipment_count int       `gorm:"column:F_equipment_count"`
	F_count           int       `gorm:"column:F_count"`
	F_create_time     time.Time `gorm:"column:F_create_time;autoCreateTime"` // 首次创建时间
	F_modify_time     time.Time `gorm:"column:F_modify_time;autoUpdateTime"` // 最后更新时间
}

func (t *T_user_hammer) TableName() string {
	return "t_user_hammer"
}

func NewUserHammer(uid int64, hammerId uint64, equipmentCount int) *T_user_hammer {
	return &T_user_hammer{
		F_uid:             uid,
		F_hammer_id:       hammerId,
		F_equipment_count: equipmentCount,
		F_count:           1,
	}
}

func Create(db *gorm.DB, tx *T_user_hammer) (err error) {
	if err := db.Table("t_user_hammer").Create(tx).Error; err != nil {
		log.Errorf("t_user_hammer create err, tx: %+v, err: %v", tx, err)
		return err
	}
	return nil
}

func CreateUserHammer(db *gorm.DB, uid int64, hammerId uint64, equipmentCount int) (*T_user_hammer, error) {
	userHammer := NewUserHammer(uid, hammerId, equipmentCount)
	if err := Create(db, userHammer); err != nil {
		return nil, err
	}
	return userHammer, nil
}

func Select(db *gorm.DB, fields map[string]interface{}) ([]*T_user_hammer, error) {
	var (
		txs = []*T_user_hammer{}
	)

	result := db.Table("t_user_hammer").Where(fields).Find(&txs)
	if err := result.Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("t_user_hammer select, err[%s], fields[%v]", err, fields)
		return nil, err
	}
	return txs, nil
}

func Update(db *gorm.DB, sFields map[string]interface{}, uFields map[string]interface{}) error {
	result := db.Table("t_user_hammer").Where(sFields).Updates(uFields)
	if err := result.Error; err != nil {
		log.Errorf("t_user_hammer SelectAndUpdate err:%v, sFields:%+v, uFields:%+v", err, sFields, uFields)
	}
	return nil
}

func UpdateCount(db *gorm.DB, uid int64, hammerId uint64, count int) error {
	sFields := map[string]interface{}{
		"F_uid":       uid,
		"F_hammer_id": hammerId,
	}
	uFields := map[string]interface{}{
		"F_count": count,
	}

	return Update(db, sFields, uFields)
}

func NotEquipHammer(db *gorm.DB, uid int64) error {
	result := db.Table("t_user_hammer").Where(`"F_equipment_count" != ?`, 0).Update("F_equipment_count", 0)
	return result.Error
}

func EquipHammer(db *gorm.DB, uid int64, hammerId uint64) error {
	sFields := map[string]interface{}{
		"F_uid":       uid,
		"F_hammer_id": hammerId,
	}
	uFields := map[string]interface{}{
		"F_equipment_count": 1,
	}
	return Update(db, sFields, uFields)
}

func SelectHammer(db *gorm.DB, uid int64) ([]*T_user_hammer, error) {
	sFields := map[string]interface{}{
		"F_uid": uid,
	}
	return Select(db, sFields)
}

func SelectHammerById(db *gorm.DB, uid int64, hammerId uint64) (*T_user_hammer, error) {
	sFields := map[string]interface{}{
		"F_uid":       uid,
		"F_hammer_id": hammerId,
	}

	res, err := Select(db, sFields)
	if err != nil {
		return nil, err
	} else if len(res) == 0 {
		return nil, nil
	}

	return res[0], nil
}

func CreateOrUpdateCount(db *gorm.DB, uid int64, hammerId uint64, count int) (*T_user_hammer, error) {
	data, err := SelectHammerById(db, uid, hammerId)
	if err != nil {
		return nil, err
	}

	if data == nil {
		data, err = CreateUserHammer(db, uid, hammerId, 0)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	data.F_count += count
	if err := UpdateCount(db, uid, hammerId, data.F_count); err != nil {
		return nil, err
	}

	return data, nil
}
