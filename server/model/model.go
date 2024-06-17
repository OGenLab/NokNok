package model

import (
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"lightspeed.2dao3.com/nokserver/config"
	"qoobing.com/gomod/log"
	"qoobing.com/gomod/redis/sentinel"
)

type Model struct {
	DB         *gorm.DB
	DbTxStatus DbTxStatus
	Redis      redis.Conn
	RedisPool  *redis.Pool
}

func NewModel() *Model {
	return NewModelDefault()
}

type DbTxStatus int

const (
	DBTX_STATUS_TX_NONE    DbTxStatus = 0
	DBTX_STATUS_TX_DOING   DbTxStatus = 1
	DBTX_STATUS_TX_SUCCESS DbTxStatus = 2
	DBTX_STATUS_TX_FAILED  DbTxStatus = 3
)

func (m *Model) Close() {
	if m.Redis != nil {
		m.Redis.Close()
	}
	if m.DB == nil {
		m.DbTxStatus = DBTX_STATUS_TX_NONE
	} else if m.DbTxStatus == DBTX_STATUS_TX_DOING {
		m.DB.Rollback()
		m.DbTxStatus = DBTX_STATUS_TX_FAILED
	}
}

func (m *Model) Begin() {
	if m.DB == nil {
		panic("unreachable code, m.DB is uninitialized")
	} else if m.DbTxStatus != DBTX_STATUS_TX_NONE {
		panic("unreachable code, begin transaction towice???")
	} else {
		m.DbTxStatus = DBTX_STATUS_TX_DOING
		m.DB = m.DB.Begin()
	}
}

func (m *Model) Commit() error {
	if m.DB == nil {
		panic("unreachable code, m.DB is uninitialized")
	} else if m.DbTxStatus != DBTX_STATUS_TX_DOING {
		return nil
	} else if err := m.DB.Commit().Error; err != nil {
		m.DbTxStatus = DBTX_STATUS_TX_FAILED
		return err
	} else {
		m.DbTxStatus = DBTX_STATUS_TX_SUCCESS
		return nil
	}
}

func NewModelDefault() *Model {
	return NewModelWithOption(OptOpenDefaultDatabase, OptOpenDefaultRedis)
}

func NewModelDefaultRedis() *Model {
	return NewModelWithOption(OptOpenDefaultRedis)
}

func NewModelDefaultDatabase() *Model {
	return NewModelWithOption(OptOpenDefaultDatabase)
}

type Option func(*Model)

func NewModelWithOption(options ...Option) *Model {
	n := 0
	m := Model{}
	log.Debugf("Start NewModelWithOption...")
	for _, option := range options {
		n++
		option(&m)
	}
	log.Debugf("Finish NewModelWithOption(with %d options)", n)
	return &m
}

// OptOpenDefaultDatabase option function for open default database
func OptOpenDefaultDatabase(m *Model) {
	if defaultDB != nil {
		m.DB = defaultDB
	} else if db, err := gorm.Open(
		postgres.New(postgres.Config{
			DSN:                  defaultDsn,
			PreferSimpleProtocol: true,
		}), &gorm.Config{
			NowFunc: func() time.Time {
				russianTimeZone, _ := time.LoadLocation("Europe/Moscow")
				return time.Now().In(russianTimeZone)
			},
		}); err != nil {
		panic("DATABASE_OPEN_ERROR")
	} else {
		if defaultDbDebug {
			m.DB = db.Debug()
		} else {
			m.DB = db
		}
		if sqlDB, err := m.DB.DB(); err != nil {
			panic("DATABASE_OPEN_ERROR")
		} else {
			sqlDB.SetMaxIdleConns(3)
			sqlDB.SetMaxOpenConns(10)
			sqlDB.SetConnMaxLifetime(time.Minute)
		}
		defaultDB = m.DB
	}
	//TODO: check cocurrent
	m.DB = m.DB.Session(&gorm.Session{QueryFields: true})
	log.Debugf("Opt for open default database done")
	return
}

// OptOpenDefaultRedisSentinelPool option function for open default redis sentinel
func OptOpenDefaultRedis(m *Model) {
	if defaultRedis != nil {
		m.Redis = defaultRedis.Get()
		m.RedisPool = defaultRedis
	} else if pool := sentinel.NewPool(defaultRds); pool == nil {
		panic("REDIS_ERROR")
	} else {
		defaultRedis = pool
		m.Redis = defaultRedis.Get()
		m.RedisPool = defaultRedis
	}

	log.Debugf("Opt for open default redis done")
	return
}

var (
	defaultDB         *gorm.DB
	defaultDbDebug    bool
	defaultDsn        string
	defaultRedis      *redis.Pool
	defaultRedisDebug bool
	defaultRds        sentinel.Config
)

func init() {
	// database init
	coreDb := config.Instance().NokDb
	arrConfStr := []string{
		"host=" + coreDb.Host,
		"port=" + strconv.Itoa(coreDb.Port),
		"user=" + coreDb.Username,
		"password=" + coreDb.Password,
		"dbname=" + coreDb.Dbname,
		coreDb.ExtraParameters,
	}
	defaultDsn = strings.Join(arrConfStr, " ")
	defaultDbDebug = coreDb.Debug

	// redis init
	coreRedis := config.Instance().NokRedis
	defaultRds = coreRedis
	defaultRedisDebug = coreRedis.Debug

	// debug
	if defaultDbDebug && len(coreDb.Password) > 5 {
		p := coreDb.Password
		l := len(p)
		secDefaultDsn := strings.Replace(defaultDsn, p[2:l-3], "*****", 1)
		log.Debugf("defaultDsn='%s'", secDefaultDsn)
	}
	if defaultRedisDebug && len(coreRedis.Password) > 5 {
		p := coreRedis.Password
		l := len(p)
		secDefaultRds := defaultRds
		secDefaultRds.Password = strings.Replace(p, p[2:l-3], "*****", 1)
		log.PrintPretty("defaultRds=", secDefaultRds)
	}
}
