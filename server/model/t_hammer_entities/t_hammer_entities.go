package t_hammer_entities

import (
	"time"

	"gorm.io/gorm"
	"lightspeed.2dao3.com/nokserver/config"
	"qoobing.com/gomod/log"
)

type T_hammer_entities struct {
	F_season      int       `gorm:"column:F_season"`
	F_entities_id uint64    `gorm:"column:F_entities_id"`
	F_count       uint64    `gorm:"column:F_count"`
	F_create_time time.Time `gorm:"column:F_create_time;autoCreateTime"` // 首次创建时间
	F_modify_time time.Time `gorm:"column:F_modify_time;autoUpdateTime"` // 最后更新时间
}

func (t *T_hammer_entities) TableName() string {
	return "t_hammer_entities"
}

func Create(db *gorm.DB, tx *T_hammer_entities) (err error) {
	if err := db.Table("t_hammer_entities").Create(tx).Error; err != nil {
		log.Errorf("t_entities_properties create err, tx: %+v, err: %v", tx, err)
		return err
	}
	return nil
}

func CreateHammerEntities(db *gorm.DB, entitiesId uint64, count uint64) (*T_hammer_entities, error) {
	data := &T_hammer_entities{
		F_season:      config.Instance().NokConfig.Season,
		F_entities_id: entitiesId,
		F_count:       count,
	}

	if err := Create(db, data); err != nil {
		return nil, err
	}
	return data, nil
}

func Select(db *gorm.DB, fields map[string]interface{}) ([]*T_hammer_entities, error) {
	var (
		txs = []*T_hammer_entities{}
	)

	result := db.Table("t_hammer_entities").Where(fields).Find(&txs)
	if err := result.Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("t_hammer_entities select, err[%s], fields[%v]", err, fields)
		return nil, err
	}
	return txs, nil
}

func SelectHammerEntities(db *gorm.DB, entitiesId uint64) (*T_hammer_entities, error) {
	sFields := map[string]interface{}{
		"F_season":      config.Instance().NokConfig.Season,
		"F_entities_id": entitiesId,
	}

	datas, err := Select(db, sFields)
	if err != nil {
		return nil, err
	} else if len(datas) == 0 {
		return nil, nil
	}
	return datas[0], nil
}

func Update(db *gorm.DB, sFields map[string]interface{}, uFields map[string]interface{}) error {
	result := db.Table("t_hammer_entities").Where(sFields).Updates(uFields)
	if err := result.Error; err != nil {
		log.Errorf("t_hammer_entities SelectAndUpdate err:%v, sFields:%+v, uFields:%+v", err, sFields, uFields)
	}
	return nil
}

func AddCount(db *gorm.DB, entitiesId uint64, count uint64) error {
	updatedCount := gorm.Expr(`"F_count" + ?`, count)
	res := db.Table("t_hammer_entities").Where(`"F_season" = ? AND "F_entities_id" = ?`, config.Instance().NokConfig.Season, entitiesId).Update("F_count", updatedCount)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func GetLeaderboard(db *gorm.DB, limit int, offset int) ([]*T_hammer_entities, error) {
	var datas = []*T_hammer_entities{}
	res := db.Table("t_hammer_entities").Where(`"F_season" = ?`, config.Instance().NokConfig.Season).Order(`"F_count" desc`).Offset(offset).Limit(limit).Find(&datas)
	if err := res.Error; err != nil {
		return nil, err
	}

	return datas, nil
}

func CreateOrUpdateCount(db *gorm.DB, entitiesId uint64, count uint64) (*T_hammer_entities, error) {
	data, err := SelectHammerEntities(db, entitiesId)
	if err != nil {
		return nil, err
	}

	if data == nil {
		data, err = CreateHammerEntities(db, entitiesId, count)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	data.F_count += count
	if err := AddCount(db, entitiesId, count); err != nil {
		return nil, err
	}
	return data, nil
}
