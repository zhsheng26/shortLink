package main

import "os"

func main() {
	setEnv()
	a := App{}
	sha1 := toSha1("http://www.baidu.com")
	println(sha1)
	a.Initialize(getEnv())
	a.Run(":8080")
}

func setEnv() {
	_ = os.Setenv("RedisAddr", "192.168.1.2:18160")
	_ = os.Setenv("RedisPwd", "")
	_ = os.Setenv("RedisDb", "3")
}
