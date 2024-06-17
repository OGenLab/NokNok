package cache

import "github.com/gomodule/redigo/redis"

func HScan(rds redis.Conn, params []interface{}) (interface{}, error) {
	val, err := rds.Do("HSCAN", params...)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func HGet(rds redis.Conn, params []interface{}) (interface{}, error) {
	val, err := rds.Do("HGET", params...)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func HDel(rds redis.Conn, params []interface{}) error {
	_, err := rds.Do("HDEL", params...)
	if err != nil {
		return err
	}
	return nil
}

func HSet(rds redis.Conn, params []interface{}) error {
	if _, err := rds.Do("HSET", params...); err != nil {
		return err
	}
	return nil
}

func HGetAll(rds redis.Conn, params []interface{}) (interface{}, error) {
	res, err := rds.Do("HGETALL", params...)
	if err != nil {
		return nil, err
	}
	return res, nil
}
