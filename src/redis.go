package main

import (
	"github.com/garyburd/redigo/redis"
	"github.com/sendwithus/lib-go"
	"strings"
	"time"
)

type ConnectionInformation struct {
	Password string
	Host     string
}

func NewPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			connectionInfo := ParseRedisConnectionString(swu.GetEnvVariable("REDIS_URL", true))
			return redis.Dial("tcp", connectionInfo.Host, redis.DialPassword(connectionInfo.Password))
		},
	}
}

func ParseRedisConnectionString(connectionString string) ConnectionInformation {
	connectionString = strings.Replace(connectionString, "redis://", "", 1)
	println("REDIS_URL: ", connectionString)

	parts := strings.Split(connectionString, "@")
	userBits := strings.Split(parts[0], ":")

	info := ConnectionInformation{
		Host:     parts[1],
		Password: userBits[1],
	}

	return info
}
