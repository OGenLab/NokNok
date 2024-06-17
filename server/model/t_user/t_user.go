package t_user

import (
	"errors"
	"time"

	"github.com/speps/go-hashids"
	"gorm.io/gorm"
	"lightspeed.2dao3.com/nokserver/internal/protocol"

	"qoobing.com/gomod/log"
)

var (
	ErrAccountNotExist = errors.New("account not exist")
)

type T_user struct {
	F_uid                   int64     `gorm:"column:F_uid"`
	F_first_name            string    `gorm:"column:F_first_name"`
	F_last_name             string    `gorm:"column:F_last_name"`
	F_img_uri               string    `gorm:"column:F_img_uri"`
	F_is_premium            bool      `gorm:"column:F_is_premium"`
	F_invited_count         int       `gorm:"column:F_invited_count"`
	F_invited_premium_count int       `gorm:"column:F_invited_premium_count"`
	F_channel               string    `gorm:"column:F_channel"`
	F_invite_uid            int64     `gorm:"column:F_invite_uid"`
	F_referral_code         string    `gorm:"column:F_referral_code"`
	F_create_time           time.Time `gorm:"column:F_create_time;autoCreateTime"` // 首次创建时间
	F_modify_time           time.Time `gorm:"column:F_modify_time;autoUpdateTime"` // 最后更新时间
}

func NewUser(uid int64, firstName, lastName, imgURI string, isPremium bool, channel string, inviteUID int64) *T_user {
	return &T_user{
		F_uid:           uid,
		F_first_name:    firstName,
		F_last_name:     lastName,
		F_img_uri:       imgURI,
		F_is_premium:    isPremium,
		F_channel:       channel,
		F_invite_uid:    inviteUID,
		F_referral_code: ID2ReferralCode(uid),
	}
}

func (t *T_user) TableName() string {
	return "t_user"
}

func Create(db *gorm.DB, tx *T_user) (err error) {
	if err := db.Table("t_user").Create(tx).Error; err != nil {
		log.Errorf("t_user create err, tx: %+v, err: %v", tx, err)
		return err
	}
	return nil
}

func CreateUser(db *gorm.DB, uid int64, firstName, lastName, imgURI string, isPremium bool, channel string, inviteUID int64) (*T_user, error) {
	user := NewUser(uid, firstName, lastName, imgURI, isPremium, channel, inviteUID)
	if err := Create(db, user); err != nil {
		return nil, err
	}

	return user, nil
}

func Select(db *gorm.DB, fields map[string]interface{}) ([]*T_user, error) {
	var (
		txs = []*T_user{}
	)

	result := db.Table("t_user").Where(fields).Find(&txs)
	if err := result.Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("t_user select, err[%s], fields[%v]", err, fields)
		return nil, err
	}
	return txs, nil
}

func Update(db *gorm.DB, sFields map[string]interface{}, uFields map[string]interface{}) error {
	result := db.Table("t_user").Where(sFields).Updates(uFields)
	if err := result.Error; err != nil {
		log.Errorf("t_user SelectAndUpdate err:%v, sFields:%+v, uFields:%+v", err, sFields, uFields)
	}
	return nil
}

func UpdateInvitedCount(db *gorm.DB, uid int64, isPremium bool) error {
	sFields := map[string]interface{}{
		"F_uid": uid,
	}
	uFields := map[string]interface{}{}
	if isPremium {
		uFields["F_invited_premium_count"] = gorm.Expr(`"F_invited_premium_count" + ?`, 1)
	} else {
		uFields["F_invited_count"] = gorm.Expr(`"F_invited_count" + ?`, 1)
	}

	return Update(db, sFields, uFields)
}

func SelectUser(db *gorm.DB, uid int64) (*T_user, error) {
	data, err := Select(db, map[string]interface{}{"F_uid": uid})
	if err != nil {
		return nil, err
	} else if len(data) == 0 {
		return nil, ErrAccountNotExist
	}
	return data[0], err
}

func SelectUserByRefferCode(db *gorm.DB, code string) (*T_user, error) {
	sFileds := map[string]interface{}{"F_referral_code": code}
	users, err := Select(db, sFileds)
	if err != nil {
		return nil, err
	} else if len(users) == 0 {
		return nil, ErrAccountNotExist
	}
	return users[0], nil
}

func ID2ReferralCode(id int64) string {
	h, _ := hashids.New()
	encodedID, err := h.EncodeInt64([]int64{id})
	if err != nil {
		log.Errorf("ID2ReferralCode error: %+v", err)
		return ""
	}

	return encodedID
}

func SelectOrCreate(db *gorm.DB, uid int64, invitingUid int64, channel string, tgUser *protocol.WebAppUser) (*T_user, error) {
	user, userErr := SelectUser(db, uid)
	if userErr != nil && !errors.Is(userErr, ErrAccountNotExist) {
		return nil, userErr
	} else if userErr == nil {
		return user, nil
	}

	user, err := CreateUser(db, uid, tgUser.FirstName, tgUser.LastName, tgUser.PhotoURL, tgUser.IsPremium, channel, invitingUid)
	if err != nil {
		return nil, err
	}

	return user, userErr
}
