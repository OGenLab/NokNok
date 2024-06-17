package cache

import (
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gomodule/redigo/redis"
	"qoobing.com/gomod/log"
)

func GetKey(rds redis.Conn, params []interface{}) (interface{}, error) {
	val, err := rds.Do("GET", params...)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func SetKey(rds redis.Conn, params []interface{}) error {
	if _, err := rds.Do("SET", params...); err != nil {
		return err
	}
	return nil
}

func DelKey(rds redis.Conn, params []interface{}) error {
	if _, err := rds.Do("DEL", params...); err != nil {
		return err
	}
	return nil
}

func NewFaucetKey(cid string, address common.Address) string {
	strs := []string{"faucet", cid, address.Hex()}
	return strings.Join(strs, ":")
}

func DelFaucetKey(rds redis.Conn, cid string, address common.Address) error {
	key := NewFaucetKey(cid, address)
	if err := DelKey(rds, []interface{}{key}); err != nil {
		return err
	}
	return nil
}

func SetFaucetKey(rds redis.Conn, cid string, address common.Address, txHash common.Hash, timestamp int64) error {
	ts := strconv.FormatInt(timestamp, 10)
	value := strings.Join([]string{txHash.Hex(), ts}, ":")
	key := NewFaucetKey(cid, address)

	if err := SetKey(rds, []interface{}{key, value, "EX", 24 * 60 * 60}); err != nil {
		return err
	}
	return nil
}

func GetFaucetKey(rds redis.Conn, cid string, address common.Address) (common.Hash, int64, error) {
	key := NewFaucetKey(cid, address)
	res, err := GetKey(rds, []interface{}{key})
	if err != nil {
		log.Errorf("cache getFaucetKey err: %v", err)
		return common.Hash{}, 0, err
	} else if res == nil {
		return common.Hash{}, 0, nil
	}

	val, err := redis.String(res, err)
	if err != nil {
		log.Errorf("faucet value to string err: %v", err)
		return common.Hash{}, 0, err
	}

	strs := strings.Split(val, ":")
	txHash := common.HexToHash(strs[0])
	timeStamp, _ := strconv.ParseInt(strs[1], 10, 64)
	return txHash, timeStamp, nil
}
