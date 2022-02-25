package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type Config struct {
	Debug       bool   `json:"debug"`
	AppID       int    `json:"appID"`
	AppSecret   string `json:"appSecret"`
	BaseAPI     string `json:"baseAPI"`
	CallbackURL string `json:"callbackURL"`
	RedisDB     struct {
		Addr     string `json:"address"`
		Password string `json:"password"`
		Db       int    `json:"db"`
	} `json:"redisDB"`
	DNS struct {
		Port     int    `json:"port"`
		Protocol string `json:"proto"`
		Domain   string `json:"domain"`
	} `json:"dns"`
}

type DeezerClient struct {
	config      Config
	HTTPClient  *http.Client
	RedisClient *redis.Client
}

func readConf(cPath string) Config {
	confFile, err := ioutil.ReadFile(cPath)
	checkErr(err)
	conf := Config{}
	err = json.Unmarshal([]byte(confFile), &conf)
	log.Debugf("[APP] Loaded config at %s", confFile)
	return conf
}

func newDeezerClient(cPath string) *DeezerClient {
	c := readConf(cPath)
	return &DeezerClient{
		config: c,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
		RedisClient: redis.NewClient(&redis.Options{
			Addr:     c.RedisDB.Addr,
			Password: c.RedisDB.Password, // no password set
			DB:       c.RedisDB.Db,       // use default DB
		}),
	}
}

func (c *DeezerClient) getQuery(q string, t string) []byte {
	url := fmt.Sprintf("%s/%s?access_token=%s", c.config.BaseAPI, q, t)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	checkErr(err)

	log.WithFields(log.Fields{
		"url": strings.Replace(url, t, "", 1),
	}).Debug("[API] getQuery")

	req.Header.Set("User-Agent", "1035-deezer-cli")
	res, err := c.HTTPClient.Do(req)
	checkErr(err)
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	checkErr(err)

	return body
}

func (c *DeezerClient) getClient(token string) deezerAccount {
	q := c.getQuery("/user/me/", token)
	u := deezerAccount{}
	err := json.Unmarshal(q, &u)
	checkErr(err)
	log.WithFields(log.Fields{
		"uid": u.ID,
	}).Debug("[API] getClient")
	return u
}

func (c *DeezerClient) getHistory(uid string) deezerHistory {
	log.WithFields(log.Fields{
		"uid": uid,
	}).Debug("[API] getHistory")
	t, err := c.RedisClient.Get(ctx, uid).Result()
	if err == redis.Nil {
		log.WithFields(log.Fields{
			"uid": uid,
		}).Debug("[REDIS] getHistory - Empty REDIS")
		return deezerHistory{}
	} else {
		q := c.getQuery("/user/me/history", t)
		hist := deezerHistory{}
		err := json.Unmarshal(q, &hist)
		checkErr(err)
		return hist
	}
}

func (c *DeezerClient) lastTitle(uid string) deezerTitle {
	hist := c.getHistory(uid).Data

	if len(hist) < 1 {
		return deezerTitle{}
	} else {
		lastMusic := hist[0]
		return lastMusic
	}
}
