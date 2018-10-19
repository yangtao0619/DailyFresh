package controllers

import (
	"errors"
	"github.com/gomodule/redigo/redis"
)

func GetRedisConnect() (conn redis.Conn, err error) {
	conn, err = redis.Dial("tcp", "172.17.0.9:6379")
	if err != nil {
		err = errors.New("redis连接错误")
	}
	return conn, err
}
