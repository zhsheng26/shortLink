package main

import "os"

func main() {
	setEnv()
	a := App{}
	a.Initialize(getEnv())
	a.Run(":8080")
}

func setEnv() {
	_ = os.Setenv("RedisAddr", "192.168.1.2:18160")
	_ = os.Setenv("RedisPwd", "")
	_ = os.Setenv("RedisDb", "3")
}
