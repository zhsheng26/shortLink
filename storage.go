package main

type Storage interface {
	//将长地址转为短地址
	Shorten(url string, exp int64) (string, error)
	ShortLinkInfo(eid string) (interface{}, error)
	UnShorten(eid string) (string, error)
}
