package main

import (
	"log"
	"os"
	"strconv"
)

type Env struct {
	S Storage
}

func getEnv() *Env {
	redisAddr := os.Getenv("RedisAddr")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisPwd := os.Getenv("RedisPwd")
	if redisPwd == "" {
		redisPwd = ""
	}
	redisDb := os.Getenv("RedisDb")
	if redisDb == "" {
		redisDb = "0"
	}
	db, err := strconv.Atoi(redisDb)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("connect to redis (addr:%s password: %s db: %d", redisAddr, redisPwd, db)
	cli := NewRedisCli(redisAddr, redisPwd, db)
	return &Env{S: cli}
}
