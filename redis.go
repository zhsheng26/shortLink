package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/mattheath/base62"
	"github.com/speps/go-hashids"
	"time"
)

const (
	UrlIdKey           = "next.url.id"
	ShortLinkKey       = "shortLink:%s:url"
	UrlHashKey         = "urlHash:%s:url"
	ShortLinkDetailKey = "shortLink:%s:detail"
)

type RedisCli struct {
	Cli *redis.Client
}

func (r *RedisCli) Shorten(url string, exp int64) (string, error) {
	urlHash := toHash(url)
	d, err := r.Cli.Get(fmt.Sprintf(UrlHashKey, urlHash)).Result()
	if err == redis.Nil {
		//no exist ,generate new
	} else if err != nil {
		return "", err
	} else {
		if d == "{}" {
			//expiration,nothing to do, generate new
		} else {
			return d, nil
		}
	}
	err = r.Cli.Incr(UrlIdKey).Err()
	if err != nil {
		return "", err
	}
	id, err := r.Cli.Get(UrlIdKey).Int64()
	if err != nil {
		return "", err
	}
	shortLink := base62.EncodeInt64(id)

	err = r.Cli.Set(fmt.Sprintf(ShortLinkKey, shortLink), url,
		time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", err
	}

	err = r.Cli.Set(fmt.Sprintf(UrlHashKey, urlHash), shortLink,
		time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", nil
	}
	detail, err := json.Marshal(&UrlDetail{
		Url:                 url,
		CreateAt:            time.Now().String(),
		ExpirationInMinutes: time.Duration(exp),
	})
	if err != nil {
		return "", err
	}
	err = r.Cli.Set(fmt.Sprintf(ShortLinkDetailKey, shortLink), detail,
		time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", err
	}
	return shortLink, nil
}

func toHash(url string) string {
	hd := hashids.NewData()
	hd.Salt = url
	hd.MinLength = 0
	h, _ := hashids.NewWithData(hd)
	r, _ := h.Encode([]int{45, 434, 1313, 99})
	return r
}

func (r *RedisCli) ShortLinkInfo(shortLink string) (interface{}, error) {
	detail, err := r.Cli.Get(fmt.Sprintf(ShortLinkDetailKey, shortLink)).Result()
	if err == redis.Nil {
		return "", NewNotFindErr(errors.New("unknown short url"))
	} else if err != nil {
		return "", err
	} else {
		return detail, nil
	}
}

func (r *RedisCli) UnShorten(shortLink string) (string, error) {
	url, err := r.Cli.Get(fmt.Sprintf(ShortLinkKey, shortLink)).Result()
	if err == redis.Nil {
		return "", NewNotFindErr(errors.New("unknown short url"))
	} else if err != nil {
		return "", err
	} else {
		return url, nil
	}
}

type UrlDetail struct {
	Url                 string        `json:"url"`
	CreateAt            string        `json:"create_at"`
	ExpirationInMinutes time.Duration `json:"expiration_in_minutes"`
}

func NewRedisCli(addr string, pwd string, db int) *RedisCli {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       db,
	})
	if _, err := client.Ping().Result(); err != nil {
		panic(err)
	}
	return &RedisCli{Cli: client}
}
