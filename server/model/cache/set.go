package cache

import "github.com/gomodule/redigo/redis"

func SAdd(rds redis.Conn, params []interface{}) error {
	_, err := rds.Do("SADD", params...) // SADD your_set_key element1 element2 element3 ...
	if err != nil {
		return err
	}
	return nil
}

func SMember(rds redis.Conn, params []interface{}) (interface{}, error) {
	val, err := rds.Do("SMEMBERS", params...) // SMEMBERS your_set_key
	if err != nil {
		return nil, err
	}
	return val, nil
}

func SDel(rds redis.Conn, params []interface{}) error {
	_, err := rds.Do("DEL", params...) // SMEMBERS your_set_key
	if err != nil {
		return err
	}
	return nil
}
