package main

import (
	"flag"
	"net/http"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/miekg/dns"
)

func checkErr(err error) {
	if err != nil {
		log.Fatalf("Error : %s", err)
	}
}

func main() {
	configPatch := flag.String("config", "config.json", "the patch to the config file")
	flag.Parse()

	cli := newDeezerClient(*configPatch)

	if cli.config.Debug {
		log.SetLevel(log.DebugLevel)
	}
	log.SetOutput(os.Stdout)

	http.HandleFunc("/callback", cli.doAuth)

	log.Infof("Oauth URL : %s", cli.getOauthURL())

	// HTTP for Deezer Auth
	go func() {
		log.Infof("[HTTP] Started")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Redis Ping
	log.Infof("[REDIS] Started connectivity check")
	err := cli.RedisClient.Set(ctx, "_ping", "success", 1*time.Second).Err()
	checkErr(err)
	log.Infof("[REDIS] Connected to REDIS")

	// DNS Listener
	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) { HandleDNSRequest(w, r, cli) })
	log.Infof("[DNS] Started")
	dnsSrv := &dns.Server{Addr: ":" + strconv.Itoa(cli.config.DNS.Port), Net: cli.config.DNS.Protocol}
	if err := dnsSrv.ListenAndServe(); err != nil {
		log.Errorf("[DNS] Server failed : %s\n", err.Error())
	}
}
