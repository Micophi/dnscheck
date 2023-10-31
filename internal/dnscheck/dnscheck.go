package dnscheck

import (
	"fmt"
	"time"

	"github.com/miekg/dns"
)

var sinkholeIp = "0.0.0.0"
var adguardRedirectIp = "94.140.14.33"
var ciraRedirectIps = []string{"75.2.78.236", "99.83.179.4", "99.83.178.7", "75.2.110.227"}
var blockedIpAnswers = append(ciraRedirectIps, sinkholeIp, adguardRedirectIp)

var maximumRetries = 20

func IsDomainBlocked(domain string, dnsServer string) (bool, time.Duration, int, bool) {
	msg := new(dns.Msg)
	msg.SetQuestion(fmt.Sprintf("%s.", domain), dns.TypeA)

	client := dns.Client{}

	retries := 0
	var err error
	for err == nil {
		in, rtt, err := client.Exchange(msg, fmt.Sprintf("%s:53", dnsServer))

		if err != nil {
			retries++
		} else if in.Answer == nil || len(in.Answer) == 0 {
			return true, rtt, retries, false
		} else {
			var ip string
			for _, element := range in.Answer {
				if _, ok := element.(*dns.A); ok {
					ip = element.(*dns.A).A.String()
					break
				}
			}
			for _, blocked := range blockedIpAnswers {
				if ip == blocked {
					return true, rtt, retries, false
				}
			}

			return false, rtt, retries, false
		}
		if retries == maximumRetries {
			return false, 0, retries, true
		}
	}
	return false, 0, retries, false
}
