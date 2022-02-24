package main

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/miekg/dns"
)

type handler struct{}

func genTxtRecord(m *dns.Msg, q dns.Question, value string) {
	rr, err := dns.NewRR(fmt.Sprintf("%s %v %s \"%s\"", q.Name, 30, "TXT", value))
	if err == nil {
		m.Answer = append(m.Answer, rr)
	}
}

func HandleDNSRequest(w dns.ResponseWriter, r *dns.Msg, c *DeezerClient) {

	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	if r.Opcode == dns.OpcodeQuery {
		for _, q := range m.Question {

			if strings.HasSuffix(q.Name, c.config.DNS.Domain) && q.Qtype == dns.TypeTXT {
				log.WithFields(log.Fields{
					"fqdn": q.Name,
					"type": q.Qtype,
				}).Debug("[DNS] Got query")

				for _, q := range m.Question {
					uid := strings.Replace(q.Name, "."+c.config.DNS.Domain, "", 1)
					lastPlayed := c.lastTitle(uid)

					if (deezerTitle{}) != lastPlayed {
						log.WithFields(log.Fields{
							"fqdn":  q.Name,
							"title": lastPlayed.Title,
						}).Debug("[DNS] Record found")

						genTxtRecord(m, q, "Last played song : "+lastPlayed.Title)
						genTxtRecord(m, q, "Author : "+lastPlayed.Artist.Name)
					} else {
						log.WithFields(log.Fields{
							"fqdn": q.Name,
						}).Debug("[DNS] No recound found")

						genTxtRecord(m, q, "Can't get this user playing song.")
						genTxtRecord(m, q, "User may not exist.")
						genTxtRecord(m, q, "If that's you, connect the app to your Deezer account :")
						genTxtRecord(m, q, c.getOauthURL())
					}

				}
				w.WriteMsg(m)
			} else {
				log.WithFields(log.Fields{
					"fqdn": q.Name,
					"type": q.Qtype,
				}).Debug("[DNS] Got not authorized query")

				m := new(dns.Msg)
				m.SetRcode(r, dns.RcodeRefused)
				w.WriteMsg(m)
			}
		}
	}

	w.WriteMsg(m)
}
