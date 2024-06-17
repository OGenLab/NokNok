package t_player

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"lightspeed.2dao3.com/nokserver/common/boost"
	"lightspeed.2dao3.com/nokserver/common/entities"
	"lightspeed.2dao3.com/nokserver/common/task"
	"lightspeed.2dao3.com/nokserver/config"
	"lightspeed.2dao3.com/nokserver/internal/protocol"
	"lightspeed.2dao3.com/nokserver/model/t_user_hammer"

	"qoobing.com/gomod/log"
)

var (
	ErrAccountNotExist = errors.New("account not exist")
	ErrMaxBoost        = errors.New("max boost")
)

type T_player struct {
	F_uid                int64     `gorm:"column:F_uid"`
	F_first_name         string    `gorm:"column:F_first_name"`
	F_last_name          string    `gorm:"column:F_last_name"`
	F_is_premium         bool      `gorm:"column:F_is_premium"`
	F_stamina            int       `gorm:"column:F_stamina"`
	F_stamina_fresh_time time.Time `gorm:"column:F_stamina_fresh_time"`
	F_boost              int       `gorm:"column:F_boost"`
	F_coins              uint64    `gorm:"column:F_coins"`
	F_hit_count          int       `gorm:"column:F_hit_count"`
	F_last_hit_time      time.Time `gorm:"column:F_last_hit_time"`
	F_create_time        time.Time `gorm:"column:F_create_time;autoCreateTime"` // 首次创建时间
	F_modify_time        time.Time `gorm:"column:F_modify_time;autoUpdateTime"` // 最后更新时间
}

func (t *T_player) TableName() string {
	return "t_player"
}

func NewPlayer(uid int64, firstName, lastName string, isPremium bool, coins uint64) *T_player {
	return &T_player{
		F_uid:                uid,
		F_first_name:         firstName,
		F_last_name:          lastName,
		F_stamina:            0,
		F_stamina_fresh_time: time.Now(),
		F_boost:              1,
		F_coins:              coins,
		F_hit_count:          0,
	}
}

func Create(db *gorm.DB, tx *T_player) (err error) {
	if err := db.Table("t_player").Create(tx).Error; err != nil {
		log.Errorf("t_player create err, tx: %+v, err: %v", tx, err)
		return err
	}
	return nil
}

func CreatePlayer(db *gorm.DB, uid int64, firstName, lastName string, isPremium bool, coins uint64) (*T_player, error) {
	player := NewPlayer(uid, firstName, lastName, isPremium, coins)
	if err := Create(db, player); err != nil {
		return nil, err
	}

	return player, nil
}

func Select(db *gorm.DB, fields map[string]interface{}) ([]*T_player, error) {
	var (
		txs = []*T_player{}
	)

	result := db.Table("t_player").Where(fields).Find(&txs)
	if err := result.Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("t_player select, err[%s], fields[%v]", err, fields)
		return nil, err
	}
	return txs, nil
}

func Update(db *gorm.DB, sFields map[string]interface{}, uFields map[string]interface{}) error {
	result := db.Table("t_player").Where(sFields).Updates(uFields)
	if err := result.Error; err != nil {
		log.Errorf("t_player SelectAndUpdate err:%v, sFields:%+v, uFields:%+v", err, sFields, uFields)
	}
	return nil
}

func SelectPlayer(db *gorm.DB, uid int64) (*T_player, error) {
	data, err := Select(db, map[string]interface{}{"F_uid": uid})
	if err != nil {
		return nil, err
	} else if len(data) == 0 {
		return nil, ErrAccountNotExist
	}
	return data[0], err
}

func AddCoins(db *gorm.DB, uid int64, coins uint64) error {
	updatedCoins := gorm.Expr(`"F_coins" + ?`, coins)
	res := db.Table("t_player").Where(`"F_uid" = ?`, uid).Update("F_coins", updatedCoins)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func NextBoost(curBoost int) (uint64, error) {
	boostConfig := boost.BoostSet[curBoost+1]
	if boostConfig == nil {
		return 0, ErrMaxBoost
	}
	return boostConfig.F_needed_coins, nil
}

func UpdateHitCount(db *gorm.DB, uid int64, hitCount int) error {
	res := db.Table("t_player").Where(`"F_uid" = ?`, uid).Updates(map[string]interface{}{
		"F_hit_count":     gorm.Expr(`"F_hit_count" + ?`, hitCount),
		"F_last_hit_time": time.Now(),
	})

	return res.Error
}

func UpdateHammerGopher(db *gorm.DB, uid int64, coins uint64, consumeStamina int, hitCount int) error {
	res := db.Table("t_player").Where(`"F_uid" = ?`, uid).Updates(map[string]interface{}{
		"F_coins":         gorm.Expr(`"F_coins" + ?`, coins),
		"F_stamina":       gorm.Expr(`"F_stamina" + ?`, consumeStamina),
		"F_hit_count":     gorm.Expr(`"F_hit_count" + ?`, hitCount),
		"F_last_hit_time": time.Now(),
	})

	return res.Error
}

func UpdateNextBoost(db *gorm.DB, uid int64, neededCoins uint64) error {
	res := db.Table("t_player").Where(`"F_uid" = ? AND "F_coins" >= ?`, uid, neededCoins).Updates(map[string]interface{}{
		"F_coins": gorm.Expr(`"F_coins" - ?`, neededCoins),
		"F_boost": gorm.Expr(`"F_boost" + ?`, 1),
	})
	if res.Error != nil || res.RowsAffected == 0 {
		return errors.New("update next boost")
	}
	return nil
}

func SelectOrCreate(db *gorm.DB, uid int64, tgUser *protocol.WebAppUser) (*T_player, error) {
	allSet := task.TaskSet.AllSet()
	player, playerErr := SelectPlayer(db, uid)
	if playerErr != nil && !errors.Is(playerErr, ErrAccountNotExist) {
		return nil, playerErr
	}
	if playerErr == nil {
		return player, nil
	}

	player, err := CreatePlayer(db, uid, tgUser.FirstName, tgUser.LastName, tgUser.IsPremium, allSet[task.BASIC_TASK_TYPE][task.BECOME_BOT_USER_ID].F_reward)
	if err != nil {
		return nil, err
	}

	_, err = t_user_hammer.CreateUserHammer(db, uid, entities.WOOD_HAMMER_ID, 1)
	if err != nil {
		return nil, err
	}

	return player, playerErr
}

func RefreshPlayer(db *gorm.DB) error {
	// 获取当前时间
	currentTime := time.Now()
	// 获取俄罗斯时区
	russianTimeZone, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return err
	}

	// 获取俄罗斯早上 8 点的时间
	russianEightAM := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), config.Instance().NokConfig.RefreshTime, 0, 0, 0, russianTimeZone)
	uFields := map[string]interface{}{
		"F_stamina":            0,
		"F_stamina_fresh_time": russianEightAM,
	}

	res := db.Where(`"F_stamina_fresh_time" < ?`, russianEightAM).Updates(uFields)
	return res.Error
}
