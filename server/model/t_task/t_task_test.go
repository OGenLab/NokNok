package t_task

// import (
// 	"testing"
// 	"time"

// 	"lightspeed.2dao3.com/nokserver/config"
// 	"lightspeed.2dao3.com/nokserver/model"
// 	"lightspeed.2dao3.com/nokserver/model/t_entities_properties"
// )

// func TestLoginGameTaskForDaily(t *testing.T) {
// 	var (
// 		m = model.NewModelDefault()
// 	)
// 	defer m.Close()

// 	task := &T_task{
// 		F_task_id:    1,
// 		F_season:     config.Instance().NokConfig.Season,
// 		F_type:       DAILY_TASK_TYPE,
// 		F_threshold:  1,
// 		F_reward:     50,
// 		F_period:     uint64(time.Duration(24 * time.Hour)),
// 		F_start_time: time.Now(),
// 		F_end_time:   time.Date(2050, time.December, 31, 23, 59, 59, 0, time.UTC),
// 	}
// 	if err := Create(m.DB, task); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestInviteTaskForDaily(t *testing.T) {
// 	var (
// 		m = model.NewModelDefault()
// 	)
// 	defer m.Close()

// 	task := &T_task{
// 		F_task_id:    2,
// 		F_season:     config.Instance().NokConfig.Season,
// 		F_type:       DAILY_TASK_TYPE,
// 		F_threshold:  1,
// 		F_reward:     100,
// 		F_period:     uint64(time.Duration(24 * time.Hour)),
// 		F_start_time: time.Now(),
// 		F_end_time:   time.Date(2050, time.December, 31, 23, 59, 59, 0, time.UTC),
// 	}
// 	if err := Create(m.DB, task); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestNokBotUserTaskForBasic(t *testing.T) {
// 	var (
// 		m = model.NewModelDefault()
// 	)
// 	defer m.Close()

// 	task := &T_task{
// 		F_task_id:    3,
// 		F_season:     config.Instance().NokConfig.Season,
// 		F_type:       BASIC_TASK_TYPE,
// 		F_threshold:  1,
// 		F_reward:     100,
// 		F_period:     0,
// 		F_start_time: time.Now(),
// 		F_end_time:   time.Date(2050, time.December, 31, 23, 59, 59, 0, time.UTC),
// 	}
// 	if err := Create(m.DB, task); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestJoinNokNewsChannelForBasic(t *testing.T) {
// 	var (
// 		m = model.NewModelDefault()
// 	)
// 	defer m.Close()

// 	task := &T_task{
// 		F_task_id:    4,
// 		F_season:     config.Instance().NokConfig.Season,
// 		F_type:       BASIC_TASK_TYPE,
// 		F_threshold:  1,
// 		F_reward:     200,
// 		F_period:     0,
// 		F_start_time: time.Now(),
// 		F_end_time:   time.Date(2050, time.December, 31, 23, 59, 59, 0, time.UTC),
// 	}
// 	if err := Create(m.DB, task); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestFollonOnXForBasic(t *testing.T) {
// 	var (
// 		m = model.NewModelDefault()
// 	)
// 	defer m.Close()

// 	task := &T_task{
// 		F_task_id:    5,
// 		F_season:     config.Instance().NokConfig.Season,
// 		F_type:       BASIC_TASK_TYPE,
// 		F_threshold:  1,
// 		F_reward:     200,
// 		F_period:     0,
// 		F_start_time: time.Now(),
// 		F_end_time:   time.Date(2050, time.December, 31, 23, 59, 59, 0, time.UTC),
// 	}
// 	if err := Create(m.DB, task); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestJoinNokGroupForBasic(t *testing.T) {
// 	var (
// 		m = model.NewModelDefault()
// 	)
// 	defer m.Close()

// 	task := &T_task{
// 		F_task_id:    6,
// 		F_season:     config.Instance().NokConfig.Season,
// 		F_type:       BASIC_TASK_TYPE,
// 		F_threshold:  1,
// 		F_reward:     200,
// 		F_period:     0,
// 		F_start_time: time.Now(),
// 		F_end_time:   time.Date(2050, time.December, 31, 23, 59, 59, 0, time.UTC),
// 	}
// 	if err := Create(m.DB, task); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestConnectWalletForBasic(t *testing.T) {
// 	var (
// 		m = model.NewModelDefault()
// 	)
// 	defer m.Close()

// 	task := &T_task{
// 		F_task_id:    7,
// 		F_season:     config.Instance().NokConfig.Season,
// 		F_type:       BASIC_TASK_TYPE,
// 		F_threshold:  1,
// 		F_reward:     200,
// 		F_period:     0,
// 		F_start_time: time.Now(),
// 		F_end_time:   time.Date(2050, time.December, 31, 23, 59, 59, 0, time.UTC),
// 	}
// 	if err := Create(m.DB, task); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestDrawTaskForDaily(t *testing.T) {
// 	var (
// 		m = model.NewModelDefault()
// 	)
// 	defer m.Close()

// 	task := &T_task{
// 		F_task_id:    8,
// 		F_season:     config.Instance().NokConfig.Season,
// 		F_type:       DRAW_TASK_TYPE,
// 		F_threshold:  3,
// 		F_reward:     t_entities_properties.HAMMER_TYPE,
// 		F_period:     0,
// 		F_start_time: time.Now(),
// 		F_end_time:   time.Date(2050, time.December, 31, 23, 59, 59, 0, time.UTC),
// 	}
// 	if err := Create(m.DB, task); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestInviteCommonTaskForDaily(t *testing.T) {
// 	var (
// 		m = model.NewModelDefault()
// 	)
// 	defer m.Close()

// 	task := &T_task{
// 		F_task_id:    INVITE_COMMON_ID,
// 		F_season:     config.Instance().NokConfig.Season,
// 		F_type:       BASIC_TASK_TYPE,
// 		F_threshold:  1,
// 		F_reward:     100,
// 		F_period:     0,
// 		F_start_time: time.Now(),
// 		F_end_time:   time.Date(2050, time.December, 31, 23, 59, 59, 0, time.UTC),
// 	}
// 	if err := Create(m.DB, task); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestInvitePremiumTaskForDaily(t *testing.T) {
// 	var (
// 		m = model.NewModelDefault()
// 	)
// 	defer m.Close()

// 	task := &T_task{
// 		F_task_id:    INVITE_PREMIUM_ID,
// 		F_season:     config.Instance().NokConfig.Season,
// 		F_type:       BASIC_TASK_TYPE,
// 		F_threshold:  1,
// 		F_reward:     500,
// 		F_period:     0,
// 		F_start_time: time.Now(),
// 		F_end_time:   time.Date(2050, time.December, 31, 23, 59, 59, 0, time.UTC),
// 	}
// 	if err := Create(m.DB, task); err != nil {
// 		t.Fatal(err)
// 	}
// }
