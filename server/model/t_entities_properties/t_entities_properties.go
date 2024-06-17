package t_entities_properties

import (
	"time"

	"gorm.io/gorm"

	"qoobing.com/gomod/log"
)

type T_entities_properties struct {
	F_entities_id uint64    `gorm:"column:F_entities_id"`
	F_season      int       `gorm:"column:F_season"`
	F_type        int       `gorm:"column:F_type"`
	F_type_ext    string    `gorm:"column:F_type_ext"`
	F_name        string    `gorm:"column:F_name"`
	F_image_url   string    `gorm:"column:F_image_url"`
	F_probability int       `gorm:"column:F_probability"`
	F_create_time time.Time `gorm:"column:F_create_time;autoCreateTime"` // 首次创建时间
	F_modify_time time.Time `gorm:"column:F_modify_time;autoUpdateTime"` // 最后更新时间
}

func (t *T_entities_properties) TableName() string {
	return "t_entities_properties"
}

func Create(db *gorm.DB, tx *T_entities_properties) (err error) {
	if err := db.Table("t_entities_properties").Create(tx).Error; err != nil {
		log.Errorf("t_entities_properties create err, tx: %+v, err: %v", tx, err)
		return err
	}
	return nil
}

func Select(db *gorm.DB, fields map[string]interface{}) ([]*T_entities_properties, error) {
	var (
		txs = []*T_entities_properties{}
	)

	result := db.Table("t_entities_properties").Where(fields).Find(&txs)
	if err := result.Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("t_entities_properties select, err[%s], fields[%v]", err, fields)
		return nil, err
	}
	return txs, nil
}

func Update(db *gorm.DB, sFields map[string]interface{}, uFields map[string]interface{}) error {
	result := db.Table("t_entities_properties").Where(sFields).Updates(uFields)
	if err := result.Error; err != nil {
		log.Errorf("t_entities_properties SelectAndUpdate err:%v, sFields:%+v, uFields:%+v", err, sFields, uFields)
	}
	return nil
}
