package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (c *DeezerClient) getOauthURL() string {
	url := fmt.Sprintf("https://connect.deezer.com/oauth/auth.php?app_id=%v&redirect_uri=%s&perms=listening_history,offline_access", c.config.AppID, c.config.CallbackURL)
	return (url)
}

func (c *DeezerClient) doAuth(w http.ResponseWriter, r *http.Request) {
	authCode, check := r.URL.Query()["code"]

	if !check || len(authCode[0]) < 1 {
		log.Debug("[HTTP] [AUTH] Missing `code` callback.")
		fmt.Fprintf(w, "ERROR - Missing `code` callback\n")
		return
	}

	url := fmt.Sprintf("https://connect.deezer.com/oauth/access_token.php?app_id=%v&secret=%s&code=%s", c.config.AppID, c.config.AppSecret, authCode[0])
	query, err := http.Get(url)
	checkErr(err)
	body, _ := ioutil.ReadAll(query.Body)

	//Not clean stuff going on here...
	token := strings.Split(string(body), "&")[0]
	token = strings.Replace(token, "access_token=", "", 1)
	if token == "wrong code" {
		fmt.Fprintf(w, "ERROR - Invalid callback code\n")
		return
	}
	uid := c.getClient(token).ID

	err = c.RedisClient.Set(ctx, string(uid), token, 0).Err()
	checkErr(err)

	log.WithFields(log.Fields{
		"uid":   uid,
		"token": token[0:5],
	}).Debug("[REDIS] Set in DB")

	fmt.Fprintf(w, "SUCCESS - Your account is now connected !\n")
	return
}
